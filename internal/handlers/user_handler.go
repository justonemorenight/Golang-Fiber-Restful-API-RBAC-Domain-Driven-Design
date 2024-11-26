package handlers

import (
	"backend-fiber/internal/errors"
	"backend-fiber/internal/models"
	"backend-fiber/internal/services"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *services.UserService
}

var validate = validator.New()

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser godoc
// @Summary Tạo user mới
// @Description Tạo một user mới với thông tin được cung cấp
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User Info"
// @Success 201 {object} models.User
// @Failure 400 {object} errors.AppError
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validate.Struct(user); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" is "+err.Tag())
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	if err := h.service.CreateUser(user); err != nil {
		return err
	}

	user.Password = ""
	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUsers godoc
// @Summary Lấy danh sách users
// @Description Lấy danh sách tất cả users
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} errors.AppError
// @Router /users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.service.GetUsers()
	if err != nil {
		return err
	}

	return c.JSON(users)
}

// GetUserByID godoc
// @Summary Lấy user bằng ID
// @Description Lấy thông tin chi tiết của user bằng ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} errors.AppError "Invalid ID"
// @Failure 404 {object} errors.AppError "User not found"
// @Failure 500 {object} errors.AppError "Internal server error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		if err == errors.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(errors.NewAppError(
				fiber.StatusNotFound,
				"User not found",
				fmt.Sprintf("User with ID %d not found", id),
			))
		}
		return err
	}

	user.Password = ""
	return c.JSON(user)
}
