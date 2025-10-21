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

type permissionService struct {
	permissionRepo      repository.PermissionRepository
	permissionRedisRepo repository.PermissionRedisRepository
	cacheTTL            time.Duration
}

// NewPermissionService creates a new instance of PermissionService
func NewPermissionService(
	permissionRepo repository.PermissionRepository,
	permissionRedisRepo repository.PermissionRedisRepository,
) service.PermissionService {
	return &permissionService{
		permissionRepo:      permissionRepo,
		permissionRedisRepo: permissionRedisRepo,
		cacheTTL:            60 * time.Minute, // 60 minutes cache TTL for permissions (longer for auth checks)
	}
}

// CreatePermission creates a new permission with validation
func (s *permissionService) CreatePermission(ctx context.Context, permission *entity.Permission) error {
	// Validate permission name
	if err := s.ValidatePermissionName(ctx, permission.Name); err != nil {
		return err
	}

	// Check if permission name is unique
	isUnique, err := s.IsPermissionNameUnique(ctx, permission.Name, nil)
	if err != nil {
		return fmt.Errorf("failed to check permission name uniqueness: %w", err)
	}
	if !isUnique {
		return errors.New("permission name already exists")
	}

	// Generate ID if not provided
	if permission.ID == "" {
		permission.ID = uuid.New().String()
	}

	// Create permission in database
	if err := s.permissionRepo.Create(ctx, permission); err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}

	// Cache the new permission
	if err := s.permissionRedisRepo.SetPermissionByID(ctx, permission); err != nil {
		log.Printf("Failed to cache new permission: %v", err)
	}
	if err := s.permissionRedisRepo.SetPermissionByName(ctx, permission); err != nil {
		log.Printf("Failed to cache new permission by name: %v", err)
	}

	// Invalidate list caches
	if err := s.permissionRedisRepo.InvalidatePermissionCache(ctx); err != nil {
		log.Printf("Failed to invalidate permission cache: %v", err)
	}

	return nil
}

// GetPermissionByID retrieves a permission by ID with caching
func (s *permissionService) GetPermissionByID(ctx context.Context, id string) (*entity.Permission, error) {
	if id == "" {
		return nil, errors.New("permission ID cannot be empty")
	}

	// Try to get from cache first
	cachedPermission, err := s.permissionRedisRepo.GetPermissionByID(ctx, id)
	if err == nil && cachedPermission != nil {
		log.Printf("Permission %s retrieved from cache", id)
		return cachedPermission, nil
	}

	// If not in cache, get from database
	permission, err := s.permissionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}
	if permission == nil {
		return nil, errors.New("permission not found")
	}

	// Cache the result
	if err := s.permissionRedisRepo.SetPermissionByID(ctx, permission); err != nil {
		log.Printf("Failed to cache permission: %v", err)
	}

	return permission, nil
}

// GetPermissionByName retrieves a permission by name with caching
func (s *permissionService) GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error) {
	if name == "" {
		return nil, errors.New("permission name cannot be empty")
	}

	// Try to get from cache first
	cachedPermission, err := s.permissionRedisRepo.GetPermissionByName(ctx, name)
	if err == nil && cachedPermission != nil {
		log.Printf("Permission %s retrieved from cache", name)
		return cachedPermission, nil
	}

	// If not in cache, get from database
	permission, err := s.permissionRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}
	if permission == nil {
		return nil, errors.New("permission not found")
	}

	// Cache the result
	if err := s.permissionRedisRepo.SetPermissionByName(ctx, permission); err != nil {
		log.Printf("Failed to cache permission: %v", err)
	}

	return permission, nil
}

// UpdatePermission updates a permission with validation
func (s *permissionService) UpdatePermission(ctx context.Context, permission *entity.Permission) error {
	if permission.ID == "" {
		return errors.New("permission ID cannot be empty")
	}

	// Validate permission name
	if err := s.ValidatePermissionName(ctx, permission.Name); err != nil {
		return err
	}

	// Check if permission exists
	existingPermission, err := s.permissionRepo.GetByID(ctx, permission.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing permission: %w", err)
	}
	if existingPermission == nil {
		return errors.New("permission not found")
	}

	// Check if new name is unique (excluding current permission)
	isUnique, err := s.IsPermissionNameUnique(ctx, permission.Name, &permission.ID)
	if err != nil {
		return fmt.Errorf("failed to check permission name uniqueness: %w", err)
	}
	if !isUnique {
		return errors.New("permission name already exists")
	}

	// Update permission in database
	if err := s.permissionRepo.Update(ctx, permission); err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}

	// Invalidate cache for this permission
	if err := s.permissionRedisRepo.InvalidatePermissionCacheByID(ctx, permission.ID); err != nil {
		log.Printf("Failed to invalidate permission cache: %v", err)
	}

	// Invalidate old name cache if name changed
	if existingPermission.Name != permission.Name {
		if err := s.permissionRedisRepo.InvalidatePermissionCache(ctx); err != nil {
			log.Printf("Failed to invalidate permission cache: %v", err)
		}
	}

	return nil
}

