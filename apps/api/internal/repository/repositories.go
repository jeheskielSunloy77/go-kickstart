package repository

import (
	"github.com/jeheskielSunloy77/go-kickstart/internal/server"
)

type Repositories struct {
	Auth              AuthRepository
	AuthSession       AuthSessionRepository
	User              UserRepository
	EmailVerification EmailVerificationRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Auth:              NewAuthRepository(s.DB.DB),
		AuthSession:       NewAuthSessionRepository(s.DB.DB),
		User:              NewUserRepository(s.DB.DB),
		EmailVerification: NewEmailVerificationRepository(s.DB.DB),
	}
}
