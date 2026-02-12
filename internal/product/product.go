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
		BrandID:     int64(request.BrandId),
		CategoryID:  int64(request.CategoryId),
		Name:        request.Name,
		Slug:        request.Slug,
		Description: helpers.DerefString(request.Description, ""),
		PriceKopeck: int32(request.PriceKopeck),
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
		Id:          int(product.ID),
		BrandId:     int(product.BrandID),
		CategoryId:  int(product.CategoryID),
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		PriceKopeck: int(product.PriceKopeck),
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
			Id:          int(product.ID),
			BrandId:     int(product.BrandID),
			CategoryId:  int(product.CategoryID),
			Name:        product.Name,
			Slug:        product.Slug,
			Description: product.Description,
			PriceKopeck: int(product.PriceKopeck),
			IsActive:    product.IsActive,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}
	return productsResponse, nil
}

func (s *Service) Get(ctx context.Context, id int) (DTO.ProductResponse, *errs.AppError) {
	product, err := s.q.GetProduct(ctx, int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO.ProductResponse{}, errs.NotFound(err)
		}

		return DTO.ProductResponse{}, errs.Internal(err)
	}

	productResponse := DTO.ProductResponse{
		Id:          int(product.ID),
		BrandId:     int(product.BrandID),
		CategoryId:  int(product.CategoryID),
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		PriceKopeck: int(product.PriceKopeck),
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
	return productResponse, nil
}

func (s *Service) Update(ctx context.Context, id int, request DTO.ProductRequest) (DTO.ProductResponse, *errs.AppError) {
	product, err := s.q.UpdateProduct(ctx, queries.UpdateProductParams{
		ID:          int64(id),
		BrandID:     int64(request.BrandId),
		CategoryID:  int64(request.CategoryId),
		Name:        request.Name,
		Slug:        request.Slug,
		Description: helpers.DerefString(request.Description, ""),
		PriceKopeck: int32(request.PriceKopeck),
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
		Id:          int(product.ID),
		BrandId:     int(product.BrandID),
		CategoryId:  int(product.CategoryID),
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		PriceKopeck: int(product.PriceKopeck),
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
	return productResponse, nil
}

func (s *Service) Delete(ctx context.Context, id int) (int, *errs.AppError) {
	rows, err := s.q.DeleteProduct(ctx, int64(id))
	if err != nil {
		return 0, errs.Internal(err)
	}
	if rows == 0 {
		return int(rows), errs.NotFound(errors.New("product not found"))
	}
	return int(rows), nil
}
