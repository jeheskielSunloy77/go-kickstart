package repository

import (
	"github.com/jeheskielSunloy77/go-kickstart/internal/server"
)

type Repositories struct {
	Auth *AuthRepository
	User *UserRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Auth: NewAuthRepository(s.DB.DB),
		User: NewUserRepository(s.DB.DB),
	}
}
