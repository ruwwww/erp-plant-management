package repository

import (
	"context"
	"server/internal/core/domain"

	"gorm.io/gorm"
)

type MediaRepository interface {
	Repository[domain.MediaAsset]
	GetByUUID(ctx context.Context, uuid string) (*domain.MediaAsset, error)
}

type MediaLinkRepository interface {
	Repository[domain.MediaLink]
	GetByEntity(ctx context.Context, entityType string, entityID int, zone string) ([]domain.MediaLink, error)
	DeleteByEntity(ctx context.Context, entityType string, entityID int, zone string) error
}

type mediaRepository struct {
	*GormRepository[domain.MediaAsset]
}

func NewMediaRepository(db *gorm.DB) MediaRepository {
	return &mediaRepository{NewGormRepository[domain.MediaAsset](db)}
}

func (r *mediaRepository) GetByUUID(ctx context.Context, uuid string) (*domain.MediaAsset, error) {
	var asset domain.MediaAsset
	err := r.DB.WithContext(ctx).Where("uuid = ?", uuid).First(&asset).Error
	return &asset, err
}

type mediaLinkRepository struct {
	*GormRepository[domain.MediaLink]
}

func NewMediaLinkRepository(db *gorm.DB) MediaLinkRepository {
	return &mediaLinkRepository{NewGormRepository[domain.MediaLink](db)}
}

func (r *mediaLinkRepository) GetByEntity(ctx context.Context, entityType string, entityID int, zone string) ([]domain.MediaLink, error) {
	var links []domain.MediaLink
	err := r.DB.WithContext(ctx).Preload("Media").Where("entity_type = ? AND entity_id = ? AND zone = ?", entityType, entityID, zone).Order("sort_order ASC").Find(&links).Error
	return links, err
}

func (r *mediaLinkRepository) DeleteByEntity(ctx context.Context, entityType string, entityID int, zone string) error {
	return r.DB.WithContext(ctx).Where("entity_type = ? AND entity_id = ? AND zone = ?", entityType, entityID, zone).Delete(&domain.MediaLink{}).Error
}
