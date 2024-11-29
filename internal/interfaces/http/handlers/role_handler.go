package handlers

import (
	"backend-fiber/internal/application/rbac"
	"backend-fiber/internal/db"
	"backend-fiber/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoleHandler struct {
	rbacService *rbac.Service
}

func NewRoleHandler(rbacService *rbac.Service) *RoleHandler {
	return &RoleHandler{
		rbacService: rbacService,
	}
}

// @Summary Create a new role
// @Description Create a new role in the system
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.SwaggerCreateRoleRequest true "Create role request"
// @Success 201 {object} models.SwaggerResponse{data=models.SwaggerRoleResponse}
// @Failure 400 {object} models.SwaggerResponse{error=models.ErrorData}
// @Failure 401 {object} models.SwaggerResponse{error=models.ErrorData}
// @Router /roles [post]
func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var req CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	role, err := h.rbacService.CreateRole(c.Context(), db.CreateRoleParams{
		Name: req.Name,
		Description: pgtype.Text{
			String: req.Description,
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(role)
}

// @Summary Get role by ID
// @Description Get role information by ID
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Role ID"
// @Success 200 {object} models.RoleResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /roles/{id} [get]
func (h *RoleHandler) GetRole(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	role, err := h.rbacService.GetRole(c.Context(), int32(id))
	if err != nil {
		return err
	}

	return c.JSON(toRoleResponse(role))
}

// @Summary List all roles
// @Description Get all roles in the system
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {array} models.RoleResponse
// @Router /roles [get]
func (h *RoleHandler) ListRoles(c *fiber.Ctx) error {
	roles, err := h.rbacService.ListRoles(c.Context())
	if err != nil {
		return err
	}

	var response []models.RoleResponse
	for _, role := range roles {
		response = append(response, toRoleResponse(role))
	}
	return c.JSON(response)
}

// @Summary Update role
// @Description Update role information
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param role body UpdateRoleRequest true "Role info"
// @Success 200 {object} db.Role
// @Router /roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	var req UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	err = h.rbacService.UpdateRole(c.Context(), db.UpdateRoleParams{
		ID:   int32(id),
		Name: req.Name,
		Description: pgtype.Text{
			String: req.Description,
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

// @Summary Delete role
// @Description Delete a role from the system
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 204 "No Content"
// @Router /roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	err = h.rbacService.DeleteRole(c.Context(), int32(id))
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Assign permission to role
// @Description Assign a permission to a role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param permission body AssignPermissionRequest true "Permission info"
// @Success 200 "OK"
// @Router /roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermission(c *fiber.Ctx) error {
	roleID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	var req AssignPermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	err = h.rbacService.AssignPermissionToRole(c.Context(), db.AssignPermissionToRoleParams{
		RoleID:       int32(roleID),
		PermissionID: req.PermissionID,
	})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

// @Summary Remove permission from role
// @Description Remove a permission from a role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param permissionId path int true "Permission ID"
// @Success 204 "No Content"
// @Router /roles/{id}/permissions/{permissionId} [delete]
func (h *RoleHandler) RemovePermission(c *fiber.Ctx) error {
	roleID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	permissionID, err := c.ParamsInt("permissionId")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid permission ID")
	}

	err = h.rbacService.RemovePermissionFromRole(c.Context(), roleID, permissionID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary List role permissions
// @Description Get all permissions assigned to a role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {array} db.Permission
// @Router /roles/{id}/permissions [get]
func (h *RoleHandler) ListRolePermissions(c *fiber.Ctx) error {
	roleID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	permissions, err := h.rbacService.GetRolePermissions(c.Context(), int32(roleID))
	if err != nil {
		return err
	}

	return c.JSON(permissions)
}

// @Summary Assign user to role
// @Description Assign a user to a role
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Role ID"
// @Param userId body int true "User ID"
// @Success 200 {object} string "User assigned to role successfully"
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /roles/{id}/users [post]
func (h *RoleHandler) AssignUser(c *fiber.Ctx) error {
	roleID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	var req AssignUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	err = h.rbacService.AssignRoleToUser(c.Context(), int32(roleID), req.UserID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

// @Summary Remove user from role
// @Description Remove a user from a role
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Role ID"
// @Param userId path int true "User ID"
// @Success 200 {object} string "User removed from role successfully"
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /roles/{id}/users/{userId} [delete]
func (h *RoleHandler) RemoveUser(c *fiber.Ctx) error {
	roleID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	userID, err := c.ParamsInt("userId")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	err = h.rbacService.RemoveUserFromRole(c.Context(), roleID, userID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary List role users
// @Description Get all users assigned to a role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {array} db.User
// @Router /roles/{id}/users [get]
func (h *RoleHandler) ListRoleUsers(c *fiber.Ctx) error {
	roleID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	users, err := h.rbacService.GetRoleUsers(c.Context(), int32(roleID))
	if err != nil {
		return err
	}

	return c.JSON(users)
}

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required" example:"admin"`
	Description string `json:"description" example:"Administrator role"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AssignPermissionRequest struct {
	PermissionID int32 `json:"permission_id"`
}

type AssignUserRequest struct {
	UserID int32 `json:"user_id"`
}

func toRoleResponse(role db.Role) models.RoleResponse {
	return models.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description.String,
	}
}
