package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/errs"
	"github.com/jeheskielSunloy77/go-kickstart/internal/server"
)

type AuthMiddleware struct {
	server *server.Server
	secret []byte
}

func NewAuthMiddleware(s *server.Server) *AuthMiddleware {
	return &AuthMiddleware{
		server: s,
		secret: []byte(s.Config.Auth.SecretKey),
	}
}

func (auth *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		authHeader := c.Get(fiber.HeaderAuthorization)
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		rawToken := strings.TrimSpace(parts[1])
		if rawToken == "" {
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errs.NewUnauthorizedError("invalid token", true)
			}
			return auth.secret, nil
		})
		if err != nil || !token.Valid {
			auth.server.Logger.Error().
				Err(err).
				Str("function", "RequireAuth").
				Str("request_id", GetRequestID(c)).
				Dur("duration", time.Since(start)).
				Msg("token validation failed")

			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		if claims.Subject == "" {
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		if _, err := uuid.Parse(claims.Subject); err != nil {
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		c.Locals("user_id", claims.Subject)

		auth.server.Logger.Info().
			Str("function", "RequireAuth").
			Str("user_id", claims.Subject).
			Str("request_id", GetRequestID(c)).
			Dur("duration", time.Since(start)).
			Msg("user authenticated successfully")

		return c.Next()
	}
}
