package handlers

import (
	"backend-fiber/internal/db"
	"backend-fiber/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         *db.User `json:"user"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponseSwagger
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  err.Error(),
		})
	}

	ip := c.IP()
	userAgent := c.Get("User-Agent")

	user, accessToken, refreshToken, err := h.userService.Login(c.Context(), req.Email, req.Password, ip, userAgent)
	if err != nil {
		return err
	}

	return c.JSON(LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}
