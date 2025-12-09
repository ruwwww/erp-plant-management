package repository

import (
	"context"
	"server/internal/core/domain"
	"time"

	"gorm.io/gorm"
)

type promotionRepository struct {
	*GormRepository[domain.Promotion]
}

func NewPromotionRepository(db *gorm.DB) PromotionRepository {
	return &promotionRepository{
		GormRepository: NewGormRepository[domain.Promotion](db),
	}
}

func (r *promotionRepository) GetActivePromotions(ctx context.Context) ([]domain.Promotion, error) {
	var promotions []domain.Promotion
	now := time.Now()
	err := r.DB.WithContext(ctx).
		Where("is_active = ? AND (starts_at IS NULL OR starts_at <= ?) AND (ends_at IS NULL OR ends_at >= ?)", true, now, now).
		Find(&promotions).Error
	return promotions, err
}