// DeletePermission deletes a permission after validation
func (s *permissionService) DeletePermission(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("permission ID cannot be empty")
	}

	// Check if permission can be deleted
	if err := s.CanDeletePermission(ctx, id); err != nil {
		return err
	}

	// Get permission before deletion for cache invalidation
	permission, err := s.permissionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get permission: %w", err)
	}
	if permission == nil {
		return errors.New("permission not found")
	}

	// Delete permission from database
	if err := s.permissionRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	// Invalidate cache
	if err := s.permissionRedisRepo.InvalidatePermissionCacheByID(ctx, id); err != nil {
		log.Printf("Failed to invalidate permission cache: %v", err)
	}

	return nil
}

// GetPermissionList retrieves all permissions with caching
func (s *permissionService) GetPermissionList(ctx context.Context, limit, offset int) ([]*entity.Permission, error) {
	// Try to get from cache first
	cachedPermissions, err := s.permissionRedisRepo.GetPermissionList(ctx, limit, offset)
	if err == nil && cachedPermissions != nil && len(cachedPermissions) > 0 {
		log.Println("Permission list retrieved from cache")
		return cachedPermissions, nil
	}

	// If not in cache, get from database
	permissions, err := s.permissionRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission list: %w", err)
	}

	// Cache the result
	if len(permissions) > 0 {
		if err := s.permissionRedisRepo.SetPermissionList(ctx, permissions, limit, offset); err != nil {
			log.Printf("Failed to cache permission list: %v", err)
		}
	}

	return permissions, nil
}

// GetPermissionCount retrieves the total number of permissions with caching
func (s *permissionService) GetPermissionCount(ctx context.Context) (int64, error) {
	// Try to get from cache first
	cachedCount, err := s.permissionRedisRepo.GetPermissionTotal(ctx)
	if err == nil && cachedCount > 0 {
		log.Println("Permission count retrieved from cache")
		return cachedCount, nil
	}

	// If not in cache, get from database
	count, err := s.permissionRepo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get permission count: %w", err)
	}

	// Cache the result
	if err := s.permissionRedisRepo.SetPermissionTotal(ctx, count); err != nil {
		log.Printf("Failed to cache permission count: %v", err)
	}

	return count, nil
}

// GetPermissionsByResource retrieves permissions by resource type with caching
func (s *permissionService) GetPermissionsByResource(ctx context.Context, resource string, limit, offset int) ([]*entity.Permission, error) {
	if resource == "" {
		return nil, errors.New("resource cannot be empty")
	}

	// Try to get from cache first
	cachedPermissions, err := s.permissionRedisRepo.GetPermissionsByResource(ctx, resource)
	if err == nil && cachedPermissions != nil {
		log.Printf("Permissions for resource %s retrieved from cache", resource)
		return cachedPermissions, nil
	}

	// If not in cache, get from database
	permissions, err := s.permissionRepo.ListByResource(ctx, resource, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions by resource: %w", err)
	}

	// Cache the result
	if err := s.permissionRedisRepo.SetPermissionsByResource(ctx, resource, permissions); err != nil {
		log.Printf("Failed to cache permissions by resource: %v", err)
	}

	return permissions, nil
}

// GetPermissionsByAction retrieves permissions by action type
func (s *permissionService) GetPermissionsByAction(ctx context.Context, action string, limit, offset int) ([]*entity.Permission, error) {
	if action == "" {
		return nil, errors.New("action cannot be empty")
	}

	permissions, err := s.permissionRepo.ListByAction(ctx, action, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions by action: %w", err)
	}

	return permissions, nil
}

