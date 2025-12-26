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

func NewRepositories(s *server.Server) *Repositories {
	var cacheClient cache.Cache
	if s.Redis != nil && s.Config.Cache.TTL > 0 {
		cacheClient = cache.NewRedisCache(s.Redis, &s.Config.Cache)
	}
	return &Repositories{
		Auth:              NewAuthRepository(s.DB.DB),
		AuthSession:       NewAuthSessionRepository(s.DB.DB),
		User:              NewUserRepository(s.DB.DB, cacheClient, s.Config.Cache.TTL),
		EmailVerification: NewEmailVerificationRepository(s.DB.DB),
	}
}
