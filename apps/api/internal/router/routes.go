package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jeheskielSunloy77/go-kickstart/internal/handler"
	"github.com/jeheskielSunloy77/go-kickstart/internal/middleware"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
)

func registerRoutes(
	r *fiber.App,
	h *handler.Handlers,
	middlewares *middleware.Middlewares,
) {
	// system routes
	r.Get("/status", h.Health.CheckHealth)
	r.Static("/static", "static")
	r.Get("/docs", h.OpenAPI.ServeOpenAPIUI)

	// versioned routes
	api := r.Group("/api/v1")

	authGroup := api.Group("/auth")
	authGroup.Post("/register", h.Auth.Register)
	authGroup.Post("/login", h.Auth.Login)
	authGroup.Post("/google", h.Auth.GoogleLogin)

	// protected routes
	protected := api.Group("", middlewares.Auth.RequireAuth())

	resource(protected, "/users", h.User.ResourceHandler)
}

func resource[T any, S model.StoreDTO[T], U model.UpdateDTO[T]](group fiber.Router, path string, h *handler.ResourceHandler[T, S, U], authMiddleware ...fiber.Handler) {
	g := group.Group(path, authMiddleware...)
	g.Get("/", h.GetMany())
	g.Get("/:id", h.GetByID())
	g.Post("/", h.Store())
	g.Delete("/:id", h.Destroy())
	g.Delete("/:id/kill", h.Kill())
	g.Patch("/:id/restore", h.Restore())
	g.Patch("/:id", h.Update())
}
