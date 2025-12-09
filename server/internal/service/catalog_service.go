package service

import (
	"context"
	"server/internal/core/domain"
	"server/internal/repository"
)

type CatalogServiceImpl struct {
	productRepo  repository.Repository[domain.Product]
	categoryRepo repository.Repository[domain.Category]
}

func NewCatalogService(productRepo repository.Repository[domain.Product], categoryRepo repository.Repository[domain.Category]) CatalogService {
	return &CatalogServiceImpl{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *CatalogServiceImpl) GetProducts(ctx context.Context, filter ProductFilterParams) ([]domain.Product, int64, error) {
	// TODO: Implement filtering and pagination
	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	return products, int64(len(products)), nil
}

func (s *CatalogServiceImpl) GetProductDetail(ctx context.Context, slug string) (*domain.Product, error) {
	return s.productRepo.FindOne(ctx, "slug = ?", slug)
}

func (s *CatalogServiceImpl) CreateProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepo.Create(ctx, product)
}

func (s *CatalogServiceImpl) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepo.Update(ctx, product)
}

func (s *CatalogServiceImpl) SoftDeleteProduct(ctx context.Context, id int) error {
	// TODO: Implement soft delete logic (set DeletedAt)
	return s.productRepo.Delete(ctx, id)
}

func (s *CatalogServiceImpl) RestoreProduct(ctx context.Context, id int) error {
	// TODO: Implement restore logic
	return nil
}
