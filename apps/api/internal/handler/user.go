package handler

import (
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"github.com/jeheskielSunloy77/go-kickstart/internal/service"
)

type UserHandler struct {
	*ResourceHandler[model.User, *model.StoreUserDTO, *model.UpdateUserDTO]
}

func NewUserHandler(h Handler, service *service.UserService) *UserHandler {
	return &UserHandler{
		ResourceHandler: NewResourceHandler("user", h, service),
	}
}
