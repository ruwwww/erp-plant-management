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
