package domain

import (
	"time"
)

type Category struct {
	ID          int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"not null;size:150" json:"name"`
	Slug        string     `gorm:"unique;not null;size:180" json:"slug"`
	ParentID    *int       `json:"parent_id"`
	Parent      *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type Product struct {
	ID          int              `gorm:"primaryKey;autoIncrement" json:"id"`
	SKU         string           `gorm:"unique;not null;size:64" json:"sku"`
	Name        string           `gorm:"not null;size:255" json:"name"`
	Slug        string           `gorm:"unique;not null;size:255" json:"slug"`
	Description *string          `json:"description"`
	CategoryID  *int             `json:"category_id"`
	Category    *Category        `gorm:"foreignKey:CategoryID" json:"category"`
	SupplierID  *int             `json:"supplier_id"`
	Supplier    *Supplier        `gorm:"foreignKey:SupplierID" json:"supplier"`
	Condition   ProductCondition `gorm:"not null;default:'NEW'" json:"condition"`
	IsActive    bool             `gorm:"not null;default:true" json:"is_active"`
	IsFeatured  bool             `gorm:"not null;default:false" json:"is_featured"`
	BasePrice   float64          `gorm:"not null;type:decimal(12,2);default:0" json:"base_price"`
	TaxClass    *string          `gorm:"size:50" json:"tax_class"`
	WeightKG    *float64         `gorm:"type:decimal(8,3)" json:"weight_kg"`
	HeightCM    *float64         `gorm:"type:decimal(8,3)" json:"height_cm"`
	WidthCM     *float64         `gorm:"type:decimal(8,3)" json:"width_cm"`
	DepthCM     *float64         `gorm:"type:decimal(8,3)" json:"depth_cm"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   *time.Time       `json:"deleted_at"`
}

type ProductVariant struct {
	ID             int                    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID      int                    `gorm:"not null" json:"product_id"`
	Product        *Product               `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	SKU            string                 `gorm:"not null;size:64" json:"sku"`
	Name           *string                `gorm:"size:255" json:"name"`
	Attributes     map[string]interface{} `gorm:"type:jsonb;not null;default:'{}'" json:"attributes"`
	Price          float64                `gorm:"not null;type:decimal(12,2);default:0" json:"price"`
	CompareAtPrice *float64               `gorm:"type:decimal(12,2)" json:"compare_at_price"`
	Barcode        *string                `gorm:"size:64" json:"barcode"`
	StockControl   bool                   `gorm:"not null;default:true" json:"stock_control"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	DeletedAt      *time.Time             `json:"deleted_at"`
}

type Tag struct {
	ID          int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string  `gorm:"unique;not null;size:100" json:"name"`
	Slug        string  `gorm:"unique;not null;size:120" json:"slug"`
	DisplayName *string `gorm:"size:150" json:"display_name"`
}

// Join Table for Product <-> Tags
type ProductTag struct {
	ProductID int `gorm:"primaryKey" json:"product_id"`
	TagID     int `gorm:"primaryKey" json:"tag_id"`
}
