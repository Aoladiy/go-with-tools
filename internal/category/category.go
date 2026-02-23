package category

import (
	"context"
	"errors"
	"fmt"
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/errs"
	"go-with-tools/internal/helpers"
	"net/http"

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

	category, err := s.q.CreateCategory(ctx, mapRequestToCreateParams(request))
	if err != nil {
		if pgErr, isUniqueViolation := errs.IsUniqueViolation(err); isUniqueViolation {
			return DTO.CategoryResponse{}, errs.UniqueViolation(err, pgErr)
		}
		if pgErr, isForeignKeyViolation := errs.IsForeignKeyViolation(err); isForeignKeyViolation {
			return DTO.CategoryResponse{}, errs.ForeignKeyViolation(err, pgErr)
		}
		return DTO.CategoryResponse{}, errs.Internal(err)
	}

	return mapCreateRowToResponse(category), nil
}

func (s *Service) GetAll(ctx context.Context) ([]DTO.CategoryResponse, *errs.AppError) {
	categories, err := s.q.GetAllCategories(ctx)
	if err != nil {
		return nil, errs.Internal(err)
	}

	categoriesResponse := make([]DTO.CategoryResponse, len(categories))
	for i, category := range categories {
		categoriesResponse[i] = mapGetAllRowToResponse(category)
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

	return mapGetRowToResponse(category), nil
}

func (s *Service) Update(ctx context.Context, id int64, request DTO.CategoryRequest) (DTO.CategoryResponse, *errs.AppError) {
	if request.ParentId != nil {
		_, appErr := s.Get(ctx, *request.ParentId)
		if appErr != nil {
			if appErr.Code == http.StatusNotFound {
				return DTO.CategoryResponse{}, errs.NotFound(fmt.Errorf("parent category with id=%d not found | %w", *request.ParentId, appErr.Unwrap()))
			}
			return DTO.CategoryResponse{}, appErr
		}
	}

	category, err := s.q.UpdateCategory(ctx, mapRequestToUpdateParams(id, request))
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

	return mapUpdateRowToResponse(category), nil
}

func (s *Service) Delete(ctx context.Context, id int64) (int, *errs.AppError) {
	var rows int64
	appErr := helpers.WithTx(ctx, s.p, s.q, func(timeout context.Context, q *queries.Queries) *errs.AppError {
		var err error
		rows, err = q.DeleteCategory(timeout, id)
		if err != nil {
			return errs.Internal(err)
		}
		if rows == 0 {
			return errs.NotFound(errors.New("category not found"))
		}
		_, err = q.DeleteProductsByCategoryId(timeout, id)
		if err != nil {
			return errs.Internal(err)
		}

		appErr := s.deleteRecursivelyByParentId(timeout, q, id)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return 0, appErr
	}

	return int(rows), nil
}

func (s *Service) deleteRecursivelyByParentId(timeout context.Context, q *queries.Queries, id int64) *errs.AppError {
	idsByParentId, err := q.GetCategoriesByParentId(timeout, &id)
	if err != nil {
		return errs.Internal(err)
	}
	for _, idByParentId := range idsByParentId {
		appErr := s.deleteRecursivelyByParentId(timeout, q, idByParentId)
		if appErr != nil {
			return appErr
		}
	}
	_, err = q.DeleteCategory(timeout, id)
	if err != nil {
		return errs.Internal(err)
	}
	_, err = q.DeleteProductsByCategoryId(timeout, id)
	if err != nil {
		return errs.Internal(err)
	}
	return nil
}
