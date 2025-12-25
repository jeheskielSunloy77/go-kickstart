package handler

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/config"
	"github.com/jeheskielSunloy77/go-kickstart/internal/errs"
	"github.com/jeheskielSunloy77/go-kickstart/internal/middleware"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"github.com/jeheskielSunloy77/go-kickstart/internal/server"
	"github.com/jeheskielSunloy77/go-kickstart/internal/service"
)

type AuthHandler struct {
	Handler
	authService service.AuthServiceInterface
}

func NewAuthHandler(h Handler, authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		Handler:     h,
		authService: authService,
	}
}

func (h *AuthHandler) Register() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, req *model.RegisterDTO) (*model.User, error) {
		result, err := h.authService.Register(c.UserContext(), req.Email, req.Username, req.Password, c.Get(fiber.HeaderUserAgent), c.IP())
		if err != nil {
			return nil, err
		}
		h.setAuthCookies(c, result)
		return result.User, nil
	}, http.StatusCreated, &model.RegisterDTO{})
}

func (h *AuthHandler) Login() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, req *model.LoginDTO) (*model.User, error) {
		identifier := req.Identifier
		if isEmail(identifier) {
			identifier = normalizeEmail(identifier)
		}

		result, err := h.authService.Login(c.UserContext(), identifier, req.Password, c.Get(fiber.HeaderUserAgent), c.IP())
		if err != nil {
			return nil, err
		}
		h.setAuthCookies(c, result)
		return result.User, nil
	}, http.StatusOK, &model.LoginDTO{})
}

func (h *AuthHandler) GoogleLogin() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, req *model.GoogleLoginDTO) (*model.User, error) {
		result, err := h.authService.LoginWithGoogle(c.UserContext(), req.IDToken, c.Get(fiber.HeaderUserAgent), c.IP())
		if err != nil {
			return nil, err
		}
		h.setAuthCookies(c, result)
		return result.User, nil
	}, http.StatusOK, &model.GoogleLoginDTO{})
}

func (h *AuthHandler) VerifyEmail() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, req *model.VerifyEmailDTO) (*model.User, error) {
		return h.authService.VerifyEmail(c.UserContext(), req.Email, req.Code)
	}, http.StatusOK, &model.VerifyEmailDTO{})
}

func (h *AuthHandler) Refresh() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) (*model.User, error) {
		refreshToken := c.Cookies(h.refreshCookieName())
		result, err := h.authService.Refresh(c.UserContext(), refreshToken, c.Get(fiber.HeaderUserAgent), c.IP())
		if err != nil {
			return nil, err
		}
		h.setAuthCookies(c, result)
		return result.User, nil
	}, http.StatusOK, &model.EmptyDTO{})
}

func (h *AuthHandler) Me() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) (*model.User, error) {
		userID, err := h.parseUserID(c)
		if err != nil {
			return nil, err
		}
		return h.authService.CurrentUser(c.UserContext(), userID)
	}, http.StatusOK, &model.EmptyDTO{})
}

func (h *AuthHandler) ResendVerification() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) (*server.Response[any], error) {
		userID, err := h.parseUserID(c)
		if err != nil {
			return nil, err
		}

		if err := h.authService.ResendVerification(c.UserContext(), userID); err != nil {
			return nil, err
		}

		resp := server.Response[any]{
			Status:  http.StatusOK,
			Success: true,
			Message: "Verification email sent if needed.",
		}
		return &resp, nil
	}, http.StatusOK, &model.EmptyDTO{})
}

func (h *AuthHandler) Logout() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) (*server.Response[any], error) {
		refreshToken := c.Cookies(h.refreshCookieName())
		if err := h.authService.Logout(c.UserContext(), refreshToken); err != nil {
			return nil, err
		}
		h.clearAuthCookies(c)

		resp := server.Response[any]{
			Status:  http.StatusOK,
			Success: true,
			Message: "Logged out successfully.",
		}
		return &resp, nil
	}, http.StatusOK, &model.EmptyDTO{})
}

func (h *AuthHandler) LogoutAll() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, _ *model.EmptyDTO) (*server.Response[any], error) {
		userID, err := h.parseUserID(c)
		if err != nil {
			return nil, err
		}

		if err := h.authService.LogoutAll(c.UserContext(), userID); err != nil {
			return nil, err
		}
		h.clearAuthCookies(c)

		resp := server.Response[any]{
			Status:  http.StatusOK,
			Success: true,
			Message: "Logged out from all sessions.",
		}
		return &resp, nil
	}, http.StatusOK, &model.EmptyDTO{})
}

func (h *AuthHandler) parseUserID(c *fiber.Ctx) (uuid.UUID, error) {
	raw := middleware.GetUserID(c)
	if raw == "" {
		return uuid.Nil, errs.NewUnauthorizedError("Unauthorized", false)
	}
	userID, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, errs.NewUnauthorizedError("Unauthorized", false)
	}
	return userID, nil
}

func (h *AuthHandler) setAuthCookies(c *fiber.Ctx, result *service.AuthResult) {
	if result == nil {
		return
	}

	sameSite := cookieSameSiteMode(h.server.Config.Auth.CookieSameSite)
	secure := h.server.Config.Primary.Env == config.EnvProduction

	accessCookie := &fiber.Cookie{
		Name:     h.accessCookieName(),
		Value:    result.Token.Token,
		Expires:  result.Token.ExpiresAt,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/",
		Domain:   h.server.Config.Auth.CookieDomain,
	}
	refreshCookie := &fiber.Cookie{
		Name:     h.refreshCookieName(),
		Value:    result.RefreshToken.Token,
		Expires:  result.RefreshToken.ExpiresAt,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/",
		Domain:   h.server.Config.Auth.CookieDomain,
	}

	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)
}

func (h *AuthHandler) clearAuthCookies(c *fiber.Ctx) {
	sameSite := cookieSameSiteMode(h.server.Config.Auth.CookieSameSite)
	secure := h.server.Config.Primary.Env == config.EnvProduction
	expired := time.Unix(0, 0)

	c.Cookie(&fiber.Cookie{
		Name:     h.accessCookieName(),
		Value:    "",
		Expires:  expired,
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/",
		Domain:   h.server.Config.Auth.CookieDomain,
	})
	c.Cookie(&fiber.Cookie{
		Name:     h.refreshCookieName(),
		Value:    "",
		Expires:  expired,
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/",
		Domain:   h.server.Config.Auth.CookieDomain,
	})
}

func (h *AuthHandler) accessCookieName() string {
	if h.server != nil && h.server.Config.Auth.AccessCookieName != "" {
		return h.server.Config.Auth.AccessCookieName
	}
	return "access_token"
}

func (h *AuthHandler) refreshCookieName() string {
	if h.server != nil && h.server.Config.Auth.RefreshCookieName != "" {
		return h.server.Config.Auth.RefreshCookieName
	}
	return "refresh_token"
}

func cookieSameSiteMode(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "strict":
		return fiber.CookieSameSiteStrictMode
	case "none":
		return fiber.CookieSameSiteNoneMode
	default:
		return fiber.CookieSameSiteLaxMode
	}
}

func isEmail(identifier string) bool {
	emailRegex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	return emailRegex.MatchString(identifier)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
