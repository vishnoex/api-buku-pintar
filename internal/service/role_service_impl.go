package service

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/internal/domain/service"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type roleService struct {
	roleRepo      repository.RoleRepository
	roleRedisRepo repository.RoleRedisRepository
	userRepo      repository.UserRepository
	cacheTTL      time.Duration
}

// NewRoleService creates a new instance of RoleService
func NewRoleService(
	roleRepo repository.RoleRepository,
	roleRedisRepo repository.RoleRedisRepository,
	userRepo repository.UserRepository,
) service.RoleService {
	return &roleService{
		roleRepo:      roleRepo,
		roleRedisRepo: roleRedisRepo,
		userRepo:      userRepo,
		cacheTTL:      30 * time.Minute, // 30 minutes cache TTL for roles
	}
}

// CreateRole creates a new role with validation
func (s *roleService) CreateRole(ctx context.Context, role *entity.Role) error {
	// Validate role name
	if err := s.ValidateRoleName(ctx, role.Name); err != nil {
		return err
	}

	// Check if role name is unique
	isUnique, err := s.IsRoleNameUnique(ctx, role.Name, nil)
	if err != nil {
		return fmt.Errorf("failed to check role name uniqueness: %w", err)
	}
	if !isUnique {
		return errors.New("role name already exists")
	}

	// Generate ID if not provided
	if role.ID == "" {
		role.ID = uuid.New().String()
	}

	// Create role in database
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	// Cache the new role
	if err := s.roleRedisRepo.SetRoleByID(ctx, role); err != nil {
		log.Printf("Failed to cache new role: %v", err)
	}
	if err := s.roleRedisRepo.SetRoleByName(ctx, role); err != nil {
		log.Printf("Failed to cache new role by name: %v", err)
	}

	// Invalidate list caches
	if err := s.roleRedisRepo.InvalidateRoleCache(ctx); err != nil {
		log.Printf("Failed to invalidate role cache: %v", err)
	}

	return nil
}

// GetRoleByID retrieves a role by ID with caching
func (s *roleService) GetRoleByID(ctx context.Context, id string) (*entity.Role, error) {
	if id == "" {
		return nil, errors.New("role ID cannot be empty")
	}

	// Try to get from cache first
	cachedRole, err := s.roleRedisRepo.GetRoleByID(ctx, id)
	if err == nil && cachedRole != nil {
		log.Printf("Role %s retrieved from cache", id)
		return cachedRole, nil
	}

	// If not in cache, get from database
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	// Cache the result
	if err := s.roleRedisRepo.SetRoleByID(ctx, role); err != nil {
		log.Printf("Failed to cache role: %v", err)
	}

	return role, nil
}

// GetRoleByName retrieves a role by name with caching
func (s *roleService) GetRoleByName(ctx context.Context, name string) (*entity.Role, error) {
	if name == "" {
		return nil, errors.New("role name cannot be empty")
	}

	// Try to get from cache first
	cachedRole, err := s.roleRedisRepo.GetRoleByName(ctx, name)
	if err == nil && cachedRole != nil {
		log.Printf("Role %s retrieved from cache", name)
		return cachedRole, nil
	}

	// If not in cache, get from database
	role, err := s.roleRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	// Cache the result
	if err := s.roleRedisRepo.SetRoleByName(ctx, role); err != nil {
		log.Printf("Failed to cache role: %v", err)
	}

	return role, nil
}

// UpdateRole updates a role with validation
func (s *roleService) UpdateRole(ctx context.Context, role *entity.Role) error {
	if role.ID == "" {
		return errors.New("role ID cannot be empty")
	}

	// Validate role name
	if err := s.ValidateRoleName(ctx, role.Name); err != nil {
		return err
	}

	// Check if role exists
	existingRole, err := s.roleRepo.GetByID(ctx, role.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing role: %w", err)
	}
	if existingRole == nil {
		return errors.New("role not found")
	}

	// Check if new name is unique (excluding current role)
	isUnique, err := s.IsRoleNameUnique(ctx, role.Name, &role.ID)
	if err != nil {
		return fmt.Errorf("failed to check role name uniqueness: %w", err)
	}
	if !isUnique {
		return errors.New("role name already exists")
	}

	// Update role in database
	if err := s.roleRepo.Update(ctx, role); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	// Invalidate cache for this role
	if err := s.roleRedisRepo.InvalidateRoleCacheByID(ctx, role.ID); err != nil {
		log.Printf("Failed to invalidate role cache: %v", err)
	}

	// Invalidate old name cache if name changed
	if existingRole.Name != role.Name {
		if err := s.roleRedisRepo.InvalidateRoleCache(ctx); err != nil {
			log.Printf("Failed to invalidate role cache: %v", err)
		}
	}

	return nil
}

