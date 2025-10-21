package redis

import (
	"buku-pintar/internal/domain/entity"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// PermissionRedisRepository implements permission caching operations
// Optimized for fast authorization checks with tiered TTL strategy
type PermissionRedisRepository struct {
	client *redis.Client
}

// NewPermissionRedisRepository creates a new permission Redis repository
func NewPermissionRedisRepository(client *redis.Client) *PermissionRedisRepository {
	return &PermissionRedisRepository{
		client: client,
	}
}

// Cache key constants
const (
	permissionByIDPrefix       = "permission:id:"
	permissionByNamePrefix     = "permission:name:"
	permissionListPrefix       = "permission:list:"
	permissionTotalKey         = "permission:total"
	permissionByResourcePrefix = "permission:resource:"
	userPermissionsPrefix      = "user:permissions:"
	userHasPermissionPrefix    = "user:has_permission:"
)

// Cache TTL constants - Tiered strategy
const (
	permissionTTL            = 24 * time.Hour  // Individual permissions rarely change
	permissionListTTL        = 1 * time.Hour   // Lists change more frequently
	userPermissionsTTL       = 30 * time.Minute // User permissions need faster invalidation
	permissionCheckTTL       = 15 * time.Minute // Boolean checks can be shorter
	resourcePermissionsTTL   = 6 * time.Hour    // Resource-based grouping is stable
)

// ========================================
// Permission Caching Operations
// ========================================

// GetPermissionByID retrieves a permission from cache by ID
func (r *PermissionRedisRepository) GetPermissionByID(ctx context.Context, id string) (*entity.Permission, error) {
	key := permissionByIDPrefix + id
	
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get permission from cache: %w", err)
	}
	
	var permission entity.Permission
	if err := json.Unmarshal([]byte(data), &permission); err != nil {
		return nil, fmt.Errorf("failed to unmarshal permission: %w", err)
	}
	
	return &permission, nil
}

// SetPermissionByID stores a permission in cache by ID
func (r *PermissionRedisRepository) SetPermissionByID(ctx context.Context, permission *entity.Permission) error {
	if permission == nil || permission.ID == "" {
		return fmt.Errorf("invalid permission: nil or empty ID")
	}
	
	key := permissionByIDPrefix + permission.ID
	
	data, err := json.Marshal(permission)
	if err != nil {
		return fmt.Errorf("failed to marshal permission: %w", err)
	}
	
	if err := r.client.Set(ctx, key, data, permissionTTL).Err(); err != nil {
		return fmt.Errorf("failed to set permission in cache: %w", err)
	}
	
	return nil
}

// GetPermissionByName retrieves a permission from cache by name
func (r *PermissionRedisRepository) GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error) {
	key := permissionByNamePrefix + name
	
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get permission from cache: %w", err)
	}
	
	var permission entity.Permission
	if err := json.Unmarshal([]byte(data), &permission); err != nil {
		return nil, fmt.Errorf("failed to unmarshal permission: %w", err)
	}
	
	return &permission, nil
}

// SetPermissionByName stores a permission in cache by name
func (r *PermissionRedisRepository) SetPermissionByName(ctx context.Context, permission *entity.Permission) error {
	if permission == nil || permission.Name == "" {
		return fmt.Errorf("invalid permission: nil or empty name")
	}
	
	key := permissionByNamePrefix + permission.Name
	
	data, err := json.Marshal(permission)
	if err != nil {
		return fmt.Errorf("failed to marshal permission: %w", err)
	}
	
	if err := r.client.Set(ctx, key, data, permissionTTL).Err(); err != nil {
		return fmt.Errorf("failed to set permission in cache: %w", err)
	}
	
	return nil
}

// ========================================
// Permission List Caching Operations
// ========================================

// GetPermissionList retrieves a permission list from cache
func (r *PermissionRedisRepository) GetPermissionList(ctx context.Context, limit, offset int) ([]*entity.Permission, error) {
	key := fmt.Sprintf("%s%d:%d", permissionListPrefix, limit, offset)
	
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get permission list from cache: %w", err)
	}
	
	var permissions []*entity.Permission
	if err := json.Unmarshal([]byte(data), &permissions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal permission list: %w", err)
	}
	
	return permissions, nil
}

