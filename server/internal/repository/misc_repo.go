package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type CashMoveRepository struct {
	*GormRepository[domain.POSCashMove]
}

func NewCashMoveRepository(db *gorm.DB) ports.CashMoveRepository {
	return &CashMoveRepository{NewGormRepository[domain.POSCashMove](db)}
}

type MovementRepository struct {
	*GormRepository[domain.StockMovement]
}

func NewMovementRepository(db *gorm.DB) ports.MovementRepository {
	return &MovementRepository{NewGormRepository[domain.StockMovement](db)}
}
func (r *MovementRepository) GetHistory(ctx context.Context, variantID, locationID int, limit int) ([]domain.StockMovement, error) {
	var moves []domain.StockMovement
	err := r.DB.WithContext(ctx).
		Where("variant_id = ? AND location_id = ?", variantID, locationID).
		Order("created_at desc").
		Limit(limit).
		Find(&moves).Error
	return moves, err
}

type InvoiceRepository struct {
	*GormRepository[domain.Invoice]
}

func NewInvoiceRepository(db *gorm.DB) ports.InvoiceRepository {
	return &InvoiceRepository{NewGormRepository[domain.Invoice](db)}
}
func (r *InvoiceRepository) FindByOrder(ctx context.Context, orderID int) (*domain.Invoice, error) {
	return r.FindOne(ctx, "sales_order_id = ?", orderID)
}

type CategoryRepository struct {
	*GormRepository[domain.Category]
}

func NewCategoryRepository(db *gorm.DB) ports.CategoryRepository {
	return &CategoryRepository{NewGormRepository[domain.Category](db)}
}
func (r *CategoryRepository) GetTree(ctx context.Context) ([]domain.Category, error) {
	// In a real app, you might build the tree structure in Go code
	// after fetching all categories, or use a recursive CTE in SQL.
	// For now, fetching all is fine.
	return r.FindAll(ctx)
}
