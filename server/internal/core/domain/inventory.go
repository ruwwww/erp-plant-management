package domain

import "time"

type InventoryLocation struct {
	ID        int      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string   `gorm:"not null;size:100" json:"name"`
	Code      string   `gorm:"not null;size:50;unique" json:"code"`
	Type      string   `gorm:"not null;size:50" json:"type"` // WAREHOUSE, STORE
	AddressID *int     `json:"address_id"`
	Address   *Address `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	IsActive  bool     `gorm:"not null;default:true" json:"is_active"`
}

type Stock struct {
	LocationID  int                `gorm:"primaryKey" json:"location_id"`
	VariantID   int                `gorm:"primaryKey" json:"variant_id"`
	Location    *InventoryLocation `gorm:"foreignKey:LocationID" json:"location"`
	Variant     *ProductVariant    `gorm:"foreignKey:VariantID" json:"variant"`
	Quantity    int                `gorm:"not null;default:0" json:"quantity"`
	SafetyStock int                `gorm:"not null;default:0" json:"safety_stock"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type StockMovement struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	LocationID     int       `gorm:"not null" json:"location_id"`
	VariantID      int       `gorm:"not null" json:"variant_id"`
	QuantityChange int       `gorm:"not null" json:"quantity_change"`
	Reason         string    `gorm:"not null;size:50" json:"reason"` // SALE, RESTOCK, ADJUSTMENT
	ReferenceType  *string   `gorm:"size:50" json:"reference_type"`  // ORDER, PO
	ReferenceID    *int      `json:"reference_id"`
	CreatedBy      *int      `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}
