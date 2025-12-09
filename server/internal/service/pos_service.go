package service

import (
	"context"
	"errors"
	"server/internal/core/domain"
	"server/internal/repository"
	"time"
)

type POSServiceImpl struct {
	sessionRepo  repository.Repository[domain.POSSession]
	cashMoveRepo repository.Repository[domain.POSCashMove]
}

func NewPOSService(sessionRepo repository.Repository[domain.POSSession], cashMoveRepo repository.Repository[domain.POSCashMove]) POSService {
	return &POSServiceImpl{
		sessionRepo:  sessionRepo,
		cashMoveRepo: cashMoveRepo,
	}
}

func (s *POSServiceImpl) OpenSession(ctx context.Context, userID int, openingFloat float64) (*domain.POSSession, error) {
	// Check if user already has active session
	existing, _ := s.GetActiveSession(ctx, userID)
	if existing != nil {
		return nil, errors.New("user already has an active session")
	}

	session := &domain.POSSession{
		UserID:      userID,
		OpeningCash: openingFloat,
		Status:      domain.SessionOpened,
		OpenedAt:    time.Now(),
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *POSServiceImpl) CloseSession(ctx context.Context, sessionID int, closingCashActual float64, note string) error {
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return err
	}

	session.ClosingCashActual = &closingCashActual
	session.Status = domain.SessionClosed
	now := time.Now()
	session.ClosedAt = &now
	// session.Note = &note // If field exists

	return s.sessionRepo.Update(ctx, session)
}

func (s *POSServiceImpl) RecordCashMove(ctx context.Context, sessionID int, amount float64, moveType domain.CashMoveType, reason string) error {
	move := &domain.POSCashMove{
		POSSessionID: sessionID,
		Amount:       amount,
		Type:         moveType,
		Reason:       &reason,
		CreatedAt:    time.Now(),
	}
	return s.cashMoveRepo.Create(ctx, move)
}

func (s *POSServiceImpl) GetActiveSession(ctx context.Context, userID int) (*domain.POSSession, error) {
	// Assuming FindOne can handle complex queries or we filter in memory (bad)
	// Ideally: "user_id = ? AND status = 'OPENED'"
	return s.sessionRepo.FindOne(ctx, "user_id = ? AND status = ?", userID, domain.SessionOpened)
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
