package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

func NewAssemblyRepository(db *gorm.DB) ports.AssemblyRepository {
	return &AssemblyRepository{NewGormRepository[domain.StockAssembly](db)}
}

func (r *AssemblyRepository) GetRecipe(ctx context.Context, variantID int) ([]domain.ProductRecipe, error) {
	var recipes []domain.ProductRecipe
	err := r.DB.WithContext(ctx).
		Where("parent_variant_id = ?", variantID).
		Preload("ChildVariant"). // Need details of components
		Find(&recipes).Error
	return recipes, err
}
