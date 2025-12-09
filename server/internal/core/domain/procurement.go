package domain

import "time"

type PurchaseOrder struct {
	ID             int                 `gorm:"primaryKey;autoIncrement" json:"id"`
	PONumber       string              `gorm:"unique;not null;size:64" json:"po_number"`
	SupplierID     int                 `gorm:"not null" json:"supplier_id"`
	Supplier       *Supplier           `gorm:"foreignKey:SupplierID" json:"supplier"`
	Status         PurchaseOrderStatus `gorm:"not null;default:'DRAFT'" json:"status"`
	ExpectedAt     *time.Time          `json:"expected_at"`
	Notes          *string             `json:"notes"`
	SubtotalAmount float64             `gorm:"not null;type:decimal(12,2);default:0" json:"subtotal_amount"`
	TaxAmount      float64             `gorm:"not null;type:decimal(12,2);default:0" json:"tax_amount"`
	ShippingAmount float64             `gorm:"not null;type:decimal(12,2);default:0" json:"shipping_amount"`
	TotalAmount    float64             `gorm:"not null;type:decimal(12,2);default:0" json:"total_amount"`
	CreatedBy      *int                `json:"created_by"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	Items          []PurchaseOrderItem `gorm:"foreignKey:PurchaseOrderID" json:"items,omitempty"`
}

type PurchaseOrderItem struct {
	ID                   int             `gorm:"primaryKey;autoIncrement" json:"id"`
	PurchaseOrderID      int             `gorm:"not null" json:"purchase_order_id"`
	VariantID            int             `gorm:"not null" json:"variant_id"`
	Variant              *ProductVariant `gorm:"foreignKey:VariantID" json:"variant"`
	QuantityOrdered      int             `gorm:"not null;default:1" json:"quantity_ordered"`
	QuantityReceived     int             `gorm:"not null;default:0" json:"quantity_received"`
	QuoteRefNumber       *string         `gorm:"size:64" json:"quote_ref_number"`
	QuoteValidUntil      *time.Time      `gorm:"type:date" json:"quote_valid_until"`
	VendorQuoteAttachUrl *string         `gorm:"size:512" json:"vendor_quote_attachment_url"`
	UnitCost             float64         `gorm:"not null;type:decimal(12,2);default:0" json:"unit_cost"`
	LineTotal            float64         `gorm:"not null;type:decimal(14,2);default:0" json:"line_total"`
}
