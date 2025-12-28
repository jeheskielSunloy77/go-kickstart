package repository

import (
	"context"
	"reflect"
	"sync"

	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"gorm.io/gorm"
)

// MockResourceRepository is a simple in-memory implementation of
// repository.ResourceRepository[T] for use in tests across entities.
type MockResourceRepository[T model.BaseModel] struct {
	mu      sync.RWMutex
	data    map[uuid.UUID]T
	deleted map[uuid.UUID]T
	cacheEn bool
}

func NewMockResourceRepository[T model.BaseModel](cacheEnabled bool) *MockResourceRepository[T] {
	return &MockResourceRepository[T]{
		data:    make(map[uuid.UUID]T),
		deleted: make(map[uuid.UUID]T),
		cacheEn: cacheEnabled,
	}
}

func (m *MockResourceRepository[T]) CacheEnabled() bool {
	return m.cacheEn
}

func (m *MockResourceRepository[T]) Store(ctx context.Context, entity *T) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := (*entity).GetID()
	m.data[id] = *entity
	// if it existed in deleted, remove tombstone
	delete(m.deleted, id)
	return nil
}

func (m *MockResourceRepository[T]) GetByID(ctx context.Context, id uuid.UUID, preloads []string) (*T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if v, ok := m.data[id]; ok {
		return &v, nil
	}
	// not found
	return nil, gorm.ErrRecordNotFound
}

func (m *MockResourceRepository[T]) GetMany(ctx context.Context, opts GetManyOptions) ([]T, int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := make([]T, 0, len(m.data))
	for _, v := range m.data {
		list = append(list, v)
	}
	return list, int64(len(list)), nil
}

func (m *MockResourceRepository[T]) Update(ctx context.Context, entity T, updates ...map[string]any) (*T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := entity.GetID()
	if _, ok := m.data[id]; !ok {
		return nil, gorm.ErrRecordNotFound
	}
	// If updates provided, apply common update keys to the entity (for tests).
	if len(updates) > 0 && updates[0] != nil {
		upd := updates[0]
		// apply known fields
		if email, ok := upd["email"].(string); ok {
			// set via field assignment using a type assertion
			if e, ok := any(&entity).(*T); ok {
				_ = e
			}
			// fallback: attempt to set via reflection for arbitrary T
			// use reflection to set Email, Username, PasswordHash when present
			setFieldIfAvailable(&entity, "Email", email)
		}
		if username, ok := upd["username"].(string); ok {
			setFieldIfAvailable(&entity, "Username", username)
		}
		if ph, ok := upd["password_hash"].(string); ok {
			setFieldIfAvailable(&entity, "PasswordHash", ph)
		}
	}
	m.data[id] = entity
	return &entity, nil
}

func (m *MockResourceRepository[T]) Destroy(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.data[id]; ok {
		// soft-delete: move to deleted map
		m.deleted[id] = v
		delete(m.data, id)
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *MockResourceRepository[T]) Kill(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	delete(m.deleted, id)
	return nil
}

func (m *MockResourceRepository[T]) Restore(ctx context.Context, id uuid.UUID) (*T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.deleted[id]; ok {
		m.data[id] = v
		delete(m.deleted, id)
		return &v, nil
	}
	return nil, gorm.ErrRecordNotFound
}

// setFieldIfAvailable tries to set a field by name on the target value pointed
// to by ptr using reflection. If the field doesn't exist or can't be set,
// the function is a no-op.
func setFieldIfAvailable[T any](ptr *T, fieldName string, value any) {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return
	}
	rv = rv.Elem()
	field := rv.FieldByName(fieldName)
	if !field.IsValid() || !field.CanSet() {
		return
	}
	val := reflect.ValueOf(value)
	if !val.IsValid() {
		return
	}
	if val.Type().AssignableTo(field.Type()) {
		field.Set(val)
		return
	}
	if val.Type().ConvertibleTo(field.Type()) {
		field.Set(val.Convert(field.Type()))
	}
}
