package handlers

import (
	appuser "backend-fiber/internal/application/user"
	models "backend-fiber/internal/models"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userService *appuser.Service
}

func NewAuthHandler(userService *appuser.Service) *AuthHandler {
	return &AuthHandler{userService: userService}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	User         models.UserResponse `json:"user"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Errors  string `json:"errors,omitempty"`
}

// Login godoc
// @Summary Login user
// @Description Login for existing user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "Invalid request body",
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "Validation error",
			Errors:  err.Error(),
		})
	}

	ip := c.IP()
	userAgent := c.Get("User-Agent")

	dbUser, accessToken, refreshToken, err := h.userService.Login(c.Context(), req.Email, req.Password, ip, userAgent)
	if err != nil {
		return err
	}

	userResponse := models.UserResponse{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}

	return c.JSON(LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userResponse,
	})
}
