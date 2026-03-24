package category

import (
	"github.com/Aoladiy/go-with-tools/internal/DTO"
	"github.com/Aoladiy/go-with-tools/internal/database/queries"
)

func mapRequestToCreateParams(request DTO.CategoryRequest) queries.CreateCategoryParams {
	return queries.CreateCategoryParams{
		Name:     request.Name,
		Slug:     request.Slug,
		ParentID: request.ParentId,
	}
}

func mapCreateRowToResponse(category queries.CreateCategoryRow) DTO.CategoryResponse {
	return DTO.CategoryResponse{
		Id:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		ParentId:  category.ParentID,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func mapGetAllRowToResponse(category queries.GetAllCategoriesRow) DTO.CategoryResponse {
	return DTO.CategoryResponse{
		Id:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		ParentId:  category.ParentID,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func mapGetRowToResponse(category queries.GetCategoryRow) DTO.CategoryResponse {
	return DTO.CategoryResponse{
		Id:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		ParentId:  category.ParentID,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func mapRequestToUpdateParams(id int64, request DTO.CategoryRequest) queries.UpdateCategoryParams {
	return queries.UpdateCategoryParams{
		ID:       id,
		Name:     request.Name,
		Slug:     request.Slug,
		ParentID: request.ParentId,
	}
}

func mapUpdateRowToResponse(category queries.UpdateCategoryRow) DTO.CategoryResponse {
	return DTO.CategoryResponse{
		Id:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		ParentId:  category.ParentID,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}
