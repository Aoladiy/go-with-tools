package server

import (
	"go-with-tools/internal/DTO"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SignUpHandler registers a new admin user and returns JWT tokens
//
//	@Summary		Admin sign up
//	@Description	Register a new admin user and receive access and refresh JWT tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		DTO.SignUpRequest	true	"Sign up credentials"
//	@Success		200		{object}	DTO.JWTResponse
//	@Failure		400		{object}	DTO.ErrorResponse	"invalid input or password too short"
//	@Failure		409		{object}	DTO.ErrorResponse	"email already exists"
//	@Failure		500		{object}	DTO.ErrorResponse
//	@Router			/admin/sign-up [post]
func (s *Server) SignUpHandler(c *gin.Context) {
	request, err := bindJson[DTO.SignUpRequest](c)
	if err != nil {
		respondError(c, err)
		return
	}
	jwtResponse, err := s.auth.SignUp(c.Request.Context(), request)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, jwtResponse)
}

// SignInHandler logins a new admin user and returns JWT tokens
//
//	@Summary		Admin login
//	@Description	login a new admin user and receive access and refresh JWT tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		DTO.SignInRequest	true	"login credentials"
//	@Success		200		{object}	DTO.JWTResponse
//	@Failure		400		{object}	DTO.ErrorResponse	"invalid input or password too short"
//	@Failure		401		{object}	DTO.ErrorResponse	"wrong credentials"
//	@Failure		500		{object}	DTO.ErrorResponse
//	@Router			/admin/sign-in [post]
func (s *Server) SignInHandler(c *gin.Context) {
	request, err := bindJson[DTO.SignInRequest](c)
	if err != nil {
		respondError(c, err)
		return
	}
	jwtResponse, err := s.auth.SignIn(c.Request.Context(), request)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, jwtResponse)
}

// TokenRefreshHandler returns new access and refresh tokens
//
//	@Summary		Token refresh
//	@Description	Returns new access and refresh tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		DTO.TokenRefreshRequest	true	"JWT access token"
//	@Success		200		{object}	DTO.JWTResponse
//	@Failure		400		{object}	DTO.ErrorResponse	"invalid input"
//	@Failure		500		{object}	DTO.ErrorResponse
//	@Router			/admin/token-refresh [post]
func (s *Server) TokenRefreshHandler(c *gin.Context) {
	request, err := bindJson[DTO.TokenRefreshRequest](c)
	if err != nil {
		respondError(c, err)
		return
	}
	jwtResponse, err := s.auth.TokenRefresh(c.Request.Context(), request)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, jwtResponse)
}

func (s *Server) SignOutHandler(c *gin.Context) {
	// TODO
}

// CreateProductHandler creates a new product
//
//	@Summary		Create product
//	@Description	Create a new product entry
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			body	body		DTO.ProductRequest	true	"Product data"
//	@Success		201		{object}	DTO.ProductResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		401		{object}	DTO.ErrorResponse
//	@Failure		409		{object}	DTO.ErrorResponse	"slug or name already exists"
//	@Failure		422		{object}	DTO.ErrorResponse	"brand or category not found"
//	@Failure		500		{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/products [post]
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

// GetAllProductHandler returns all products
//
//	@Summary		List products
//	@Description	Get a list of all products
//	@Tags			products
//	@Produce		json
//	@Success		200	{array}		DTO.ProductResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Failure		500	{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/products [get]
func (s *Server) GetAllProductHandler(c *gin.Context) {
	products, err := s.product.GetAll(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, nonNilSlice(products))
}

