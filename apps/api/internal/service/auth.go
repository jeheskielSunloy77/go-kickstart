package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/config"
	"github.com/jeheskielSunloy77/go-kickstart/internal/errs"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"github.com/jeheskielSunloy77/go-kickstart/internal/repository"
	"github.com/jeheskielSunloy77/go-kickstart/internal/sqlerr"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

var (
	minPasswordLength = 8
	emailRegex        = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

type AuthService struct {
	repo           repository.AuthRepositoryInterface
	secretKey      []byte
	accessTokenTTL time.Duration
	googleClientID string
}

type AuthToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type AuthResult struct {
	User  *model.User `json:"user"`
	Token AuthToken   `json:"token"`
}

type AuthServiceInterface interface {
	Register(ctx context.Context, email, username, password string) (*AuthResult, error)
	Login(ctx context.Context, identifier, password string) (*AuthResult, error)
	LoginWithGoogle(ctx context.Context, idToken string) (*AuthResult, error)
}

func NewAuthService(cfg *config.AuthConfig, repo repository.AuthRepositoryInterface) *AuthService {
	return &AuthService{
		repo:           repo,
		secretKey:      []byte(cfg.SecretKey),
		accessTokenTTL: cfg.AccessTokenTTL,
		googleClientID: cfg.GoogleClientID,
	}
}

func (s *AuthService) Register(ctx context.Context, email, username, password string) (*AuthResult, error) {
	if len(password) < minPasswordLength {
		return nil, errs.NewBadRequestError(
			fmt.Sprintf("Password must be at least %d characters", minPasswordLength),
			true,
			nil,
			[]errs.FieldError{{Field: "password", Error: "too short"}},
			nil,
		)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.NewInternalServerError()
	}

	user := &model.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(passwordHash),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, sqlerr.HandleError(err)
	}

	token, exp, err := s.generateToken(user.ID)
	if err != nil {
		return nil, errs.NewInternalServerError()
	}

	return &AuthResult{
		User:  user,
		Token: AuthToken{Token: token, ExpiresAt: exp},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, identifier, password string) (*AuthResult, error) {
	user, err := s.lookupUser(ctx, identifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewUnauthorizedError("Invalid credentials", true)
		}
		return nil, sqlerr.HandleError(err)
	}

	if user.PasswordHash == "" {
		return nil, errs.NewUnauthorizedError("Password login not available for this account", true)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errs.NewUnauthorizedError("Invalid credentials", true)
	}

	now := time.Now().UTC()
	_ = s.repo.UpdateLoginAt(ctx, user.ID, now)

	token, exp, err := s.generateToken(user.ID)
	if err != nil {
		return nil, errs.NewInternalServerError()
	}

	return &AuthResult{
		User:  user,
		Token: AuthToken{Token: token, ExpiresAt: exp},
	}, nil
}

func (s *AuthService) LoginWithGoogle(ctx context.Context, idToken string) (*AuthResult, error) {
	if s.googleClientID == "" {
		return nil, errs.NewBadRequestError("Google login is not configured", false, nil, nil, nil)
	}

	payload, err := idtoken.Validate(ctx, idToken, s.googleClientID)
	if err != nil {
		return nil, errs.NewUnauthorizedError("Invalid Google token", false)
	}

	subject := payload.Subject
	emailClaim, _ := payload.Claims["email"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)
	if emailClaim == "" || !emailVerified {
		return nil, errs.NewUnauthorizedError("Google account email is not verified", true)
	}

	user, findErr := s.repo.GetByGoogleID(ctx, subject)
	if findErr != nil && !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return nil, sqlerr.HandleError(findErr)
	}

	if user == nil {
		// Try to link existing account by email
		user, findErr = s.repo.GetByEmail(ctx, emailClaim)
		switch {
		case findErr == nil:
			user.GoogleID = &subject
			if err := s.repo.Save(ctx, user); err != nil {
				return nil, sqlerr.HandleError(err)
			}
		case errors.Is(findErr, gorm.ErrRecordNotFound):
			username := deriveUsername(emailClaim)
			user = &model.User{
				Email:    emailClaim,
				Username: username,
				GoogleID: &subject,
			}
			if err := s.repo.CreateUser(ctx, user); err != nil {
				return nil, sqlerr.HandleError(err)
			}
		default:
			return nil, sqlerr.HandleError(findErr)
		}
	}

	now := time.Now().UTC()
	_ = s.repo.UpdateLoginAt(ctx, user.ID, now)

	token, exp, err := s.generateToken(user.ID)
	if err != nil {
		return nil, errs.NewInternalServerError()
	}

	return &AuthResult{
		User:  user,
		Token: AuthToken{Token: token, ExpiresAt: exp},
	}, nil
}

func (s *AuthService) lookupUser(ctx context.Context, identifier string) (*model.User, error) {
	if emailRegex.MatchString(identifier) {
		return s.repo.GetByEmail(ctx, identifier)
	}
	return s.repo.GetByUsername(ctx, identifier)
}

func (s *AuthService) generateToken(userID uuid.UUID) (string, time.Time, error) {
	exp := time.Now().Add(s.accessTokenTTL)
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

func deriveUsername(email string) string {
	parts := regexp.MustCompile("@").Split(email, 2)
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return fmt.Sprintf("user-%s", uuid.New().String()[:8])
}
