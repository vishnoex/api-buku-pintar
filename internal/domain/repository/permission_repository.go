package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// PermissionRepository defines the interface for permission data operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type PermissionRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, permission *entity.Permission) error
	GetByID(ctx context.Context, id string) (*entity.Permission, error)
	GetByName(ctx context.Context, name string) (*entity.Permission, error)
	Update(ctx context.Context, permission *entity.Permission) error
	Delete(ctx context.Context, id string) error
	
	// List operations
	List(ctx context.Context, limit, offset int) ([]*entity.Permission, error)
	Count(ctx context.Context) (int64, error)
	
	// Search and filter operations
	ListByResource(ctx context.Context, resource string, limit, offset int) ([]*entity.Permission, error)
	ListByAction(ctx context.Context, action string, limit, offset int) ([]*entity.Permission, error)
	ListByResourceAndAction(ctx context.Context, resource, action string, limit, offset int) ([]*entity.Permission, error)
	
	// Permission-Role operations
	GetRolesByPermissionID(ctx context.Context, permissionID string, limit, offset int) ([]*entity.Role, error)
	CountRolesByPermissionID(ctx context.Context, permissionID string) (int64, error)
	
	// User permission check operations
	GetPermissionsByUserID(ctx context.Context, userID string) ([]*entity.Permission, error)
	HasPermission(ctx context.Context, userID, permissionName string) (bool, error)
	HasPermissions(ctx context.Context, userID string, permissionNames []string) (bool, error)
	
	// Bulk operations
	CreateBulk(ctx context.Context, permissions []*entity.Permission) error
	GetByNames(ctx context.Context, names []string) ([]*entity.Permission, error)
}

// PermissionRedisRepository defines the interface for permission Redis operations
// Used for caching to improve performance of authorization checks
type PermissionRedisRepository interface {
	// Permission caching
	GetPermissionByID(ctx context.Context, id string) (*entity.Permission, error)
	SetPermissionByID(ctx context.Context, permission *entity.Permission) error
	GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error)
	SetPermissionByName(ctx context.Context, permission *entity.Permission) error
	
	// Permission list caching
	GetPermissionList(ctx context.Context, limit, offset int) ([]*entity.Permission, error)
	SetPermissionList(ctx context.Context, permissions []*entity.Permission, limit, offset int) error
	GetPermissionTotal(ctx context.Context) (int64, error)
	SetPermissionTotal(ctx context.Context, count int64) error
	
	// Permission by resource/action caching
	GetPermissionsByResource(ctx context.Context, resource string) ([]*entity.Permission, error)
	SetPermissionsByResource(ctx context.Context, resource string, permissions []*entity.Permission) error
	
	// User permission caching (critical for fast authorization)
	GetUserPermissions(ctx context.Context, userID string) ([]*entity.Permission, error)
	SetUserPermissions(ctx context.Context, userID string, permissions []*entity.Permission) error
	
	// Permission check caching (boolean results)
	GetUserHasPermission(ctx context.Context, userID, permissionName string) (*bool, error)
	SetUserHasPermission(ctx context.Context, userID, permissionName string, hasPermission bool) error
	
	// Cache invalidation
	InvalidatePermissionCache(ctx context.Context) error
	InvalidatePermissionCacheByID(ctx context.Context, permissionID string) error
	InvalidateUserPermissionsCache(ctx context.Context, userID string) error
	InvalidateUserPermissionCheck(ctx context.Context, userID, permissionName string) error
}
