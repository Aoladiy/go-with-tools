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
	product, err := s.q.WithTx(tx).CreateProduct(timeout, mapRequestToCreateParams(request))
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

	err = tx.Commit(timeout)
	if err != nil {
		return DTO.ProductResponse{}, errs.Internal(err)
	}
	return mapCreateRowToResponse(product), nil
}

func (s *Service) GetAll(ctx context.Context) ([]DTO.ProductResponse, *errs.AppError) {
	products, err := s.q.GetAllProducts(ctx)
	if err != nil {
		return nil, errs.Internal(err)
	}

	productsResponse := make([]DTO.ProductResponse, len(products))
	for i, product := range products {
		productsResponse[i] = mapGetAllRowToResponse(product)
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

	return mapGetRowToResponse(product), nil
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
	qtx := s.q.WithTx(tx)

	oldProduct, err := qtx.GetProduct(timeout, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO.ProductResponse{}, errs.NotFound(err)
		}

		return DTO.ProductResponse{}, errs.Internal(err)
	}

	product, err := qtx.UpdateProduct(timeout, mapRequestToUpdateParams(id, request))
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

	return mapUpdateRowToResponse(product), nil
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
