package DTO

import "time"

type BrandResponse struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryResponse struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ParentId  *int64    `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductResponse struct {
	Id          int64     `json:"id"`
	BrandId     int64     `json:"brand_id"`
	CategoryId  int64     `json:"category_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	PriceKopeck int32     `json:"price_kopeck"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
