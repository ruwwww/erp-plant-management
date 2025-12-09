package dto

import "time"

type ProductFilterParams struct {
	Search        string
	CategorySlug  string
	Tags          []string // Tag slugs for filtering
	MinPrice      float64
	MaxPrice      float64
	IsActive      *bool // nil = all (admin), true = active only (customer)
	Page          int
	Limit         int
	CreatedAfter  *time.Time // gte (Start Date)
	CreatedBefore *time.Time // lte (End Date)
}

type OrderFilterParams struct {
	DateFrom *time.Time
	DateTo   *time.Time
	Status   string
	Search   string // Order Number or Customer Name
	Page     int
	Limit    int
}

// Tags
type CreateTagRequest struct {
	Name        string  `json:"name" validate:"required"`
	Slug        string  `json:"slug" validate:"required"`
	DisplayName *string `json:"display_name"`
}

type UpdateProductTagsRequest struct {
	TagIDs []int `json:"tag_ids" validate:"required"`
}
