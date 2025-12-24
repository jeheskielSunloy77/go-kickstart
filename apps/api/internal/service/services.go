package service

import (
	"github.com/jeheskielSunloy77/go-kickstart/internal/lib/job"
	"github.com/jeheskielSunloy77/go-kickstart/internal/repository"
	"github.com/jeheskielSunloy77/go-kickstart/internal/server"
)

type Services struct {
	Auth *AuthService
	User *UserService
	Job  *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(&s.Config.Auth, repos.Auth)
	userService := NewUserService(repos.User)

	return &Services{
		Job:  s.Job,
		Auth: authService,
		User: userService,
	}, nil
}
