package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	*GormRepository[domain.Stock]
}

func NewInventoryRepository(db *gorm.DB) ports.InventoryRepository {
	return &InventoryRepository{NewGormRepository[domain.Stock](db)}
}

func (r *InventoryRepository) GetStock(ctx context.Context, variantID, locationID int) (*domain.Stock, error) {
	var stock domain.Stock
	// Use Find to avoid "Record Not Found" error, we want to return 0 quantity if not found,
	// but strictly speaking FindOne handles not found by returning error.
	// Let's stick to standard behavior.
	err := r.DB.WithContext(ctx).
		Where("variant_id = ? AND location_id = ?", variantID, locationID).
		First(&stock).Error
	return &stock, err
}

// UpdateStockAtomic performs SQL-level arithmetic to avoid race conditions
// UPDATE stocks SET quantity_on_hand = quantity_on_hand + X WHERE ...
func (r *InventoryRepository) UpdateStockAtomic(ctx context.Context, variantID, locationID, qtyChange int) error {
	// GORM's UpdateColumn skips Hooks (UpdatedAt), so be careful,
	// but it is the fastest and safest way for atomic updates.
	return r.DB.WithContext(ctx).
		Model(&domain.Stock{}).
		Where("variant_id = ? AND location_id = ?", variantID, locationID).
		UpdateColumn("quantity_on_hand", gorm.Expr("quantity_on_hand + ?", qtyChange)).
		Error
}
