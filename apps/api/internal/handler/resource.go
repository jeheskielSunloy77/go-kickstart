package handler

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/jeheskielSunloy77/go-kickstart/internal/lib/utils"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"github.com/jeheskielSunloy77/go-kickstart/internal/repository"
	"github.com/jeheskielSunloy77/go-kickstart/internal/service"
)

type emptyRequest struct{}

func (r emptyRequest) Validate() error { return nil }

type ResourceHandler[T any, S model.StoreDTO[T], U model.UpdateDTO[T]] struct {
	Handler
	service service.ResourceServiceInterface[T, S, U]
}

func NewResourceHandler[T any, S model.StoreDTO[T], U model.UpdateDTO[T]](base Handler, service service.ResourceServiceInterface[T, S, U]) *ResourceHandler[T, S, U] {
	return &ResourceHandler[T, S, U]{
		Handler: base,
		service: service,
	}
}

func (h *ResourceHandler[T, S, U]) Update() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, dto U) (*T, error) {
		id, err := utils.ParseUUIDParam(c.Params("id"))
		if err != nil {
			return nil, err
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), h.server.Config.Server.WriteTimeout)
		defer cancel()

		return h.service.Update(ctx, id, dto)
	}, http.StatusOK, newDTO[U]())
}

func (h *ResourceHandler[T, S, U]) GetByID() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ emptyRequest) (*T, error) {
		id, err := utils.ParseUUIDParam(c.Params("id"))
		if err != nil {
			return nil, err
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), h.server.Config.Server.ReadTimeout)
		defer cancel()

		preloads := repository.ParsePreloads(c.Query("preloads"))
		return h.service.GetByID(ctx, id, preloads)
	}, http.StatusOK, emptyRequest{})
}

func (h *ResourceHandler[T, S, U]) GetMany() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ emptyRequest) (model.PaginatedResponse[T], error) {
		ctx, cancel := context.WithTimeout(c.UserContext(), h.server.Config.Server.ReadTimeout)
		defer cancel()

		options := repository.NewGetManyOptionsFromRequest(c)
		entities, total, err := h.service.GetMany(ctx, options)
		if err != nil {
			return model.PaginatedResponse[T]{}, err
		}

		resp := model.PaginatedResponse[T]{
			Data:       entities,
			Page:       options.Offset/options.Limit + 1,
			Limit:      options.Limit,
			Total:      int(total),
			TotalPages: int((total + int64(options.Limit) - 1) / int64(options.Limit)),
		}
		return resp, nil
	}, http.StatusOK, emptyRequest{})
}

func (h *ResourceHandler[T, S, U]) Destroy() fiber.Handler {
	return HandleNoContent(h.Handler, func(c *fiber.Ctx, _ emptyRequest) error {
		id, err := utils.ParseUUIDParam(c.Params("id"))
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), h.server.Config.Server.WriteTimeout)
		defer cancel()

		return h.service.Destroy(ctx, id)
	}, http.StatusNoContent, emptyRequest{})
}

func (h *ResourceHandler[T, S, U]) Kill() fiber.Handler {
	return HandleNoContent(h.Handler, func(c *fiber.Ctx, _ emptyRequest) error {
		id, err := utils.ParseUUIDParam(c.Params("id"))
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), h.server.Config.Server.WriteTimeout)
		defer cancel()

		return h.service.Kill(ctx, id)
	}, http.StatusNoContent, emptyRequest{})
}

func (h *ResourceHandler[T, S, U]) Restore() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ emptyRequest) (*T, error) {
		id, err := utils.ParseUUIDParam(c.Params("id"))
		if err != nil {
			return nil, err
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), h.server.Config.Server.WriteTimeout)
		defer cancel()

		preloads := repository.ParsePreloads(c.Query("preloads"))
		return h.service.Restore(ctx, id, preloads)
	}, http.StatusOK, emptyRequest{})
}

func (h *ResourceHandler[T, S, U]) Store() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, dto S) (*T, error) {
		ctx, cancel := context.WithTimeout(c.UserContext(), h.server.Config.Server.WriteTimeout)
		defer cancel()

		return h.service.Store(ctx, dto)
	}, http.StatusCreated, newDTO[S]())
}

func newDTO[T any]() T {
	var dto T
	t := reflect.TypeOf(dto)
	if t != nil && t.Kind() == reflect.Pointer {
		return reflect.New(t.Elem()).Interface().(T)
	}
	return dto
}