// GetProductHandler returns a product by ID
//
//	@Summary		Get product
//	@Description	Fetch a single product by its ID
//	@Tags			products
//	@Produce		json
//	@Param			id	path		int	true	"Product ID"
//	@Success		200	{object}	DTO.ProductResponse
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Failure		500	{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/products/{id} [get]
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

// UpdateProductHandler updates an existing product
//
//	@Summary		Update product
//	@Description	Update product data by ID
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Product ID"
//	@Param			body	body		DTO.ProductRequest	true	"Updated product data"
//	@Success		200		{object}	DTO.ProductResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		401		{object}	DTO.ErrorResponse
//	@Failure		404		{object}	DTO.ErrorResponse
//	@Failure		409		{object}	DTO.ErrorResponse	"slug or name already exists"
//	@Failure		422		{object}	DTO.ErrorResponse	"brand or category not found"
//	@Failure		500		{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/products/{id} [put]
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

// DeleteProductHandler deletes a product by ID
//
//	@Summary		Delete product
//	@Description	Permanently delete a product by its ID
//	@Tags			products
//	@Produce		json
//	@Param			id	path	int	true	"Product ID"
//	@Success		204
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/products/{id} [delete]
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

// GetProductPriceHistory returns price history for a product
//
//	@Summary		Product price history
//	@Description	Get the full price change history for a product by ID
//	@Tags			products
//	@Produce		json
//	@Param			id	path		int	true	"Product ID"
//	@Success		200	{array}		queries.ProductPriceHistory
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/products/{id}/priceHistory [get]
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

// CreateBrandHandler creates a new brand
//
//	@Summary		Create brand
//	@Description	Create a new brand entry
//	@Tags			brands
//	@Accept			json
//	@Produce		json
//	@Param			body	body		DTO.BrandRequest	true	"Brand data"
//	@Success		201		{object}	DTO.BrandResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		401		{object}	DTO.ErrorResponse
//	@Failure		409		{object}	DTO.ErrorResponse	"name or slug already exists"
//	@Failure		500		{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/brands [post]
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

// GetAllBrandHandler returns all brands
//
//	@Summary		List brands
//	@Description	Get a list of all brands
//	@Tags			brands
//	@Produce		json
//	@Success		200	{array}		DTO.BrandResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/brands [get]
func (s *Server) GetAllBrandHandler(c *gin.Context) {
	brands, err := s.brand.GetAll(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, nonNilSlice(brands))
}

// GetBrandHandler returns a brand by ID
//
//	@Summary		Get brand
//	@Description	Fetch a single brand by its ID
//	@Tags			brands
//	@Produce		json
//	@Param			id	path		int	true	"Brand ID"
//	@Success		200	{object}	DTO.BrandResponse
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/brands/{id} [get]
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

// UpdateBrandHandler updates an existing brand
//
//	@Summary		Update brand
//	@Description	Update brand data by ID
//	@Tags			brands
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Brand ID"
//	@Param			body	body		DTO.BrandRequest	true	"Updated brand data"
//	@Success		200		{object}	DTO.BrandResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		401		{object}	DTO.ErrorResponse
//	@Failure		404		{object}	DTO.ErrorResponse
//	@Failure		409		{object}	DTO.ErrorResponse	"name or slug already exists"
//	@Failure		500		{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/brands/{id} [put]
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

// DeleteBrandHandler deletes a brand by ID
//
//	@Summary		Delete brand
//	@Description	Permanently delete a brand by its ID
//	@Tags			brands
//	@Produce		json
//	@Param			id	path	int	true	"Brand ID"
//	@Success		204
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/brands/{id} [delete]
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

// CreateCategoryHandler creates a new category
//
//	@Summary		Create category
//	@Description	Create a new category entry
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			body	body		DTO.CategoryRequest	true	"Category data"
//	@Success		201		{object}	DTO.CategoryResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		401		{object}	DTO.ErrorResponse
//	@Failure		409		{object}	DTO.ErrorResponse	"slug already exists"
//	@Failure		422		{object}	DTO.ErrorResponse	"parent category not found"
//	@Failure		500		{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/categories [post]
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

// GetAllCategoryHandler returns all categories
//
//	@Summary		List categories
//	@Description	Get a list of all categories
//	@Tags			categories
//	@Produce		json
//	@Success		200	{array}		DTO.CategoryResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/categories [get]
func (s *Server) GetAllCategoryHandler(c *gin.Context) {
	categories, err := s.category.GetAll(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, nonNilSlice(categories))
}

// GetCategoryHandler returns a category by ID
//
//	@Summary		Get category
//	@Description	Fetch a single category by its ID
//	@Tags			categories
//	@Produce		json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	DTO.CategoryResponse
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/categories/{id} [get]
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

// UpdateCategoryHandler updates an existing category
//
//	@Summary		Update category
//	@Description	Update category data by ID
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Category ID"
//	@Param			body	body		DTO.CategoryRequest	true	"Updated category data"
//	@Success		200		{object}	DTO.CategoryResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		401		{object}	DTO.ErrorResponse
//	@Failure		404		{object}	DTO.ErrorResponse
//	@Failure		409		{object}	DTO.ErrorResponse	"slug already exists"
//	@Failure		422		{object}	DTO.ErrorResponse	"parent category not found"
//	@Failure		500		{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/categories/{id} [put]
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

// DeleteCategoryHandler deletes a category by ID
//
//	@Summary		Delete category
//	@Description	Permanently delete a category by its ID
//	@Tags			categories
//	@Produce		json
//	@Param			id	path	int	true	"Category ID"
//	@Success		204
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/categories/{id} [delete]
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
