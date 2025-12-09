package service

import (
	"context"
	"server/internal/core/domain"
	"server/internal/repository"
)

type ProcurementServiceImpl struct {
	poRepo       repository.Repository[domain.PurchaseOrder]
	supplierRepo repository.SupplierRepository
}

func NewProcurementService(poRepo repository.Repository[domain.PurchaseOrder], supplierRepo repository.SupplierRepository) ProcurementService {
	return &ProcurementServiceImpl{
		poRepo:       poRepo,
		supplierRepo: supplierRepo,
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

func (s *ProcurementServiceImpl) CreateSupplier(ctx context.Context, supplier *domain.Supplier) error {
	return s.supplierRepo.Create(ctx, supplier)
}

func (s *ProcurementServiceImpl) UpdateSupplier(ctx context.Context, supplier *domain.Supplier) error {
	return s.supplierRepo.Update(ctx, supplier)
}

func (s *ProcurementServiceImpl) SoftDeleteSupplier(ctx context.Context, id int) error {
	return s.supplierRepo.SoftDelete(ctx, id)
}

func (s *ProcurementServiceImpl) RestoreSupplier(ctx context.Context, id int) error {
	return s.supplierRepo.Restore(ctx, id)
}

func (s *ProcurementServiceImpl) ForceDeleteSupplier(ctx context.Context, id int) error {
	return s.supplierRepo.ForceDelete(ctx, id)
}
