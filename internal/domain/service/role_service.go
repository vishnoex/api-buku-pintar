package service

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// RoleService defines the interface for role business operations
// Clean Architecture: Service layer, orchestrates business logic and validation
type RoleService interface {
	// Role CRUD operations
	CreateRole(ctx context.Context, role *entity.Role) error
	GetRoleByID(ctx context.Context, id string) (*entity.Role, error)
	GetRoleByName(ctx context.Context, name string) (*entity.Role, error)
	UpdateRole(ctx context.Context, role *entity.Role) error
	DeleteRole(ctx context.Context, id string) error
	
	// Role list operations
	GetRoleList(ctx context.Context, limit, offset int) ([]*entity.Role, error)
	GetRoleCount(ctx context.Context) (int64, error)
	
	// Role-Permission operations
	GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*entity.Permission, error)
	AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error
	AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error
	RemoveAllPermissionsFromRole(ctx context.Context, roleID string) error
	
	// Role with complete information
	GetRoleWithPermissions(ctx context.Context, roleID string) (*entity.RoleWithPermissions, error)
	
	// User-Role operations
	GetUsersByRoleID(ctx context.Context, roleID string, limit, offset int) ([]*entity.User, error)
	CountUsersByRoleID(ctx context.Context, roleID string) (int64, error)
	AssignRoleToUser(ctx context.Context, userID, roleID string) error
	RemoveRoleFromUser(ctx context.Context, userID string) error
	
	// Validation
	ValidateRoleName(ctx context.Context, name string) error
	IsRoleNameUnique(ctx context.Context, name string, excludeID *string) (bool, error)
	CanDeleteRole(ctx context.Context, roleID string) error
}
