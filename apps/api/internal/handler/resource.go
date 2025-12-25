package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jeheskielSunloy77/go-kickstart/internal/lib/utils"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"github.com/jeheskielSunloy77/go-kickstart/internal/repository"
	"github.com/jeheskielSunloy77/go-kickstart/internal/server"
	"github.com/jeheskielSunloy77/go-kickstart/internal/service"
)

type ResourceHandler[T any, S model.StoreDTO[T], U model.UpdateDTO[T]] struct {
	Handler
	resourceName string
	service      service.ResourceServiceInterface[T, S, U]
}

func NewResourceHandler[T any, S model.StoreDTO[T], U model.UpdateDTO[T]](resourceName string, base Handler, service service.ResourceServiceInterface[T, S, U]) *ResourceHandler[T, S, U] {
	return &ResourceHandler[T, S, U]{
		resourceName: resourceName,
		Handler:      base,
		service:      service,
	}
}

func (h *ResourceHandler[T, S, U]) Update() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, dto U) (*T, error) {
		id, err := utils.ParseUUIDParam(c.Params("id"))
		if err != nil {
			return nil, err
		}

		return h.service.Update(c.UserContext(), id, dto)
	}, http.StatusOK, model.NewDTO[U]())
}

func (h *ResourceHandler[T, S, U]) GetByID() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) (*T, error) {
		id, err := utils.ParseUUIDParam(c.Params("id"))
		if err != nil {
			return nil, err
		}

		preloads := repository.ParsePreloads(c.Query("preloads"))
		return h.service.GetByID(c.UserContext(), id, preloads)
	}, http.StatusOK, &model.EmptyDTO{})
}

func (h *ResourceHandler[T, S, U]) GetMany() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) (server.PaginatedResponse[T], error) {
		options := repository.NewGetManyOptionsFromRequest(c)
		entities, total, err := h.service.GetMany(c.UserContext(), options)
		if err != nil {
			return server.PaginatedResponse[T]{}, err
		}

		resp := server.NewPaginatedResponse("Successfully fetched "+h.resourceName+"s!", entities, total, options.Limit, options.Offset)
		return resp, nil
	}, http.StatusOK, &model.EmptyDTO{})
}

func (h *ResourceHandler[T, S, U]) Destroy() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) (*server.Response[T], error) {
		id, err := utils.ParseUUIDParam(c.Params("id"))
		if err != nil {
			return nil, err
		}

		err = h.service.Destroy(c.UserContext(), id)
		if err != nil {
			return nil, err
		}

		resp := server.Response[T]{
			Status:  http.StatusNoContent,
			Success: true,
			Message: "Successfully deleted " + h.resourceName + "!",
		}

		return &resp, nil
	}, http.StatusNoContent, &model.EmptyDTO{})
}

func (h *ResourceHandler[T, S, U]) Kill() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) (*server.Response[T], error) {
		id, err := utils.ParseUUIDParam(c.Params("id"))
		if err != nil {
			return nil, err
		}

		err = h.service.Kill(c.UserContext(), id)
		if err != nil {
			return nil, err
		}

		resp := server.Response[T]{
			Status:  http.StatusNoContent,
			Success: true,
			Message: "Successfully permanently deleted " + h.resourceName + "!",
		}

		return &resp, nil
	}, http.StatusNoContent, &model.EmptyDTO{})
}

func (h *ResourceHandler[T, S, U]) Restore() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) (*server.Response[T], error) {
		id, err := utils.ParseUUIDParam(c.Params("id"))
		if err != nil {
			return nil, err
		}

		preloads := repository.ParsePreloads(c.Query("preloads"))
		entity, err := h.service.Restore(c.UserContext(), id, preloads)
		if err != nil {
			return nil, err
		}

		resp := server.Response[T]{
			Status:  http.StatusOK,
			Success: true,
			Message: "Successfully restored " + h.resourceName + "!",
			Data:    entity,
		}

		return &resp, nil
	}, http.StatusOK, &model.EmptyDTO{})
}

func (h *ResourceHandler[T, S, U]) Store() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, dto S) (*server.Response[T], error) {
		entity, err := h.service.Store(c.UserContext(), dto)
		if err != nil {
			return nil, err
		}

		resp := server.Response[T]{
			Status:  http.StatusCreated,
			Success: true,
			Message: "Successfully created " + h.resourceName + "!",
			Data:    entity,
		}

		return &resp, nil
	}, http.StatusCreated, model.NewDTO[S]())
}
