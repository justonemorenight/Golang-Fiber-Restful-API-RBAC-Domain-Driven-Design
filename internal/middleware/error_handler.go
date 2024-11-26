package middleware

import (
	"backend-fiber/internal/errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Mặc định là internal server error
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	detail := ""

	// Check error type
	if e, ok := err.(*errors.AppError); ok {
		code = e.Code
		message = e.Message
		detail = e.Detail
	} else if e, ok := err.(*fiber.Error); ok {
		// Handle Fiber's built-in errors
		code = e.Code
		message = e.Message
	}

	// Log error ở đây nếu cần
	// logger.Error(err)

	// Return error response
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    code,
			"message": message,
			"detail":  detail,
		},
	})
}
