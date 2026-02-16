package category

import (
	"context"
	"errors"
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/errs"

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
	rows, err := s.q.DeleteCategory(ctx, id)
	if err != nil {
		return 0, errs.Internal(err)
	}
	if rows == 0 {
		return int(rows), errs.NotFound(errors.New("category not found"))
	}
	return int(rows), nil
}
