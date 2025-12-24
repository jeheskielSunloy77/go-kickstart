package repository

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Save(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *model.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *AuthRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "LOWER(email) = ?", strings.ToLower(email)).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "LOWER(username) = ?", strings.ToLower(username)).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "google_id = ?", googleID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) UpdateLoginAt(ctx context.Context, id uuid.UUID, ts time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("last_login_at", ts).
		Error
}
