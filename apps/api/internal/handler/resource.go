package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jeheskielSunloy77/go-kickstart/internal/lib/utils"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"github.com/jeheskielSunloy77/go-kickstart/internal/repository"
	"github.com/jeheskielSunloy77/go-kickstart/internal/service"
)

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
	return Handle(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) (model.PaginatedResponse[T], error) {
		options := repository.NewGetManyOptionsFromRequest(c)
		entities, total, err := h.service.GetMany(c.UserContext(), options)
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
	}, http.StatusOK, &model.EmptyDTO{})
}

func (h *ResourceHandler[T, S, U]) Destroy() fiber.Handler {
	return HandleNoContent(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) error {
		id, err := utils.ParseUUIDParam(c.Params("id"))
		if err != nil {
			return err
		}

		return h.service.Destroy(c.UserContext(), id)
	}, http.StatusNoContent, &model.EmptyDTO{})
}

func (h *ResourceHandler[T, S, U]) Kill() fiber.Handler {
	return HandleNoContent(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) error {
		id, err := utils.ParseUUIDParam(c.Params("id"))
		if err != nil {
			return err
		}

		return h.service.Kill(c.UserContext(), id)
	}, http.StatusNoContent, &model.EmptyDTO{})
}

func (h *ResourceHandler[T, S, U]) Restore() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) (*T, error) {
		id, err := utils.ParseUUIDParam(c.Params("id"))
		if err != nil {
			return nil, err
		}

		preloads := repository.ParsePreloads(c.Query("preloads"))
		return h.service.Restore(c.UserContext(), id, preloads)
	}, http.StatusOK, &model.EmptyDTO{})
}

func (h *ResourceHandler[T, S, U]) Store() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, dto S) (*T, error) {
		return h.service.Store(c.UserContext(), dto)
	}, http.StatusCreated, model.NewDTO[S]())
}
