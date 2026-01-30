package handler

import (
	"github.com/jeheskielSunloy77/go-kickstart/internal/application"
	"github.com/jeheskielSunloy77/go-kickstart/internal/infrastructure/server"
)

type Handlers struct {
	Health  *HealthHandler
	Auth    *AuthHandler
	User    *UserHandler
	OpenAPI *OpenAPIHandler
}

func NewHandlers(s *server.Server, services *application.Services) *Handlers {
	h := NewHandler(s)

	return &Handlers{
		Health:  NewHealthHandler(h),
		Auth:    NewAuthHandler(h, services.Auth),
		User:    NewUserHandler(h, services.User),
		OpenAPI: NewOpenAPIHandler(h),
	}
}
