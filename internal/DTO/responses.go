package DTO

import "time"

type BrandResponse struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type CategoryResponse struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	ParentId  *int       `json:"parent_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
