package service

import (
	"context"
	"server/internal/core/domain"
	"server/internal/repository"
)

type OrderServiceImpl struct {
	orderRepo repository.Repository[domain.SalesOrder]
}

func NewOrderService(orderRepo repository.Repository[domain.SalesOrder]) OrderService {
	return &OrderServiceImpl{
		orderRepo: orderRepo,
	}
}

func (s *OrderServiceImpl) PlaceOrder(ctx context.Context, order *domain.SalesOrder) error {
	return s.orderRepo.Create(ctx, order)
}

func (s *OrderServiceImpl) CancelOrder(ctx context.Context, orderID int, reason string) error {
	// TODO: Implement cancel logic
	return nil
}

func (s *OrderServiceImpl) ProcessReturn(ctx context.Context, orderID int, items []domain.Return) error {
	// TODO: Implement return logic
	return nil
}
