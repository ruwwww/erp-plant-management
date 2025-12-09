package handlers

import (
	"server/internal/dto"
	"server/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate request (simplified, ideally use a validator library)
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}

	accessToken, refreshToken, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400, // 24 hours
		// User summary would ideally come from service or decoded token
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	// TODO: Implement refresh logic in service
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// TODO: Implement logout logic (e.g. blacklist token)
	return c.SendStatus(fiber.StatusOK)
}

func (h *AuthHandler) RequestPasswordReset(c *fiber.Ctx) error {
	// TODO: Implement request password reset logic in service
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *AuthHandler) ConfirmPasswordReset(c *fiber.Ctx) error {
	// TODO: Implement confirm password reset logic in service
	return c.SendStatus(fiber.StatusNotImplemented)
}
