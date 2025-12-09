package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type purchaseOrderRepository struct {
	*GormRepository[domain.PurchaseOrder]
}

func NewPurchaseOrderRepository(db *gorm.DB) PurchaseOrderRepository {
	return &purchaseOrderRepository{NewGormRepository[domain.PurchaseOrder](db)}
}

func (r *purchaseOrderRepository) GetFullPO(ctx context.Context, id int) (*domain.PurchaseOrder, error) {
	var po domain.PurchaseOrder
	err := r.DB.WithContext(ctx).
		Preload("Items").
		Preload("Supplier").
		First(&po, id).Error
	return &po, err
}
