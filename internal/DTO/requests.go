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
