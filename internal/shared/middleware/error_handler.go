package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go-architecture/internal/shared/errors"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default error
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	var details map[string]interface{}

	// Check if it's an AppError
	if appErr, ok := err.(*errors.AppError); ok {
		code = appErr.Code
		message = appErr.Message
		details = appErr.Details
	} else if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": message,
			"details": details,
		},
	})
}
