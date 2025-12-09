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
	return nil, nil
}

func (s *ProcurementServiceImpl) CreatePO(ctx context.Context, po *domain.PurchaseOrder) error {
	return s.poRepo.Create(ctx, po)
}

func (s *ProcurementServiceImpl) ReceivePO(ctx context.Context, poID int, receivedItems map[int]int) error {
	// TODO: Implement logic
	return nil
}
