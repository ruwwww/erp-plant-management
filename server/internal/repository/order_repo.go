package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type orderRepository struct {
	*GormRepository[domain.SalesOrder]
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{NewGormRepository[domain.SalesOrder](db)}
}

func (r *orderRepository) GetFullOrder(ctx context.Context, orderNumber string) (*domain.SalesOrder, error) {
	var order domain.SalesOrder
	err := r.DB.WithContext(ctx).
		Preload("Items").
		Preload("Customer").
		Preload("Customer.User"). // Nested Preload
		Where("order_number = ?", orderNumber).
		First(&order).Error
	return &order, err
}

func (r *orderRepository) GetByPOSSession(ctx context.Context, sessionID int) ([]domain.SalesOrder, error) {
	var orders []domain.SalesOrder
	err := r.DB.WithContext(ctx).
		Preload("Items").
		Where("pos_session_id = ?", sessionID).
		Find(&orders).Error
	return orders, err
}
