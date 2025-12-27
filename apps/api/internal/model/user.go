package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents an application user with local and federated auth support.
type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`

	Email           string     `json:"email" gorm:"uniqueIndex;not null"`
	Username        string     `json:"username" gorm:"not null"`
	PasswordHash    string     `json:"-"`
	GoogleID        *string    `json:"googleId,omitempty" gorm:"uniqueIndex"`
	EmailVerifiedAt *time.Time `json:"emailVerifiedAt,omitempty"`
	LastLoginAt     *time.Time `json:"lastLoginAt,omitempty"`
	IsAdmin         bool       `json:"isAdmin" gorm:"not null;default:false"`
}

func (m User) GetID() uuid.UUID {
	return m.ID
}

type StoreUserDTO struct {
	Email    string  `json:"email" validate:"required,email"`
	Username string  `json:"username" validate:"required,min=3,max=50"`
	Password string  `json:"password" validate:"min=8,max=128"`
	GoogleID *string `json:"googleId" validate:"omitempty"`
}

func (d *StoreUserDTO) Validate() error {
	return validator.New().Struct(d)
}

func (d *StoreUserDTO) ToModel() *User {
	return &User{
		Email:    d.Email,
		Username: d.Username,
		GoogleID: d.GoogleID,
	}
}

type UpdateUserDTO struct {
	Email    *string `json:"email" validate:"omitempty,email"`
	Username *string `json:"username" validate:"omitempty,min=3,max=50"`
	Password *string `json:"password" validate:"omitempty,min=8,max=128"`
}

func (d *UpdateUserDTO) ToModel() *User {
	user := &User{}
	if d.Email != nil {
		user.Email = *d.Email
	}
	if d.Username != nil {
		user.Username = *d.Username
	}
	if d.Password != nil {
		user.PasswordHash = *d.Password
	}
	return user
}

func (d *UpdateUserDTO) ToMap() map[string]any {
	updates := make(map[string]any)
	if d.Email != nil {
		updates["email"] = *d.Email
	}
	if d.Username != nil {
		updates["username"] = *d.Username
	}
	if d.Password != nil {
		updates["password_hash"] = *d.Password
	}
	return updates
}

func (d *UpdateUserDTO) Validate() error {
	return validator.New().Struct(d)
}
