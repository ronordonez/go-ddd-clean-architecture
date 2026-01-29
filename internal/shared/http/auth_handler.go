package http

import (
	"net/http"
	"time"

	"go-architecture/internal/shared/auth"
	"go-architecture/internal/shared/config"

	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// Simple in-memory auth for demo purposes. Replace with real auth in production.
func LoginHandler(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}
		if req.Username == "" || req.Password == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "username and password required"})
		}

		// Demo logic: username admin/password admin123 => admin role
		var role string
		var userID string
		var email string
		if req.Username == "admin" && req.Password == "admin123" {
			role = "admin"
			userID = "1"
			email = "admin@example.com"
		} else if req.Password == "password" {
			role = "user"
			userID = "2"
			email = req.Username + "@example.com"
		} else {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}

		token, exp, err := auth.GenerateToken(userID, email, role, cfg.JWT.Secret, cfg.JWT.Expiration)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
		}

		return c.Status(http.StatusOK).JSON(LoginResponse{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresAt:   exp,
		})
	}
}
