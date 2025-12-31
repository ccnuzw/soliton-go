package orm

import (
	"context"
	"errors"
	"fmt"

	"github.com/soliton-go/framework/ddd"
	"gorm.io/gorm"
)

// Repository is a generic interface for repositories.
type Repository[T ddd.Entity, ID ddd.ID] interface {
	Find(ctx context.Context, id ID) (T, error)
	FindAll(ctx context.Context) ([]T, error)
	Save(ctx context.Context, entity T) error
	Delete(ctx context.Context, id ID) error
}

// GormRepository is a generic implementation of Repository using GORM.
type GormRepository[T ddd.Entity, ID ddd.ID] struct {
	db *gorm.DB
}

// NewGormRepository creates a new GormRepository.
func NewGormRepository[T ddd.Entity, ID ddd.ID](db *gorm.DB) *GormRepository[T, ID] {
	return &GormRepository[T, ID]{db: db}
}

func (r *GormRepository[T, ID]) Find(ctx context.Context, id ID) (T, error) {
	var entity T
	result := r.db.WithContext(ctx).First(&entity, "id = ?", id.String())
	if result.Error != nil {
		var zero T
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return zero, fmt.Errorf("entity not found: %w", result.Error)
		}
		return zero, result.Error
	}
	return entity, nil
}

func (r *GormRepository[T, ID]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	result := r.db.WithContext(ctx).Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	return entities, nil
}

func (r *GormRepository[T, ID]) Save(ctx context.Context, entity T) error {
	// If it's an AggregateRoot, we should dispatch events here.
	// But first, let's just save.
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *GormRepository[T, ID]) Delete(ctx context.Context, id ID) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, "id = ?", id.String()).Error
}
