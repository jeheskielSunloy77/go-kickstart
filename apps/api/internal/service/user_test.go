package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/errs"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"github.com/jeheskielSunloy77/go-kickstart/internal/repository"
	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/require"
)

type mockUserRepo struct {
	storeFn   func(ctx context.Context, entity *model.User) error
	getByIDFn func(ctx context.Context, id uuid.UUID, preloads []string) (*model.User, error)
	getManyFn func(ctx context.Context, opts repository.GetManyOptions) ([]model.User, int64, error)
	updateFn  func(ctx context.Context, entity model.User, updates ...map[string]any) (*model.User, error)
	destroyFn func(ctx context.Context, id uuid.UUID) error
	killFn    func(ctx context.Context, id uuid.UUID) error
	restoreFn func(ctx context.Context, id uuid.UUID) (*model.User, error)
}

func (m *mockUserRepo) CacheEnabled() bool {
	return false
}

func (m *mockUserRepo) Store(ctx context.Context, entity *model.User) error {
	if m.storeFn != nil {
		return m.storeFn(ctx, entity)
	}
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID, preloads []string) (*model.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id, preloads)
	}
	return nil, nil
}

func (m *mockUserRepo) GetMany(ctx context.Context, opts repository.GetManyOptions) ([]model.User, int64, error) {
	if m.getManyFn != nil {
		return m.getManyFn(ctx, opts)
	}
	return nil, 0, nil
}

func (m *mockUserRepo) Update(ctx context.Context, entity model.User, updates ...map[string]any) (*model.User, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, entity, updates...)
	}
	return &entity, nil
}

func (m *mockUserRepo) Destroy(ctx context.Context, id uuid.UUID) error {
	if m.destroyFn != nil {
		return m.destroyFn(ctx, id)
	}
	return nil
}

func (m *mockUserRepo) Kill(ctx context.Context, id uuid.UUID) error {
	if m.killFn != nil {
		return m.killFn(ctx, id)
	}
	return nil
}

func (m *mockUserRepo) Restore(ctx context.Context, id uuid.UUID) (*model.User, error) {
	if m.restoreFn != nil {
		return m.restoreFn(ctx, id)
	}
	return nil, nil
}

func newUserServiceWithRepo(repo *mockUserRepo) UserService {
	return &userService{
		ResourceService: &resourceService[model.User, *model.StoreUserDTO, *model.UpdateUserDTO]{
			resourceName: "user",
			repo:         repo,
		},
		repo: repo,
	}
}

func ptrString(v string) *string {
	return &v
}

// Ensures Store hashes passwords before persisting users.
func TestUserServiceStore_HashesPassword(t *testing.T) {
	ctx := context.Background()
	var stored *model.User

	repo := &mockUserRepo{
		storeFn: func(_ context.Context, user *model.User) error {
			stored = user
			return nil
		},
	}

	svc := newUserServiceWithRepo(repo)

	user, err := svc.Store(ctx, &model.StoreUserDTO{
		Email:    "user@example.com",
		Username: "user",
		Password: "password123",
	})
	require.NoError(t, err)
	require.NotNil(t, stored)
	require.Equal(t, stored, user)
	require.NotEmpty(t, stored.PasswordHash)
	require.NotEqual(t, "password123", stored.PasswordHash)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(stored.PasswordHash), []byte("password123")))
}

// Ensures Update normalizes inputs and hashes passwords before updating.
func TestUserServiceUpdate_NormalizesAndHashes(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	var capturedUpdates map[string]any

	existing := &model.User{ID: id, Email: "old@example.com", Username: "old"}

	repo := &mockUserRepo{
		getByIDFn: func(_ context.Context, gotID uuid.UUID, _ []string) (*model.User, error) {
			require.Equal(t, id, gotID)
			return existing, nil
		},
		updateFn: func(_ context.Context, entity model.User, updates ...map[string]any) (*model.User, error) {
			capturedUpdates = updates[0]
			if email, ok := capturedUpdates["email"].(string); ok {
				entity.Email = email
			}
			if username, ok := capturedUpdates["username"].(string); ok {
				entity.Username = username
			}
			if hash, ok := capturedUpdates["password_hash"].(string); ok {
				entity.PasswordHash = hash
			}
			return &entity, nil
		},
	}

	svc := newUserServiceWithRepo(repo)

	updated, err := svc.Update(ctx, id, &model.UpdateUserDTO{
		Email:    ptrString("  TEST@EXAMPLE.COM "),
		Username: ptrString("  Alice  "),
		Password: ptrString("password123"),
	})
	require.NoError(t, err)
	require.NotNil(t, updated)

	require.Equal(t, "test@example.com", capturedUpdates["email"])
	require.Equal(t, "Alice", capturedUpdates["username"])
	require.NotEmpty(t, capturedUpdates["password_hash"])
	hash := capturedUpdates["password_hash"].(string)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(hash), []byte("password123")))
}

// Ensures Update returns the existing entity when no meaningful updates are provided.
func TestUserServiceUpdate_EmptyUpdates(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	updateCalled := false

	existing := &model.User{ID: id, Email: "old@example.com", Username: "old"}

	repo := &mockUserRepo{
		getByIDFn: func(_ context.Context, gotID uuid.UUID, _ []string) (*model.User, error) {
			require.Equal(t, id, gotID)
			return existing, nil
		},
		updateFn: func(_ context.Context, entity model.User, updates ...map[string]any) (*model.User, error) {
			updateCalled = true
			return &entity, nil
		},
	}

	svc := newUserServiceWithRepo(repo)

	updated, err := svc.Update(ctx, id, &model.UpdateUserDTO{
		Email: ptrString("   "),
	})
	require.NoError(t, err)
	require.False(t, updateCalled)
	require.Equal(t, existing, updated)
}

// Ensures Update rejects short passwords before performing repository lookups.
func TestUserServiceUpdate_PasswordTooShort(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	getCalled := false

	repo := &mockUserRepo{
		getByIDFn: func(_ context.Context, gotID uuid.UUID, _ []string) (*model.User, error) {
			getCalled = true
			return &model.User{ID: gotID}, nil
		},
	}

	svc := newUserServiceWithRepo(repo)

	_, err := svc.Update(ctx, id, &model.UpdateUserDTO{
		Password: ptrString("short"),
	})
	require.Error(t, err)
	require.False(t, getCalled)

	var httpErr *errs.ErrorResponse
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, http.StatusBadRequest, httpErr.Status)
}
