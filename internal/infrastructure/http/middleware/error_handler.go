package middleware

import (
	apperrors "backend-fiber/internal/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case *apperrors.AppError:
		return c.Status(e.Code).JSON(fiber.Map{
			"success": false,
			"error":   e,
		})
	case *fiber.Error:
		return c.Status(e.Code).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    e.Code,
				"message": e.Message,
			},
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Internal Server Error",
			},
		})
	}
}
