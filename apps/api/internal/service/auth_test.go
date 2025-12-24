package service

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/config"
	"github.com/jeheskielSunloy77/go-kickstart/internal/errs"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/stretchr/testify/require"
)

type mockAuthRepo struct {
	createUserFn            func(ctx context.Context, user *model.User) error
	getByEmailFn            func(ctx context.Context, email string) (*model.User, error)
	getByUsernameFn         func(ctx context.Context, username string) (*model.User, error)
	getByGoogleIDFn         func(ctx context.Context, googleID string) (*model.User, error)
	saveFn                  func(ctx context.Context, user *model.User) error
	updateLoginAtFn         func(ctx context.Context, id uuid.UUID, ts time.Time) error
	updateEmailVerifiedAtFn func(ctx context.Context, id uuid.UUID, ts time.Time) error
}

type mockVerificationRepo struct {
	createFn       func(ctx context.Context, verification *model.EmailVerification) error
	getActiveFn    func(ctx context.Context, userID uuid.UUID, codeHash string, now time.Time) (*model.EmailVerification, error)
	expireActiveFn func(ctx context.Context, userID uuid.UUID, now time.Time) error
	markVerifiedFn func(ctx context.Context, id uuid.UUID, verifiedAt time.Time) error
}

func (m *mockAuthRepo) Save(ctx context.Context, user *model.User) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, user)
	}
	return nil
}

func (m *mockAuthRepo) CreateUser(ctx context.Context, user *model.User) error {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, user)
	}
	return nil
}

func (m *mockAuthRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	if m.getByUsernameFn != nil {
		return m.getByUsernameFn(ctx, username)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) GetByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	if m.getByGoogleIDFn != nil {
		return m.getByGoogleIDFn(ctx, googleID)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) UpdateLoginAt(ctx context.Context, id uuid.UUID, ts time.Time) error {
	if m.updateLoginAtFn != nil {
		return m.updateLoginAtFn(ctx, id, ts)
	}
	return nil
}

func (m *mockAuthRepo) UpdateEmailVerifiedAt(ctx context.Context, id uuid.UUID, ts time.Time) error {
	if m.updateEmailVerifiedAtFn != nil {
		return m.updateEmailVerifiedAtFn(ctx, id, ts)
	}
	return nil
}

func (m *mockVerificationRepo) Create(ctx context.Context, verification *model.EmailVerification) error {
	if m.createFn != nil {
		return m.createFn(ctx, verification)
	}
	return nil
}

func (m *mockVerificationRepo) GetActiveByUserIDAndCodeHash(ctx context.Context, userID uuid.UUID, codeHash string, now time.Time) (*model.EmailVerification, error) {
	if m.getActiveFn != nil {
		return m.getActiveFn(ctx, userID, codeHash, now)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockVerificationRepo) ExpireActiveByUserID(ctx context.Context, userID uuid.UUID, now time.Time) error {
	if m.expireActiveFn != nil {
		return m.expireActiveFn(ctx, userID, now)
	}
	return nil
}

func (m *mockVerificationRepo) MarkVerified(ctx context.Context, id uuid.UUID, verifiedAt time.Time) error {
	if m.markVerifiedFn != nil {
		return m.markVerifiedFn(ctx, id, verifiedAt)
	}
	return nil
}

// Ensures Register hashes passwords and returns a signed token tied to the user ID.
func TestAuthServiceRegister_HashesPasswordAndReturnsToken(t *testing.T) {
	secret := "test-secret"
	ttl := 15 * time.Minute
	ctx := context.Background()
	var createdUser *model.User

	repo := &mockAuthRepo{
		createUserFn: func(_ context.Context, user *model.User) error {
			if user.ID == uuid.Nil {
				user.ID = uuid.New()
			}
			createdUser = user
			return nil
		},
	}

	svc := NewAuthService(&config.AuthConfig{SecretKey: secret, AccessTokenTTL: ttl}, repo, nil, nil, nil)

	result, err := svc.Register(ctx, "user@example.com", "user", "password123")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, createdUser)

	require.NotEmpty(t, result.User.PasswordHash)
	require.NotEqual(t, "password123", result.User.PasswordHash)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(result.User.PasswordHash), []byte("password123")))

	claims := &jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(result.Token.Token, claims, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid)
	require.Equal(t, result.User.ID.String(), claims.Subject)
}

// Ensures Register rejects short passwords before hitting the repository.
func TestAuthServiceRegister_ShortPassword(t *testing.T) {
	ctx := context.Background()
	called := false

	repo := &mockAuthRepo{
		createUserFn: func(_ context.Context, user *model.User) error {
			called = true
			return nil
		},
	}

	svc := NewAuthService(&config.AuthConfig{SecretKey: "test"}, repo, nil, nil, nil)

	_, err := svc.Register(ctx, "user@example.com", "user", "short")
	require.Error(t, err)
	require.False(t, called)

	var httpErr *errs.HTTPError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, http.StatusBadRequest, httpErr.Status)
}

