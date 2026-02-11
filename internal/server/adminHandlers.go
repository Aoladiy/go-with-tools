package server

import (
	"go-with-tools/internal/DTO"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) LoginHandler(c *gin.Context) {
	// TODO
}

func (s *Server) LogoutHandler(c *gin.Context) {
	// TODO
}

func (s *Server) CreateProductHandler(c *gin.Context) {
	// TODO
}

func (s *Server) GetAllProductHandler(c *gin.Context) {
	// TODO
}

func (s *Server) GetProductHandler(c *gin.Context) {
	// TODO
}

func (s *Server) UpdateProductHandler(c *gin.Context) {
	// TODO
}

func (s *Server) DeleteProductHandler(c *gin.Context) {
	// TODO
}

func (s *Server) CreateBrandHandler(c *gin.Context) {
	request, err := bindJson[DTO.CreateBrandRequest](c)
	if err != nil {
		respondError(c, err)
		return
	}
	brand, err := s.brand.Create(c.Request.Context(), request)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, brand)
}

func (s *Server) GetAllBrandHandler(c *gin.Context) {
	brands, err := s.brand.GetAll(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, nonNilSlice(brands))
}

func (s *Server) GetBrandHandler(c *gin.Context) {
	id, err := getIntPathParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	brand, err := s.brand.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, brand)
}

func (s *Server) UpdateBrandHandler(c *gin.Context) {
	id, err := getIntPathParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	request, err := bindJson[DTO.UpdateBrandRequest](c)
	if err != nil {
		respondError(c, err)
		return
	}
	brand, err := s.brand.Update(c.Request.Context(), id, request)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, brand)
}

func (s *Server) DeleteBrandHandler(c *gin.Context) {
	id, err := getIntPathParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	_, err = s.brand.Delete(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *Server) CreateCategoryHandler(c *gin.Context) {
	// TODO
}

func (s *Server) GetAllCategoryHandler(c *gin.Context) {
	// TODO
}

func (s *Server) GetCategoryHandler(c *gin.Context) {
	// TODO
}

func (s *Server) UpdateCategoryHandler(c *gin.Context) {
	// TODO
}

func (s *Server) DeleteCategoryHandler(c *gin.Context) {
	// TODO
}

func (s *Server) CreateInventoryMovementHandler(c *gin.Context) {
	// TODO
}
