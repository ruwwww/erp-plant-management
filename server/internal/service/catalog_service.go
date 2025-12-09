package service

import (
	"context"
	"server/internal/core/domain"
	"server/internal/repository"
)

type CatalogService struct {
	productRepo  repository.Repository[domain.Product]
	categoryRepo repository.Repository[domain.Category]
}

func NewCatalogService(productRepo repository.Repository[domain.Product], categoryRepo repository.Repository[domain.Category]) *CatalogService {
	return &CatalogService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *CatalogService) GetProducts(ctx context.Context) ([]domain.Product, error) {
	return s.productRepo.FindAll(ctx)
}

func (s *CatalogService) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	return s.productRepo.FindByID(ctx, id)
}

func (s *CatalogService) CreateProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepo.Create(ctx, product)
}
