package handlers

import (
	apperrors "backend-fiber/internal/errors"
	"backend-fiber/internal/models"
	"backend-fiber/internal/services"
	"errors"
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

// Thêm struct mới cho request
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required" example:"johndoe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"secret123"`
}

// CreateUser godoc
// @Summary Tạo user mới
// @Description Tạo một user mới với thông tin được cung cấp
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User Info"
// @Success 201 {object} models.User
// @Failure 400 {object} apperrors.AppError
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	req := new(CreateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validate.Struct(req); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" is "+err.Tag())
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	// Chuyển đổi từ request sang model
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
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
// @Failure 500 {object} apperrors.AppError
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
// @Failure 400 {object} apperrors.AppError "Invalid ID"
// @Failure 404 {object} apperrors.AppError "User not found"
// @Failure 500 {object} apperrors.AppError "Internal server error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(apperrors.NewAppError(
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
