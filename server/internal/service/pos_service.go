package service

import (
	"context"
	"server/internal/core/domain"
	"server/internal/repository"
)

type POSServiceImpl struct {
	sessionRepo repository.Repository[domain.POSSession]
}

func NewPOSService(sessionRepo repository.Repository[domain.POSSession]) POSService {
	return &POSServiceImpl{
		sessionRepo: sessionRepo,
	}
}

func (s *POSServiceImpl) OpenSession(ctx context.Context, userID int, openingFloat float64) (*domain.POSSession, error) {
	// TODO: Implement logic
	return nil, nil
}

func (s *POSServiceImpl) CloseSession(ctx context.Context, sessionID int, closingCashActual float64, note string) error {
	// TODO: Implement logic
	return nil
}

func (s *POSServiceImpl) RecordCashMove(ctx context.Context, sessionID int, amount float64, moveType domain.CashMoveType, reason string) error {
	// TODO: Implement logic
	return nil
}

func (s *POSServiceImpl) GetActiveSession(ctx context.Context, userID int) (*domain.POSSession, error) {
	// TODO: Implement logic
	return nil, nil
}

func (s *POSServiceImpl) GetSessionDetails(ctx context.Context, sessionID int) (*domain.POSSession, error) {
	return nil, nil
}

func (s *POSServiceImpl) GetCashMoves(ctx context.Context, sessionID int) ([]domain.POSCashMove, error) {
	return nil, nil
}

func (s *POSServiceImpl) ScanProduct(ctx context.Context, barcode string) (*domain.ProductVariant, int, error) {
	return nil, 0, nil
}

func (s *POSServiceImpl) SearchCustomer(ctx context.Context, query string) ([]domain.User, error) {
	return nil, nil
}

func (s *POSServiceImpl) OverridePrice(ctx context.Context, variantID int, newPrice float64, managerPIN string) error {
	return nil
}

func (s *POSServiceImpl) VoidOrder(ctx context.Context, orderID int, managerPIN string) error {
	return nil
}

func (s *POSServiceImpl) PrintReceipt(ctx context.Context, orderID int) error {
	return nil
}
