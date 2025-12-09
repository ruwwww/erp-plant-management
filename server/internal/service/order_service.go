package service

import (
	"context"
	"errors"
	"server/internal/core/domain"
	"server/internal/repository"

	"gorm.io/gorm"
)

type OrderServiceImpl struct {
	orderRepo        repository.OrderRepository
	InventoryService InventoryService
	db               *gorm.DB // <-- Perbaikan Arsitektur: Inject DB
}

func NewOrderService(orderRepo repository.OrderRepository, inventoryService InventoryService, db *gorm.DB) OrderService {
	return &OrderServiceImpl{
		orderRepo:        orderRepo,
		InventoryService: inventoryService,
		db:               db,
	}
}

func (s *OrderServiceImpl) GetByPOSSession(ctx context.Context, sessionID int) ([]domain.SalesOrder, error) {
	return s.orderRepo.GetByPOSSession(ctx, sessionID)
}

func (s *OrderServiceImpl) PlaceOrder(ctx context.Context, order *domain.SalesOrder) error {
	// Logic expansion: Here is where you should eventually put:
	// 1. Validate Stock (call inventoryService)
	// 2. Calculate Totals (call cartService)
	return s.orderRepo.Create(ctx, order)
}

func (s *OrderServiceImpl) GetOrder(ctx context.Context, orderNumber string) (*domain.SalesOrder, error) {
	return s.orderRepo.GetFullOrder(ctx, orderNumber)
}

func (s *OrderServiceImpl) GetCustomerHistory(ctx context.Context, userID int, page, limit int) ([]domain.SalesOrder, error) {
	// TODO: Add Pagination logic to Repo
	return s.orderRepo.Find(ctx, "customer_id = ?", userID)
}
func (s *OrderServiceImpl) CancelOrder(ctx context.Context, orderID int, reason string) error {
	order, err := s.orderRepo.GetFullOrderByID(ctx, orderID) // Perbaikan: Ganti method call
	if err != nil {
		return err
	}

	if order.Status == domain.OrderShipped || order.Status == domain.OrderCompleted {
		return errors.New("Order cannot be cancelled: Item has already been shipped or completed")
	}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil || tx.Error != nil {
			tx.Rollback()
		}
	}()

	for _, item := range order.Items {
		cmd := StockMoveCmd{
			LocationID:    1, // TODO: Tentukan Lokasi Gudang
			VariantID:     item.VariantID,
			QtyChange:     item.Quantity, // Quantity positif = penambahan stok
			Reason:        domain.ReasonReturn,
			ReferenceID:   orderID,
			ReferenceType: "sales_orders",
			// UserID:        userID, // TODO: Dapatkan dari context/auth
		}

		if err := s.InventoryService.ExecuteMovement(ctx, cmd); err != nil {
			tx.Rollback()
			return err
		}
	}

	order.Status = domain.OrderCancelled

	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *OrderServiceImpl) ProcessReturn(ctx context.Context, orderID int, items []domain.Return) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	order.Status = domain.OrderReturned // Simplification
	return s.orderRepo.Update(ctx, order)
}

func (s *OrderServiceImpl) GetOrderList(ctx context.Context, filter OrderFilterParams) ([]domain.SalesOrder, int64, error) {
	// If your Generic Repo supports dynamic filtering, map it here.
	// Otherwise, you rely on a specific Repo method.

	// For now, let's assume we map it to a generic Map for the Repo
	repoFilter := make(map[string]interface{})

	if filter.Status != "" {
		repoFilter["status"] = filter.Status
	}
	if filter.Search != "" {
		repoFilter["order_number"] = filter.Search // Simplification
	}

	// Warning: Generic FindAll usually doesn't return Total Count or handle Pagination.
	// Ideally: return s.orderRepo.Search(ctx, repoFilter, filter.Page, filter.Limit)

	// Temporary fallback to FindAll to keep code compiling:
	orders, err := s.orderRepo.FindAll(ctx)
	return orders, int64(len(orders)), err
}

func (s *OrderServiceImpl) SubmitReview(ctx context.Context, review *domain.Review) error {
	return nil // Placeholder
}

// -------------------------------------------------------
// 4. INTERFACE MISMATCH
// These are implemented but likely NOT in your Interface.
// Add them to 'interfaces.go' or keep them as internal helpers.
// -------------------------------------------------------

func (s *OrderServiceImpl) SoftDeleteOrder(ctx context.Context, orderID int) error {
	return s.orderRepo.SoftDelete(ctx, orderID)
}

func (s *OrderServiceImpl) RestoreOrder(ctx context.Context, orderID int) error {
	return s.orderRepo.Restore(ctx, orderID)
}

func (s *OrderServiceImpl) ForceDeleteOrder(ctx context.Context, orderID int) error {
	return s.orderRepo.ForceDelete(ctx, orderID)
}
