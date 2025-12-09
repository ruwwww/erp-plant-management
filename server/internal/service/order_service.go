package service

import (
	"context"
	"server/internal/core/domain"
	"server/internal/repository"
)

type OrderService struct {
	orderRepo repository.Repository[domain.SalesOrder]
}

func NewOrderService(orderRepo repository.Repository[domain.SalesOrder]) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *domain.SalesOrder) error {
	return s.orderRepo.Create(ctx, order)
}

func (s *OrderService) GetOrders(ctx context.Context) ([]domain.SalesOrder, error) {
	return s.orderRepo.FindAll(ctx)
}