// GetPermissionsByResourceAndAction retrieves permissions by resource and action
func (s *permissionService) GetPermissionsByResourceAndAction(ctx context.Context, resource, action string) ([]*entity.Permission, error) {
	if resource == "" || action == "" {
		return nil, errors.New("resource and action cannot be empty")
	}

	permissions, err := s.permissionRepo.ListByResourceAndAction(ctx, resource, action, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	return permissions, nil
}

// GetRolesByPermissionID retrieves all roles that have a specific permission
func (s *permissionService) GetRolesByPermissionID(ctx context.Context, permissionID string, limit, offset int) ([]*entity.Role, error) {
	if permissionID == "" {
		return nil, errors.New("permission ID cannot be empty")
	}

	roles, err := s.permissionRepo.GetRolesByPermissionID(ctx, permissionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	return roles, nil
}

// CountRolesByPermissionID counts the number of roles that have a specific permission
func (s *permissionService) CountRolesByPermissionID(ctx context.Context, permissionID string) (int64, error) {
	if permissionID == "" {
		return 0, errors.New("permission ID cannot be empty")
	}

	count, err := s.permissionRepo.CountRolesByPermissionID(ctx, permissionID)
	if err != nil {
		return 0, fmt.Errorf("failed to count roles: %w", err)
	}

	return count, nil
}

// GetPermissionWithRoles retrieves a permission with all its roles
func (s *permissionService) GetPermissionWithRoles(ctx context.Context, permissionID string) (*entity.PermissionWithRoles, error) {
	if permissionID == "" {
		return nil, errors.New("permission ID cannot be empty")
	}

	// Get permission
	permission, err := s.GetPermissionByID(ctx, permissionID)
	if err != nil {
		return nil, err
	}

	// Get roles
	roles, err := s.GetRolesByPermissionID(ctx, permissionID, 1000, 0)
	if err != nil {
		return nil, err
	}

	// Convert []*Role to []Role
	roleValues := make([]entity.Role, len(roles))
	for i, r := range roles {
		roleValues[i] = *r
	}

	return &entity.PermissionWithRoles{
		Permission: *permission,
		Roles:      roleValues,
	}, nil
}

// GetUserPermissions retrieves all permissions for a user (CRITICAL for authorization)
func (s *permissionService) GetUserPermissions(ctx context.Context, userID string) ([]*entity.Permission, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	// Try to get from cache first (CRITICAL for performance)
	cachedPermissions, err := s.permissionRedisRepo.GetUserPermissions(ctx, userID)
	if err == nil && cachedPermissions != nil {
		log.Printf("User %s permissions retrieved from cache", userID)
		return cachedPermissions, nil
	}

	// If not in cache, get from database
	permissions, err := s.permissionRepo.GetPermissionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Cache the result (important for authorization performance)
	if err := s.permissionRedisRepo.SetUserPermissions(ctx, userID, permissions); err != nil {
		log.Printf("Failed to cache user permissions: %v", err)
	}

	return permissions, nil
}

// GetUserPermissionNames retrieves just the permission names for a user
func (s *permissionService) GetUserPermissionNames(ctx context.Context, userID string) ([]string, error) {
	permissions, err := s.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(permissions))
	for i, p := range permissions {
		names[i] = p.Name
	}

	return names, nil
}

// HasPermission checks if a user has a specific permission (CRITICAL for authorization)
func (s *permissionService) HasPermission(ctx context.Context, userID, permissionName string) (bool, error) {
	if userID == "" || permissionName == "" {
		return false, errors.New("user ID and permission name cannot be empty")
	}

	// Try to get from cache first (CRITICAL for performance)
	cachedResult, err := s.permissionRedisRepo.GetUserHasPermission(ctx, userID, permissionName)
	if err == nil && cachedResult != nil {
		log.Printf("User %s permission check for %s retrieved from cache", userID, permissionName)
		return *cachedResult, nil
	}

	// If not in cache, check database
	hasPermission, err := s.permissionRepo.HasPermission(ctx, userID, permissionName)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	// Cache the result (important for authorization performance)
	if err := s.permissionRedisRepo.SetUserHasPermission(ctx, userID, permissionName, hasPermission); err != nil {
		log.Printf("Failed to cache permission check: %v", err)
	}

	return hasPermission, nil
}

// HasPermissions checks if a user has ALL specified permissions (AND logic)
func (s *permissionService) HasPermissions(ctx context.Context, userID string, permissionNames []string) (bool, error) {
	if userID == "" {
		return false, errors.New("user ID cannot be empty")
	}
	if len(permissionNames) == 0 {
		return false, errors.New("permission names cannot be empty")
	}

	hasPermissions, err := s.permissionRepo.HasPermissions(ctx, userID, permissionNames)
	if err != nil {
		return false, fmt.Errorf("failed to check permissions: %w", err)
	}

	return hasPermissions, nil
}

// HasAnyPermission checks if a user has ANY of the specified permissions (OR logic)
func (s *permissionService) HasAnyPermission(ctx context.Context, userID string, permissionNames []string) (bool, error) {
	if userID == "" {
		return false, errors.New("user ID cannot be empty")
	}
	if len(permissionNames) == 0 {
		return false, errors.New("permission names cannot be empty")
	}

	// Check each permission individually
	for _, permissionName := range permissionNames {
		hasPermission, err := s.HasPermission(ctx, userID, permissionName)
		if err != nil {
			return false, err
		}
		if hasPermission {
			return true, nil
		}
	}

	return false, nil
}

// GetUserPermissionsForResource retrieves all user permissions for a specific resource
func (s *permissionService) GetUserPermissionsForResource(ctx context.Context, userID, resource string) ([]*entity.Permission, error) {
	if userID == "" || resource == "" {
		return nil, errors.New("user ID and resource cannot be empty")
	}

	// Get all user permissions
	allPermissions, err := s.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Filter by resource
	resourcePermissions := make([]*entity.Permission, 0)
	for _, p := range allPermissions {
		if p.Resource == resource {
			resourcePermissions = append(resourcePermissions, p)
		}
	}

	return resourcePermissions, nil
}

// CanUserPerformAction checks if a user can perform a specific action on a resource
func (s *permissionService) CanUserPerformAction(ctx context.Context, userID, resource, action string) (bool, error) {
	if userID == "" || resource == "" || action == "" {
		return false, errors.New("user ID, resource, and action cannot be empty")
	}

	// Generate permission name
	permissionName := s.GeneratePermissionName(resource, action)

	// Check if user has this permission
	return s.HasPermission(ctx, userID, permissionName)
}

// CreatePermissionsBulk creates multiple permissions in bulk (for seeding)
func (s *permissionService) CreatePermissionsBulk(ctx context.Context, permissions []*entity.Permission) error {
	if len(permissions) == 0 {
		return errors.New("permissions cannot be empty")
	}

	// Validate all permissions
	for _, permission := range permissions {
		if err := s.ValidatePermissionName(ctx, permission.Name); err != nil {
			return fmt.Errorf("invalid permission %s: %w", permission.Name, err)
		}

		// Generate ID if not provided
		if permission.ID == "" {
			permission.ID = uuid.New().String()
		}
	}

	// Create permissions in bulk
	if err := s.permissionRepo.CreateBulk(ctx, permissions); err != nil {
		return fmt.Errorf("failed to create permissions in bulk: %w", err)
	}

	// Invalidate cache
	if err := s.permissionRedisRepo.InvalidatePermissionCache(ctx); err != nil {
		log.Printf("Failed to invalidate permission cache: %v", err)
	}

	return nil
}

// GetPermissionsByNames retrieves multiple permissions by their names
func (s *permissionService) GetPermissionsByNames(ctx context.Context, names []string) ([]*entity.Permission, error) {
	if len(names) == 0 {
		return []*entity.Permission{}, nil
	}

	permissions, err := s.permissionRepo.GetByNames(ctx, names)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions by names: %w", err)
	}

	return permissions, nil
}

