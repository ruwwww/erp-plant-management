package dto

import "time"

type ProductFilterParams struct {
	Search        string
	CategorySlug  string
	MinPrice      float64
	MaxPrice      float64
	IsActive      *bool // nil = all (admin), true = active only (customer)
	Page          int
	Limit         int
	CreatedAfter  *time.Time // gte (Start Date)
	CreatedBefore *time.Time // lte (End Date)
}
