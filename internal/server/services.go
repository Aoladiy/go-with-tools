package server

import (
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/helpers"

	"github.com/gin-gonic/gin"
)

func (s *Server) createBrandService(c *gin.Context) (queries.Brand, *helpers.CustomErr) {
	var request DTO.CreateBrandRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return queries.Brand{}, &helpers.ErrBadRequest
	}
	brand, err := s.q.CreateBrand(c.Request.Context(), queries.CreateBrandParams{
		Name: request.Name,
		Slug: request.Slug,
	})
	if err != nil {
		return queries.Brand{}, &helpers.ErrInternalServerError
	}

	return brand, nil
}
