package dto

// Sessions
type OpenSessionRequest struct {
	OpeningCash float64 `json:"opening_cash" validate:"gte=0"`
}

type CloseSessionRequest struct {
	ClosingCashActual float64 `json:"closing_cash_actual" validate:"gte=0"`
	Note              string  `json:"note"`
}

type CashMoveRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Type   string  `json:"type" validate:"required,oneof=ADD DROP"`
	Reason string  `json:"reason" validate:"required"`
}

// Orders
type CreatePOSOrderRequest struct {
	POSSessionID     int               `json:"pos_session_id" validate:"required"`
	CustomerID       *int              `json:"customer_id"` // Optional (Walk-in)
	Items            []CartItemRequest `json:"items" validate:"required,dive"`
	PaymentMethod    string            `json:"payment_method" validate:"required"` // CASH, EDC, QRIS
	DiscountOverride *float64          `json:"discount_override"`                  // Manager only
}

type OverridePriceRequest struct {
	VariantID  int     `json:"variant_id" validate:"required"`
	NewPrice   float64 `json:"new_price" validate:"required,gte=0"`
	Reason     string  `json:"reason" validate:"required"`
	ManagerPIN string  `json:"manager_pin" validate:"required"`
}

type CustomerSearchRequest struct {
	Query string `json:"query" validate:"required,min=3"` // Name or Phone
}
