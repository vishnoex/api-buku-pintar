package helper

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"context"
	"fmt"
	"strings"
)

// PermissionChecker provides utility functions for permission checking
type PermissionChecker struct {
	permissionService service.PermissionService
	roleService       service.RoleService
}

// NewPermissionChecker creates a new PermissionChecker instance
func NewPermissionChecker(
	permissionService service.PermissionService,
	roleService service.RoleService,
) *PermissionChecker {
	return &PermissionChecker{
		permissionService: permissionService,
		roleService:       roleService,
	}
}

// ============================================================================
// Basic Permission Checks
// ============================================================================

// CanUser checks if a user has a specific permission
func (pc *PermissionChecker) CanUser(ctx context.Context, userID, permissionName string) (bool, error) {
	return pc.permissionService.HasPermission(ctx, userID, permissionName)
}

// CanUserAll checks if a user has all specified permissions (AND logic)
func (pc *PermissionChecker) CanUserAll(ctx context.Context, userID string, permissions ...string) (bool, error) {
	return pc.permissionService.HasPermissions(ctx, userID, permissions)
}

// CanUserAny checks if a user has any of the specified permissions (OR logic)
func (pc *PermissionChecker) CanUserAny(ctx context.Context, userID string, permissions ...string) (bool, error) {
	return pc.permissionService.HasAnyPermission(ctx, userID, permissions)
}

// ============================================================================
// Resource-Based Permission Checks
// ============================================================================

// CanCreate checks if a user can create a resource
func (pc *PermissionChecker) CanCreate(ctx context.Context, userID, resource string) (bool, error) {
	return pc.permissionService.CanUserPerformAction(ctx, userID, resource, string(entity.ActionCreate))
}

// CanRead checks if a user can read a resource
func (pc *PermissionChecker) CanRead(ctx context.Context, userID, resource string) (bool, error) {
	return pc.permissionService.CanUserPerformAction(ctx, userID, resource, string(entity.ActionRead))
}

// CanUpdate checks if a user can update a resource
func (pc *PermissionChecker) CanUpdate(ctx context.Context, userID, resource string) (bool, error) {
	return pc.permissionService.CanUserPerformAction(ctx, userID, resource, string(entity.ActionUpdate))
}

// CanDelete checks if a user can delete a resource
func (pc *PermissionChecker) CanDelete(ctx context.Context, userID, resource string) (bool, error) {
	return pc.permissionService.CanUserPerformAction(ctx, userID, resource, string(entity.ActionDelete))
}

// CanList checks if a user can list resources
func (pc *PermissionChecker) CanList(ctx context.Context, userID, resource string) (bool, error) {
	return pc.permissionService.CanUserPerformAction(ctx, userID, resource, string(entity.ActionList))
}

// CanManage checks if a user can manage a resource (admin-level access)
func (pc *PermissionChecker) CanManage(ctx context.Context, userID, resource string) (bool, error) {
	return pc.permissionService.CanUserPerformAction(ctx, userID, resource, string(entity.ActionManage))
}

// ============================================================================
// Resource-Specific Permission Checks
// ============================================================================

// Ebook Permissions
func (pc *PermissionChecker) CanCreateEbook(ctx context.Context, userID string) (bool, error) {
	return pc.CanCreate(ctx, userID, string(entity.ResourceEbook))
}

func (pc *PermissionChecker) CanReadEbook(ctx context.Context, userID string) (bool, error) {
	return pc.CanRead(ctx, userID, string(entity.ResourceEbook))
}

func (pc *PermissionChecker) CanUpdateEbook(ctx context.Context, userID string) (bool, error) {
	return pc.CanUpdate(ctx, userID, string(entity.ResourceEbook))
}

func (pc *PermissionChecker) CanDeleteEbook(ctx context.Context, userID string) (bool, error) {
	return pc.CanDelete(ctx, userID, string(entity.ResourceEbook))
}

func (pc *PermissionChecker) CanManageEbook(ctx context.Context, userID string) (bool, error) {
	return pc.CanManage(ctx, userID, string(entity.ResourceEbook))
}

