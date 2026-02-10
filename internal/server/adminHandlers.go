package server

import (
	"errors"
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/database/queries"
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
	request, err := bindJson[DTO.CreateBrandRequest](c)
	if err != nil {
		fail(c, "error binding json while creating brand", err)
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
	brands, err := s.q.GetAllBrands(c.Request.Context())
	if err != nil {
		log.Println("error getting brand", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nonNilSlice(brands))
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
}

func (s *Server) UpdateBrandHandler(c *gin.Context) {
	var updateBrandRequest DTO.UpdateBrandRequest
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
