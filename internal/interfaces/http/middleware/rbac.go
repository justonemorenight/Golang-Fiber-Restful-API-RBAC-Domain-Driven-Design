package middleware

import (
	"backend-fiber/internal/application/rbac"

	"github.com/gofiber/fiber/v2"
)

type RBACMiddleware struct {
	rbacService *rbac.Service
}

func NewRBACMiddleware(rbacService *rbac.Service) *RBACMiddleware {
	return &RBACMiddleware{
		rbacService: rbacService,
	}
}

func (m *RBACMiddleware) RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(int32)

		if permission == "users.read_self" && c.Path() != "/api/v1/users/profile" {
			paramID, err := c.ParamsInt("id")
			if err != nil {
				return fiber.ErrBadRequest
			}

			if int32(paramID) != userID {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Không có quyền xem thông tin của user khác",
				})
			}
		}

		hasPermission, err := m.rbacService.CheckPermission(c.Context(), userID, permission)
		if err != nil {
			return err
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Không có quyền truy cập",
			})
		}

		return c.Next()
	}
}
