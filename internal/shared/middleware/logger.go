package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go-architecture/internal/shared/logger"
)

func RequestLogger(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		log.Info("HTTP Request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration", duration.Milliseconds(),
			"ip", c.IP(),
		)

		return err
	}
}
