package service

import (
	"context"

	"github.com/soliton-go/framework/ddd"
	"github.com/soliton-go/framework/orm"
)

// Service defines standard CRUD operations.
type Service[T ddd.Entity, ID ddd.ID] interface {
	Get(ctx context.Context, id ID) (T, error)
	List(ctx context.Context) ([]T, error)
	Create(ctx context.Context, entity T) error
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id ID) error
}

// BaseService is a generic implementation of Service.
type BaseService[T ddd.Entity, ID ddd.ID] struct {
	Repo orm.Repository[T, ID]
}

// NewBaseService creates a new BaseService.
func NewBaseService[T ddd.Entity, ID ddd.ID](repo orm.Repository[T, ID]) *BaseService[T, ID] {
	return &BaseService[T, ID]{Repo: repo}
}

func (s *BaseService[T, ID]) Get(ctx context.Context, id ID) (T, error) {
	return s.Repo.Find(ctx, id)
}

func (s *BaseService[T, ID]) List(ctx context.Context) ([]T, error) {
	return s.Repo.FindAll(ctx)
}

func (s *BaseService[T, ID]) Create(ctx context.Context, entity T) error {
	return s.Repo.Save(ctx, entity)
}

func (s *BaseService[T, ID]) Update(ctx context.Context, entity T) error {
	return s.Repo.Save(ctx, entity)
}

func (s *BaseService[T, ID]) Delete(ctx context.Context, id ID) error {
	return s.Repo.Delete(ctx, id)
}
