package service

import (
	"context"
	"errors"
	"server/internal/core/domain"
	"server/internal/repository"
)

type FulfillmentServiceImpl struct {
	orderRepo repository.Repository[domain.SalesOrder]
}

func NewFulfillmentService(orderRepo repository.Repository[domain.SalesOrder]) FulfillmentService {
	return &FulfillmentServiceImpl{
		orderRepo: orderRepo,
	}
}

func (s *FulfillmentServiceImpl) GetQueue(ctx context.Context) ([]domain.SalesOrder, error) {
	// Get orders with ShipmentStatus == ReadyToPack
	// Since generic repo doesn't support filtering, this is simplified
	// In real implementation, use a custom method in OrderRepository
	return s.orderRepo.FindAll(ctx)
}

func (s *FulfillmentServiceImpl) PackOrder(ctx context.Context, orderID int) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.ShipmentStatus != domain.ShipmentReadyToPack {
		return errors.New("order not ready to pack")
	}

	order.ShipmentStatus = domain.ShipmentShipped
	return s.orderRepo.Update(ctx, order)
}

func (s *FulfillmentServiceImpl) ShipOrder(ctx context.Context, orderID int, carrier, trackingNumber string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	order.Carrier = &carrier
	order.TrackingNumber = &trackingNumber
	order.ShipmentStatus = domain.ShipmentShipped
	return s.orderRepo.Update(ctx, order)
}