// Ensures Login returns unauthorized when the user lookup fails.
func TestAuthServiceLogin_UserNotFound(t *testing.T) {
	ctx := context.Background()

	repo := &mockAuthRepo{
		getByEmailFn: func(_ context.Context, email string) (*model.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	svc := NewAuthService(&config.AuthConfig{SecretKey: "test"}, repo, nil, nil, nil)

	_, err := svc.Login(ctx, "user@example.com", "password123")
	require.Error(t, err)

	var httpErr *errs.HTTPError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, http.StatusUnauthorized, httpErr.Status)
}

// Ensures Login rejects invalid passwords without updating login timestamps.
func TestAuthServiceLogin_PasswordMismatch(t *testing.T) {
	ctx := context.Background()
	called := false

	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	repo := &mockAuthRepo{
		getByUsernameFn: func(_ context.Context, username string) (*model.User, error) {
			return &model.User{ID: uuid.New(), Username: username, PasswordHash: string(hash)}, nil
		},
		updateLoginAtFn: func(_ context.Context, id uuid.UUID, ts time.Time) error {
			called = true
			return nil
		},
	}

	svc := NewAuthService(&config.AuthConfig{SecretKey: "test"}, repo, nil, nil, nil)

	_, err = svc.Login(ctx, "user", "wrong-password")
	require.Error(t, err)
	require.False(t, called)

	var httpErr *errs.HTTPError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, http.StatusUnauthorized, httpErr.Status)
}

// Ensures Login updates login timestamps and returns a valid token on success.
func TestAuthServiceLogin_Success(t *testing.T) {
	secret := "test-secret"
	ctx := context.Background()
	called := false

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	userID := uuid.New()
	repo := &mockAuthRepo{
		getByEmailFn: func(_ context.Context, email string) (*model.User, error) {
			return &model.User{ID: userID, Email: email, PasswordHash: string(hash)}, nil
		},
		updateLoginAtFn: func(_ context.Context, id uuid.UUID, ts time.Time) error {
			called = true
			return nil
		},
	}

	svc := NewAuthService(&config.AuthConfig{SecretKey: secret, AccessTokenTTL: time.Minute}, repo, nil, nil, nil)

	result, err := svc.Login(ctx, "user@example.com", "password123")
	require.NoError(t, err)
	require.True(t, called)
	require.NotNil(t, result)

	claims := &jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(result.Token.Token, claims, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid)
	require.Equal(t, userID.String(), claims.Subject)
}

// Ensures LoginWithGoogle fails fast when Google auth is not configured.
func TestAuthServiceLoginWithGoogle_ConfigMissing(t *testing.T) {
	ctx := context.Background()

	svc := NewAuthService(&config.AuthConfig{SecretKey: "test"}, &mockAuthRepo{}, nil, nil, nil)

	_, err := svc.LoginWithGoogle(ctx, "token")
	require.Error(t, err)

	var httpErr *errs.HTTPError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, http.StatusBadRequest, httpErr.Status)
}

// Ensures VerifyEmail marks the user as verified when the code is valid.
func TestAuthServiceVerifyEmail_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	code := "123456"
	codeHash := hashVerificationCode(code)
	verifiedCalled := false

	repo := &mockAuthRepo{
		getByEmailFn: func(_ context.Context, email string) (*model.User, error) {
			return &model.User{ID: userID, Email: email}, nil
		},
		updateEmailVerifiedAtFn: func(_ context.Context, id uuid.UUID, ts time.Time) error {
			verifiedCalled = true
			return nil
		},
	}

	verificationRepo := &mockVerificationRepo{
		getActiveFn: func(_ context.Context, id uuid.UUID, hash string, now time.Time) (*model.EmailVerification, error) {
			require.Equal(t, userID, id)
			require.Equal(t, codeHash, hash)
			return &model.EmailVerification{ID: uuid.New(), UserID: id}, nil
		},
		markVerifiedFn: func(_ context.Context, id uuid.UUID, verifiedAt time.Time) error {
			return nil
		},
	}

	svc := NewAuthService(&config.AuthConfig{SecretKey: "test"}, repo, verificationRepo, nil, nil)

	user, err := svc.VerifyEmail(ctx, "user@example.com", code)
	require.NoError(t, err)
	require.True(t, verifiedCalled)
	require.NotNil(t, user.EmailVerifiedAt)
}

// Ensures VerifyEmail rejects invalid codes.
func TestAuthServiceVerifyEmail_InvalidCode(t *testing.T) {
	ctx := context.Background()

	repo := &mockAuthRepo{
		getByEmailFn: func(_ context.Context, email string) (*model.User, error) {
			return &model.User{ID: uuid.New(), Email: email}, nil
		},
	}

	verificationRepo := &mockVerificationRepo{
		getActiveFn: func(_ context.Context, id uuid.UUID, hash string, now time.Time) (*model.EmailVerification, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	svc := NewAuthService(&config.AuthConfig{SecretKey: "test"}, repo, verificationRepo, nil, nil)

	_, err := svc.VerifyEmail(ctx, "user@example.com", "bad-code")
	require.Error(t, err)

	var httpErr *errs.HTTPError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, http.StatusBadRequest, httpErr.Status)
}
