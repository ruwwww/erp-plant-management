package service

import (
	"context"
	"errors"
	"server/internal/core/domain"
	"server/internal/repository"
)

type ProcurementServiceImpl struct {
	poRepo       repository.Repository[domain.PurchaseOrder]
	supplierRepo repository.SupplierRepository
	inventorySvc InventoryService
}

func NewProcurementService(poRepo repository.Repository[domain.PurchaseOrder], supplierRepo repository.SupplierRepository, inventorySvc InventoryService) ProcurementService {
	return &ProcurementServiceImpl{
		poRepo:       poRepo,
		supplierRepo: supplierRepo,
		inventorySvc: inventorySvc,
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

	if po.Status == domain.POCompleted || po.Status == domain.POCancelled {
		return errors.New("cannot receive items for completed or cancelled PO")
	}

	allReceived := true
	for i := range po.Items {
		item := &po.Items[i]
		if qty, ok := receivedItems[item.ID]; ok {
			item.QuantityReceived += qty
			if item.QuantityReceived < item.QuantityOrdered {
				allReceived = false
			}
			// Create stock movement for received quantity
			err = s.inventorySvc.ExecuteMovement(ctx, StockMoveCmd{
				LocationID:    1, // TODO: configurable default location
				VariantID:     item.VariantID,
				QtyChange:     qty,
				Reason:        domain.ReasonPurchase,
				ReferenceID:   poID,
				ReferenceType: "PURCHASE_ORDER",
				UserID:        0, // TODO: get from context
			})
			if err != nil {
				return err
			}
		} else if item.QuantityReceived < item.QuantityOrdered {
			allReceived = false
		}
	}

	if allReceived {
		po.Status = domain.POCompleted
	} else {
		po.Status = domain.POPartiallyReceived
	}

	return s.poRepo.Update(ctx, po)
}

func (s *ProcurementServiceImpl) GetSuppliers(ctx context.Context) ([]domain.Supplier, error) {
	return s.supplierRepo.FindAll(ctx)
}

func (s *ProcurementServiceImpl) GetSupplier(ctx context.Context, id int) (*domain.Supplier, error) {
	return s.supplierRepo.FindByID(ctx, id)
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
