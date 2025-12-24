package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/errs"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"github.com/jeheskielSunloy77/go-kickstart/internal/service"

	"github.com/stretchr/testify/require"
)

type stubAuthService struct {
	registerFn    func(ctx context.Context, email, username, password string) (*service.AuthResult, error)
	loginFn       func(ctx context.Context, identifier, password string) (*service.AuthResult, error)
	googleLoginFn func(ctx context.Context, idToken string) (*service.AuthResult, error)
}

func (s *stubAuthService) Register(ctx context.Context, email, username, password string) (*service.AuthResult, error) {
	if s.registerFn != nil {
		return s.registerFn(ctx, email, username, password)
	}
	return nil, nil
}

func (s *stubAuthService) Login(ctx context.Context, identifier, password string) (*service.AuthResult, error) {
	if s.loginFn != nil {
		return s.loginFn(ctx, identifier, password)
	}
	return nil, nil
}

func (s *stubAuthService) LoginWithGoogle(ctx context.Context, idToken string) (*service.AuthResult, error) {
	if s.googleLoginFn != nil {
		return s.googleLoginFn(ctx, idToken)
	}
	return nil, nil
}

// Ensures Register returns validation errors without invoking the service.
func TestAuthHandlerRegister_ValidationError(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	called := false
	authService := &stubAuthService{
		registerFn: func(ctx context.Context, email, username, password string) (*service.AuthResult, error) {
			called = true
			return nil, nil
		},
	}

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Post("/register", h.Register)

	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(mustJSON(t, map[string]any{})))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.False(t, called)
}

// Ensures Register returns a 201 response with the auth payload on success.
func TestAuthHandlerRegister_Success(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	userID := uuid.New()
	authService := &stubAuthService{
		registerFn: func(ctx context.Context, email, username, password string) (*service.AuthResult, error) {
			return &service.AuthResult{
				User:  &model.User{ID: userID, Email: email, Username: username},
				Token: service.AuthToken{Token: "token", ExpiresAt: time.Now().Add(time.Hour)},
			}, nil
		},
	}

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Post("/register", h.Register)

	body := mustJSON(t, map[string]any{
		"email":    "user@example.com",
		"username": "user",
		"password": "password123",
	})

	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var got service.AuthResult
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.Equal(t, userID, got.User.ID)
	require.Equal(t, "user@example.com", got.User.Email)
}

// Ensures Login normalizes email identifiers and maps auth errors to HTTP responses.
func TestAuthHandlerLogin_NormalizesEmail(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	var gotIdentifier string
	authService := &stubAuthService{
		loginFn: func(ctx context.Context, identifier, password string) (*service.AuthResult, error) {
			gotIdentifier = identifier
			return nil, errs.NewUnauthorizedError("Invalid credentials", true)
		},
	}

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Post("/login", h.Login)

	body := mustJSON(t, map[string]any{
		"identifier": "USER@Example.COM",
		"password":   "password123",
	})

	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	fmt.Printf("comparing user@example.com to %s", gotIdentifier)
	require.Equal(t, "user@example.com", gotIdentifier)
}