// Article Permissions
func (pc *PermissionChecker) CanCreateArticle(ctx context.Context, userID string) (bool, error) {
	return pc.CanCreate(ctx, userID, string(entity.ResourceArticle))
}

func (pc *PermissionChecker) CanReadArticle(ctx context.Context, userID string) (bool, error) {
	return pc.CanRead(ctx, userID, string(entity.ResourceArticle))
}

func (pc *PermissionChecker) CanUpdateArticle(ctx context.Context, userID string) (bool, error) {
	return pc.CanUpdate(ctx, userID, string(entity.ResourceArticle))
}

func (pc *PermissionChecker) CanDeleteArticle(ctx context.Context, userID string) (bool, error) {
	return pc.CanDelete(ctx, userID, string(entity.ResourceArticle))
}

func (pc *PermissionChecker) CanManageArticle(ctx context.Context, userID string) (bool, error) {
	return pc.CanManage(ctx, userID, string(entity.ResourceArticle))
}

// User Permissions
func (pc *PermissionChecker) CanCreateUser(ctx context.Context, userID string) (bool, error) {
	return pc.CanCreate(ctx, userID, string(entity.ResourceUser))
}

func (pc *PermissionChecker) CanReadUser(ctx context.Context, userID string) (bool, error) {
	return pc.CanRead(ctx, userID, string(entity.ResourceUser))
}

func (pc *PermissionChecker) CanUpdateUser(ctx context.Context, userID string) (bool, error) {
	return pc.CanUpdate(ctx, userID, string(entity.ResourceUser))
}

func (pc *PermissionChecker) CanDeleteUser(ctx context.Context, userID string) (bool, error) {
	return pc.CanDelete(ctx, userID, string(entity.ResourceUser))
}

func (pc *PermissionChecker) CanManageUsers(ctx context.Context, userID string) (bool, error) {
	return pc.CanManage(ctx, userID, string(entity.ResourceUser))
}

// Payment Permissions
func (pc *PermissionChecker) CanReadPayment(ctx context.Context, userID string) (bool, error) {
	return pc.CanRead(ctx, userID, string(entity.ResourcePayment))
}

func (pc *PermissionChecker) CanManagePayments(ctx context.Context, userID string) (bool, error) {
	return pc.CanManage(ctx, userID, string(entity.ResourcePayment))
}

// Banner Permissions
func (pc *PermissionChecker) CanCreateBanner(ctx context.Context, userID string) (bool, error) {
	return pc.CanCreate(ctx, userID, string(entity.ResourceBanner))
}

func (pc *PermissionChecker) CanUpdateBanner(ctx context.Context, userID string) (bool, error) {
	return pc.CanUpdate(ctx, userID, string(entity.ResourceBanner))
}

func (pc *PermissionChecker) CanDeleteBanner(ctx context.Context, userID string) (bool, error) {
	return pc.CanDelete(ctx, userID, string(entity.ResourceBanner))
}

// Category Permissions
func (pc *PermissionChecker) CanCreateCategory(ctx context.Context, userID string) (bool, error) {
	return pc.CanCreate(ctx, userID, string(entity.ResourceCategory))
}

func (pc *PermissionChecker) CanUpdateCategory(ctx context.Context, userID string) (bool, error) {
	return pc.CanUpdate(ctx, userID, string(entity.ResourceCategory))
}

func (pc *PermissionChecker) CanDeleteCategory(ctx context.Context, userID string) (bool, error) {
	return pc.CanDelete(ctx, userID, string(entity.ResourceCategory))
}

// ============================================================================
// Role-Based Checks
// ============================================================================

// IsAdmin checks if a user has the admin role
func (pc *PermissionChecker) IsAdmin(ctx context.Context, userID string) (bool, error) {
	return pc.hasRole(ctx, userID, string(entity.RoleTypeAdmin))
}

// IsEditor checks if a user has the editor role
func (pc *PermissionChecker) IsEditor(ctx context.Context, userID string) (bool, error) {
	return pc.hasRole(ctx, userID, string(entity.RoleTypeEditor))
}

// IsReader checks if a user has the reader role
func (pc *PermissionChecker) IsReader(ctx context.Context, userID string) (bool, error) {
	return pc.hasRole(ctx, userID, string(entity.RoleTypeReader))
}

