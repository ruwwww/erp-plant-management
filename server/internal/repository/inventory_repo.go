package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type inventoryRepository struct {
	*GormRepository[domain.Stock]
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{NewGormRepository[domain.Stock](db)}
}

func (r *inventoryRepository) GetStock(ctx context.Context, variantID, locationID int) (*domain.Stock, error) {
	var stock domain.Stock

	err := r.DB.WithContext(ctx).
		Where("variant_id = ? AND location_id = ?", variantID, locationID).
		First(&stock).Error
	return &stock, err
}

func (r *inventoryRepository) UpdateStockAtomic(ctx context.Context, variantID, locationID, qtyChange int) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Stock{}).
		Where("variant_id = ? AND location_id = ?", variantID, locationID).
		UpdateColumn("quantity_on_hand", gorm.Expr("quantity_on_hand + ?", qtyChange)).
		Error
}
