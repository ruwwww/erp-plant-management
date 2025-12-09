package middleware

import (
	"os"
	"server/internal/core/domain"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// 1. Protect: Validates JWT and injects UserID/Role into Context
func Protect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"})
		}

		// Remove "Bearer " prefix
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Inject into context for Handlers to use
		// JWT numbers are float64
		if sub, ok := claims["sub"].(float64); ok {
			c.Locals("userID", int(sub))
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token subject"})
		}

		c.Locals("role", claims["role"])

		return c.Next()
	}
}

// 2. Authorize: Checks if user has one of the required roles
func Authorize(allowedRoles ...domain.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleStr, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Role not found in token"})
		}

		userRole := domain.UserRole(userRoleStr)

		// Allow if user has ANY of the allowed roles
		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied: Insufficient privileges"})
	}
}
