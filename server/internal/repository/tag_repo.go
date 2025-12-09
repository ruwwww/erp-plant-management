package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type tagRepository struct {
	*GormRepository[domain.Tag]
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{
		GormRepository: NewGormRepository[domain.Tag](db),
	}
}

func (r *tagRepository) FindBySlug(ctx context.Context, slug string) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.DB.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error
	return &tag, err
}