// SetPermissionList stores a permission list in cache
func (r *PermissionRedisRepository) SetPermissionList(ctx context.Context, permissions []*entity.Permission, limit, offset int) error {
	if permissions == nil {
		return fmt.Errorf("invalid permissions: nil")
	}
	
	key := fmt.Sprintf("%s%d:%d", permissionListPrefix, limit, offset)
	
	data, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal permission list: %w", err)
	}
	
	if err := r.client.Set(ctx, key, data, permissionListTTL).Err(); err != nil {
		return fmt.Errorf("failed to set permission list in cache: %w", err)
	}
	
	return nil
}

// GetPermissionTotal retrieves the total permission count from cache
func (r *PermissionRedisRepository) GetPermissionTotal(ctx context.Context) (int64, error) {
	count, err := r.client.Get(ctx, permissionTotalKey).Int64()
	if err == redis.Nil {
		return 0, nil // Cache miss
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get permission total from cache: %w", err)
	}
	
	return count, nil
}

// SetPermissionTotal stores the total permission count in cache
func (r *PermissionRedisRepository) SetPermissionTotal(ctx context.Context, count int64) error {
	if err := r.client.Set(ctx, permissionTotalKey, count, permissionListTTL).Err(); err != nil {
		return fmt.Errorf("failed to set permission total in cache: %w", err)
	}
	
	return nil
}

// ========================================
// Permission by Resource Caching Operations
// ========================================

// GetPermissionsByResource retrieves permissions by resource from cache
func (r *PermissionRedisRepository) GetPermissionsByResource(ctx context.Context, resource string) ([]*entity.Permission, error) {
	key := permissionByResourcePrefix + resource
	
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions by resource from cache: %w", err)
	}
	
	var permissions []*entity.Permission
	if err := json.Unmarshal([]byte(data), &permissions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal permissions: %w", err)
	}
	
	return permissions, nil
}

// SetPermissionsByResource stores permissions by resource in cache
func (r *PermissionRedisRepository) SetPermissionsByResource(ctx context.Context, resource string, permissions []*entity.Permission) error {
	if permissions == nil {
		return fmt.Errorf("invalid permissions: nil")
	}
	
	key := permissionByResourcePrefix + resource
	
	data, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal permissions: %w", err)
	}
	
	if err := r.client.Set(ctx, key, data, resourcePermissionsTTL).Err(); err != nil {
		return fmt.Errorf("failed to set permissions by resource in cache: %w", err)
	}
	
	return nil
}

// ========================================
// User Permission Caching Operations
// Critical for fast authorization checks
// ========================================

// GetUserPermissions retrieves user permissions from cache
func (r *PermissionRedisRepository) GetUserPermissions(ctx context.Context, userID string) ([]*entity.Permission, error) {
	key := userPermissionsPrefix + userID
	
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions from cache: %w", err)
	}
	
	var permissions []*entity.Permission
	if err := json.Unmarshal([]byte(data), &permissions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user permissions: %w", err)
	}
	
	return permissions, nil
}

// SetUserPermissions stores user permissions in cache
func (r *PermissionRedisRepository) SetUserPermissions(ctx context.Context, userID string, permissions []*entity.Permission) error {
	if permissions == nil {
		return fmt.Errorf("invalid permissions: nil")
	}
	
	key := userPermissionsPrefix + userID
	
	data, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal user permissions: %w", err)
	}
	
	if err := r.client.Set(ctx, key, data, userPermissionsTTL).Err(); err != nil {
		return fmt.Errorf("failed to set user permissions in cache: %w", err)
	}
	
	return nil
}

// ========================================
// Permission Check Caching Operations
// Optimized for boolean authorization checks
// ========================================

// GetUserHasPermission retrieves a permission check result from cache
func (r *PermissionRedisRepository) GetUserHasPermission(ctx context.Context, userID, permissionName string) (*bool, error) {
	key := fmt.Sprintf("%s%s:%s", userHasPermissionPrefix, userID, permissionName)
	
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get permission check from cache: %w", err)
	}
	
	hasPermission := val == "1"
	return &hasPermission, nil
}

// SetUserHasPermission stores a permission check result in cache
func (r *PermissionRedisRepository) SetUserHasPermission(ctx context.Context, userID, permissionName string, hasPermission bool) error {
	key := fmt.Sprintf("%s%s:%s", userHasPermissionPrefix, userID, permissionName)
	
	val := "0"
	if hasPermission {
		val = "1"
	}
	
	if err := r.client.Set(ctx, key, val, permissionCheckTTL).Err(); err != nil {
		return fmt.Errorf("failed to set permission check in cache: %w", err)
	}
	
	return nil
}

