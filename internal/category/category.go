package category

import (
	"context"
	"errors"
	"fmt"
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/errs"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	q *queries.Queries
	p *pgxpool.Pool
}

func New(q *queries.Queries, p *pgxpool.Pool) *Service {
	return &Service{q: q, p: p}
}

func (s *Service) Create(ctx context.Context, request DTO.CategoryRequest) (DTO.CategoryResponse, *errs.AppError) {
	if request.ParentId != nil {
		_, appErr := s.Get(ctx, *request.ParentId)
		if appErr != nil {
			if appErr.Code == http.StatusNotFound {
				return DTO.CategoryResponse{}, errs.NotFound(fmt.Errorf("parent category with id=%d not found | %w", *request.ParentId, appErr.Unwrap()))
			}
			return DTO.CategoryResponse{}, appErr
		}
	}

	category, err := s.q.CreateCategory(ctx, queries.CreateCategoryParams{
		Name:     request.Name,
		Slug:     request.Slug,
		ParentID: request.ParentId,
	})
	if err != nil {
		if pgErr, isUniqueViolation := errs.IsUniqueViolation(err); isUniqueViolation {
			return DTO.CategoryResponse{}, errs.UniqueViolation(err, pgErr)
		}
		if pgErr, isForeignKeyViolation := errs.IsForeignKeyViolation(err); isForeignKeyViolation {
			return DTO.CategoryResponse{}, errs.ForeignKeyViolation(err, pgErr)
		}
		return DTO.CategoryResponse{}, errs.Internal(err)
	}

	categoryResponse := DTO.CategoryResponse{
		Id:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		ParentId:  category.ParentID,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
	return categoryResponse, nil
}

func (s *Service) GetAll(ctx context.Context) ([]DTO.CategoryResponse, *errs.AppError) {
	categories, err := s.q.GetAllCategories(ctx)
	if err != nil {
		return nil, errs.Internal(err)
	}

	categoriesResponse := make([]DTO.CategoryResponse, len(categories))
	for i, category := range categories {
		categoriesResponse[i] = DTO.CategoryResponse{
			Id:        category.ID,
			Name:      category.Name,
			Slug:      category.Slug,
			ParentId:  category.ParentID,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		}
	}
	return categoriesResponse, nil
}

func (s *Service) Get(ctx context.Context, id int64) (DTO.CategoryResponse, *errs.AppError) {
	category, err := s.q.GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO.CategoryResponse{}, errs.NotFound(err)
		}

		return DTO.CategoryResponse{}, errs.Internal(err)
	}

	categoryResponse := DTO.CategoryResponse{
		Id:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		ParentId:  category.ParentID,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
	return categoryResponse, nil
}

func (s *Service) Update(ctx context.Context, id int64, request DTO.CategoryRequest) (DTO.CategoryResponse, *errs.AppError) {
	if request.ParentId != nil {
		_, appErr := s.Get(ctx, *request.ParentId)
		if appErr != nil {
			if appErr.Code == http.StatusNotFound {
				return DTO.CategoryResponse{}, errs.NotFound(fmt.Errorf("parent category with id=%d not found | %w", id, appErr.Unwrap()))
			}
			return DTO.CategoryResponse{}, appErr
		}
	}

	category, err := s.q.UpdateCategory(ctx, queries.UpdateCategoryParams{
		ID:       id,
		Name:     request.Name,
		Slug:     request.Slug,
		ParentID: request.ParentId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO.CategoryResponse{}, errs.NotFound(err)
		}
		if pgErr, isUniqueViolation := errs.IsUniqueViolation(err); isUniqueViolation {
			return DTO.CategoryResponse{}, errs.UniqueViolation(err, pgErr)
		}
		if pgErr, isForeignKeyViolation := errs.IsForeignKeyViolation(err); isForeignKeyViolation {
			return DTO.CategoryResponse{}, errs.ForeignKeyViolation(err, pgErr)
		}
		return DTO.CategoryResponse{}, errs.Internal(err)
	}

	categoryResponse := DTO.CategoryResponse{
		Id:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		ParentId:  category.ParentID,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
	return categoryResponse, nil
}

func (s *Service) Delete(ctx context.Context, id int64) (int, *errs.AppError) {
	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	tx, err := s.p.Begin(timeout)
	if err != nil {
		return 0, errs.Internal(err)
	}
	defer tx.Rollback(timeout)

	rows, err := s.q.WithTx(tx).DeleteCategory(timeout, id)
	if err != nil {
		return 0, errs.Internal(err)
	}
	if rows == 0 {
		return int(rows), errs.NotFound(errors.New("category not found"))
	}
	_, err = s.q.WithTx(tx).DeleteProductsByCategoryId(timeout, id)
	if err != nil {
		return 0, errs.Internal(err)
	}

	appErr := s.deleteRecursivelyByParentId(timeout, tx, id)
	if appErr != nil {
		return 0, appErr
	}

	err = tx.Commit(timeout)
	if err != nil {
		return 0, errs.Internal(err)
	}
	return int(rows), nil
}

func (s *Service) deleteRecursivelyByParentId(timeout context.Context, tx pgx.Tx, id int64) *errs.AppError {
	idsByParentId, err := s.q.WithTx(tx).GetCategoriesByParentId(timeout, &id)
	if err != nil {
		return errs.Internal(err)
	}
	for _, idByParentId := range idsByParentId {
		appErr := s.deleteRecursivelyByParentId(timeout, tx, idByParentId)
		if appErr != nil {
			return appErr
		}
	}
	_, err = s.q.WithTx(tx).DeleteCategory(timeout, id)
	if err != nil {
		return errs.Internal(err)
	}
	_, err = s.q.WithTx(tx).DeleteProductsByCategoryId(timeout, id)
	if err != nil {
		return errs.Internal(err)
	}
	return nil
}
