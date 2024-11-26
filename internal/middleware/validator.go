package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New() // Thêm khai báo validator

func ValidateRequest(model interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(model); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Dữ liệu không hợp lệ",
			})
		}

		if err := validate.Struct(model); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": err.Error(),
			})
		}

		return c.Next()
	}
}
