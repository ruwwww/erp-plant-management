package domain

import "time"

type ProductRecipe struct {
	ID int `gorm:"primaryKey;autoIncrement" json:"id"`

	ParentVariantID int             `gorm:"not null" json:"parent_variant_id"`
	ParentVariant   *ProductVariant `gorm:"foreignKey:ParentVariantID" json:"parent_variant,omitempty"`

	ChildVariantID int             `gorm:"not null" json:"child_variant_id"`
	ChildVariant   *ProductVariant `gorm:"foreignKey:ChildVariantID" json:"child_variant,omitempty"`

	QuantityNeeded float64 `gorm:"not null;type:decimal(8,3);default:1" json:"quantity_needed"`

	CreatedAt time.Time `gorm:"not null;default:current_timestamp" json:"created_at"`
}

func (ProductRecipe) TableName() string {
	return "product_recipes"
}

type StockAssembly struct {
	ID int `gorm:"primaryKey;autoIncrement" json:"id"`

	AssemblyNumber string `gorm:"unique;not null;size:64" json:"assembly_number"`

	VariantID int             `gorm:"not null" json:"variant_id"`
	Variant   *ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`

	QuantityProduced int `gorm:"not null" json:"quantity_produced"`

	TotalCost float64 `gorm:"not null;type:decimal(12,2)" json:"total_cost"`

	CreatedBy *int  `json:"created_by"`
	User      *User `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`

	CreatedAt time.Time `gorm:"not null;default:current_timestamp" json:"created_at"`
}

func (StockAssembly) TableName() string {
	return "stock_assemblies"
}
