package rbac

import (
	"backend-fiber/internal/db"
	"context"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

func (s *Service) CheckPermission(ctx context.Context, userID int32, requiredPermission string) (bool, error) {
	// Get user permissions
	permissions, err := s.queries.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	// Check if user has the required permission
	for _, p := range permissions {
		if p.Name == requiredPermission {
			return true, nil
		}
	}

	return false, nil
}

func (s *Service) AssignRoleToUser(ctx context.Context, userID int32, roleID int32) error {
	return s.queries.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: userID,
		RoleID: roleID,
	})
}

func (s *Service) GetUserRoles(ctx context.Context, userID int32) ([]db.Role, error) {
	return s.queries.GetUserRoles(ctx, userID)
}

func (s *Service) GetUserPermissions(ctx context.Context, userID int32) ([]db.Permission, error) {
	return s.queries.GetUserPermissions(ctx, userID)
}

func (s *Service) CreateRole(ctx context.Context, params db.CreateRoleParams) (db.Role, error) {
	return s.queries.CreateRole(ctx, params)
}

func (s *Service) GetRole(ctx context.Context, id int32) (db.Role, error) {
	return s.queries.GetRole(ctx, id)
}

func (s *Service) ListRoles(ctx context.Context) ([]db.Role, error) {
	return s.queries.ListRoles(ctx)
}

func (s *Service) UpdateRole(ctx context.Context, params db.UpdateRoleParams) error {
	return s.queries.UpdateRole(ctx, params)
}

func (s *Service) DeleteRole(ctx context.Context, id int32) error {
	return s.queries.DeleteRole(ctx, id)
}

func (s *Service) AssignPermissionToRole(ctx context.Context, params db.AssignPermissionToRoleParams) error {
	return s.queries.AssignPermissionToRole(ctx, params)
}

func (s *Service) RemovePermissionFromRole(ctx context.Context, roleID, permissionID int) error {
	return s.queries.RemovePermissionFromRole(ctx, db.RemovePermissionFromRoleParams{
		RoleID:       int32(roleID),
		PermissionID: int32(permissionID),
	})
}

func (s *Service) GetRolePermissions(ctx context.Context, roleID int32) ([]db.Permission, error) {
	return s.queries.GetRolePermissions(ctx, roleID)
}

func (s *Service) RemoveUserFromRole(ctx context.Context, roleID, userID int) error {
	return s.queries.RemoveUserFromRole(ctx, db.RemoveUserFromRoleParams{
		RoleID: int32(roleID),
		UserID: int32(userID),
	})
}

func (s *Service) GetRoleUsers(ctx context.Context, roleID int32) ([]db.User, error) {
	return s.queries.GetRoleUsers(ctx, roleID)
}
