package server

import (
	"context"
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/errs"
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
