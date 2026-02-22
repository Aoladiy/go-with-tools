package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	front := apiV1.Group("/front")
	front.GET("/health", s.healthHandler)
	front.GET("/hello", s.HelloWorldHandler)
	front.GET("/categories", s.CategoriesHandler)
	front.GET("/brands", s.BrandsHandler)
	front.GET("/products", s.ProductsHandler)

	apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(2)))

	admin := apiV1.Group("/admin")
	admin.POST("/sign-up", s.SignUpHandler)
	admin.POST("/sign-in", s.SignInHandler)
	admin.POST("/token-refresh", s.TokenRefreshHandler)
	admin.POST("/sign-out", s.SignOutHandler)

	products := admin.Group("/products")
	products.Use(AuthByJWT())
	products.POST("", s.CreateProductHandler)
	products.GET("", s.GetAllProductHandler)
	products.GET("/:id", s.GetProductHandler)
	products.PUT("/:id", s.UpdateProductHandler)
	products.DELETE("/:id", s.DeleteProductHandler)
	products.GET("/:id/priceHistory", s.GetProductPriceHistory)

	brands := admin.Group("/brands")
	brands.Use(AuthByJWT())
	brands.POST("", s.CreateBrandHandler)
	brands.GET("", s.GetAllBrandHandler)
	brands.GET("/:id", s.GetBrandHandler)
	brands.PUT("/:id", s.UpdateBrandHandler)
	brands.DELETE("/:id", s.DeleteBrandHandler)

	categories := admin.Group("/categories")
	categories.Use(AuthByJWT())
	categories.POST("", s.CreateCategoryHandler)
	categories.GET("", s.GetAllCategoryHandler)
	categories.GET("/:id", s.GetCategoryHandler)
	categories.PUT("/:id", s.UpdateCategoryHandler)
	categories.DELETE("/:id", s.DeleteCategoryHandler)

	admin.POST("/inventory/adjustments", s.CreateInventoryMovementHandler).Use(AuthByJWT())

	return r
}
