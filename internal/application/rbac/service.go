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
