package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type cashMoveRepository struct {
	*GormRepository[domain.POSCashMove]
}

func NewCashMoveRepository(db *gorm.DB) CashMoveRepository {
	return &cashMoveRepository{
		GormRepository: NewGormRepository[domain.POSCashMove](db),
	}
}

func (r *cashMoveRepository) GetBySession(ctx context.Context, sessionID int) ([]domain.POSCashMove, error) {
	var moves []domain.POSCashMove
	err := r.DB.WithContext(ctx).
		Where("pos_session_id = ?", sessionID).
		Find(&moves).Error
	return moves, err
}

type movementRepository struct {
	*GormRepository[domain.StockMovement]
}

func NewMovementRepository(db *gorm.DB) MovementRepository {
	return &movementRepository{NewGormRepository[domain.StockMovement](db)}
}
func (r *movementRepository) GetHistory(ctx context.Context, variantID, locationID int, limit int) ([]domain.StockMovement, error) {
	var moves []domain.StockMovement
	err := r.DB.WithContext(ctx).
		Where("variant_id = ? AND location_id = ?", variantID, locationID).
		Order("created_at desc").
		Limit(limit).
		Find(&moves).Error
	return moves, err
}

type invoiceRepository struct {
	*GormRepository[domain.Invoice]
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{NewGormRepository[domain.Invoice](db)}
}
func (r *invoiceRepository) FindByOrder(ctx context.Context, orderID int) (*domain.Invoice, error) {
	return r.FindOne(ctx, "sales_order_id = ?", orderID)
}

type categoryRepository struct {
	*GormRepository[domain.Category]
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{NewGormRepository[domain.Category](db)}
}
func (r *categoryRepository) GetTree(ctx context.Context) ([]domain.Category, error) {
	// In a real app, you might build the tree structure in Go code
	// after fetching all categories, or use a recursive CTE in SQL.
	// For now, fetching all is fine.
	return r.FindAll(ctx)
}
