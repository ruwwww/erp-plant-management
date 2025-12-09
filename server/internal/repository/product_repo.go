package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type productRepository struct {
	*GormRepository[domain.Product]
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{NewGormRepository[domain.Product](db)}
}

func (r *productRepository) GetFullProduct(ctx context.Context, slug string) (*domain.Product, error) {
	var product domain.Product
	// Preload necessary relations
	err := r.DB.WithContext(ctx).
		Preload("Category").
		Preload("Supplier").
		Where("slug = ?", slug).
		First(&product).Error
	return &product, err
}

func (r *productRepository) Search(ctx context.Context, filter map[string]interface{}, page, limit int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.DB.WithContext(ctx).Model(&domain.Product{}).Preload("Category").Preload("Supplier")

	if search, ok := filter["search"].(string); ok && search != "" {
		query = query.Where("name ILIKE ? OR sku ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if catSlug, ok := filter["category_slug"].(string); ok && catSlug != "" {
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

func (r *productRepository) SoftDelete(ctx context.Context, id int) error {
	var product domain.Product
	return r.DB.WithContext(ctx).Delete(&product, id).Error
}

func (r *productRepository) Restore(ctx context.Context, id int) error {
	var product domain.Product
	// Find the record (including soft deleted)
	if err := r.DB.WithContext(ctx).Unscoped().First(&product, id).Error; err != nil {
		return err
	}
	// Update DeletedAt to nil
	return r.DB.WithContext(ctx).Unscoped().Model(&product).Update("DeletedAt", nil).Error
}

// ForceDelete permanently deletes a product if only it has no sales records
func (r *productRepository) ForceDelete(ctx context.Context, id int) error {
	var product domain.Product

	return r.DB.WithContext(ctx).Unscoped().Delete(&product).Error
}
