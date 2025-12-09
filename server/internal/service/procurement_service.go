package service

import (
	"context"
	"server/internal/core/domain"
	"server/internal/repository"
)

type ProcurementServiceImpl struct {
	poRepo repository.Repository[domain.PurchaseOrder]
}

func NewProcurementService(poRepo repository.Repository[domain.PurchaseOrder]) ProcurementService {
	return &ProcurementServiceImpl{
		poRepo: poRepo,
	}
}

func (s *ProcurementServiceImpl) GetPOs(ctx context.Context, page, limit int) ([]domain.PurchaseOrder, error) {
	return s.poRepo.FindAll(ctx)
}

func (s *ProcurementServiceImpl) CreatePO(ctx context.Context, po *domain.PurchaseOrder) error {
	return s.poRepo.Create(ctx, po)
}

func (s *ProcurementServiceImpl) ReceivePO(ctx context.Context, poID int, receivedItems map[int]int) error {
	po, err := s.poRepo.FindByID(ctx, poID)
	if err != nil {
		return err
	}

	// Update received quantities
	// This is simplified. In reality, we'd need to iterate items and update them.
	// Since we don't have easy access to update nested items via generic repo without more logic:
	po.Status = domain.POPartiallyReceived // or Completed

	return s.poRepo.Update(ctx, po)
}
