package repository

import (
	"context"

	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id any) (*T, error)
	FindOne(ctx context.Context, condition interface{}, args ...interface{}) (*T, error)
	FindAll(ctx context.Context) ([]T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id any) error
}

type GormRepository[T any] struct {
	DB *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{DB: db}
}

func (r *GormRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *GormRepository[T]) FindByID(ctx context.Context, id any) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).First(&entity, id).Error
	return &entity, err
}

func (r *GormRepository[T]) FindOne(ctx context.Context, condition interface{}, args ...interface{}) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).Where(condition, args...).First(&entity).Error
	return &entity, err
}

func (r *GormRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *GormRepository[T]) Delete(ctx context.Context, id any) error {
	var entity T
	return r.DB.WithContext(ctx).Delete(&entity, id).Error
}
