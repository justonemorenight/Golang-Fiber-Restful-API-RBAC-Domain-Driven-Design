package handlers

import (
	"backend-fiber/internal/application/user"
	customerrors "backend-fiber/internal/pkg/errors"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *user.Service
}

var validate = validator.New()

func NewUserHandler(service *user.Service) *UserHandler {
	return &UserHandler{service: service}
}

// Thêm struct mới cho request
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required" example:"johndoe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"secret123"`
}

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User information"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /register [post]
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

	// Gọi service với context và các tham số
	user, err := h.service.CreateUser(c.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}

	// Ẩn password trước khi trả về
	user.Password = ""
	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve all users from the system
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.service.GetUsers(c.Context())
	if err != nil {
		return err
	}

	// Ẩn password của tất cả users
	for i := range users {
		users[i].Password = ""
	}
	return c.JSON(users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security ApiKeyAuth
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	user, err := h.service.GetUserByID(c.Context(), int32(id))
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(customerrors.NewAppError(
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