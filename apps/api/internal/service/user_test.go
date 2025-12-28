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

func newUserServiceWithRepo(repo repository.UserRepository) UserService {
	return &userService{
		ResourceService: NewResourceService[model.User, *model.StoreUserDTO, *model.UpdateUserDTO]("user", repo),
		repo:            repo,
	}
}

// trackingRepo wraps a UserRepository and records whether GetByID was called.
type trackingRepo struct {
	repository.UserRepository
	called *bool
}

func (t *trackingRepo) GetByID(ctx context.Context, id uuid.UUID, preloads []string) (*model.User, error) {
	if t.called != nil {
		*t.called = true
	}
	return t.UserRepository.GetByID(ctx, id, preloads)
}

func ptrString(v string) *string {
	return &v
}

// Ensures Store hashes passwords before persisting users.
func TestUserServiceStore_HashesPassword(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewMockResourceRepository[model.User](false)

	svc := newUserServiceWithRepo(repo)

	user, err := svc.Store(ctx, &model.StoreUserDTO{
		Email:    "user@example.com",
		Username: "user",
		Password: "password123",
	})
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotEmpty(t, user.PasswordHash)
	require.NotEqual(t, "password123", user.PasswordHash)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123")))
}

// Ensures Update normalizes inputs and hashes passwords before updating.
func TestUserServiceUpdate_NormalizesAndHashes(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	existing := &model.User{ID: id, Email: "old@example.com", Username: "old"}

	repo := repository.NewMockResourceRepository[model.User](false)
	// pre-populate existing user
	require.NoError(t, repo.Store(ctx, existing))

	svc := newUserServiceWithRepo(repo)

	updated, err := svc.Update(ctx, id, &model.UpdateUserDTO{
		Email:    ptrString("  TEST@EXAMPLE.COM "),
		Username: ptrString("  Alice  "),
		Password: ptrString("password123"),
	})
	require.NoError(t, err)
	require.NotNil(t, updated)

	require.Equal(t, "test@example.com", updated.Email)
	require.Equal(t, "Alice", updated.Username)
	require.NotEmpty(t, updated.PasswordHash)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(updated.PasswordHash), []byte("password123")))
}

// Ensures Update returns the existing entity when no meaningful updates are provided.
func TestUserServiceUpdate_EmptyUpdates(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	existing := &model.User{ID: id, Email: "old@example.com", Username: "old"}

	repo := repository.NewMockResourceRepository[model.User](false)
	require.NoError(t, repo.Store(ctx, existing))

	svc := newUserServiceWithRepo(repo)

	updated, err := svc.Update(ctx, id, &model.UpdateUserDTO{
		Email: ptrString("   "),
	})
	require.NoError(t, err)
	require.Equal(t, existing.Email, updated.Email)
	require.Equal(t, existing.Username, updated.Username)
}

// Ensures Update rejects short passwords before performing repository lookups.
func TestUserServiceUpdate_PasswordTooShort(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	getCalled := false

	// wrap underlying mock to track whether GetByID is called
	base := repository.NewMockResourceRepository[model.User](false)
	trepo := &trackingRepo{UserRepository: base, called: &getCalled}

	svc := newUserServiceWithRepo(trepo)

	_, err := svc.Update(ctx, id, &model.UpdateUserDTO{
		Password: ptrString("short"),
	})
	require.Error(t, err)
	require.False(t, getCalled)

	var httpErr *errs.ErrorResponse
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, http.StatusBadRequest, httpErr.Status)
}
