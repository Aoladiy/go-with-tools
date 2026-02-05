package server

import (
	"go-with-tools/internal/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HelloWorldHandler(c *gin.Context) {
	brands, err := s.q.GetAllBrands(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, helpers.NonNil(brands))
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) BrandsHandler(c *gin.Context) {
	// TODO
}

func (s *Server) CategoriesHandler(c *gin.Context) {
	// TODO
}

func (s *Server) ProductsHandler(c *gin.Context) {
	// TODO
}
