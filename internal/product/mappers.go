package product

import (
	"github.com/Aoladiy/go-with-tools/internal/DTO"
	"github.com/Aoladiy/go-with-tools/internal/database/queries"
	"github.com/Aoladiy/go-with-tools/internal/helpers"
)

func mapRequestToCreateParams(request DTO.ProductRequest) queries.CreateProductParams {
	return queries.CreateProductParams{
		BrandID:     request.BrandId,
		CategoryID:  request.CategoryId,
		Name:        request.Name,
		Slug:        request.Slug,
		Description: helpers.DerefString(request.Description, ""),
		PriceKopeck: request.PriceKopeck,
		IsActive:    helpers.DerefBool(request.IsActive, true),
	}
}

func mapCreateRowToResponse(product queries.CreateProductRow) DTO.ProductResponse {
	return DTO.ProductResponse{
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

func mapGetAllRowToResponse(product queries.GetAllProductsRow) DTO.ProductResponse {
	return DTO.ProductResponse{
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

func mapGetRowToResponse(product queries.GetProductRow) DTO.ProductResponse {
	return DTO.ProductResponse{
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

func mapRequestToUpdateParams(id int64, request DTO.ProductRequest) queries.UpdateProductParams {
	return queries.UpdateProductParams{
		ID:          id,
		BrandID:     request.BrandId,
		CategoryID:  request.CategoryId,
		Name:        request.Name,
		Slug:        request.Slug,
		Description: helpers.DerefString(request.Description, ""),
		PriceKopeck: request.PriceKopeck,
		IsActive:    helpers.DerefBool(request.IsActive, true),
	}
}

func mapUpdateRowToResponse(product queries.UpdateProductRow) DTO.ProductResponse {
	return DTO.ProductResponse{
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
