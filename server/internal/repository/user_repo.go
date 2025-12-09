package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	*GormRepository[domain.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{NewGormRepository[domain.User](db)}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.FindOne(ctx, "email = ?", email)
}
