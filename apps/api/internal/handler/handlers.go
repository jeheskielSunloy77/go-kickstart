package handler

import (
	"github.com/jeheskielSunloy77/go-kickstart/internal/server"
	"github.com/jeheskielSunloy77/go-kickstart/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	Auth    *AuthHandler
	User    *UserHandler
	OpenAPI *OpenAPIHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	h := NewHandler(s)

	return &Handlers{
		Health:  NewHealthHandler(h),
		Auth:    NewAuthHandler(h, services.Auth),
		User:    NewUserHandler(h, services.User),
		OpenAPI: NewOpenAPIHandler(h),
	}
}
