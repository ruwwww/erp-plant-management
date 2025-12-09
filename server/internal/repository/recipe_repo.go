package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type RecipeRepository interface {
	Repository[domain.ProductRecipe]
	GetRecipe(ctx context.Context, variantID int) ([]domain.ProductRecipe, error)
}

type recipeRepository struct {
	*GormRepository[domain.ProductRecipe]
}

func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{NewGormRepository[domain.ProductRecipe](db)}
}

func (r *recipeRepository) GetRecipe(ctx context.Context, variantID int) ([]domain.ProductRecipe, error) {
	var recipes []domain.ProductRecipe
	err := r.DB.WithContext(ctx).
		Where("parent_variant_id = ?", variantID).
		Preload("ChildVariant").
		Find(&recipes).Error
	return recipes, err
}
