package repository

import (
	"context"
	"server/internal/core/domain"
	"server/internal/dto"

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

	err := r.DB.WithContext(ctx).
		Preload("Category").
		Preload("Supplier").
		Where("slug = ?", slug).
		First(&product).Error
	return &product, err
}

func (r *productRepository) Search(ctx context.Context, filter dto.ProductFilterParams) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.DB.WithContext(ctx).Model(&domain.Product{}).
		Preload("Category").
		Preload("Supplier")

	if filter.Search != "" {
		query = query.Where("name ILIKE ? OR sku ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	if filter.CategorySlug != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.slug = ?", filter.CategorySlug)
	}

	if filter.MinPrice > 0 {
		query = query.Where("base_price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("base_price <= ?", filter.MaxPrice)
	}

	if filter.CreatedAfter != nil {
		query = query.Where("created_at >= ?", filter.CreatedAfter)
	}

	if filter.CreatedBefore != nil {
		query = query.Where("created_at <= ?", filter.CreatedBefore)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}

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

	if err := r.DB.WithContext(ctx).Unscoped().First(&product, id).Error; err != nil {
		return err
	}

	return r.DB.WithContext(ctx).Unscoped().Model(&product).Update("DeletedAt", nil).Error
}

func (r *productRepository) ForceDelete(ctx context.Context, id int) error {
	var product domain.Product

	return r.DB.WithContext(ctx).Unscoped().Delete(&product, id).Error
}
