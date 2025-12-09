package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type PurchaseOrderRepository struct {
	*GormRepository[domain.PurchaseOrder]
}

func NewPurchaseOrderRepository(db *gorm.DB) ports.PurchaseOrderRepository {
	return &PurchaseOrderRepository{NewGormRepository[domain.PurchaseOrder](db)}
}

func (r *PurchaseOrderRepository) GetFullPO(ctx context.Context, id int) (*domain.PurchaseOrder, error) {
	var po domain.PurchaseOrder
	err := r.DB.WithContext(ctx).
		Preload("Items").
		Preload("Supplier").
		First(&po, id).Error
	return &po, err
}