// DeleteRole deletes a role after validation
func (s *roleService) DeleteRole(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("role ID cannot be empty")
	}

	// Check if role can be deleted
	if err := s.CanDeleteRole(ctx, id); err != nil {
		return err
	}

	// Get role before deletion for cache invalidation
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Delete role from database
	if err := s.roleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	// Invalidate cache
	if err := s.roleRedisRepo.InvalidateRoleCacheByID(ctx, id); err != nil {
		log.Printf("Failed to invalidate role cache: %v", err)
	}

	return nil
}

// GetRoleList retrieves all roles with caching
func (s *roleService) GetRoleList(ctx context.Context, limit, offset int) ([]*entity.Role, error) {
	// Try to get from cache first
	cachedRoles, err := s.roleRedisRepo.GetRoleList(ctx, limit, offset)
	if err == nil && cachedRoles != nil && len(cachedRoles) > 0 {
		log.Println("Role list retrieved from cache")
		return cachedRoles, nil
	}

	// If not in cache, get from database
	roles, err := s.roleRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get role list: %w", err)
	}

	// Cache the result
	if len(roles) > 0 {
		if err := s.roleRedisRepo.SetRoleList(ctx, roles, limit, offset); err != nil {
			log.Printf("Failed to cache role list: %v", err)
		}
	}

	return roles, nil
}

// GetRoleCount retrieves the total number of roles with caching
func (s *roleService) GetRoleCount(ctx context.Context) (int64, error) {
	// Try to get from cache first
	cachedCount, err := s.roleRedisRepo.GetRoleTotal(ctx)
	if err == nil && cachedCount > 0 {
		log.Println("Role count retrieved from cache")
		return cachedCount, nil
	}

	// If not in cache, get from database
	count, err := s.roleRepo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get role count: %w", err)
	}

	// Cache the result
	if err := s.roleRedisRepo.SetRoleTotal(ctx, count); err != nil {
		log.Printf("Failed to cache role count: %v", err)
	}

	return count, nil
}

// GetPermissionsByRoleID retrieves all permissions for a role with caching
func (s *roleService) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*entity.Permission, error) {
	if roleID == "" {
		return nil, errors.New("role ID cannot be empty")
	}

	// Try to get from cache first
	cachedPermissions, err := s.roleRedisRepo.GetPermissionsByRoleID(ctx, roleID)
	if err == nil && cachedPermissions != nil {
		log.Printf("Permissions for role %s retrieved from cache", roleID)
		return cachedPermissions, nil
	}

	// If not in cache, get from database
	permissions, err := s.roleRepo.GetPermissionsByRoleID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	// Cache the result
	if err := s.roleRedisRepo.SetPermissionsByRoleID(ctx, roleID, permissions); err != nil {
		log.Printf("Failed to cache permissions: %v", err)
	}

	return permissions, nil
}

// AssignPermissionToRole assigns a permission to a role
func (s *roleService) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	if roleID == "" || permissionID == "" {
		return errors.New("role ID and permission ID cannot be empty")
	}

	// Verify role exists
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Assign permission
	if err := s.roleRepo.AssignPermissionToRole(ctx, roleID, permissionID); err != nil {
		return fmt.Errorf("failed to assign permission: %w", err)
	}

	// Invalidate permission cache for this role
	if err := s.roleRedisRepo.InvalidateRoleCacheByID(ctx, roleID); err != nil {
		log.Printf("Failed to invalidate role cache: %v", err)
	}

	return nil
}

// RemovePermissionFromRole removes a permission from a role
func (s *roleService) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	if roleID == "" || permissionID == "" {
		return errors.New("role ID and permission ID cannot be empty")
	}

	// Remove permission
	if err := s.roleRepo.RemovePermissionFromRole(ctx, roleID, permissionID); err != nil {
		return fmt.Errorf("failed to remove permission: %w", err)
	}

	// Invalidate permission cache for this role
	if err := s.roleRedisRepo.InvalidateRoleCacheByID(ctx, roleID); err != nil {
		log.Printf("Failed to invalidate role cache: %v", err)
	}

	return nil
}

// AssignPermissionsToRole assigns multiple permissions to a role
func (s *roleService) AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error {
	if roleID == "" {
		return errors.New("role ID cannot be empty")
	}
	if len(permissionIDs) == 0 {
		return errors.New("permission IDs cannot be empty")
	}

	// Verify role exists
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Assign permissions in bulk
	if err := s.roleRepo.AssignPermissionsToRole(ctx, roleID, permissionIDs); err != nil {
		return fmt.Errorf("failed to assign permissions: %w", err)
	}

	// Invalidate permission cache for this role
	if err := s.roleRedisRepo.InvalidateRoleCacheByID(ctx, roleID); err != nil {
		log.Printf("Failed to invalidate role cache: %v", err)
	}

	return nil
}

