package service

import (
	"context"
	"errors"
	"server/internal/core/domain"
	"server/internal/repository"
	"time"

	"gorm.io/gorm"
)

type InventoryServiceImpl struct {
	stockRepo    repository.Repository[domain.Stock]
	movementRepo repository.Repository[domain.StockMovement]
	db           *gorm.DB // Needed for transaction
}

func NewInventoryService(stockRepo repository.Repository[domain.Stock], movementRepo repository.Repository[domain.StockMovement], db *gorm.DB) InventoryService {
	return &InventoryServiceImpl{
		stockRepo:    stockRepo,
		movementRepo: movementRepo,
		db:           db,
	}
}

func (s *InventoryServiceImpl) TransferStock(ctx context.Context, variantID, qty, fromLocID, toLocID, userID int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Deduct from Source
		deductCmd := StockMoveCmd{
			LocationID:    fromLocID,
			VariantID:     variantID,
			QtyChange:     -qty,
			Reason:        domain.ReasonTransfer,
			ReferenceType: "TRANSFER",
			UserID:        userID,
		}
		if err := s.executeMovementTx(tx, deductCmd); err != nil {
			return err
		}

		// 2. Add to Destination
		addCmd := StockMoveCmd{
			LocationID:    toLocID,
			VariantID:     variantID,
			QtyChange:     qty,
			Reason:        domain.ReasonTransfer,
			ReferenceType: "TRANSFER",
			UserID:        userID,
		}
		if err := s.executeMovementTx(tx, addCmd); err != nil {
			return err
		}

		return nil
	})
}

func (s *InventoryServiceImpl) GetStockLevel(ctx context.Context, variantID, locationID int) (int, error) {
	stock, err := s.stockRepo.FindOne(ctx, "location_id = ? AND variant_id = ?", locationID, variantID)
	if err != nil {
		return 0, err
	}
	return stock.Quantity, nil
}

func (s *InventoryServiceImpl) BulkAdjustStock(ctx context.Context, cmds []StockMoveCmd) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, cmd := range cmds {
			if err := s.executeMovementTx(tx, cmd); err != nil {
				return err
			}
		}
		return nil
	})
}

// Helper to execute movement within a transaction
func (s *InventoryServiceImpl) executeMovementTx(tx *gorm.DB, cmd StockMoveCmd) error {
	// 1. Create Movement
	movement := &domain.StockMovement{
		LocationID:     cmd.LocationID,
		VariantID:      cmd.VariantID,
		QuantityChange: cmd.QtyChange,
		Reason:         string(cmd.Reason),
		ReferenceType:  &cmd.ReferenceType,
		ReferenceID:    &cmd.ReferenceID,
		CreatedBy:      &cmd.UserID,
		CreatedAt:      time.Now(),
	}
	if err := tx.Create(movement).Error; err != nil {
		return err
	}

	// 2. Update Stock
	var stock domain.Stock
	err := tx.Where("location_id = ? AND variant_id = ?", cmd.LocationID, cmd.VariantID).First(&stock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if cmd.QtyChange < 0 {
				return errors.New("insufficient stock")
			}
			stock = domain.Stock{
				LocationID: cmd.LocationID,
				VariantID:  cmd.VariantID,
				Quantity:   cmd.QtyChange,
				UpdatedAt:  time.Now(),
			}
			if err := tx.Create(&stock).Error; err != nil {
				return err
			}
			return nil
		}
		return err
	}

	newQty := stock.Quantity + cmd.QtyChange
	if newQty < 0 {
		return errors.New("insufficient stock")
	}

	stock.Quantity = newQty
	stock.UpdatedAt = time.Now()
	if err := tx.Save(&stock).Error; err != nil {
		return err
	}

	return nil
}

func (s *InventoryServiceImpl) ExecuteMovement(ctx context.Context, cmd StockMoveCmd) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.executeMovementTx(tx, cmd)
	})
}

func (s *InventoryServiceImpl) GetMovements(ctx context.Context, variantID, locationID, page, limit int) ([]domain.StockMovement, error) {
	return nil, nil
}

func (s *InventoryServiceImpl) GetLocations(ctx context.Context) ([]domain.InventoryLocation, error) {
	return nil, nil
}

func (s *InventoryServiceImpl) CreateLocation(ctx context.Context, loc *domain.InventoryLocation) error {
	return nil
}

func (s *InventoryServiceImpl) ExportStockSnapshot(ctx context.Context) ([]byte, error) {
	return nil, nil
}
