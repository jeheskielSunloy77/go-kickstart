package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/errs"
	"github.com/jeheskielSunloy77/go-kickstart/internal/middleware"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"github.com/jeheskielSunloy77/go-kickstart/internal/server"
	"github.com/jeheskielSunloy77/go-kickstart/internal/service"

	"github.com/stretchr/testify/require"
)

type stubAuthService struct {
	registerFn           func(ctx context.Context, email, username, password, userAgent, ipAddress string) (*service.AuthResult, error)
	loginFn              func(ctx context.Context, identifier, password, userAgent, ipAddress string) (*service.AuthResult, error)
	startGoogleAuthFn    func(ctx context.Context) (*service.GoogleAuthStart, error)
	completeGoogleAuthFn func(ctx context.Context, code, state, stateCookie, userAgent, ipAddress string) (*service.AuthResult, error)
	verifyEmailFn        func(ctx context.Context, email, code string) (*model.User, error)
	refreshFn            func(ctx context.Context, refreshToken, userAgent, ipAddress string) (*service.AuthResult, error)
	logoutFn             func(ctx context.Context, refreshToken string) error
	logoutAllFn          func(ctx context.Context, userID uuid.UUID) error
	currentUserFn        func(ctx context.Context, userID uuid.UUID) (*model.User, error)
	resendVerificationFn func(ctx context.Context, userID uuid.UUID) error
}

func (s *stubAuthService) Register(ctx context.Context, email, username, password, userAgent, ipAddress string) (*service.AuthResult, error) {
	if s.registerFn != nil {
		return s.registerFn(ctx, email, username, password, userAgent, ipAddress)
	}
	return nil, nil
}

func (s *stubAuthService) Login(ctx context.Context, identifier, password, userAgent, ipAddress string) (*service.AuthResult, error) {
	if s.loginFn != nil {
		return s.loginFn(ctx, identifier, password, userAgent, ipAddress)
	}
	return nil, nil
}

func (s *stubAuthService) StartGoogleAuth(ctx context.Context) (*service.GoogleAuthStart, error) {
	if s.startGoogleAuthFn != nil {
		return s.startGoogleAuthFn(ctx)
	}
	return nil, nil
}

func (s *stubAuthService) CompleteGoogleAuth(ctx context.Context, code, state, stateCookie, userAgent, ipAddress string) (*service.AuthResult, error) {
	if s.completeGoogleAuthFn != nil {
		return s.completeGoogleAuthFn(ctx, code, state, stateCookie, userAgent, ipAddress)
	}
	return nil, nil
}

func (s *stubAuthService) VerifyEmail(ctx context.Context, email, code string) (*model.User, error) {
	if s.verifyEmailFn != nil {
		return s.verifyEmailFn(ctx, email, code)
	}
	return nil, nil
}

func (s *stubAuthService) Refresh(ctx context.Context, refreshToken, userAgent, ipAddress string) (*service.AuthResult, error) {
	if s.refreshFn != nil {
		return s.refreshFn(ctx, refreshToken, userAgent, ipAddress)
	}
	return nil, nil
}

func (s *stubAuthService) Logout(ctx context.Context, refreshToken string) error {
	if s.logoutFn != nil {
		return s.logoutFn(ctx, refreshToken)
	}
	return nil
}

func (s *stubAuthService) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	if s.logoutAllFn != nil {
		return s.logoutAllFn(ctx, userID)
	}
	return nil
}

func (s *stubAuthService) CurrentUser(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	if s.currentUserFn != nil {
		return s.currentUserFn(ctx, userID)
	}
	return nil, nil
}

func (s *stubAuthService) ResendVerification(ctx context.Context, userID uuid.UUID) error {
	if s.resendVerificationFn != nil {
		return s.resendVerificationFn(ctx, userID)
	}
	return nil
}

// Ensures Register returns validation errors without invoking the service.
func TestAuthHandlerRegister_ValidationError(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	called := false
	authService := &stubAuthService{
		registerFn: func(ctx context.Context, email, username, password, userAgent, ipAddress string) (*service.AuthResult, error) {
			called = true
			return nil, nil
		},
	}

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Post("/register", h.Register())

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
		registerFn: func(ctx context.Context, email, username, password, userAgent, ipAddress string) (*service.AuthResult, error) {
			return &service.AuthResult{
				User:         &model.User{ID: userID, Email: email, Username: username},
				Token:        service.AuthToken{Token: "token", ExpiresAt: time.Now().Add(time.Hour)},
				RefreshToken: service.AuthToken{Token: "refresh", ExpiresAt: time.Now().Add(24 * time.Hour)},
			}, nil
		},
	}

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Post("/register", h.Register())

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

	var got model.User
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.Equal(t, userID, got.ID)
	require.Equal(t, "user@example.com", got.Email)
}

