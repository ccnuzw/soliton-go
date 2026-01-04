package orm

import (
	"context"
	"errors"
	"fmt"
	"reflect"

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
	dest := any(&entity)
	entityType := reflect.TypeOf(entity)
	if entityType != nil && entityType.Kind() == reflect.Ptr {
		ptr := reflect.New(entityType.Elem()).Interface()
		typed, ok := ptr.(T)
		if !ok {
			var zero T
			return zero, fmt.Errorf("failed to allocate entity")
		}
		entity = typed
		dest = entity
	}
	result := r.db.WithContext(ctx).First(dest, "id = ?", id.String())
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
	model := any(&entity)
	entityType := reflect.TypeOf(entity)
	if entityType != nil && entityType.Kind() == reflect.Ptr {
		model = reflect.New(entityType.Elem()).Interface()
	}
	return r.db.WithContext(ctx).Delete(model, "id = ?", id.String()).Error
}
