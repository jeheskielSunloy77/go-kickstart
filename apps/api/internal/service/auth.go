package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jeheskielSunloy77/go-kickstart/internal/config"
	"github.com/jeheskielSunloy77/go-kickstart/internal/errs"
	"github.com/jeheskielSunloy77/go-kickstart/internal/lib/job"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"github.com/jeheskielSunloy77/go-kickstart/internal/repository"
	"github.com/jeheskielSunloy77/go-kickstart/internal/sqlerr"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

var (
	minPasswordLength = 8
	emailRegex        = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

type AuthService struct {
	repo                 repository.AuthRepositoryInterface
	verificationRepo     repository.EmailVerificationRepositoryInterface
	taskEnqueuer         TaskEnqueuer
	logger               *zerolog.Logger
	secretKey            []byte
	accessTokenTTL       time.Duration
	googleClientID       string
	emailVerificationTTL time.Duration
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
	VerifyEmail(ctx context.Context, email, code string) (*model.User, error)
}

type TaskEnqueuer interface {
	EnqueueContext(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

func NewAuthService(cfg *config.AuthConfig, repo repository.AuthRepositoryInterface, verificationRepo repository.EmailVerificationRepositoryInterface, taskEnqueuer TaskEnqueuer, logger *zerolog.Logger) *AuthService {
	return &AuthService{
		repo:                 repo,
		verificationRepo:     verificationRepo,
		taskEnqueuer:         taskEnqueuer,
		logger:               logger,
		secretKey:            []byte(cfg.SecretKey),
		accessTokenTTL:       cfg.AccessTokenTTL,
		googleClientID:       cfg.GoogleClientID,
		emailVerificationTTL: cfg.EmailVerificationTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, email, username, password string) (*AuthResult, error) {
	if len(password) < minPasswordLength {
		return nil, errs.NewBadRequestError(
			fmt.Sprintf("Password must be at least %d characters", minPasswordLength),
			true,
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

	if err := s.queueEmailVerification(ctx, user); err != nil {
		s.logVerificationQueueError(err)
	}

	token, exp, err := s.generateToken(user)
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

	token, exp, err := s.generateToken(user)
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
		return nil, errs.NewBadRequestError("Google login is not configured", false, nil, nil)
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
	if user.EmailVerifiedAt == nil {
		_ = s.repo.UpdateEmailVerifiedAt(ctx, user.ID, now)
		user.EmailVerifiedAt = &now
	}

	token, exp, err := s.generateToken(user)
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

func (s *AuthService) VerifyEmail(ctx context.Context, email, code string) (*model.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, invalidVerificationError()
		}
		return nil, sqlerr.HandleError(err)
	}

	if user.EmailVerifiedAt != nil {
		return user, nil
	}

	if s.verificationRepo == nil {
		return nil, errs.NewInternalServerError()
	}

	codeHash := hashVerificationCode(code)
	now := time.Now().UTC()
	verification, err := s.verificationRepo.GetActiveByUserIDAndCodeHash(ctx, user.ID, codeHash, now)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, invalidVerificationError()
		}
		return nil, sqlerr.HandleError(err)
	}

	if err := s.verificationRepo.MarkVerified(ctx, verification.ID, now); err != nil {
		return nil, sqlerr.HandleError(err)
	}

	if err := s.repo.UpdateEmailVerifiedAt(ctx, user.ID, now); err != nil {
		return nil, sqlerr.HandleError(err)
	}

	user.EmailVerifiedAt = &now
	return user, nil
}

func (s *AuthService) generateToken(user *model.User) (string, time.Time, error) {
	if user == nil {
		return "", time.Time{}, errs.NewInternalServerError()
	}

	exp := time.Now().Add(s.accessTokenTTL)
	claims := model.AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
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

func (s *AuthService) queueEmailVerification(ctx context.Context, user *model.User) error {
	if user == nil || user.Email == "" || user.EmailVerifiedAt != nil || s.verificationRepo == nil {
		return nil
	}

	code, err := generateVerificationCode()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	ttl := s.emailVerificationTTL
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	if err := s.verificationRepo.ExpireActiveByUserID(ctx, user.ID, now); err != nil {
		return err
	}

	verification := &model.EmailVerification{
		UserID:    user.ID,
		Email:     user.Email,
		CodeHash:  hashVerificationCode(code),
		ExpiresAt: now.Add(ttl),
	}
	if err := s.verificationRepo.Create(ctx, verification); err != nil {
		return err
	}

	if s.taskEnqueuer == nil {
		return nil
	}

	expiresInMinutes := int(ttl.Minutes())
	if expiresInMinutes <= 0 {
		expiresInMinutes = 1
	}
	task, err := job.NewEmailVerificationTask(job.EmailVerificationPayload{
		To:               user.Email,
		Username:         user.Username,
		Code:             code,
		ExpiresInMinutes: expiresInMinutes,
	})
	if err != nil {
		return err
	}

	_, err = s.taskEnqueuer.EnqueueContext(ctx, task)
	return err
}

func (s *AuthService) logVerificationQueueError(err error) {
	if err == nil || s.logger == nil {
		return
	}
	s.logger.Error().Err(err).Msg("failed to queue email verification")
}

func generateVerificationCode() (string, error) {
	const codeLength = 6
	const maxDigit = 10

	code := make([]byte, 0, codeLength)
	for range codeLength {
		n, err := rand.Int(rand.Reader, big.NewInt(maxDigit))
		if err != nil {
			return "", err
		}
		code = append(code, byte('0'+n.Int64()))
	}

	return string(code), nil
}

func hashVerificationCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}

func invalidVerificationError() *errs.ErrorResponse {
	return errs.NewBadRequestError(
		"Invalid or expired verification code",
		true,
		[]errs.FieldError{{Field: "code", Error: "invalid or expired"}},
		nil,
	)
}
