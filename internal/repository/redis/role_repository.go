package redis

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type roleRedisRepository struct {
	client *redis.Client
}

// NewRoleRedisRepository creates a new instance of RoleRedisRepository
func NewRoleRedisRepository(client *redis.Client) repository.RoleRedisRepository {
	return &roleRedisRepository{
		client: client,
	}
}

// GetRoleByID retrieves a role from cache by ID
func (r *roleRedisRepository) GetRoleByID(ctx context.Context, id string) (*entity.Role, error) {
	key := fmt.Sprintf("role:id:%s", id)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var role entity.Role
	err = json.Unmarshal([]byte(data), &role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// SetRoleByID stores a role in cache by ID
func (r *roleRedisRepository) SetRoleByID(ctx context.Context, role *entity.Role) error {
	key := fmt.Sprintf("role:id:%s", role.ID)

	data, err := json.Marshal(role)
	if err != nil {
		return err
	}

	// Cache for 30 minutes (roles don't change often)
	return r.client.Set(ctx, key, data, 30*time.Minute).Err()
}

// GetRoleByName retrieves a role from cache by name
func (r *roleRedisRepository) GetRoleByName(ctx context.Context, name string) (*entity.Role, error) {
	key := fmt.Sprintf("role:name:%s", name)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var role entity.Role
	err = json.Unmarshal([]byte(data), &role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// SetRoleByName stores a role in cache by name
func (r *roleRedisRepository) SetRoleByName(ctx context.Context, role *entity.Role) error {
	key := fmt.Sprintf("role:name:%s", role.Name)

	data, err := json.Marshal(role)
	if err != nil {
		return err
	}

	// Cache for 30 minutes
	return r.client.Set(ctx, key, data, 30*time.Minute).Err()
}

// GetRoleList retrieves a list of roles from cache
func (r *roleRedisRepository) GetRoleList(ctx context.Context, limit, offset int) ([]*entity.Role, error) {
	key := fmt.Sprintf("role:list:%d:%d", limit, offset)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var roles []*entity.Role
	err = json.Unmarshal([]byte(data), &roles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// SetRoleList stores a list of roles in cache
func (r *roleRedisRepository) SetRoleList(ctx context.Context, roles []*entity.Role, limit, offset int) error {
	key := fmt.Sprintf("role:list:%d:%d", limit, offset)

	data, err := json.Marshal(roles)
	if err != nil {
		return err
	}

	// Cache for 15 minutes
	return r.client.Set(ctx, key, data, 15*time.Minute).Err()
}

// GetRoleTotal retrieves the total count of roles from cache
func (r *roleRedisRepository) GetRoleTotal(ctx context.Context) (int64, error) {
	key := "role:count:total"

	count, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Cache miss
		}
		return 0, err
	}

	return count, nil
}

// SetRoleTotal stores the total count of roles in cache
func (r *roleRedisRepository) SetRoleTotal(ctx context.Context, count int64) error {
	key := "role:count:total"

	// Cache for 15 minutes
	return r.client.Set(ctx, key, count, 15*time.Minute).Err()
}

// GetPermissionsByRoleID retrieves permissions for a role from cache
func (r *roleRedisRepository) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*entity.Permission, error) {
	key := fmt.Sprintf("role:permissions:%s", roleID)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var permissions []*entity.Permission
	err = json.Unmarshal([]byte(data), &permissions)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// SetPermissionsByRoleID stores permissions for a role in cache
func (r *roleRedisRepository) SetPermissionsByRoleID(ctx context.Context, roleID string, permissions []*entity.Permission) error {
	key := fmt.Sprintf("role:permissions:%s", roleID)

	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	// Cache for 1 hour (permissions are critical for authorization)
	return r.client.Set(ctx, key, data, 1*time.Hour).Err()
}

// GetUserPermissions retrieves all permissions for a user from cache
// This is the most critical cache for authorization checks
func (r *roleRedisRepository) GetUserPermissions(ctx context.Context, userID string) ([]*entity.Permission, error) {
	key := fmt.Sprintf("user:permissions:%s", userID)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var permissions []*entity.Permission
	err = json.Unmarshal([]byte(data), &permissions)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// SetUserPermissions stores all permissions for a user in cache
// This cache is critical for fast authorization checks on every request
func (r *roleRedisRepository) SetUserPermissions(ctx context.Context, userID string, permissions []*entity.Permission) error {
	key := fmt.Sprintf("user:permissions:%s", userID)

	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	// Cache for 1 hour (critical for performance, invalidate when role changes)
	return r.client.Set(ctx, key, data, 1*time.Hour).Err()
}

// InvalidateRoleCache invalidates all role-related cache
func (r *roleRedisRepository) InvalidateRoleCache(ctx context.Context) error {
	// Get all keys matching role pattern
	pattern := "role:*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	// Delete all role-related cache
	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}

// InvalidateRoleCacheByID invalidates cache for a specific role
func (r *roleRedisRepository) InvalidateRoleCacheByID(ctx context.Context, roleID string) error {
	// Keys to invalidate
	keys := []string{
		fmt.Sprintf("role:id:%s", roleID),
		fmt.Sprintf("role:permissions:%s", roleID),
	}

	// Also invalidate list caches (since role data might have changed)
	listPattern := "role:list:*"
	listKeys, err := r.client.Keys(ctx, listPattern).Result()
	if err != nil {
		return err
	}
	keys = append(keys, listKeys...)

	// Invalidate count cache
	keys = append(keys, "role:count:total")

	// Delete all identified keys
	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}

// InvalidateUserPermissionsCache invalidates permission cache for a specific user
// Call this when a user's role changes
func (r *roleRedisRepository) InvalidateUserPermissionsCache(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user:permissions:%s", userID)
	return r.client.Del(ctx, key).Err()
}
