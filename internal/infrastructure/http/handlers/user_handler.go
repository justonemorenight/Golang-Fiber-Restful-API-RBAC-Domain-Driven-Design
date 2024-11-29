package handlers

import (
	apperrors "backend-fiber/internal/pkg/errors"
	"backend-fiber/internal/pkg/validator"

	"github.com/gofiber/fiber"
)

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserHandler struct {
	// TODO: add necessary dependencies
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	req := new(CreateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return apperrors.ErrBadRequest
	}

	if err := validator.ValidateStruct(req); err != nil {
		return apperrors.ErrBadRequest
	}

	return c.Status(fiber.StatusCreated).JSON(req)
}
