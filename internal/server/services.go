package server

import (
	"context"
	"errors"
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/errs"

	"github.com/jackc/pgx/v5"
)

func (s *Server) createBrandService(ctx context.Context, request DTO.CreateBrandRequest) (queries.Brand, *errs.AppError) {
	brand, err := s.q.CreateBrand(ctx, queries.CreateBrandParams{
		Name: request.Name,
		Slug: request.Slug,
	})
	if err != nil {
		return queries.Brand{}, errs.Internal(err)
	}

	return brand, nil
}

func (s *Server) GetAllBrandService(ctx context.Context) ([]queries.Brand, *errs.AppError) {
	brands, err := s.q.GetAllBrands(ctx)
	if err != nil {
		return nil, errs.Internal(err)
	}
	return brands, nil
}

func (s *Server) GetBrandService(ctx context.Context, id int) (queries.Brand, *errs.AppError) {
	brand, err := s.q.GetBrand(ctx, int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return queries.Brand{}, errs.NotFound(err)
		}

		return queries.Brand{}, errs.Internal(err)
	}
	return brand, nil
}

func (s *Server) UpdateBrandService(ctx context.Context, id int, request DTO.UpdateBrandRequest) (queries.Brand, *errs.AppError) {
	brand, err := s.q.UpdateBrand(ctx, queries.UpdateBrandParams{
		ID:   int64(id),
		Name: request.Name,
		Slug: request.Slug,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return queries.Brand{}, errs.NotFound(err)
		}

		return queries.Brand{}, errs.Internal(err)
	}
	return brand, nil
}

func (s *Server) DeleteBrandService(ctx context.Context, id int) (int, *errs.AppError) {
	rows, err := s.q.DeleteBrand(ctx, int64(id))
	if err != nil {
		return 0, errs.Internal(err)
	}
	if rows == 0 {
		return int(rows), errs.NotFound(err)
	}
	return int(rows), nil
}
