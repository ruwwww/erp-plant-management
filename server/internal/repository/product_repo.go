package repository

import (
	"context"
	"server/internal/core/domain"
	"server/internal/core/ports"

	"gorm.io/gorm"
)

type ProductRepository struct {
	*GormRepository[domain.Product]
}

func NewProductRepository(db *gorm.DB) ports.ProductRepository {
	return &ProductRepository{NewGormRepository[domain.Product](db)}
}

func (r *ProductRepository) GetFullProduct(ctx context.Context, slug string) (*domain.Product, error) {
	var product domain.Product
	// Preload necessary relations
	err := r.DB.WithContext(ctx).
		Preload("Category").
		Preload("Supplier").
		Where("slug = ?", slug).
		First(&product).Error
	return &product, err
}

func (r *ProductRepository) Search(ctx context.Context, filter map[string]interface{}, page, limit int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.DB.WithContext(ctx).Model(&domain.Product{})

	// Dynamic Filtering
	if search, ok := filter["search"].(string); ok && search != "" {
		query = query.Where("name ILIKE ? OR sku ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if catSlug, ok := filter["category_slug"].(string); ok && catSlug != "" {
		// Sub-query logic or simple join could go here
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.slug = ?", catSlug)
	}
	if min, ok := filter["min_price"].(float64); ok {
		query = query.Where("base_price >= ?", min)
	}
	if active, ok := filter["is_active"].(bool); ok {
		query = query.Where("is_active = ?", active)
	}

	query.Count(&total)

	// Pagination
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&products).Error

	return products, total, err
}
