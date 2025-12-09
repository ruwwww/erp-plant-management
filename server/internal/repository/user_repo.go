package repository

import (
	"context"
	"server/internal/core/domain"
	"server/internal/core/ports"

	"gorm.io/gorm"
)

type UserRepository struct {
	*GormRepository[domain.User]
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &UserRepository{NewGormRepository[domain.User](db)}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.FindOne(ctx, "email = ?", email)
}
