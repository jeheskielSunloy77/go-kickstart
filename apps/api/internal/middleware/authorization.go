package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jeheskielSunloy77/go-kickstart/internal/errs"
	"github.com/jeheskielSunloy77/go-kickstart/internal/service"
)

type AuthorizationMiddleware struct {
	authorizer AuthorizationEnforcer
}

type AuthorizationEnforcer interface {
	Enforce(ctx context.Context, sub service.AuthorizationSubject, obj service.AuthorizationObject, act string) (bool, error)
}

func NewAuthorizationMiddleware(authorizer AuthorizationEnforcer) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{authorizer: authorizer}
}

func (am *AuthorizationMiddleware) RequireAuthorization() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := GetUserID(c)
		if userID == "" {
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		isAdmin := GetUserIsAdmin(c)
		if isAdmin {
			return c.Next()
		}

		if am.authorizer == nil {
			return errs.NewInternalServerError()
		}

		subject := service.AuthorizationSubject{
			ID:      userID,
			Email:   GetUserEmail(c),
			IsAdmin: isAdmin,
		}
		obj := service.AuthorizationObject{
			Route:  routePattern(c),
			Path:   c.Path(),
			Params: c.AllParams(),
			Query:  c.Queries(),
		}

		allowed, err := am.authorizer.Enforce(c.UserContext(), subject, obj, c.Method())
		if err != nil {
			logger := GetLogger(c)
			logger.Error().Err(err).Msg("authorization check failed")
			return errs.NewInternalServerError()
		}
		if !allowed {
			return errs.NewForbiddenError("Forbidden", false)
		}

		return c.Next()
	}
}

func routePattern(c *fiber.Ctx) string {
	if c.Route() != nil && c.Route().Path != "" {
		return c.Route().Path
	}
	return c.Path()
}
