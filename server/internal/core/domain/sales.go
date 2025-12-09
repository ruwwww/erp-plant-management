package domain

import (
	"time"
)

type SalesOrder struct {
	ID                      int                    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNumber             string                 `gorm:"unique;not null;size:64" json:"order_number"`
	CustomerID              *int                   `json:"customer_id"`
	Customer                *Customer              `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	GuestEmail              *string                `gorm:"size:320" json:"guest_email"`
	Channel                 OrderChannel           `gorm:"not null;default:'WEB'" json:"channel"`
	Status                  OrderStatus            `gorm:"not null;default:'PENDING'" json:"status"`
	PaymentStatus           PaymentStatus          `gorm:"not null;default:'UNPAID'" json:"payment_status"`
	ShipmentStatus          ShipmentStatus         `gorm:"not null;default:'PENDING'" json:"shipment_status"`
	POSSessionID            *int                   `json:"pos_session_id"`
	ShippingAddressSnapshot map[string]interface{} `gorm:"type:jsonb" json:"shipping_address_snapshot"`
	BillingAddressSnapshot  map[string]interface{} `gorm:"type:jsonb" json:"billing_address_snapshot"`
	ExpiresAt               *time.Time             `json:"expires_at"`
	SubtotalAmount          float64                `gorm:"not null;type:decimal(12,2);default:0" json:"subtotal_amount"`
	ShippingAmount          float64                `gorm:"not null;type:decimal(12,2);default:0" json:"shipping_amount"`
	TaxAmount               float64                `gorm:"not null;type:decimal(12,2);default:0" json:"tax_amount"`
	DiscountAmount          float64                `gorm:"not null;type:decimal(12,2);default:0" json:"discount_amount"`
	TotalAmount             float64                `gorm:"not null;type:decimal(14,2);default:0" json:"total_amount"`
	PlacedAt                time.Time              `gorm:"not null;default:current_timestamp" json:"placed_at"`
	CreatedBy               *int                   `json:"created_by"`
	CreatedAt               time.Time              `json:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at"`
	DeletedAt               *time.Time             `json:"deleted_at"`
	Items                   []SalesOrderItem       `gorm:"foreignKey:SalesOrderID" json:"items,omitempty"`
}

type SalesOrderItem struct {
	ID           int             `gorm:"primaryKey;autoIncrement" json:"id"`
	SalesOrderID int             `gorm:"not null" json:"sales_order_id"`
	VariantID    int             `gorm:"not null" json:"variant_id"`
	Variant      *ProductVariant `gorm:"foreignKey:VariantID" json:"variant"`
	ProductName  string          `gorm:"not null;size:255" json:"product_name"`
	SKU          string          `gorm:"not null;size:64" json:"sku"`
	Quantity     int             `gorm:"not null;default:1" json:"quantity"`
	UnitPrice    float64         `gorm:"not null;type:decimal(12,2);default:0" json:"unit_price"`
	TaxRate      float64         `gorm:"not null;type:decimal(5,4);default:0" json:"tax_rate"`
	TaxAmount    float64         `gorm:"not null;type:decimal(12,2);default:0" json:"tax_amount"`
	LineTotal    float64         `gorm:"not null;type:decimal(14,2);default:0" json:"line_total"`
}

type Shipment struct {
	ID             int            `gorm:"primaryKey;autoIncrement" json:"id"`
	SalesOrderID   int            `gorm:"not null" json:"sales_order_id"`
	SalesOrder     *SalesOrder    `gorm:"foreignKey:SalesOrderID" json:"sales_order"`
	Carrier        *string        `gorm:"size:100" json:"carrier"`
	TrackingNumber *string        `gorm:"unique;size:150" json:"tracking_number"`
	Status         ShipmentStatus `gorm:"not null;default:'PENDING'" json:"status"`
	ShippedAt      *time.Time     `json:"shipped_at"`
	DeliveredAt    *time.Time     `json:"delivered_at"`
	CreatedAt      time.Time      `json:"created_at"`
}

type Return struct {
	ID           int         `gorm:"primaryKey;autoIncrement" json:"id"`
	SalesOrderID int         `gorm:"not null" json:"sales_order_id"`
	SalesOrder   *SalesOrder `gorm:"foreignKey:SalesOrderID" json:"sales_order"`
	Reason       *string     `gorm:"size:255" json:"reason"`
	Status       string      `gorm:"not null;default:'REQUESTED';size:50" json:"status"`
	RequestedAt  time.Time   `gorm:"not null;default:current_timestamp" json:"requested_at"`
	ProcessedAt  *time.Time  `json:"processed_at"`
	RefundAmount *float64    `gorm:"type:decimal(12,2);default:0" json:"refund_amount"`
}
