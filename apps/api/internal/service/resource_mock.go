package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"github.com/jeheskielSunloy77/go-kickstart/internal/repository"
)

// MockResourceService is a generic mock for ResourceService interfaces in tests.
type MockResourceService[T model.BaseModel, S model.StoreDTO[T], U model.UpdateDTO[T]] struct {
	StoreFn   func(ctx context.Context, dto S) (*T, error)
	GetByIDFn func(ctx context.Context, id uuid.UUID, preloads []string) (*T, error)
	GetManyFn func(ctx context.Context, opts repository.GetManyOptions) ([]T, int64, error)
	DestroyFn func(ctx context.Context, id uuid.UUID) error
	KillFn    func(ctx context.Context, id uuid.UUID) error
	UpdateFn  func(ctx context.Context, id uuid.UUID, dto U) (*T, error)
	RestoreFn func(ctx context.Context, id uuid.UUID, preloads []string) (*T, error)
}

func NewMockResourceService[T model.BaseModel, S model.StoreDTO[T], U model.UpdateDTO[T]]() *MockResourceService[T, S, U] {
	return &MockResourceService[T, S, U]{}
}

func (m *MockResourceService[T, S, U]) Store(ctx context.Context, dto S) (*T, error) {
	if m.StoreFn != nil {
		return m.StoreFn(ctx, dto)
	}
	return nil, nil
}

func (m *MockResourceService[T, S, U]) GetByID(ctx context.Context, id uuid.UUID, preloads []string) (*T, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id, preloads)
	}
	return nil, nil
}

func (m *MockResourceService[T, S, U]) GetMany(ctx context.Context, opts repository.GetManyOptions) ([]T, int64, error) {
	if m.GetManyFn != nil {
		return m.GetManyFn(ctx, opts)
	}
	return nil, 0, nil
}

func (m *MockResourceService[T, S, U]) Destroy(ctx context.Context, id uuid.UUID) error {
	if m.DestroyFn != nil {
		return m.DestroyFn(ctx, id)
	}
	return nil
}

func (m *MockResourceService[T, S, U]) Kill(ctx context.Context, id uuid.UUID) error {
	if m.KillFn != nil {
		return m.KillFn(ctx, id)
	}
	return nil
}

func (m *MockResourceService[T, S, U]) Update(ctx context.Context, id uuid.UUID, dto U) (*T, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, dto)
	}
	return nil, nil
}

func (m *MockResourceService[T, S, U]) Restore(ctx context.Context, id uuid.UUID, preloads []string) (*T, error) {
	if m.RestoreFn != nil {
		return m.RestoreFn(ctx, id, preloads)
	}
	return nil, nil
}
