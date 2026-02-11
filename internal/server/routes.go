package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))
	r.Use(LogErrors())
	apiV1 := r.Group("/api/v1")
	apiV1.GET("/health", s.healthHandler)
	apiV1.GET("/hello", s.HelloWorldHandler)
	apiV1.GET("/categories", s.CategoriesHandler)
	apiV1.GET("/brands", s.BrandsHandler)
	apiV1.GET("/products", s.ProductsHandler)

	admin := apiV1.Group("/admin")
	admin.Use(AuthByJWT())
	admin.POST("/login", s.LoginHandler)
	admin.POST("/logout", s.LogoutHandler)

	products := admin.Group("/products")
	products.POST("", s.CreateProductHandler)
	products.GET("", s.GetAllProductHandler)
	products.GET("/:id", s.GetProductHandler)
	products.PUT("/:id", s.UpdateProductHandler)
	products.DELETE("/:id", s.DeleteProductHandler)

	brands := admin.Group("/brands")
	brands.POST("", s.CreateBrandHandler)
	brands.GET("", s.GetAllBrandHandler)
	brands.GET("/:id", s.GetBrandHandler)
	brands.PUT("/:id", s.UpdateBrandHandler)
	brands.DELETE("/:id", s.DeleteBrandHandler)

	categories := admin.Group("/categories")
	categories.POST("", s.CreateCategoryHandler)
	categories.GET("", s.GetAllCategoryHandler)
	categories.GET("/:id", s.GetCategoryHandler)
	categories.PUT("/:id", s.UpdateCategoryHandler)
	categories.DELETE("/:id", s.DeleteCategoryHandler)

	admin.POST("/inventory/adjustments", s.CreateInventoryMovementHandler)

	return r
}
