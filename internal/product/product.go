package product

import (
	"context"
	"errors"
	"fmt"
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/errs"
	"go-with-tools/internal/helpers"
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

func (s *Service) Create(ctx context.Context, request DTO.ProductRequest) (DTO.ProductResponse, *errs.AppError) {
	_, err := s.q.GetBrand(ctx, request.BrandId)
	if err != nil {
		return DTO.ProductResponse{}, errs.NotFound(fmt.Errorf("brand with id=%d not found | %w", request.BrandId, err))
	}
	_, err = s.q.GetCategory(ctx, request.CategoryId)
	if err != nil {
		return DTO.ProductResponse{}, errs.NotFound(fmt.Errorf("category with id=%d not found | %w", request.CategoryId, err))
	}

	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	tx, err := s.p.Begin(timeout)
	if err != nil {
		return DTO.ProductResponse{}, errs.Internal(err)
	}
	defer tx.Rollback(timeout)
	product, err := s.q.WithTx(tx).CreateProduct(timeout, queries.CreateProductParams{ //TODO handle case with non empty brand's and\or category's deleted_at
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

	appErr := s.createPriceHistory(timeout, tx, product.ID, 0, product.PriceKopeck)
	if appErr != nil {
		return DTO.ProductResponse{}, appErr
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
	err = tx.Commit(timeout)
	if err != nil {
		return DTO.ProductResponse{}, errs.Internal(err)
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
	_, err := s.q.GetBrand(ctx, request.BrandId)
	if err != nil {
		return DTO.ProductResponse{}, errs.NotFound(fmt.Errorf("brand with id=%d not found | %w", request.BrandId, err))
	}
	_, err = s.q.GetCategory(ctx, request.CategoryId)
	if err != nil {
		return DTO.ProductResponse{}, errs.NotFound(fmt.Errorf("category with id=%d not found | %w", request.CategoryId, err))
	}

	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	tx, err := s.p.Begin(timeout)
	if err != nil {
		return DTO.ProductResponse{}, errs.Internal(err)
	}
	defer tx.Rollback(timeout)

	oldProduct, err := s.q.WithTx(tx).GetProduct(timeout, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO.ProductResponse{}, errs.NotFound(err)
		}

		return DTO.ProductResponse{}, errs.Internal(err)
	}

	product, err := s.q.WithTx(tx).UpdateProduct(timeout, queries.UpdateProductParams{ //TODO handle case with non empty brand's and\or category's deleted_at
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

	if oldProduct.PriceKopeck != product.PriceKopeck {
		appErr := s.createPriceHistory(timeout, tx, product.ID, oldProduct.PriceKopeck, product.PriceKopeck)
		if appErr != nil {
			return DTO.ProductResponse{}, appErr
		}
	}

	err = tx.Commit(timeout)
	if err != nil {
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

func (s *Service) GetPriceHistory(ctx context.Context, id int64) ([]queries.ProductPriceHistory, *errs.AppError) {
	byProductId, err := s.q.GetProductPriceHistoryByProductId(ctx, id)
	if err != nil {
		return nil, errs.Internal(err)
	}
	return byProductId, nil
}

func (s *Service) createPriceHistory(timeout context.Context, tx pgx.Tx, productId int64, oldPrice, newPrice int32) *errs.AppError {
	id, err := helpers.SafeGetUserID(timeout)
	if err != nil {
		return errs.Internal(err)
	}
	_, err = s.q.WithTx(tx).CreateProductPriceHistory(timeout, queries.CreateProductPriceHistoryParams{
		ProductID:      productId,
		OldPriceKopeck: oldPrice,
		NewPriceKopeck: newPrice,
		UpdatedBy:      id,
	})
	if err != nil {
		return errs.Internal(err)
	}
	return nil
}
