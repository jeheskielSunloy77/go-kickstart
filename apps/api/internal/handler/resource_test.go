package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"github.com/jeheskielSunloy77/go-kickstart/internal/repository"

	"github.com/stretchr/testify/require"
)

type stubUserService struct {
	storeFn   func(ctx context.Context, dto *model.StoreUserDTO) (*model.User, error)
	getByIDFn func(ctx context.Context, id uuid.UUID, preloads []string) (*model.User, error)
	getManyFn func(ctx context.Context, opts repository.GetManyOptions) ([]model.User, int64, error)
	destroyFn func(ctx context.Context, id uuid.UUID) error
	killFn    func(ctx context.Context, id uuid.UUID) error
	updateFn  func(ctx context.Context, id uuid.UUID, dto *model.UpdateUserDTO) (*model.User, error)
	restoreFn func(ctx context.Context, id uuid.UUID, preloads []string) (*model.User, error)
}

func (s *stubUserService) Store(ctx context.Context, dto *model.StoreUserDTO) (*model.User, error) {
	if s.storeFn != nil {
		return s.storeFn(ctx, dto)
	}
	return nil, nil
}

func (s *stubUserService) GetByID(ctx context.Context, id uuid.UUID, preloads []string) (*model.User, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id, preloads)
	}
	return nil, nil
}

func (s *stubUserService) GetMany(ctx context.Context, opts repository.GetManyOptions) ([]model.User, int64, error) {
	if s.getManyFn != nil {
		return s.getManyFn(ctx, opts)
	}
	return nil, 0, nil
}

func (s *stubUserService) Destroy(ctx context.Context, id uuid.UUID) error {
	if s.destroyFn != nil {
		return s.destroyFn(ctx, id)
	}
	return nil
}

func (s *stubUserService) Kill(ctx context.Context, id uuid.UUID) error {
	if s.killFn != nil {
		return s.killFn(ctx, id)
	}
	return nil
}

func (s *stubUserService) Update(ctx context.Context, id uuid.UUID, dto *model.UpdateUserDTO) (*model.User, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, dto)
	}
	return nil, nil
}

func (s *stubUserService) Restore(ctx context.Context, id uuid.UUID, preloads []string) (*model.User, error) {
	if s.restoreFn != nil {
		return s.restoreFn(ctx, id, preloads)
	}
	return nil, nil
}

// Ensures Store returns validation errors without calling the service.
func TestResourceHandlerStore_ValidationError(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	called := false
	service := &stubUserService{
		storeFn: func(ctx context.Context, dto *model.StoreUserDTO) (*model.User, error) {
			called = true
			return nil, nil
		},
	}

	h := NewResourceHandler(NewHandler(srv), service)
	app.Post("/users", h.Store())

	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(mustJSON(t, map[string]any{})))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.False(t, called)
}

// Ensures Store parses the payload and returns the created user with 201 status.
func TestResourceHandlerStore_Success(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	userID := uuid.New()
	service := &stubUserService{
		storeFn: func(ctx context.Context, dto *model.StoreUserDTO) (*model.User, error) {
			require.Equal(t, "user@example.com", dto.Email)
			require.Equal(t, "user", dto.Username)
			return &model.User{ID: userID, Email: dto.Email, Username: dto.Username}, nil
		},
	}

	h := NewResourceHandler(NewHandler(srv), service)
	app.Post("/users", h.Store())

	body := mustJSON(t, map[string]any{
		"email":    "user@example.com",
		"username": "user",
		"password": "password123",
	})

	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var got model.User
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.Equal(t, userID, got.ID)
	require.Equal(t, "user@example.com", got.Email)
	require.Equal(t, "user", got.Username)
}

// Ensures GetMany uses request pagination values and returns computed page metadata.
func TestResourceHandlerGetMany_Paginates(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	var captured repository.GetManyOptions
	service := &stubUserService{
		getManyFn: func(ctx context.Context, opts repository.GetManyOptions) ([]model.User, int64, error) {
			captured = opts
			return []model.User{{ID: uuid.New()}, {ID: uuid.New()}}, 5, nil
		},
	}

	h := NewResourceHandler(NewHandler(srv), service)
	app.Get("/users", h.GetMany())

	req, err := http.NewRequest(http.MethodGet, "/users?limit=2&offset=2", nil)
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var got model.PaginatedResponse[model.User]
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.Equal(t, 2, captured.Limit)
	require.Equal(t, 2, captured.Offset)
	require.Equal(t, 2, got.Page)
	require.Equal(t, 2, got.Limit)
	require.Equal(t, 5, got.Total)
	require.Equal(t, 3, got.TotalPages)
}
