package DTO

type BrandRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CategoryRequest struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ParentId *int   `json:"parent_id,omitempty"`
}

type ProductRequest struct {
	BrandId     int     `json:"brand_id"`
	CategoryId  int     `json:"category_id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
	PriceKopeck int     `json:"price_kopeck"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
