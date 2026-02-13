package brand

import (
	"context"
	"errors"
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/errs"

	"github.com/jackc/pgx/v5"
)

type Service struct {
	q queries.Querier
}

func New(q queries.Querier) *Service {
	return &Service{q: q}
}

func (s *Service) Create(ctx context.Context, request DTO.BrandRequest) (DTO.BrandResponse, *errs.AppError) {
	brand, err := s.q.CreateBrand(ctx, queries.CreateBrandParams{
		Name: request.Name,
		Slug: request.Slug,
	})
	if err != nil {
		if pgErr, isUniqueViolation := errs.IsUniqueViolation(err); isUniqueViolation {
			return DTO.BrandResponse{}, errs.UniqueViolation(err, pgErr)
		}
		return DTO.BrandResponse{}, errs.Internal(err)
	}

	brandResponse := DTO.BrandResponse{
		Id:        brand.ID,
		Name:      brand.Name,
		Slug:      brand.Slug,
		CreatedAt: brand.CreatedAt,
		UpdatedAt: brand.UpdatedAt,
	}
	return brandResponse, nil
}

func (s *Service) GetAll(ctx context.Context) ([]DTO.BrandResponse, *errs.AppError) {
	brands, err := s.q.GetAllBrands(ctx)
	if err != nil {
		return nil, errs.Internal(err)
	}

	brandsResponse := make([]DTO.BrandResponse, len(brands))
	for i, brand := range brands {
		brandsResponse[i] = DTO.BrandResponse{
			Id:        brand.ID,
			Name:      brand.Name,
			Slug:      brand.Slug,
			CreatedAt: brand.CreatedAt,
			UpdatedAt: brand.UpdatedAt,
		}
	}
	return brandsResponse, nil
}

func (s *Service) Get(ctx context.Context, id int64) (DTO.BrandResponse, *errs.AppError) {
	brand, err := s.q.GetBrand(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO.BrandResponse{}, errs.NotFound(err)
		}

		return DTO.BrandResponse{}, errs.Internal(err)
	}

	brandResponse := DTO.BrandResponse{
		Id:        brand.ID,
		Name:      brand.Name,
		Slug:      brand.Slug,
		CreatedAt: brand.CreatedAt,
		UpdatedAt: brand.UpdatedAt,
	}
	return brandResponse, nil
}

func (s *Service) Update(ctx context.Context, id int64, request DTO.BrandRequest) (DTO.BrandResponse, *errs.AppError) {
	brand, err := s.q.UpdateBrand(ctx, queries.UpdateBrandParams{
		ID:   id,
		Name: request.Name,
		Slug: request.Slug,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO.BrandResponse{}, errs.NotFound(err)
		}

		return DTO.BrandResponse{}, errs.Internal(err)
	}

	brandResponse := DTO.BrandResponse{
		Id:        brand.ID,
		Name:      brand.Name,
		Slug:      brand.Slug,
		CreatedAt: brand.CreatedAt,
		UpdatedAt: brand.UpdatedAt,
	}
	return brandResponse, nil
}

func (s *Service) Delete(ctx context.Context, id int64) (int, *errs.AppError) {
	rows, err := s.q.DeleteBrand(ctx, id)
	if err != nil {
		return 0, errs.Internal(err)
	}
	if rows == 0 {
		return int(rows), errs.NotFound(errors.New("brand not found"))
	}
	return int(rows), nil
}
