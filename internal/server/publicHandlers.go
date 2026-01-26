package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
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
