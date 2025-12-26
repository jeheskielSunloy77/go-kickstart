package handler

import (
	"net/http"
	"net/url"
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
	return func(c *fiber.Ctx) error {
		start, err := h.authService.StartGoogleAuth(c.UserContext())
		if err != nil {
			return h.redirectGoogleFailure(c, err)
		}

		h.setGoogleStateCookie(c, start)
		return c.Redirect(start.AuthURL, http.StatusFound)
	}
}

func (h *AuthHandler) GoogleCallback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		code := c.Query("code")
		state := c.Query("state")
		stateCookie := c.Cookies(googleStateCookieName)

		result, err := h.authService.CompleteGoogleAuth(
			c.UserContext(),
			code,
			state,
			stateCookie,
			c.Get(fiber.HeaderUserAgent),
			c.IP(),
		)

		h.clearGoogleStateCookie(c)

		if err != nil {
			return h.redirectGoogleFailure(c, err)
		}

		h.setAuthCookies(c, result)

		return c.Redirect(h.server.Config.Auth.GoogleSuccessRedirectURL, http.StatusFound)
	}
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

	sameSite := string(h.server.Config.Auth.CookieSameSite)
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
	sameSite := string(h.server.Config.Auth.CookieSameSite)
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

const googleStateCookieName = "google_auth_state"

func (h *AuthHandler) setGoogleStateCookie(c *fiber.Ctx, start *service.GoogleAuthStart) {
	if start == nil {
		return
	}

	sameSite := string(h.server.Config.Auth.CookieSameSite)
	secure := h.server.Config.Primary.Env == config.EnvProduction

	c.Cookie(&fiber.Cookie{
		Name:     googleStateCookieName,
		Value:    start.StateCookie,
		Expires:  start.StateExpiresAt,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/api/v1/auth/google",
		Domain:   h.server.Config.Auth.CookieDomain,
	})
}

func (h *AuthHandler) clearGoogleStateCookie(c *fiber.Ctx) {
	sameSite := string(h.server.Config.Auth.CookieSameSite)
	secure := h.server.Config.Primary.Env == config.EnvProduction
	expired := time.Unix(0, 0)

	c.Cookie(&fiber.Cookie{
		Name:     googleStateCookieName,
		Value:    "",
		Expires:  expired,
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/api/v1/auth/google",
		Domain:   h.server.Config.Auth.CookieDomain,
	})
}

func (h *AuthHandler) defaultAuthOrigin() string {
	if h.server == nil || h.server.Config == nil {
		return ""
	}

	if len(h.server.Config.Server.CORSAllowedOrigins) == 0 {
		return ""
	}

	origin := strings.TrimSpace(h.server.Config.Server.CORSAllowedOrigins[0])
	if origin == "" || origin == "*" {
		return ""
	}

	return origin
}

func (h *AuthHandler) redirectGoogleFailure(c *fiber.Ctx, err error) error {
	middleware.GetLogger(c).Warn().Err(err).Msg("google auth failed")
	redirectURL := appendQueryParam(h.server.Config.Auth.GoogleFailureRedirectURL, "error", "google_auth_failed")
	return c.Redirect(redirectURL, http.StatusFound)
}

func appendQueryParam(rawURL, key, value string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	query := parsed.Query()
	query.Set(key, value)
	parsed.RawQuery = query.Encode()

	return parsed.String()
}

func isEmail(identifier string) bool {
	emailRegex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	return emailRegex.MatchString(identifier)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