// ValidatePermissionName validates the permission name format
func (s *permissionService) ValidatePermissionName(ctx context.Context, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("permission name cannot be empty")
	}

	// Permission name should be in format "resource:action"
	parts := strings.Split(name, ":")
	if len(parts) != 2 {
		return errors.New("permission name must be in format 'resource:action'")
	}

	resource := strings.TrimSpace(parts[0])
	action := strings.TrimSpace(parts[1])

	if resource == "" || action == "" {
		return errors.New("resource and action cannot be empty")
	}

	if len(name) > 100 {
		return errors.New("permission name must be at most 100 characters")
	}

	return nil
}

// IsPermissionNameUnique checks if a permission name is unique
func (s *permissionService) IsPermissionNameUnique(ctx context.Context, name string, excludeID *string) (bool, error) {
	existingPermission, err := s.permissionRepo.GetByName(ctx, name)
	if err != nil {
		return false, fmt.Errorf("failed to check permission name: %w", err)
	}

	// If no permission with this name exists, it's unique
	if existingPermission == nil {
		return true, nil
	}

	// If we're excluding a specific ID (for updates), check if it's the same permission
	if excludeID != nil && existingPermission.ID == *excludeID {
		return true, nil
	}

	return false, nil
}

// CanDeletePermission checks if a permission can be deleted
func (s *permissionService) CanDeletePermission(ctx context.Context, permissionID string) error {
	// Check if permission has any roles
	roleCount, err := s.permissionRepo.CountRolesByPermissionID(ctx, permissionID)
	if err != nil {
		return fmt.Errorf("failed to count roles: %w", err)
	}

	if roleCount > 0 {
		return fmt.Errorf("cannot delete permission: %d role(s) have this permission assigned", roleCount)
	}

	return nil
}

// GeneratePermissionName generates a permission name from resource and action
func (s *permissionService) GeneratePermissionName(resource, action string) string {
	return fmt.Sprintf("%s:%s", strings.ToLower(strings.TrimSpace(resource)), strings.ToLower(strings.TrimSpace(action)))
}

// ParsePermissionName extracts resource and action from a permission name
func (s *permissionService) ParsePermissionName(permissionName string) (resource, action string, err error) {
	parts := strings.Split(permissionName, ":")
	if len(parts) != 2 {
		return "", "", errors.New("invalid permission name format, expected 'resource:action'")
	}

	resource = strings.TrimSpace(parts[0])
	action = strings.TrimSpace(parts[1])

	if resource == "" || action == "" {
		return "", "", errors.New("resource and action cannot be empty")
	}

	return resource, action, nil
}
