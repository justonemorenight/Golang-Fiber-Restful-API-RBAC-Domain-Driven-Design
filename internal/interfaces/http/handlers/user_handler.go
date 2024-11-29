package handlers

import (
	"backend-fiber/internal/application/user"
	customerrors "backend-fiber/internal/pkg/errors"
	"errors"
	"fmt"
	"time"

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

// Add new struct for request
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required" example:"johndoe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"secret123"`
}

// UserResponse represents the user response for swagger documentation
type SwaggerUserResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User info"
// @Success 201 {object} SwaggerUserResponse
// @Failure 400 {object} ErrorResponse
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

	dbUser, err := h.service.CreateUser(c.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}

	response := SwaggerUserResponse{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetUsers godoc
// @Summary Get all users
// @Description Get all users in the system
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} SwaggerUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	dbUsers, err := h.service.GetUsers(c.Context())
	if err != nil {
		return err
	}

	var response []SwaggerUserResponse
	for _, user := range dbUsers {
		response = append(response, SwaggerUserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
		})
	}
	return c.JSON(response)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} SwaggerUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	dbUser, err := h.service.GetUserByID(c.Context(), int32(id))
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

	response := SwaggerUserResponse{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}
	return c.JSON(response)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile information
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int32)

	dbUser, err := h.service.GetUserByID(c.Context(), userID)
	if err != nil {
		return err
	}

	response := SwaggerUserResponse{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}
	return c.JSON(response)
}
