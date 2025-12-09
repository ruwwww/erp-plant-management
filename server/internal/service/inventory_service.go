package service

import (
	"context"
	"errors"
	"server/internal/core/domain"
	"server/internal/repository"
	"time"

	"gorm.io/gorm"
)

type InventoryService struct {
	stockRepo    repository.Repository[domain.Stock]
	movementRepo repository.Repository[domain.StockMovement]
	db           *gorm.DB // Needed for transaction
}

func NewInventoryService(stockRepo repository.Repository[domain.Stock], movementRepo repository.Repository[domain.StockMovement], db *gorm.DB) *InventoryService {
	return &InventoryService{
		stockRepo:    stockRepo,
		movementRepo: movementRepo,
		db:           db,
	}
}

func (s *InventoryService) GetStock(ctx context.Context, locationID, variantID int) (*domain.Stock, error) {
	return s.stockRepo.FindOne(ctx, "location_id = ? AND variant_id = ?", locationID, variantID)
}

func (s *InventoryService) AdjustStock(ctx context.Context, locationID, variantID, qtyChange int, reason string, refType *string, refID *int, userID *int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create Movement
		movement := &domain.StockMovement{
			LocationID:     locationID,
			VariantID:      variantID,
			QuantityChange: qtyChange,
			Reason:         reason,
			ReferenceType:  refType,
			ReferenceID:    refID,
			CreatedBy:      userID,
			CreatedAt:      time.Now(),
		}
		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		// 2. Update Stock
		// Check if stock record exists
		var stock domain.Stock
		err := tx.Where("location_id = ? AND variant_id = ?", locationID, variantID).First(&stock).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create new stock record if adding stock
				if qtyChange < 0 {
					return errors.New("insufficient stock")
				}
				stock = domain.Stock{
					LocationID: locationID,
					VariantID:  variantID,
					Quantity:   qtyChange,
					UpdatedAt:  time.Now(),
				}
				if err := tx.Create(&stock).Error; err != nil {
					return err
				}
				return nil
			}
			return err
		}

		// Update existing
		newQty := stock.Quantity + qtyChange
		if newQty < 0 {
			return errors.New("insufficient stock")
		}

		stock.Quantity = newQty
		stock.UpdatedAt = time.Now()
		if err := tx.Save(&stock).Error; err != nil {
			return err
		}

		return nil
	})
}
