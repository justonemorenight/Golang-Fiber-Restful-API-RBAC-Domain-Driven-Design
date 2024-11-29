package middleware

import (
	"backend-fiber/internal/auth"
	"backend-fiber/internal/pkg/config"
	apperrors "backend-fiber/internal/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	cfg := config.GetConfig()

	return func(c *fiber.Ctx) error {
		token := auth.ExtractBearerToken(c.Get("Authorization"))
		if token == "" {
			return apperrors.ErrUnauthorized
		}

		claims, err := auth.ValidateToken(token, cfg.JWTSecret)
		if err != nil {
			return apperrors.ErrUnauthorized
		}

		if claims == nil {
			return apperrors.ErrUnauthorized
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}
