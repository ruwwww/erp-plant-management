package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type categoryRepository struct {
	*GormRepository[domain.Category]
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{NewGormRepository[domain.Category](db)}
}

func (r *categoryRepository) GetTree(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	// For now, just return all categories.
	// A proper tree implementation would require a recursive query or building the tree in memory.
	// Given the struct doesn't have Children, we just return the flat list.
	err := r.DB.WithContext(ctx).Preload("Parent").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) SoftDelete(ctx context.Context, id int) error {
	var category domain.Category
	return r.DB.WithContext(ctx).Delete(&category, id).Error
}

func (r *categoryRepository) Restore(ctx context.Context, id int) error {
	var category domain.Category
	if err := r.DB.WithContext(ctx).Unscoped().First(&category, id).Error; err != nil {
		return err
	}
	return r.DB.WithContext(ctx).Unscoped().Model(&category).Update("DeletedAt", nil).Error
}

func (r *categoryRepository) ForceDelete(ctx context.Context, id int) error {
	var category domain.Category
	return r.DB.WithContext(ctx).Unscoped().Delete(&category, id).Error
}
