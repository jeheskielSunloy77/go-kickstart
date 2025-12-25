package model

import (
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

func (d *RegisterDTO) Validate() error {
	return validator.New().Struct(d)
}

type LoginDTO struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

func (d *LoginDTO) Validate() error {
	return validator.New().Struct(d)
}

type GoogleLoginDTO struct {
	IDToken string `json:"idToken" validate:"required"`
}

func (d *GoogleLoginDTO) Validate() error {
	return validator.New().Struct(d)
}

type VerifyEmailDTO struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,min=4,max=10"`
}

func (d *VerifyEmailDTO) Validate() error {
	return validator.New().Struct(d)
}

type AuthClaims struct {
	jwt.RegisteredClaims
	Email   string `json:"email,omitempty"`
	IsAdmin bool   `json:"is_admin"`
}
