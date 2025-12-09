package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type posSessionRepository struct {
	*GormRepository[domain.POSSession]
}

func NewPOSSessionRepository(db *gorm.DB) POSSessionRepository {
	return &posSessionRepository{NewGormRepository[domain.POSSession](db)}
}

func (r *posSessionRepository) FindActiveSession(ctx context.Context, userID int) (*domain.POSSession, error) {
	// Active means NOT CLOSED
	return r.FindOne(ctx, "user_id = ? AND status != ?", userID, domain.SessionClosed)
}
