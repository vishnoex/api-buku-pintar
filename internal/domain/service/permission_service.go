package service

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// PermissionService defines the interface for permission business operations
// Clean Architecture: Service layer, orchestrates business logic and authorization checks
type PermissionService interface {
	// Permission CRUD operations
	CreatePermission(ctx context.Context, permission *entity.Permission) error
	GetPermissionByID(ctx context.Context, id string) (*entity.Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error)
	UpdatePermission(ctx context.Context, permission *entity.Permission) error
	DeletePermission(ctx context.Context, id string) error
	
	// Permission list operations
	GetPermissionList(ctx context.Context, limit, offset int) ([]*entity.Permission, error)
	GetPermissionCount(ctx context.Context) (int64, error)
	
	// Permission search and filter
	GetPermissionsByResource(ctx context.Context, resource string, limit, offset int) ([]*entity.Permission, error)
	GetPermissionsByAction(ctx context.Context, action string, limit, offset int) ([]*entity.Permission, error)
	GetPermissionsByResourceAndAction(ctx context.Context, resource, action string) ([]*entity.Permission, error)
	
	// Permission-Role operations
	GetRolesByPermissionID(ctx context.Context, permissionID string, limit, offset int) ([]*entity.Role, error)
	CountRolesByPermissionID(ctx context.Context, permissionID string) (int64, error)
	GetPermissionWithRoles(ctx context.Context, permissionID string) (*entity.PermissionWithRoles, error)
	
	// User permission operations (Critical for authorization)
	GetUserPermissions(ctx context.Context, userID string) ([]*entity.Permission, error)
	GetUserPermissionNames(ctx context.Context, userID string) ([]string, error)
	HasPermission(ctx context.Context, userID, permissionName string) (bool, error)
	HasPermissions(ctx context.Context, userID string, permissionNames []string) (bool, error)
	HasAnyPermission(ctx context.Context, userID string, permissionNames []string) (bool, error)
	
	// Permission by resource operations
	GetUserPermissionsForResource(ctx context.Context, userID, resource string) ([]*entity.Permission, error)
	CanUserPerformAction(ctx context.Context, userID, resource, action string) (bool, error)
	
	// Bulk operations
	CreatePermissionsBulk(ctx context.Context, permissions []*entity.Permission) error
	GetPermissionsByNames(ctx context.Context, names []string) ([]*entity.Permission, error)
	
	// Validation
	ValidatePermissionName(ctx context.Context, name string) error
	IsPermissionNameUnique(ctx context.Context, name string, excludeID *string) (bool, error)
	CanDeletePermission(ctx context.Context, permissionID string) error
	
	// Permission management helpers
	GeneratePermissionName(resource, action string) string
	ParsePermissionName(permissionName string) (resource, action string, err error)
}
