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

type OrderFilterParams struct {
	DateFrom *time.Time
	DateTo   *time.Time
	Status   string
	Search   string // Order Number or Customer Name
	Page     int
	Limit    int
}
