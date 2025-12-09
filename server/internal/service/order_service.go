package service

import (
	"context"
	"server/internal/core/domain"
	"server/internal/repository"
)

type OrderServiceImpl struct {
	orderRepo repository.Repository[domain.SalesOrder]
}

func NewOrderService(orderRepo repository.Repository[domain.SalesOrder]) *OrderServiceImpl {
	return &OrderServiceImpl{
		orderRepo: orderRepo,
	}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, order *domain.SalesOrder) error {
	return s.orderRepo.Create(ctx, order)
}

func (s *OrderServiceImpl) GetOrders(ctx context.Context) ([]domain.SalesOrder, error) {
	return s.orderRepo.FindAll(ctx)
}