// IsPremium checks if a user has the premium role
func (pc *PermissionChecker) IsPremium(ctx context.Context, userID string) (bool, error) {
	return pc.hasRole(ctx, userID, string(entity.RoleTypePremium))
}

// HasRole checks if a user has a specific role
func (pc *PermissionChecker) HasRole(ctx context.Context, userID, roleName string) (bool, error) {
	return pc.hasRole(ctx, userID, roleName)
}

// HasAnyRole checks if a user has any of the specified roles
func (pc *PermissionChecker) HasAnyRole(ctx context.Context, userID string, roleNames ...string) (bool, error) {
	role, err := pc.getUserRole(ctx, userID)
	if err != nil || role == nil {
		return false, err
	}

	for _, roleName := range roleNames {
		if strings.EqualFold(role.Name, roleName) {
			return true, nil
		}
	}

	return false, nil
}

// hasRole is a helper function to check if a user has a specific role
func (pc *PermissionChecker) hasRole(ctx context.Context, userID, roleName string) (bool, error) {
	role, err := pc.getUserRole(ctx, userID)
	if err != nil || role == nil {
		return false, err
	}

	return strings.EqualFold(role.Name, roleName), nil
}

// getUserRole retrieves the user's role
func (pc *PermissionChecker) getUserRole(ctx context.Context, userID string) (*entity.Role, error) {
	permissions, err := pc.permissionService.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	// If user has permissions, get their role through the role service
	// This requires the user's role ID, so we'll need to fetch it differently
	// For now, we'll use the role service directly with the user ID
	
	// Get all roles and find which one belongs to this user
	// This is a workaround - ideally we'd have a GetRoleByUserID method
	roles, err := pc.roleService.GetRoleList(ctx, 100, 0)
	if err != nil {
		return nil, err
	}

	// If we have permissions, the user must have a role
	// This is a simplified implementation - in practice, you'd want a more direct query
	if len(permissions) == 0 {
		return nil, nil
	}

	// Return the first role that has any of these permissions
	// This is a simplified approach - you may want to add a GetUserRole method to RoleService
	for _, role := range roles {
		rolePerms, err := pc.roleService.GetPermissionsByRoleID(ctx, role.ID)
		if err != nil {
			continue
		}
		
		// Check if any user permission matches role permissions
		for _, userPerm := range permissions {
			for _, rolePerm := range rolePerms {
				if userPerm.ID == rolePerm.ID {
					return role, nil
				}
			}
		}
	}

	return nil, nil
}

// ============================================================================
// Combined Permission Checks
// ============================================================================

// CanModifyContent checks if a user can create, update, or delete content (articles or ebooks)
func (pc *PermissionChecker) CanModifyContent(ctx context.Context, userID, resourceType string) (bool, error) {
	permissions := []string{
		fmt.Sprintf("%s:%s", resourceType, entity.ActionCreate),
		fmt.Sprintf("%s:%s", resourceType, entity.ActionUpdate),
		fmt.Sprintf("%s:%s", resourceType, entity.ActionDelete),
	}
	return pc.CanUserAny(ctx, userID, permissions...)
}

// CanAccessPremiumContent checks if a user can access premium content
func (pc *PermissionChecker) CanAccessPremiumContent(ctx context.Context, userID string) (bool, error) {
	// Check if user is premium or admin
	isPremium, err := pc.IsPremium(ctx, userID)
	if err != nil {
		return false, err
	}
	if isPremium {
		return true, nil
	}

	return pc.IsAdmin(ctx, userID)
}

// CanManageContent checks if a user can fully manage content (admin-level access)
func (pc *PermissionChecker) CanManageContent(ctx context.Context, userID string) (bool, error) {
	permissions := []string{
		entity.PermissionEbookManage,
		entity.PermissionArticleManage,
	}
	return pc.CanUserAny(ctx, userID, permissions...)
}

// CanManageSystem checks if a user has system-level management permissions
func (pc *PermissionChecker) CanManageSystem(ctx context.Context, userID string) (bool, error) {
	permissions := []string{
		entity.PermissionUserManage,
		entity.PermissionRoleManage,
	}
	return pc.CanUserAny(ctx, userID, permissions...)
}

