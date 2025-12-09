package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type variantRepository struct {
	*GormRepository[domain.ProductVariant]
}

func NewVariantRepository(db *gorm.DB) VariantRepository {
	return &variantRepository{NewGormRepository[domain.ProductVariant](db)}
}

func (r *variantRepository) Restore(ctx context.Context, id int) error {
	var variant domain.ProductVariant
	if err := r.DB.WithContext(ctx).Unscoped().First(&variant, id).Error; err != nil {
		return err
	}
	return r.DB.WithContext(ctx).Unscoped().Model(&variant).Update("DeletedAt", nil).Error
}

func (r *variantRepository) ForceDelete(ctx context.Context, id int) error {
	var variant domain.ProductVariant
	return r.DB.WithContext(ctx).Unscoped().Delete(&variant, id).Error
}
