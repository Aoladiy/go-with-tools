package DTO

type BrandRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CategoryRequest struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ParentId *int64 `json:"parent_id,omitempty"`
}

type ProductRequest struct {
	BrandId     int64   `json:"brand_id"`
	CategoryId  int64   `json:"category_id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
	PriceKopeck int32   `json:"price_kopeck"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
