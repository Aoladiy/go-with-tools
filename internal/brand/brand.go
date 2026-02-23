package brand

import (
	"context"
	"errors"
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/errs"
	"go-with-tools/internal/helpers"

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

func (s *Service) Create(ctx context.Context, request DTO.BrandRequest) (DTO.BrandResponse, *errs.AppError) {
	brand, err := s.q.CreateBrand(ctx, mapRequestToCreateParams(request))
	if err != nil {
		if pgErr, isUniqueViolation := errs.IsUniqueViolation(err); isUniqueViolation {
			return DTO.BrandResponse{}, errs.UniqueViolation(err, pgErr)
		}
		return DTO.BrandResponse{}, errs.Internal(err)
	}

	return mapCreateRowToResponse(brand), nil
}

func (s *Service) GetAll(ctx context.Context) ([]DTO.BrandResponse, *errs.AppError) {
	brands, err := s.q.GetAllBrands(ctx)
	if err != nil {
		return nil, errs.Internal(err)
	}

	brandsResponse := make([]DTO.BrandResponse, len(brands))
	for i, brand := range brands {
		brandsResponse[i] = mapGetAllRowToResponse(brand)
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

	return mapGetRowToResponse(brand), nil
}

func (s *Service) Update(ctx context.Context, id int64, request DTO.BrandRequest) (DTO.BrandResponse, *errs.AppError) {
	brand, err := s.q.UpdateBrand(ctx, mapRequestToUpdateParams(id, request))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO.BrandResponse{}, errs.NotFound(err)
		}

		return DTO.BrandResponse{}, errs.Internal(err)
	}

	return mapUpdateRowToResponse(brand), nil
}

func (s *Service) Delete(ctx context.Context, id int64) (int, *errs.AppError) {
	var rows int64
	appErr := helpers.WithTx(ctx, s.p, s.q, func(timeout context.Context, q *queries.Queries) *errs.AppError {
		var err error
		rows, err = q.DeleteBrand(timeout, id)
		if err != nil {
			return errs.Internal(err)
		}
		if rows == 0 {
			return errs.NotFound(errors.New("brand not found"))
		}

		_, err = q.DeleteProductsByBrandId(timeout, id)
		if err != nil {
			return errs.Internal(err)
		}
		return nil
	})
	if appErr != nil {
		return 0, appErr
	}

	return int(rows), nil
}