// Ensures Login normalizes email identifiers and maps auth errors to HTTP responses.
func TestAuthHandlerLogin_NormalizesEmail(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	var gotIdentifier string
	authService := &stubAuthService{
		loginFn: func(ctx context.Context, identifier, password, userAgent, ipAddress string) (*service.AuthResult, error) {
			gotIdentifier = identifier
			return nil, errs.NewUnauthorizedError("Invalid credentials", true)
		},
	}

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Post("/login", h.Login())

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

// Ensures Google login start redirects to the provider and sets a state cookie.
func TestAuthHandlerGoogleLogin_RedirectsToProvider(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	authService := &stubAuthService{
		startGoogleAuthFn: func(ctx context.Context) (*service.GoogleAuthStart, error) {
			return &service.GoogleAuthStart{
				AuthURL:        "https://accounts.google.com/o/oauth2/auth?state=abc",
				StateCookie:    "cookie-value",
				StateExpiresAt: time.Now().Add(10 * time.Minute),
			}, nil
		},
	}

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Get("/google", h.GoogleLogin())

	req, err := http.NewRequest(http.MethodGet, "/google", nil)
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, resp.StatusCode)
	require.Equal(t, "https://accounts.google.com/o/oauth2/auth?state=abc", resp.Header.Get("Location"))

	found := false
	for _, header := range resp.Header["Set-Cookie"] {
		if strings.HasPrefix(header, googleStateCookieName+"=") {
			found = true
			break
		}
	}
	require.True(t, found)
}

// Ensures Google callback sets auth cookies and redirects to the web success URL.
func TestAuthHandlerGoogleCallback_RedirectsToSuccess(t *testing.T) {
	srv := newTestServer()
	srv.Config.Auth.GoogleSuccessRedirectURL = "http://localhost:3000/auth/me"
	app := newTestApp(srv)

	userID := uuid.New()
	authService := &stubAuthService{
		completeGoogleAuthFn: func(ctx context.Context, code, state, stateCookie, userAgent, ipAddress string) (*service.AuthResult, error) {
			require.Equal(t, "code", code)
			require.Equal(t, "state", state)
			require.Equal(t, "cookie-value", stateCookie)
			return &service.AuthResult{
				User:         &model.User{ID: userID, Email: "user@example.com"},
				Token:        service.AuthToken{Token: "token", ExpiresAt: time.Now().Add(time.Hour)},
				RefreshToken: service.AuthToken{Token: "refresh", ExpiresAt: time.Now().Add(24 * time.Hour)},
			}, nil
		},
	}

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Get("/google/callback", h.GoogleCallback())

	req, err := http.NewRequest(http.MethodGet, "/google/callback?code=code&state=state", nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: googleStateCookieName, Value: "cookie-value"})

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, resp.StatusCode)
	require.Equal(t, "http://localhost:3000/auth/me", resp.Header.Get("Location"))

	cookieHeaders := strings.Join(resp.Header["Set-Cookie"], "; ")
	require.Contains(t, cookieHeaders, "access_token=")
	require.Contains(t, cookieHeaders, "refresh_token=")
}

// Ensures VerifyEmail validates required fields.
func TestAuthHandlerVerifyEmail_ValidationError(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	called := false
	authService := &stubAuthService{
		verifyEmailFn: func(ctx context.Context, email, code string) (*model.User, error) {
			called = true
			return nil, nil
		},
	}

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Post("/verify-email", h.VerifyEmail())

	req, err := http.NewRequest(http.MethodPost, "/verify-email", bytes.NewReader(mustJSON(t, map[string]any{})))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.False(t, called)
}

