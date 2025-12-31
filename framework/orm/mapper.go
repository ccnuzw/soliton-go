package orm

import (
	"context"

	"gorm.io/gorm"
)

// SQLMapper provides raw SQL capabilities similar to MyBatis.
// It complements the Repository pattern by allowing complex or dynamic SQL execution.
type SQLMapper[T any] interface {
	// SelectOne executes a raw SQL query and maps the result to a single entity.
	SelectOne(ctx context.Context, sql string, args ...interface{}) (*T, error)
	// SelectList executes a raw SQL query and maps the result to a slice of entities.
	SelectList(ctx context.Context, sql string, args ...interface{}) ([]*T, error)
	// Exec executes a raw SQL statement (INSERT, UPDATE, DELETE).
	Exec(ctx context.Context, sql string, args ...interface{}) error
	// Count executes a raw SQL query to count rows.
	Count(ctx context.Context, sql string, args ...interface{}) (int64, error)
}

// GormMapper is the GORM-based implementation of SQLMapper.
type GormMapper[T any] struct {
	db *gorm.DB
}

// NewGormMapper creates a new GormMapper.
func NewGormMapper[T any](db *gorm.DB) *GormMapper[T] {
	return &GormMapper[T]{db: db}
}

func (m *GormMapper[T]) SelectOne(ctx context.Context, sql string, args ...interface{}) (*T, error) {
	var entity T
	err := m.db.WithContext(ctx).Raw(sql, args...).Scan(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (m *GormMapper[T]) SelectList(ctx context.Context, sql string, args ...interface{}) ([]*T, error) {
	var entities []*T
	err := m.db.WithContext(ctx).Raw(sql, args...).Scan(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (m *GormMapper[T]) Exec(ctx context.Context, sql string, args ...interface{}) error {
	return m.db.WithContext(ctx).Exec(sql, args...).Error
}

func (m *GormMapper[T]) Count(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Raw(sql, args...).Scan(&count).Error
	return count, err
}
