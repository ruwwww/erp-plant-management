package dto

// --- Catalog ---
type ProductFilterQuery struct {
	CategorySlug string  `query:"category"`
	MinPrice     float64 `query:"min_price"`
	MaxPrice     float64 `query:"max_price"`
	Sort         string  `query:"sort"` // "price_asc", "newest"
	Search       string  `query:"q"`
	Page         int     `query:"page"`
	Limit        int     `query:"limit"`
}

// --- Cart & Checkout ---
type CartSyncRequest struct {
	Items []CartItemRequest `json:"items" validate:"required,dive"`
}

type CartItemRequest struct {
	VariantID int `json:"variant_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

type ApplyCouponRequest struct {
	Code string `json:"code" validate:"required"`
}

type CheckoutPreviewRequest struct {
	Items             []CartItemRequest `json:"items" validate:"required,dive"`
	CouponCode        string            `json:"coupon_code"`
	ShippingAddressID int               `json:"shipping_address_id" validate:"required"`
}

type CheckoutPlaceRequest struct {
	Items             []CartItemRequest `json:"items" validate:"required,dive"`
	CouponCode        string            `json:"coupon_code"`
	ShippingAddressID int               `json:"shipping_address_id" validate:"required"`
	BillingAddressID  int               `json:"billing_address_id"` // Optional, defaults to shipping
	PaymentMethod     string            `json:"payment_method" validate:"required,oneof=STRIPE MIDTRANS"`
}

// Response
type CheckoutPreviewResponse struct {
	Subtotal              float64 `json:"subtotal"`
	DiscountAmount        float64 `json:"discount_amount"`
	TaxAmount             float64 `json:"tax_amount"`
	ShippingAmount        float64 `json:"shipping_amount"`
	TotalAmount           float64 `json:"total_amount"`
	EstimatedDeliveryDate string  `json:"estimated_delivery_date"`
}
