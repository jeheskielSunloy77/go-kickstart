package handler

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
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
	return Handle(h.Handler, func(c *fiber.Ctx, req *model.RegisterDTO) (*service.AuthResult, error) {
		return h.authService.Register(c.UserContext(), req.Email, req.Username, req.Password)
	}, http.StatusCreated, &model.RegisterDTO{})
}

func (h *AuthHandler) Login() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, req *model.LoginDTO) (*service.AuthResult, error) {
		identifier := req.Identifier
		if isEmail(identifier) {
			identifier = normalizeEmail(identifier)
		}

		return h.authService.Login(c.UserContext(), identifier, req.Password)
	}, http.StatusOK, &model.LoginDTO{})
}

func (h *AuthHandler) GoogleLogin() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, req *model.GoogleLoginDTO) (*service.AuthResult, error) {
		return h.authService.LoginWithGoogle(c.UserContext(), req.IDToken)
	}, http.StatusOK, &model.GoogleLoginDTO{})
}

func (h *AuthHandler) VerifyEmail() fiber.Handler {
	return Handle(h.Handler, func(c *fiber.Ctx, req *model.VerifyEmailDTO) (*model.User, error) {
		return h.authService.VerifyEmail(c.UserContext(), req.Email, req.Code)
	}, http.StatusOK, &model.VerifyEmailDTO{})
}

func isEmail(identifier string) bool {
	emailRegex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	return emailRegex.MatchString(identifier)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
