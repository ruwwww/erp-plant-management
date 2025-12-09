package service

import (
	"context"
	"encoding/json"
	"errors"
	"server/internal/core/domain"
	"server/internal/dto"
	"server/internal/repository"

	"gorm.io/gorm"
)

type CatalogServiceImpl struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	variantRepo  repository.VariantRepository
	tagRepo      repository.TagRepository
	db           *gorm.DB
}

func NewCatalogService(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	variantRepo repository.VariantRepository,
	tagRepo repository.TagRepository,
	db *gorm.DB,

) CatalogService {
	return &CatalogServiceImpl{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		variantRepo:  variantRepo,
		tagRepo:      tagRepo,
		db:           db,
	}
}

func (s *CatalogServiceImpl) GetProducts(ctx context.Context, filter dto.ProductFilterParams) ([]domain.Product, int64, error) {
	return s.productRepo.Search(ctx, filter)
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

func (s *CatalogServiceImpl) UpdateProduct(ctx context.Context, id int, req dto.UpdateProductRequest) error {

	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Name != nil {
		product.Name = *req.Name
	}

	if req.Description != nil {
		product.Description = req.Description
	}

	if req.BasePrice != nil {
		if *req.BasePrice < 0 {
			return errors.New("price cannot be negative")
		}
		product.BasePrice = *req.BasePrice
	}

	if req.WeightKG != nil {
		product.WeightKG = req.WeightKG
	}

	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}

	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

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
	return s.productRepo.SoftDelete(ctx, id)
}

func (s *CatalogServiceImpl) RestoreProduct(ctx context.Context, id int) error {
	return s.productRepo.Restore(ctx, id)
}

func (s *CatalogServiceImpl) ForceDeleteProduct(ctx context.Context, id int) error {
	return s.productRepo.ForceDelete(ctx, id)
}

func (s *CatalogServiceImpl) SoftDeleteVariant(ctx context.Context, id int) error {
	return s.variantRepo.SoftDelete(ctx, id)
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

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&products).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *CatalogServiceImpl) ExportProducts(ctx context.Context) ([]byte, error) {
	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return json.Marshal(products)
}

// Tags
func (s *CatalogServiceImpl) GetTags(ctx context.Context) ([]domain.Tag, error) {
	return s.tagRepo.FindAll(ctx)
}

func (s *CatalogServiceImpl) CreateTag(ctx context.Context, tag *domain.Tag) error {
	// Basic validation
	if tag.Name == "" {
		return errors.New("tag name is required")
	}
	if tag.Slug == "" {
		return errors.New("tag slug is required")
	}
	return s.tagRepo.Create(ctx, tag)
}

func (s *CatalogServiceImpl) GetTagBySlug(ctx context.Context, slug string) (*domain.Tag, error) {
	return s.tagRepo.FindBySlug(ctx, slug)
}

func (s *CatalogServiceImpl) UpdateProductTags(ctx context.Context, productID int, tagIDs []int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Clear existing tags
		if err := tx.Where("product_id = ?", productID).Delete(&domain.ProductTag{}).Error; err != nil {
			return err
		}
		// Add new tags
		for _, tagID := range tagIDs {
			if err := tx.Create(&domain.ProductTag{ProductID: productID, TagID: tagID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
