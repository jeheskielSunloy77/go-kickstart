package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"gorm.io/gorm"
)

type AuthSessionRepositoryInterface interface {
	Create(ctx context.Context, session *model.AuthSession) error
	GetByRefreshTokenHash(ctx context.Context, hash string) (*model.AuthSession, error)
	RevokeByID(ctx context.Context, id uuid.UUID, revokedAt time.Time) error
	RevokeByUserID(ctx context.Context, userID uuid.UUID, revokedAt time.Time) error
}

type AuthSessionRepository struct {
	db *gorm.DB
}

func NewAuthSessionRepository(db *gorm.DB) *AuthSessionRepository {
	return &AuthSessionRepository{db: db}
}

func (r *AuthSessionRepository) Create(ctx context.Context, session *model.AuthSession) error {
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *AuthSessionRepository) GetByRefreshTokenHash(ctx context.Context, hash string) (*model.AuthSession, error) {
	var session model.AuthSession
	if err := r.db.WithContext(ctx).First(&session, "refresh_token_hash = ?", hash).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AuthSessionRepository) RevokeByID(ctx context.Context, id uuid.UUID, revokedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.AuthSession{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"revoked_at": revokedAt,
		}).
		Error
}

func (r *AuthSessionRepository) RevokeByUserID(ctx context.Context, userID uuid.UUID, revokedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.AuthSession{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Updates(map[string]any{
			"revoked_at": revokedAt,
		}).
		Error
}
