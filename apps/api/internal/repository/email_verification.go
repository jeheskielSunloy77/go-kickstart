package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"gorm.io/gorm"
)

type EmailVerificationRepository interface {
	Create(ctx context.Context, verification *model.EmailVerification) error
	GetActiveByUserIDAndCodeHash(ctx context.Context, userID uuid.UUID, codeHash string, now time.Time) (*model.EmailVerification, error)
	ExpireActiveByUserID(ctx context.Context, userID uuid.UUID, now time.Time) error
	MarkVerified(ctx context.Context, id uuid.UUID, verifiedAt time.Time) error
}

type emailVerificationRepository struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) EmailVerificationRepository {
	return &emailVerificationRepository{db: db}
}

func (r *emailVerificationRepository) Create(ctx context.Context, verification *model.EmailVerification) error {
	if verification.ID == uuid.Nil {
		verification.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(verification).Error
}

func (r *emailVerificationRepository) GetActiveByUserIDAndCodeHash(ctx context.Context, userID uuid.UUID, codeHash string, now time.Time) (*model.EmailVerification, error) {
	var verification model.EmailVerification
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND code_hash = ? AND verified_at IS NULL AND expires_at > ?", userID, codeHash, now).
		Order("created_at desc").
		First(&verification).
		Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

func (r *emailVerificationRepository) ExpireActiveByUserID(ctx context.Context, userID uuid.UUID, now time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.EmailVerification{}).
		Where("user_id = ? AND verified_at IS NULL AND expires_at > ?", userID, now).
		Update("expires_at", now).
		Error
}

func (r *emailVerificationRepository) MarkVerified(ctx context.Context, id uuid.UUID, verifiedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.EmailVerification{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"verified_at": verifiedAt,
			"updated_at":  verifiedAt,
		}).
		Error
}
