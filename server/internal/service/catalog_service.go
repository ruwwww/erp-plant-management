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

func NewCatalogService(productRepo repository.Repository[domain.Product], categoryRepo repository.Repository[domain.Category]) *CatalogServiceImpl {
	return &CatalogServiceImpl{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *CatalogServiceImpl) GetProducts(ctx context.Context) ([]domain.Product, error) {
	return s.productRepo.FindAll(ctx)
}

func (s *CatalogServiceImpl) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	return s.productRepo.FindByID(ctx, id)
}

func (s *CatalogServiceImpl) CreateProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepo.Create(ctx, product)
}
