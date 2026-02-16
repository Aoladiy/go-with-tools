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
	request, err := bindJson[DTO.ProductRequest](c)
	if err != nil {
		respondError(c, err)
		return
	}
	product, err := s.product.Create(c.Request.Context(), request)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (s *Server) GetAllProductHandler(c *gin.Context) {
	products, err := s.product.GetAll(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, nonNilSlice(products))
}

func (s *Server) GetProductHandler(c *gin.Context) {
	id, err := getInt64PathParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	product, err := s.product.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, product)
}

func (s *Server) UpdateProductHandler(c *gin.Context) {
	id, err := getInt64PathParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	request, err := bindJson[DTO.ProductRequest](c)
	if err != nil {
		respondError(c, err)
		return
	}
	product, err := s.product.Update(c.Request.Context(), id, request)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, product)
}

func (s *Server) DeleteProductHandler(c *gin.Context) {
	id, err := getInt64PathParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	_, err = s.product.Delete(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *Server) GetProductPriceHistory(c *gin.Context) {
	id, err := getInt64PathParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	priceHistory, err := s.product.GetPriceHistory(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, nonNilSlice(priceHistory))
}

func (s *Server) CreateBrandHandler(c *gin.Context) {
	request, err := bindJson[DTO.BrandRequest](c)
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
	id, err := getInt64PathParam(c, "id")
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
	id, err := getInt64PathParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	request, err := bindJson[DTO.BrandRequest](c)
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
	id, err := getInt64PathParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	_, err = s.brand.Delete(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *Server) CreateCategoryHandler(c *gin.Context) {
	request, err := bindJson[DTO.CategoryRequest](c)
	if err != nil {
		respondError(c, err)
		return
	}
	category, err := s.category.Create(c.Request.Context(), request)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, category)
}

func (s *Server) GetAllCategoryHandler(c *gin.Context) {
	categories, err := s.category.GetAll(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, nonNilSlice(categories))
}

func (s *Server) GetCategoryHandler(c *gin.Context) {
	id, err := getInt64PathParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	category, err := s.category.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, category)
}

func (s *Server) UpdateCategoryHandler(c *gin.Context) {
	id, err := getInt64PathParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	request, err := bindJson[DTO.CategoryRequest](c)
	if err != nil {
		respondError(c, err)
		return
	}
	category, err := s.category.Update(c.Request.Context(), id, request)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, category)
}

func (s *Server) DeleteCategoryHandler(c *gin.Context) {
	id, err := getInt64PathParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	_, err = s.category.Delete(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *Server) CreateInventoryMovementHandler(c *gin.Context) {
	// TODO
}
