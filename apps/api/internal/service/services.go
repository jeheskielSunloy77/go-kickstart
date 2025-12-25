package service

import (
	"github.com/jeheskielSunloy77/go-kickstart/internal/lib/job"
	"github.com/jeheskielSunloy77/go-kickstart/internal/repository"
	"github.com/jeheskielSunloy77/go-kickstart/internal/server"
)

type Services struct {
	Auth          *AuthService
	User          *UserService
	Authorization *AuthorizationService
	Job           *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	var enqueuer TaskEnqueuer
	if s.Job != nil {
		enqueuer = s.Job.Client
	}
	authService := NewAuthService(&s.Config.Auth, repos.Auth, repos.EmailVerification, enqueuer, s.Logger)
	userService := NewUserService(repos.User)
	authorizationService, err := NewAuthorizationService(s.DB.DB, s.Logger)
	if err != nil {
		return nil, err
	}

	return &Services{
		Job:           s.Job,
		Auth:          authService,
		User:          userService,
		Authorization: authorizationService,
	}, nil
}
