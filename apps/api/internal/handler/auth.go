package handler

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jeheskielSunloy77/go-kickstart/internal/service"
	"github.com/jeheskielSunloy77/go-kickstart/internal/validation"
)

type AuthHandler struct {
	Handler
	authService *service.AuthService
}

type registerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

type loginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type googleLoginRequest struct {
	IDToken string `json:"idToken" validate:"required"`
}

func (r registerRequest) Validate() error {
	return validator.New().Struct(r)
}

func (r loginRequest) Validate() error {
	return validator.New().Struct(r)
}

func (r googleLoginRequest) Validate() error {
	return validator.New().Struct(r)
}

func NewAuthHandler(h Handler, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		Handler:     h,
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	result, err := h.authService.Register(ctx, req.Email, req.Username, req.Password)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(result)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	identifier := req.Identifier
	if isEmail(identifier) {
		identifier = normalizeEmail(identifier)
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	result, err := h.authService.Login(ctx, identifier, req.Password)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}

func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	var req googleLoginRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	result, err := h.authService.LoginWithGoogle(ctx, req.IDToken)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}

func isEmail(identifier string) bool {
	emailRegex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	return emailRegex.MatchString(identifier)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
