package model

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID  `json:"userId" gorm:"type:uuid;not null;index"`
	Email      string     `json:"email" gorm:"not null"`
	CodeHash   string     `json:"-" gorm:"not null"`
	ExpiresAt  time.Time  `json:"expiresAt" gorm:"not null"`
	VerifiedAt *time.Time `json:"verifiedAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

func (EmailVerification) TableName() string {
	return "email_verifications"
}
