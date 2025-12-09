package repository

import (
	"context"
	"server/internal/core/domain"
	"server/internal/dto"

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

func (r *orderRepository) GetFullOrderByID(ctx context.Context, id int) (*domain.SalesOrder, error) {
	var order domain.SalesOrder
	err := r.DB.WithContext(ctx).
		Preload("Items").
		Preload("Customer").
		Preload("Customer.User").
		First(&order, id).Error
	return &order, err
}

func (r *orderRepository) Search(ctx context.Context, filter dto.OrderFilterParams) ([]domain.SalesOrder, int64, error) {
	var orders []domain.SalesOrder
	var total int64

	query := r.DB.WithContext(ctx).Model(&domain.SalesOrder{}).
		Preload("Items").
		Preload("Customer")

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		// Search by Order Number, Customer Company Name, User Name, or Guest Email
		query = query.
			Joins("LEFT JOIN customers ON customers.id = sales_orders.customer_id").
			Joins("LEFT JOIN users ON users.id = customers.user_id").
			Where("sales_orders.order_number ILIKE ? OR customers.company_name ILIKE ? OR users.first_name ILIKE ? OR users.last_name ILIKE ? OR sales_orders.guest_email ILIKE ?",
				search, search, search, search, search)
	}

	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", filter.DateTo)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 10
	} // Default limit

	offset := (page - 1) * limit

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) GetByPOSSession(ctx context.Context, sessionID int) ([]domain.SalesOrder, error) {
	var orders []domain.SalesOrder
	err := r.DB.WithContext(ctx).
		Preload("Items").
		Where("pos_session_id = ?", sessionID).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) SoftDelete(ctx context.Context, id int) error {
	var order domain.SalesOrder
	return r.DB.WithContext(ctx).Delete(&order, id).Error
}

func (r *orderRepository) Restore(ctx context.Context, id int) error {
	var order domain.SalesOrder
	if err := r.DB.WithContext(ctx).Unscoped().First(&order, id).Error; err != nil {
		return err
	}
	return r.DB.WithContext(ctx).Unscoped().Model(&order).Update("DeletedAt", nil).Error
}

func (r *orderRepository) ForceDelete(ctx context.Context, id int) error {
	var order domain.SalesOrder
	return r.DB.WithContext(ctx).Unscoped().Delete(&order, id).Error
}
