package product

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

func (s *Service) Create(ctx context.Context, request DTO.ProductRequest) (DTO.ProductResponse, *errs.AppError) {
	product, err := s.q.CreateProduct(ctx, queries.CreateProductParams{
		BrandID:     request.BrandId,
		CategoryID:  request.CategoryId,
		Name:        request.Name,
		Slug:        request.Slug,
		Description: helpers.DerefString(request.Description, ""),
		PriceKopeck: request.PriceKopeck,
		IsActive:    helpers.DerefBool(request.IsActive, true),
	})
	if err != nil {
		if pgErr, isUniqueViolation := errs.IsUniqueViolation(err); isUniqueViolation {
			return DTO.ProductResponse{}, errs.UniqueViolation(err, pgErr)
		}
		if pgErr, isForeignKeyViolation := errs.IsForeignKeyViolation(err); isForeignKeyViolation {
			return DTO.ProductResponse{}, errs.ForeignKeyViolation(err, pgErr)
		}
		return DTO.ProductResponse{}, errs.Internal(err)
	}

	productResponse := DTO.ProductResponse{
		Id:          product.ID,
		BrandId:     product.BrandID,
		CategoryId:  product.CategoryID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		PriceKopeck: product.PriceKopeck,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
	return productResponse, nil
}

func (s *Service) GetAll(ctx context.Context) ([]DTO.ProductResponse, *errs.AppError) {
	products, err := s.q.GetAllProducts(ctx)
	if err != nil {
		return nil, errs.Internal(err)
	}

	productsResponse := make([]DTO.ProductResponse, len(products))
	for i, product := range products {
		productsResponse[i] = DTO.ProductResponse{
			Id:          product.ID,
			BrandId:     product.BrandID,
			CategoryId:  product.CategoryID,
			Name:        product.Name,
			Slug:        product.Slug,
			Description: product.Description,
			PriceKopeck: product.PriceKopeck,
			IsActive:    product.IsActive,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}
	return productsResponse, nil
}

func (s *Service) Get(ctx context.Context, id int64) (DTO.ProductResponse, *errs.AppError) {
	product, err := s.q.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO.ProductResponse{}, errs.NotFound(err)
		}

		return DTO.ProductResponse{}, errs.Internal(err)
	}

	productResponse := DTO.ProductResponse{
		Id:          product.ID,
		BrandId:     product.BrandID,
		CategoryId:  product.CategoryID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		PriceKopeck: product.PriceKopeck,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
	return productResponse, nil
}

func (s *Service) Update(ctx context.Context, id int64, request DTO.ProductRequest) (DTO.ProductResponse, *errs.AppError) {
	product, err := s.q.UpdateProduct(ctx, queries.UpdateProductParams{
		ID:          id,
		BrandID:     request.BrandId,
		CategoryID:  request.CategoryId,
		Name:        request.Name,
		Slug:        request.Slug,
		Description: helpers.DerefString(request.Description, ""),
		PriceKopeck: request.PriceKopeck,
		IsActive:    helpers.DerefBool(request.IsActive, true),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO.ProductResponse{}, errs.NotFound(err)
		}
		if pgErr, isUniqueViolation := errs.IsUniqueViolation(err); isUniqueViolation {
			return DTO.ProductResponse{}, errs.UniqueViolation(err, pgErr)
		}
		if pgErr, isForeignKeyViolation := errs.IsForeignKeyViolation(err); isForeignKeyViolation {
			return DTO.ProductResponse{}, errs.ForeignKeyViolation(err, pgErr)
		}

		return DTO.ProductResponse{}, errs.Internal(err)
	}

	productResponse := DTO.ProductResponse{
		Id:          product.ID,
		BrandId:     product.BrandID,
		CategoryId:  product.CategoryID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		PriceKopeck: product.PriceKopeck,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
	return productResponse, nil
}

func (s *Service) Delete(ctx context.Context, id int64) (int, *errs.AppError) {
	rows, err := s.q.DeleteProduct(ctx, id)
	if err != nil {
		return 0, errs.Internal(err)
	}
	if rows == 0 {
		return int(rows), errs.NotFound(errors.New("product not found"))
	}
	return int(rows), nil
}
