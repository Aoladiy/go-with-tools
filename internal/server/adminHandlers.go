package server

import (
	"errors"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/helpers"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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
	var request CreateBrandRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	brand, err := s.q.CreateBrand(c.Request.Context(), queries.CreateBrandParams{
		Name: request.Name,
		Slug: request.Slug,
	})
	if err != nil {
		log.Println("error creating brand", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, brand)
	return
}

func (s *Server) GetAllBrandHandler(c *gin.Context) {
	brands, err := s.q.GetAllBrands(c.Request.Context())
	if err != nil {
		log.Println("error getting brand", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, helpers.NonNil(brands))
	return
}

func (s *Server) GetBrandHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("error getting brand", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	brand, err := s.q.GetBrand(c.Request.Context(), int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{})
		} else {
			log.Println("error getting brand", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, brand)
	return
}

func (s *Server) UpdateBrandHandler(c *gin.Context) {
	var updateBrandRequest UpdateBrandRequest
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("error deleting brand", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = c.ShouldBindJSON(&updateBrandRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	brand, err := s.q.UpdateBrand(c.Request.Context(), queries.UpdateBrandParams{
		ID:   int64(id),
		Name: updateBrandRequest.Name,
		Slug: updateBrandRequest.Slug,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, brand)
	return
}

func (s *Server) DeleteBrandHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("error deleting brand", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rows, err := s.q.DeleteBrand(c.Request.Context(), int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{})
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
