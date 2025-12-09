package service

import (
	"context"
	"encoding/json"
	"server/internal/core/domain"
	"server/internal/repository"

	"gorm.io/gorm"
)

type CatalogServiceImpl struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	variantRepo  repository.VariantRepository
	db           *gorm.DB // Needed for transaction

}

func NewCatalogService(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	variantRepo repository.VariantRepository,
	db *gorm.DB,

) CatalogService {
	return &CatalogServiceImpl{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		variantRepo:  variantRepo,
		db:           db,
	}
}

func (s *CatalogServiceImpl) GetProducts(ctx context.Context, filter ProductFilterParams) ([]domain.Product, int64, error) {
	// TODO: Implement proper filtering
	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	return products, int64(len(products)), nil
}

func (s *CatalogServiceImpl) GetProductDetail(ctx context.Context, slug string) (*domain.Product, error) {
	return s.productRepo.FindOne(ctx, "slug = ?", slug)
}

func (s *CatalogServiceImpl) GetCategories(ctx context.Context) ([]domain.Category, error) {
	return s.categoryRepo.FindAll(ctx)
}

func (s *CatalogServiceImpl) GetVariants(ctx context.Context, productID int) ([]domain.ProductVariant, error) {
	return s.variantRepo.Find(ctx, "product_id = ?", productID)
}

func (s *CatalogServiceImpl) CreateProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepo.Create(ctx, product)
}

func (s *CatalogServiceImpl) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepo.Update(ctx, product)
}

func (s *CatalogServiceImpl) UpdateVariants(ctx context.Context, productID int, variants []domain.ProductVariant) error {
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var keepIDs []int
	for _, v := range variants {
		if v.ID != 0 {
			keepIDs = append(keepIDs, v.ID)
		}
	}

	if err := tx.Where("product_id = ? AND id NOT IN ?", productID, keepIDs).
		Delete(&domain.ProductVariant{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, v := range variants {
		v.ProductID = productID

		if v.ID != 0 {

			if err := tx.Save(&v).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.Create(&v).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (s *CatalogServiceImpl) SoftDeleteProduct(ctx context.Context, id int) error {
	// TODO: Implement soft delete logic (set DeletedAt)
	return s.productRepo.Delete(ctx, id)
}

func (s *CatalogServiceImpl) RestoreProduct(ctx context.Context, id int) error {
	return s.productRepo.Restore(ctx, id)
}

func (s *CatalogServiceImpl) ForceDeleteProduct(ctx context.Context, id int) error {
	return s.productRepo.ForceDelete(ctx, id)
}

func (s *CatalogServiceImpl) SoftDeleteVariant(ctx context.Context, id int) error {
	return s.variantRepo.Delete(ctx, id)
}

func (s *CatalogServiceImpl) RestoreVariant(ctx context.Context, id int) error {
	return s.variantRepo.Restore(ctx, id)
}

func (s *CatalogServiceImpl) ForceDeleteVariant(ctx context.Context, id int) error {
	return s.variantRepo.ForceDelete(ctx, id)
}

func (s *CatalogServiceImpl) ImportProducts(ctx context.Context, data []byte) error {
	var products []domain.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return err
	}
	for _, p := range products {
		if err := s.productRepo.Create(ctx, &p); err != nil {
			return err
		}
	}
	return nil
}

func (s *CatalogServiceImpl) ExportProducts(ctx context.Context) ([]byte, error) {
	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return json.Marshal(products)
}
