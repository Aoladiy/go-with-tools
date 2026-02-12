package DTO

import "time"

type BrandResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ParentId  *int      `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductResponse struct {
	Id          int       `json:"id"`
	BrandId     int       `json:"brand_id"`
	CategoryId  int       `json:"category_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	PriceKopeck int       `json:"price_kopeck"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
