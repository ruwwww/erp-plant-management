package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type assemblyRepository struct {
	*GormRepository[domain.StockAssembly]
}

func NewAssemblyRepository(db *gorm.DB) AssemblyRepository {
	return &assemblyRepository{NewGormRepository[domain.StockAssembly](db)}
}

func (r *assemblyRepository) GetRecipe(ctx context.Context, variantID int) ([]domain.ProductRecipe, error) {
	var recipes []domain.ProductRecipe
	err := r.DB.WithContext(ctx).
		Where("parent_variant_id = ?", variantID).
		Preload("ChildVariant"). // Need details of components
		Find(&recipes).Error
	return recipes, err
}

func (r *assemblyRepository) GetAllRecipes(ctx context.Context) ([]domain.ProductRecipe, error) {
	var recipes []domain.ProductRecipe
	err := r.DB.WithContext(ctx).
		Preload("ParentVariant").
		Preload("ChildVariant").
		Find(&recipes).Error
	return recipes, err
}

func (r *assemblyRepository) GetLogs(ctx context.Context, page, limit int) ([]domain.StockAssembly, error) {
	offset := (page - 1) * limit
	var logs []domain.StockAssembly
	err := r.DB.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, err
}