// RemoveAllPermissionsFromRole removes all permissions from a role
func (s *roleService) RemoveAllPermissionsFromRole(ctx context.Context, roleID string) error {
	if roleID == "" {
		return errors.New("role ID cannot be empty")
	}

	// Remove all permissions
	if err := s.roleRepo.RemoveAllPermissionsFromRole(ctx, roleID); err != nil {
		return fmt.Errorf("failed to remove permissions: %w", err)
	}

	// Invalidate permission cache for this role
	if err := s.roleRedisRepo.InvalidateRoleCacheByID(ctx, roleID); err != nil {
		log.Printf("Failed to invalidate role cache: %v", err)
	}

	return nil
}

// GetRoleWithPermissions retrieves a role with all its permissions
func (s *roleService) GetRoleWithPermissions(ctx context.Context, roleID string) (*entity.RoleWithPermissions, error) {
	if roleID == "" {
		return nil, errors.New("role ID cannot be empty")
	}

	// Get role
	role, err := s.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// Get permissions
	permissions, err := s.GetPermissionsByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// Convert []*Permission to []Permission
	permissionValues := make([]entity.Permission, len(permissions))
	for i, p := range permissions {
		permissionValues[i] = *p
	}

	return &entity.RoleWithPermissions{
		Role:        *role,
		Permissions: permissionValues,
	}, nil
}

// GetUsersByRoleID retrieves all users with a specific role
func (s *roleService) GetUsersByRoleID(ctx context.Context, roleID string, limit, offset int) ([]*entity.User, error) {
	if roleID == "" {
		return nil, errors.New("role ID cannot be empty")
	}

	users, err := s.roleRepo.GetUsersByRoleID(ctx, roleID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

// CountUsersByRoleID counts the number of users with a specific role
func (s *roleService) CountUsersByRoleID(ctx context.Context, roleID string) (int64, error) {
	if roleID == "" {
		return 0, errors.New("role ID cannot be empty")
	}

	count, err := s.roleRepo.CountUsersByRoleID(ctx, roleID)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// AssignRoleToUser assigns a role to a user
func (s *roleService) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	if userID == "" || roleID == "" {
		return errors.New("user ID and role ID cannot be empty")
	}

	// Verify role exists
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Update user's role
	user.RoleID = &roleID
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	// Invalidate user permission cache
	if err := s.roleRedisRepo.InvalidateUserPermissionsCache(ctx, userID); err != nil {
		log.Printf("Failed to invalidate user permissions cache: %v", err)
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (s *roleService) RemoveRoleFromUser(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user ID cannot be empty")
	}

	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Remove role from user
	user.RoleID = nil
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	// Invalidate user permission cache
	if err := s.roleRedisRepo.InvalidateUserPermissionsCache(ctx, userID); err != nil {
		log.Printf("Failed to invalidate user permissions cache: %v", err)
	}

	return nil
}

// ValidateRoleName validates the role name format
func (s *roleService) ValidateRoleName(ctx context.Context, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("role name cannot be empty")
	}
	if len(name) < 3 {
		return errors.New("role name must be at least 3 characters")
	}
	if len(name) > 50 {
		return errors.New("role name must be at most 50 characters")
	}
	// Check for valid characters (alphanumeric, spaces, hyphens, underscores)
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == ' ' || char == '-' || char == '_') {
			return errors.New("role name contains invalid characters")
		}
	}
	return nil
}

// IsRoleNameUnique checks if a role name is unique
func (s *roleService) IsRoleNameUnique(ctx context.Context, name string, excludeID *string) (bool, error) {
	existingRole, err := s.roleRepo.GetByName(ctx, name)
	if err != nil {
		return false, fmt.Errorf("failed to check role name: %w", err)
	}

	// If no role with this name exists, it's unique
	if existingRole == nil {
		return true, nil
	}

	// If we're excluding a specific ID (for updates), check if it's the same role
	if excludeID != nil && existingRole.ID == *excludeID {
		return true, nil
	}

	return false, nil
}

// CanDeleteRole checks if a role can be deleted
func (s *roleService) CanDeleteRole(ctx context.Context, roleID string) error {
	// Check if role has any users
	userCount, err := s.roleRepo.CountUsersByRoleID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if userCount > 0 {
		return fmt.Errorf("cannot delete role: %d user(s) are assigned to this role", userCount)
	}

	return nil
}
