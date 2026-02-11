package category

import (
	"context"
	"errors"
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/errs"
	"go-with-tools/internal/helpers"

	"github.com/jackc/pgx/v5"
)

type Service struct {
	q queries.Querier
}

func New(q queries.Querier) *Service {
	return &Service{q: q}
}

func (s *Service) Create(ctx context.Context, request DTO.CategoryRequest) (DTO.CategoryResponse, *errs.AppError) {
	category, err := s.q.CreateCategory(ctx, queries.CreateCategoryParams{
		Name:     request.Name,
		Slug:     request.Slug,
		ParentID: helpers.ToPgInt8(request.ParentId),
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
		Id:        int(category.ID),
		Name:      category.Name,
		Slug:      category.Slug,
		ParentId:  helpers.ParsePgInt8(category.ParentID),
		CreatedAt: helpers.ParsePgTimestamptz(category.CreatedAt),
		UpdatedAt: helpers.ParsePgTimestamptz(category.UpdatedAt),
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
			Id:        int(category.ID),
			Name:      category.Name,
			Slug:      category.Slug,
			ParentId:  helpers.ParsePgInt8(category.ParentID),
			CreatedAt: helpers.ParsePgTimestamptz(category.CreatedAt),
			UpdatedAt: helpers.ParsePgTimestamptz(category.UpdatedAt),
		}
	}
	return categoriesResponse, nil
}

func (s *Service) Get(ctx context.Context, id int) (DTO.CategoryResponse, *errs.AppError) {
	category, err := s.q.GetCategory(ctx, int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO.CategoryResponse{}, errs.NotFound(err)
		}

		return DTO.CategoryResponse{}, errs.Internal(err)
	}

	categoryResponse := DTO.CategoryResponse{
		Id:        int(category.ID),
		Name:      category.Name,
		Slug:      category.Slug,
		ParentId:  helpers.ParsePgInt8(category.ParentID),
		CreatedAt: helpers.ParsePgTimestamptz(category.CreatedAt),
		UpdatedAt: helpers.ParsePgTimestamptz(category.UpdatedAt),
	}
	return categoryResponse, nil
}

func (s *Service) Update(ctx context.Context, id int, request DTO.CategoryRequest) (DTO.CategoryResponse, *errs.AppError) {
	category, err := s.q.UpdateCategory(ctx, queries.UpdateCategoryParams{
		ID:       int64(id),
		Name:     request.Name,
		Slug:     request.Slug,
		ParentID: helpers.ToPgInt8(request.ParentId),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO.CategoryResponse{}, errs.NotFound(err)
		}

		return DTO.CategoryResponse{}, errs.Internal(err)
	}

	categoryResponse := DTO.CategoryResponse{
		Id:        int(category.ID),
		Name:      category.Name,
		Slug:      category.Slug,
		ParentId:  helpers.ParsePgInt8(category.ParentID),
		CreatedAt: helpers.ParsePgTimestamptz(category.CreatedAt),
		UpdatedAt: helpers.ParsePgTimestamptz(category.UpdatedAt),
	}
	return categoryResponse, nil
}

func (s *Service) Delete(ctx context.Context, id int) (int, *errs.AppError) {
	rows, err := s.q.DeleteCategory(ctx, int64(id))
	if err != nil {
		return 0, errs.Internal(err)
	}
	if rows == 0 {
		return int(rows), errs.NotFound(errors.New("category not found"))
	}
	return int(rows), nil
}
