package middleware

import (
	"backend-fiber/internal/auth"
	"backend-fiber/internal/pkg/config"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	cfg := config.GetConfig()

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token format",
			})
		}

		// Verify token with secret key from config
		token, err := jwt.ParseWithClaims(parts[1], &auth.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Add log for debugging
			log.Printf("Using secret key: %s", cfg.JWTSecret)
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			log.Printf("Token validation error: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
				"error":   err.Error(),
			})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token is not valid",
			})
		}

		claims := token.Claims.(*auth.JWTClaims)
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}
