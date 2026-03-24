package brand

import (
	"github.com/Aoladiy/go-with-tools/internal/DTO"
	"github.com/Aoladiy/go-with-tools/internal/database/queries"
)

func mapRequestToCreateParams(request DTO.BrandRequest) queries.CreateBrandParams {
	return queries.CreateBrandParams{
		Name: request.Name,
		Slug: request.Slug,
	}
}

func mapCreateRowToResponse(brand queries.CreateBrandRow) DTO.BrandResponse {
	return DTO.BrandResponse{
		Id:        brand.ID,
		Name:      brand.Name,
		Slug:      brand.Slug,
		CreatedAt: brand.CreatedAt,
		UpdatedAt: brand.UpdatedAt,
	}
}

func mapGetAllRowToResponse(brand queries.GetAllBrandsRow) DTO.BrandResponse {
	return DTO.BrandResponse{
		Id:        brand.ID,
		Name:      brand.Name,
		Slug:      brand.Slug,
		CreatedAt: brand.CreatedAt,
		UpdatedAt: brand.UpdatedAt,
	}
}

func mapGetRowToResponse(brand queries.GetBrandRow) DTO.BrandResponse {
	return DTO.BrandResponse{
		Id:        brand.ID,
		Name:      brand.Name,
		Slug:      brand.Slug,
		CreatedAt: brand.CreatedAt,
		UpdatedAt: brand.UpdatedAt,
	}
}

func mapRequestToUpdateParams(id int64, request DTO.BrandRequest) queries.UpdateBrandParams {
	return queries.UpdateBrandParams{
		ID:   id,
		Name: request.Name,
		Slug: request.Slug,
	}
}

func mapUpdateRowToResponse(brand queries.UpdateBrandRow) DTO.BrandResponse {
	return DTO.BrandResponse{
		Id:        brand.ID,
		Name:      brand.Name,
		Slug:      brand.Slug,
		CreatedAt: brand.CreatedAt,
		UpdatedAt: brand.UpdatedAt,
	}
}