// ========================================
// Cache Invalidation Operations
// Surgical invalidation for data consistency
// ========================================

// InvalidatePermissionCache invalidates all permission-related caches
// Use sparingly - only when permissions are bulk updated
func (r *PermissionRedisRepository) InvalidatePermissionCache(ctx context.Context) error {
	patterns := []string{
		permissionByIDPrefix + "*",
		permissionByNamePrefix + "*",
		permissionListPrefix + "*",
		permissionTotalKey,
		permissionByResourcePrefix + "*",
		userPermissionsPrefix + "*",
		userHasPermissionPrefix + "*",
	}
	
	for _, pattern := range patterns {
		keys, err := r.client.Keys(ctx, pattern).Result()
		if err != nil {
			return fmt.Errorf("failed to get keys for pattern %s: %w", pattern, err)
		}
		
		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("failed to delete keys for pattern %s: %w", pattern, err)
			}
		}
	}
	
	return nil
}

// InvalidatePermissionCacheByID invalidates cache for a specific permission
// Use when a single permission is updated/deleted
func (r *PermissionRedisRepository) InvalidatePermissionCacheByID(ctx context.Context, permissionID string) error {
	// Get permission name first to invalidate name-based cache
	permission, err := r.GetPermissionByID(ctx, permissionID)
	if err != nil {
		return fmt.Errorf("failed to get permission for invalidation: %w", err)
	}
	
	keys := []string{
		permissionByIDPrefix + permissionID,
	}
	
	if permission != nil && permission.Name != "" {
		keys = append(keys, permissionByNamePrefix+permission.Name)
		
		// Also invalidate resource-based cache
		if permission.Resource != "" {
			resourceKeys, err := r.client.Keys(ctx, permissionByResourcePrefix+permission.Resource).Result()
			if err == nil && len(resourceKeys) > 0 {
				keys = append(keys, resourceKeys...)
			}
		}
	}
	
	// Invalidate permission lists
	listKeys, err := r.client.Keys(ctx, permissionListPrefix+"*").Result()
	if err == nil && len(listKeys) > 0 {
		keys = append(keys, listKeys...)
	}
	
	// Invalidate total count
	keys = append(keys, permissionTotalKey)
	
	// Invalidate all user permission caches (affected by this permission change)
	userPermKeys, err := r.client.Keys(ctx, userPermissionsPrefix+"*").Result()
	if err == nil && len(userPermKeys) > 0 {
		keys = append(keys, userPermKeys...)
	}
	
	// Invalidate permission checks
	checkKeys, err := r.client.Keys(ctx, userHasPermissionPrefix+"*").Result()
	if err == nil && len(checkKeys) > 0 {
		keys = append(keys, checkKeys...)
	}
	
	if len(keys) > 0 {
		if err := r.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to invalidate permission cache: %w", err)
		}
	}
	
	return nil
}

// InvalidateUserPermissionsCache invalidates permission cache for a specific user
// Use when user's role changes or role-permission assignments change
func (r *PermissionRedisRepository) InvalidateUserPermissionsCache(ctx context.Context, userID string) error {
	patterns := []string{
		userPermissionsPrefix + userID,
		userHasPermissionPrefix + userID + ":*",
	}
	
	var keysToDelete []string
	for _, pattern := range patterns {
		keys, err := r.client.Keys(ctx, pattern).Result()
		if err != nil {
			return fmt.Errorf("failed to get keys for pattern %s: %w", pattern, err)
		}
		keysToDelete = append(keysToDelete, keys...)
	}
	
	if len(keysToDelete) > 0 {
		if err := r.client.Del(ctx, keysToDelete...).Err(); err != nil {
			return fmt.Errorf("failed to invalidate user permissions cache: %w", err)
		}
	}
	
	return nil
}

// InvalidateUserPermissionCheck invalidates a specific permission check for a user
// Use for granular invalidation when you know the exact permission that changed
func (r *PermissionRedisRepository) InvalidateUserPermissionCheck(ctx context.Context, userID, permissionName string) error {
	key := fmt.Sprintf("%s%s:%s", userHasPermissionPrefix, userID, permissionName)
	
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to invalidate user permission check: %w", err)
	}
	
	return nil
}