// Ensures VerifyEmail returns a user on success.
func TestAuthHandlerVerifyEmail_Success(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	userID := uuid.New()
	authService := &stubAuthService{
		verifyEmailFn: func(ctx context.Context, email, code string) (*model.User, error) {
			return &model.User{ID: userID, Email: email, Username: "user"}, nil
		},
	}

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Post("/verify-email", h.VerifyEmail())

	body := mustJSON(t, map[string]any{
		"email": "user@example.com",
		"code":  "123456",
	})

	req, err := http.NewRequest(http.MethodPost, "/verify-email", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var got model.User
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.Equal(t, userID, got.ID)
	require.Equal(t, "user@example.com", got.Email)
}

// Ensures Refresh pulls the refresh cookie and sets new auth cookies.
func TestAuthHandlerRefresh_Success(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	userID := uuid.New()
	refreshToken := "refresh-token"
	var gotToken string
	authService := &stubAuthService{
		refreshFn: func(ctx context.Context, token, userAgent, ipAddress string) (*service.AuthResult, error) {
			gotToken = token
			return &service.AuthResult{
				User:         &model.User{ID: userID, Email: "user@example.com"},
				Token:        service.AuthToken{Token: "access", ExpiresAt: time.Now().Add(time.Hour)},
				RefreshToken: service.AuthToken{Token: "refresh-new", ExpiresAt: time.Now().Add(24 * time.Hour)},
			}, nil
		},
	}

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Post("/refresh", h.Refresh())

	req, err := http.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(mustJSON(t, map[string]any{})))
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, refreshToken, gotToken)

	var got model.User
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.Equal(t, userID, got.ID)

	require.NotNil(t, cookieByName(resp.Cookies(), "access_token"))
	require.NotNil(t, cookieByName(resp.Cookies(), "refresh_token"))
}

// Ensures Logout clears cookies and forwards the refresh token.
func TestAuthHandlerLogout_Success(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	refreshToken := "refresh-token"
	var gotToken string
	authService := &stubAuthService{
		logoutFn: func(ctx context.Context, token string) error {
			gotToken = token
			return nil
		},
	}

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Post("/logout", h.Logout())

	req, err := http.NewRequest(http.MethodPost, "/logout", bytes.NewReader(mustJSON(t, map[string]any{})))
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, refreshToken, gotToken)

	var got server.Response[any]
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.Equal(t, "Logged out successfully.", got.Message)

	accessCookie := cookieByName(resp.Cookies(), "access_token")
	refreshCookie := cookieByName(resp.Cookies(), "refresh_token")
	require.NotNil(t, accessCookie)
	require.NotNil(t, refreshCookie)
	require.Empty(t, accessCookie.Value)
	require.Empty(t, refreshCookie.Value)
}

// Ensures Me rejects requests without a user ID in context.
func TestAuthHandlerMe_MissingUserID(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	called := false
	authService := &stubAuthService{
		currentUserFn: func(ctx context.Context, userID uuid.UUID) (*model.User, error) {
			called = true
			return nil, nil
		},
	}

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Get("/me", h.Me())

	req, err := http.NewRequest(http.MethodGet, "/me", nil)
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	require.False(t, called)
}

// Ensures Me returns the current user when authenticated.
func TestAuthHandlerMe_Success(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	userID := uuid.New()
	authService := &stubAuthService{
		currentUserFn: func(ctx context.Context, id uuid.UUID) (*model.User, error) {
			require.Equal(t, userID, id)
			return &model.User{ID: userID, Email: "user@example.com"}, nil
		},
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.UserIDKey, userID.String())
		return c.Next()
	})

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Get("/me", h.Me())

	req, err := http.NewRequest(http.MethodGet, "/me", nil)
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var got model.User
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.Equal(t, userID, got.ID)
}

// Ensures ResendVerification uses the user ID from context.
func TestAuthHandlerResendVerification_Success(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	userID := uuid.New()
	var gotID uuid.UUID
	authService := &stubAuthService{
		resendVerificationFn: func(ctx context.Context, id uuid.UUID) error {
			gotID = id
			return nil
		},
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.UserIDKey, userID.String())
		return c.Next()
	})

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Post("/resend-verification", h.ResendVerification())

	req, err := http.NewRequest(http.MethodPost, "/resend-verification", bytes.NewReader(mustJSON(t, map[string]any{})))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, userID, gotID)

	var got server.Response[any]
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.Equal(t, "Verification email sent if needed.", got.Message)
}

// Ensures LogoutAll revokes sessions and clears cookies.
func TestAuthHandlerLogoutAll_Success(t *testing.T) {
	srv := newTestServer()
	app := newTestApp(srv)

	userID := uuid.New()
	var gotID uuid.UUID
	authService := &stubAuthService{
		logoutAllFn: func(ctx context.Context, id uuid.UUID) error {
			gotID = id
			return nil
		},
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.UserIDKey, userID.String())
		return c.Next()
	})

	h := NewAuthHandler(NewHandler(srv), authService)
	app.Post("/logout-all", h.LogoutAll())

	req, err := http.NewRequest(http.MethodPost, "/logout-all", bytes.NewReader(mustJSON(t, map[string]any{})))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, userID, gotID)

	var got server.Response[any]
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.Equal(t, "Logged out from all sessions.", got.Message)

	require.NotNil(t, cookieByName(resp.Cookies(), "access_token"))
	require.NotNil(t, cookieByName(resp.Cookies(), "refresh_token"))
}

func cookieByName(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
