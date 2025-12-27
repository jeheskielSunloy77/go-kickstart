package repository

import (
	"github.com/jeheskielSunloy77/go-kickstart/internal/lib/cache"
	"github.com/jeheskielSunloy77/go-kickstart/internal/server"
)

type Repositories struct {
	Auth              AuthRepository
	AuthSession       AuthSessionRepository
	User              UserRepository
	EmailVerification EmailVerificationRepository
}

func NewRepositories(s *server.Server, cacheClient cache.Cache) *Repositories {
	return &Repositories{
		Auth:              NewAuthRepository(s.DB.DB),
		AuthSession:       NewAuthSessionRepository(s.DB.DB),
		User:              NewUserRepository(s.Config, s.DB.DB, cacheClient),
		EmailVerification: NewEmailVerificationRepository(s.DB.DB),
	}
}
