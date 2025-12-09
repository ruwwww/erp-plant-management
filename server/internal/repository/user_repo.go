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

func (r *userRepository) SoftDelete(ctx context.Context, id int) error {
	var user domain.User
	return r.DB.WithContext(ctx).Delete(&user, id).Error
}

func (r *userRepository) Restore(ctx context.Context, id int) error {
	var user domain.User
	if err := r.DB.WithContext(ctx).Unscoped().First(&user, id).Error; err != nil {
		return err
	}
	return r.DB.WithContext(ctx).Unscoped().Model(&user).Update("DeletedAt", nil).Error
}

func (r *userRepository) ForceDelete(ctx context.Context, id int) error {
	var user domain.User
	return r.DB.WithContext(ctx).Unscoped().Delete(&user, id).Error
}
