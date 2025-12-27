package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/errs"
	"github.com/jeheskielSunloy77/go-kickstart/internal/repository"
	"gorm.io/gorm"

	"github.com/stretchr/testify/require"
)

type testEntity struct {
	ID   uuid.UUID
	Name string
}

func (m testEntity) GetID() uuid.UUID {
	return m.ID
}

type testStoreDTO struct {
	Name string
}

func (d testStoreDTO) Validate() error { return nil }

func (d testStoreDTO) ToModel() *testEntity {
	return &testEntity{Name: d.Name}
}

type testUpdateDTO struct {
	Name *string
}

func (d testUpdateDTO) Validate() error { return nil }

func (d testUpdateDTO) ToModel() *testEntity { return &testEntity{} }

func (d testUpdateDTO) ToMap() map[string]any {
	updates := make(map[string]any)
	if d.Name != nil {
		updates["name"] = *d.Name
	}
	return updates
}

type mockResourceRepo struct {
	storeFn   func(ctx context.Context, entity *testEntity) error
	getByIDFn func(ctx context.Context, id uuid.UUID, preloads []string) (*testEntity, error)
	getManyFn func(ctx context.Context, opts repository.GetManyOptions) ([]testEntity, int64, error)
	updateFn  func(ctx context.Context, entity testEntity, updates ...map[string]any) (*testEntity, error)
	destroyFn func(ctx context.Context, id uuid.UUID) error
	killFn    func(ctx context.Context, id uuid.UUID) error
	restoreFn func(ctx context.Context, id uuid.UUID) (*testEntity, error)
}

func (m *mockResourceRepo) CacheEnabled() bool {
	return false
}

func (m *mockResourceRepo) Store(ctx context.Context, entity *testEntity) error {
	if m.storeFn != nil {
		return m.storeFn(ctx, entity)
	}
	return nil
}

func (m *mockResourceRepo) GetByID(ctx context.Context, id uuid.UUID, preloads []string) (*testEntity, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id, preloads)
	}
	return nil, nil
}

func (m *mockResourceRepo) GetMany(ctx context.Context, opts repository.GetManyOptions) ([]testEntity, int64, error) {
	if m.getManyFn != nil {
		return m.getManyFn(ctx, opts)
	}
	return nil, 0, nil
}

func (m *mockResourceRepo) Update(ctx context.Context, entity testEntity, updates ...map[string]any) (*testEntity, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, entity, updates...)
	}
	return &entity, nil
}

func (m *mockResourceRepo) Destroy(ctx context.Context, id uuid.UUID) error {
	if m.destroyFn != nil {
		return m.destroyFn(ctx, id)
	}
	return nil
}

func (m *mockResourceRepo) Kill(ctx context.Context, id uuid.UUID) error {
	if m.killFn != nil {
		return m.killFn(ctx, id)
	}
	return nil
}

func (m *mockResourceRepo) Restore(ctx context.Context, id uuid.UUID) (*testEntity, error) {
	if m.restoreFn != nil {
		return m.restoreFn(ctx, id)
	}
	return nil, nil
}

// Ensures GetByID maps not-found repository errors to HTTP 404 responses.
func TestResourceServiceGetByID_NotFound(t *testing.T) {
	ctx := context.Background()

	repo := &mockResourceRepo{
		getByIDFn: func(_ context.Context, id uuid.UUID, _ []string) (*testEntity, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	svc := NewResourceService[testEntity, testStoreDTO, testUpdateDTO]("widget", repo)

	_, err := svc.GetByID(ctx, uuid.New(), nil)
	require.Error(t, err)

	var httpErr *errs.ErrorResponse
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, http.StatusNotFound, httpErr.Status)
}

// Ensures Update returns the existing entity when there are no updates.
func TestResourceServiceUpdate_NoUpdates(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	updateCalled := false

	repo := &mockResourceRepo{
		getByIDFn: func(_ context.Context, gotID uuid.UUID, _ []string) (*testEntity, error) {
			require.Equal(t, id, gotID)
			return &testEntity{ID: id, Name: "current"}, nil
		},
		updateFn: func(_ context.Context, entity testEntity, updates ...map[string]any) (*testEntity, error) {
			updateCalled = true
			return &entity, nil
		},
	}

	svc := NewResourceService[testEntity, testStoreDTO, testUpdateDTO]("widget", repo)

	updated, err := svc.Update(ctx, id, testUpdateDTO{})
	require.NoError(t, err)
	require.False(t, updateCalled)
	require.NotNil(t, updated)
	require.Equal(t, "current", updated.Name)
}
