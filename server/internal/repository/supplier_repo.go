package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type supplierRepository struct {
	*GormRepository[domain.Supplier]
}

func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &supplierRepository{NewGormRepository[domain.Supplier](db)}
}

func (r *supplierRepository) Restore(ctx context.Context, id int) error {
	var supplier domain.Supplier
	if err := r.DB.WithContext(ctx).Unscoped().First(&supplier, id).Error; err != nil {
		return err
	}
	return r.DB.WithContext(ctx).Unscoped().Model(&supplier).Update("DeletedAt", nil).Error
}

func (r *supplierRepository) ForceDelete(ctx context.Context, id int) error {
	var supplier domain.Supplier
	return r.DB.WithContext(ctx).Unscoped().Delete(&supplier, id).Error
}
