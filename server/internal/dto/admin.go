package dto

import (
	"server/internal/core/domain"
)

// --- Products ---
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3"`
	SKU         string  `json:"sku" validate:"required,alphanum"`
	Slug        string  `json:"slug" validate:"required,slug"`
	Description string  `json:"description"`
	CategoryID  int     `json:"category_id" validate:"required"`
	SupplierID  int     `json:"supplier_id"`
	BasePrice   float64 `json:"base_price" validate:"required,gte=0"`
	WeightKG    float64 `json:"weight_kg"`
	// Initial Variant
	StockControl bool `json:"stock_control"`
}

func (r *CreateProductRequest) ToDomain() *domain.Product {
	p := &domain.Product{
		Name:      r.Name,
		SKU:       r.SKU,
		Slug:      r.Slug,
		BasePrice: r.BasePrice,
		IsActive:  true,
	}

	if r.Description != "" {
		desc := r.Description
		p.Description = &desc
	}

	if r.WeightKG > 0 {
		val := r.WeightKG
		p.WeightKG = &val
	}

	if r.CategoryID != 0 {
		id := r.CategoryID
		p.CategoryID = &id
	}
	if r.SupplierID != 0 {
		id := r.SupplierID
		p.SupplierID = &id
	}

	return p
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	BasePrice   *float64 `json:"base_price"`
	CategoryID  *int     `json:"category_id"`
	IsActive    *bool    `json:"is_active"`
	WeightKG    *float64 `json:"weight_kg"`
}

type ProductVariantRequest struct {
	SKU        string                 `json:"sku" validate:"required"`
	Name       string                 `json:"name"`
	Price      float64                `json:"price" validate:"gte=0"`
	Attributes map[string]interface{} `json:"attributes"`
}

// --- Users (Admin Manage) ---
type CreateStaffRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
	Role      string `json:"role" validate:"required,oneof=ADMIN MANAGER STAFF"`
	Password  string `json:"password" validate:"required,min=8"`
}

type UpdateUserStatusRequest struct {
	IsActive bool `json:"is_active"` // Ban/Unban
}

type AssignRoleRequest struct {
	Roles []string `json:"roles" validate:"required"` // RBAC
}

type AdminResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// --- Customers (CRM) ---
type CreateSegmentRequest struct {
	Name       string                 `json:"name" validate:"required"`
	Conditions map[string]interface{} `json:"conditions" validate:"required"` // e.g., {"total_spent_gt": 1000}
}

type TriggerCampaignRequest struct {
	SegmentID int    `json:"segment_id" validate:"required"`
	Subject   string `json:"subject" validate:"required"`
	Body      string `json:"body" validate:"required"` // HTML allowed
}

// --- Promotions ---
type CreatePromotionRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Code        string                 `json:"code" validate:"alphanum"`
	IsExclusive bool                   `json:"is_exclusive"`
	Conditions  map[string]interface{} `json:"conditions"` // e.g., {"min_cart_total": 50000}
	Actions     map[string]interface{} `json:"actions"`    // e.g., {"percent_off": 10}
	StartsAt    string                 `json:"starts_at"`
	EndsAt      string                 `json:"ends_at"`
}

// --- Media ---
type LinkMediaRequest struct {
	MediaIDs   []int  `json:"media_ids" validate:"required"`
	EntityType string `json:"entity_type" validate:"required,oneof=products categories"`
	EntityID   int    `json:"entity_id" validate:"required"`
	Zone       string `json:"zone" validate:"required"` // e.g. "gallery"
}

type EditMediaRequest struct {
	AltText string `json:"alt_text"`
	Title   string `json:"title"`
}

// --- Data Import/Export ---
// Usually these rely on Multipart Forms (CSV file) rather than JSON body
// But if JSON import is supported:
type BulkImportProductsRequest struct {
	Items []CreateProductRequest `json:"items" validate:"dive"`
}

// --- Suppliers ---
type CreateSupplierRequest struct {
	Name        string `json:"name" validate:"required"`
	ContactName string `json:"contact_name"`
	Email       string `json:"email" validate:"omitempty,email"`
	Phone       string `json:"phone"`
	// Address fields could be nested or separate
}

type UpdateSupplierRequest struct {
	Name        string `json:"name"`
	ContactName string `json:"contact_name"`
	Email       string `json:"email" validate:"omitempty,email"`
	Phone       string `json:"phone"`
}