// ============================================================================
// Ownership-Based Checks
// ============================================================================

// CanModifyOwnContent checks if a user can modify their own content
// This combines ownership check with permission check
func (pc *PermissionChecker) CanModifyOwnContent(ctx context.Context, userID, resourceType, ownerID string) (bool, error) {
	// If user is the owner, check for basic edit permissions
	if userID == ownerID {
		return pc.CanUpdate(ctx, userID, resourceType)
	}

	// If not the owner, check for admin-level permissions
	return pc.CanManage(ctx, userID, resourceType)
}

// CanDeleteOwnContent checks if a user can delete their own content
func (pc *PermissionChecker) CanDeleteOwnContent(ctx context.Context, userID, resourceType, ownerID string) (bool, error) {
	// If user is the owner, check for delete permissions
	if userID == ownerID {
		return pc.CanDelete(ctx, userID, resourceType)
	}

	// If not the owner, check for admin-level permissions
	return pc.CanManage(ctx, userID, resourceType)
}

// ============================================================================
// Utility Functions
// ============================================================================

// GetUserPermissions retrieves all permissions for a user
func (pc *PermissionChecker) GetUserPermissions(ctx context.Context, userID string) ([]*entity.Permission, error) {
	return pc.permissionService.GetUserPermissions(ctx, userID)
}

// GetUserPermissionNames retrieves all permission names for a user
func (pc *PermissionChecker) GetUserPermissionNames(ctx context.Context, userID string) ([]string, error) {
	permissions, err := pc.permissionService.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(permissions))
	for i, perm := range permissions {
		names[i] = perm.Name
	}

	return names, nil
}

// GetUserRole retrieves the user's role
func (pc *PermissionChecker) GetUserRole(ctx context.Context, userID string) (*entity.Role, error) {
	return pc.getUserRole(ctx, userID)
}

// GetUserRoleName retrieves the user's role name
func (pc *PermissionChecker) GetUserRoleName(ctx context.Context, userID string) (string, error) {
	role, err := pc.getUserRole(ctx, userID)
	if err != nil {
		return "", err
	}
	if role == nil {
		return "", nil
	}
	return role.Name, nil
}

// ============================================================================
// Batch Permission Checks
// ============================================================================

// CheckMultiplePermissions checks multiple permissions and returns a map of results
func (pc *PermissionChecker) CheckMultiplePermissions(ctx context.Context, userID string, permissions []string) (map[string]bool, error) {
	results := make(map[string]bool)

	for _, permission := range permissions {
		hasPermission, err := pc.CanUser(ctx, userID, permission)
		if err != nil {
			return nil, err
		}
		results[permission] = hasPermission
	}

	return results, nil
}

// CheckResourceActions checks all actions for a resource and returns a map of results
func (pc *PermissionChecker) CheckResourceActions(ctx context.Context, userID, resource string) (map[string]bool, error) {
	actions := []string{
		string(entity.ActionCreate),
		string(entity.ActionRead),
		string(entity.ActionUpdate),
		string(entity.ActionDelete),
		string(entity.ActionList),
		string(entity.ActionManage),
	}

	results := make(map[string]bool)

	for _, action := range actions {
		canPerform, err := pc.permissionService.CanUserPerformAction(ctx, userID, resource, action)
		if err != nil {
			return nil, err
		}
		results[action] = canPerform
	}

	return results, nil
}

// ============================================================================
// Permission Helpers
// ============================================================================

// BuildPermissionName creates a permission name from resource and action
func (pc *PermissionChecker) BuildPermissionName(resource, action string) string {
	return fmt.Sprintf("%s:%s", resource, action)
}

// ParsePermissionName splits a permission name into resource and action
func (pc *PermissionChecker) ParsePermissionName(permissionName string) (resource, action string) {
	parts := strings.Split(permissionName, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return permissionName, ""
}

// ValidatePermissionFormat checks if a permission name follows the expected format
func (pc *PermissionChecker) ValidatePermissionFormat(permissionName string) bool {
	parts := strings.Split(permissionName, ":")
	return len(parts) == 2 && parts[0] != "" && parts[1] != ""
}
