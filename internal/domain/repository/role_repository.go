package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// RoleRepository defines the interface for role data operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type RoleRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, role *entity.Role) error
	GetByID(ctx context.Context, id string) (*entity.Role, error)
	GetByName(ctx context.Context, name string) (*entity.Role, error)
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id string) error
	
	// List operations
	List(ctx context.Context, limit, offset int) ([]*entity.Role, error)
	Count(ctx context.Context) (int64, error)
	
	// Role-Permission operations
	GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*entity.Permission, error)
	AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error
	
	// Bulk operations
	AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error
	RemoveAllPermissionsFromRole(ctx context.Context, roleID string) error
	
	// User-Role operations
	GetUsersByRoleID(ctx context.Context, roleID string, limit, offset int) ([]*entity.User, error)
	CountUsersByRoleID(ctx context.Context, roleID string) (int64, error)
}

// RoleRedisRepository defines the interface for role Redis operations
// Used for caching to improve performance of RBAC checks
type RoleRedisRepository interface {
	// Role caching
	GetRoleByID(ctx context.Context, id string) (*entity.Role, error)
	SetRoleByID(ctx context.Context, role *entity.Role) error
	GetRoleByName(ctx context.Context, name string) (*entity.Role, error)
	SetRoleByName(ctx context.Context, role *entity.Role) error
	
	// Role list caching
	GetRoleList(ctx context.Context, limit, offset int) ([]*entity.Role, error)
	SetRoleList(ctx context.Context, roles []*entity.Role, limit, offset int) error
	GetRoleTotal(ctx context.Context) (int64, error)
	SetRoleTotal(ctx context.Context, count int64) error
	
	// Permission caching by role
	GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*entity.Permission, error)
	SetPermissionsByRoleID(ctx context.Context, roleID string, permissions []*entity.Permission) error
	
	// User's role and permissions caching (for fast authorization checks)
	GetUserPermissions(ctx context.Context, userID string) ([]*entity.Permission, error)
	SetUserPermissions(ctx context.Context, userID string, permissions []*entity.Permission) error
	
	// Cache invalidation
	InvalidateRoleCache(ctx context.Context) error
	InvalidateRoleCacheByID(ctx context.Context, roleID string) error
	InvalidateUserPermissionsCache(ctx context.Context, userID string) error
}
