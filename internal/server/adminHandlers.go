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
		fail(c, "error binding json", err)
		return
	}
	brand, err := s.createBrandService(c.Request.Context(), request)
	if err != nil {
		fail(c, "error creating brand", err)
		return
	}
	c.JSON(http.StatusCreated, brand)
}

func (s *Server) GetAllBrandHandler(c *gin.Context) {
	brands, err := s.GetAllBrandService(c.Request.Context())
	if err != nil {
		fail(c, "cannot get all brands", err)
		return
	}
	c.JSON(http.StatusOK, nonNilSlice(brands))
}

func (s *Server) GetBrandHandler(c *gin.Context) {
	id, err := getIntPathParam(c, "id")
	if err != nil {
		fail(c, "error getting path param", err)
		return
	}
	brand, err := s.GetBrandService(c.Request.Context(), id)
	if err != nil {
		fail(c, "cannot get brand by id", err)
		return
	}
	c.JSON(http.StatusOK, brand)
}

func (s *Server) UpdateBrandHandler(c *gin.Context) {
	id, err := getIntPathParam(c, "id")
	if err != nil {
		fail(c, "error getting path param", err)
		return
	}
	request, err := bindJson[DTO.UpdateBrandRequest](c)
	if err != nil {
		fail(c, "error binding json", err)
		return
	}
	brand, err := s.UpdateBrandService(c.Request.Context(), id, request)
	if err != nil {
		fail(c, "error updating brand", err)
		return
	}
	c.JSON(http.StatusOK, brand)
}

func (s *Server) DeleteBrandHandler(c *gin.Context) {
	id, err := getIntPathParam(c, "id")
	if err != nil {
		fail(c, "error getting path param", err)
		return
	}
	_, err = s.DeleteBrandService(c.Request.Context(), id)
	if err != nil {
		fail(c, "error deleting brand", err)
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
