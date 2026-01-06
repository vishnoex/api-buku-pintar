package middleware

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Context keys for permission middleware
type contextKey string

const (
	PermissionsContextKey contextKey = "user_permissions"
)

// PermissionMiddleware handles fine-grained permission-based access control
// Separates permission logic from role logic for better modularity
type PermissionMiddleware struct {
	permissionService service.PermissionService
	roleService       service.RoleService
	auditLog          *PermissionAuditLogger
	debugMode         bool
}

// PermissionMiddlewareConfig holds configuration for permission middleware
type PermissionMiddlewareConfig struct {
	EnableAuditLog bool
	EnableDebug    bool
	AuditLogger    *PermissionAuditLogger
}

// NewPermissionMiddleware creates a new instance of PermissionMiddleware
func NewPermissionMiddleware(
	permissionService service.PermissionService,
	roleService service.RoleService,
	config *PermissionMiddlewareConfig,
) *PermissionMiddleware {
	if config == nil {
		config = &PermissionMiddlewareConfig{
			EnableAuditLog: true,
			EnableDebug:    false,
		}
	}

	if config.AuditLogger == nil && config.EnableAuditLog {
		config.AuditLogger = NewPermissionAuditLogger()
	}

	return &PermissionMiddleware{
		permissionService: permissionService,
		roleService:       roleService,
		auditLog:          config.AuditLogger,
		debugMode:         config.EnableDebug,
	}
}

// CheckPermission creates a middleware that checks if user has a specific permission
// This is the core permission checking middleware
func (m *PermissionMiddleware) CheckPermission(permissionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			
			// Get user from context
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				m.logPermissionCheck(r.Context(), "", permissionName, false, "user_not_found", time.Since(startTime))
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required", permissionName)
				return
			}

			// Check permission
			hasPermission, err := m.permissionService.HasPermission(r.Context(), user.ID, permissionName)
			if err != nil {
				m.logPermissionCheck(r.Context(), user.ID, permissionName, false, fmt.Sprintf("check_error: %v", err), time.Since(startTime))
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify permission", permissionName)
				return
			}

			if !hasPermission {
				m.logPermissionCheck(r.Context(), user.ID, permissionName, false, "permission_denied", time.Since(startTime))
				m.respondWithError(w, http.StatusForbidden, 
					fmt.Sprintf("Permission denied: requires '%s'", permissionName), permissionName)
				return
			}

			// Permission granted
			m.logPermissionCheck(r.Context(), user.ID, permissionName, true, "granted", time.Since(startTime))
			
			// Add permission info to context for debugging
			if m.debugMode {
				ctx := context.WithValue(r.Context(), "checked_permission", permissionName)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CheckRole creates a middleware that checks if user has a specific role
// Convenience method for role-based checks
func (m *PermissionMiddleware) CheckRole(roleName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				m.logPermissionCheck(r.Context(), "", fmt.Sprintf("role:%s", roleName), false, "user_not_found", time.Since(startTime))
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required", roleName)
				return
			}

			// Check if user has a role assigned
			if user.RoleID == nil {
				m.logPermissionCheck(r.Context(), user.ID, fmt.Sprintf("role:%s", roleName), false, "no_role_assigned", time.Since(startTime))
				m.respondWithError(w, http.StatusForbidden, "No role assigned to user", roleName)
				return
			}

			// Get user's role from role service
			role, err := m.roleService.GetRoleByID(r.Context(), *user.RoleID)
			if err != nil {
				m.logPermissionCheck(r.Context(), user.ID, fmt.Sprintf("role:%s", roleName), false, fmt.Sprintf("check_error: %v", err), time.Since(startTime))
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify role", roleName)
				return
			}

			hasRole := role != nil && strings.EqualFold(role.Name, roleName)
			if !hasRole {
				m.logPermissionCheck(r.Context(), user.ID, fmt.Sprintf("role:%s", roleName), false, "role_denied", time.Since(startTime))
				m.respondWithError(w, http.StatusForbidden, 
					fmt.Sprintf("Access denied: requires '%s' role", roleName), roleName)
				return
			}

			m.logPermissionCheck(r.Context(), user.ID, fmt.Sprintf("role:%s", roleName), true, "granted", time.Since(startTime))
			next.ServeHTTP(w, r)
		})
	}
}

// CheckOwnership creates a middleware that verifies resource ownership
// Extracts resource ID from URL path and checks if user owns the resource
func (m *PermissionMiddleware) CheckOwnership(resourceType string, resourceIDExtractor ResourceIDExtractor) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				m.logPermissionCheck(r.Context(), "", fmt.Sprintf("ownership:%s", resourceType), false, "user_not_found", time.Since(startTime))
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required", resourceType)
				return
			}

			// Extract resource ID from request
			resourceID := resourceIDExtractor(r)
			if resourceID == "" {
				m.logPermissionCheck(r.Context(), user.ID, fmt.Sprintf("ownership:%s", resourceType), false, "resource_id_not_found", time.Since(startTime))
				m.respondWithError(w, http.StatusBadRequest, "Resource ID not found", resourceType)
				return
			}

			// Check ownership
			isOwner, err := m.checkResourceOwnership(r.Context(), user.ID, resourceType, resourceID)
			if err != nil {
				m.logPermissionCheck(r.Context(), user.ID, fmt.Sprintf("ownership:%s:%s", resourceType, resourceID), false, fmt.Sprintf("check_error: %v", err), time.Since(startTime))
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify ownership", resourceType)
				return
			}

			if !isOwner {
				m.logPermissionCheck(r.Context(), user.ID, fmt.Sprintf("ownership:%s:%s", resourceType, resourceID), false, "not_owner", time.Since(startTime))
				m.respondWithError(w, http.StatusForbidden, 
					fmt.Sprintf("Access denied: you don't own this %s", resourceType), resourceType)
				return
			}

			m.logPermissionCheck(r.Context(), user.ID, fmt.Sprintf("ownership:%s:%s", resourceType, resourceID), true, "owner", time.Since(startTime))
			
			// Add ownership info to context
			ctx := context.WithValue(r.Context(), "resource_owner", true)
			ctx = context.WithValue(ctx, "resource_id", resourceID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CheckPermissionOrOwnership allows access if user has permission OR owns the resource
// Useful for endpoints where owners or admins can access
func (m *PermissionMiddleware) CheckPermissionOrOwnership(
	permissionName string,
	resourceType string,
	resourceIDExtractor ResourceIDExtractor,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				m.logPermissionCheck(r.Context(), "", permissionName, false, "user_not_found", time.Since(startTime))
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required", permissionName)
				return
			}

			// First, check if user has the permission (e.g., admin bypass)
			hasPermission, err := m.permissionService.HasPermission(r.Context(), user.ID, permissionName)
			if err == nil && hasPermission {
				m.logPermissionCheck(r.Context(), user.ID, permissionName, true, "granted_by_permission", time.Since(startTime))
				next.ServeHTTP(w, r)
				return
			}

			// If no permission, check ownership
			resourceID := resourceIDExtractor(r)
			if resourceID == "" {
				m.logPermissionCheck(r.Context(), user.ID, permissionName, false, "resource_id_not_found", time.Since(startTime))
				m.respondWithError(w, http.StatusBadRequest, "Resource ID not found", permissionName)
				return
			}

			isOwner, err := m.checkResourceOwnership(r.Context(), user.ID, resourceType, resourceID)
			if err != nil {
				m.logPermissionCheck(r.Context(), user.ID, fmt.Sprintf("%s|ownership:%s", permissionName, resourceType), false, fmt.Sprintf("check_error: %v", err), time.Since(startTime))
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify access", permissionName)
				return
			}

			if !isOwner {
				m.logPermissionCheck(r.Context(), user.ID, fmt.Sprintf("%s|ownership:%s", permissionName, resourceType), false, "access_denied", time.Since(startTime))
				m.respondWithError(w, http.StatusForbidden, 
					"Access denied: requires permission or resource ownership", permissionName)
				return
			}

			m.logPermissionCheck(r.Context(), user.ID, fmt.Sprintf("%s|ownership:%s", permissionName, resourceType), true, "granted_by_ownership", time.Since(startTime))
			ctx := context.WithValue(r.Context(), "resource_owner", true)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CheckAllPermissions requires user to have ALL specified permissions (AND logic)
func (m *PermissionMiddleware) CheckAllPermissions(permissionNames ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			permissionsStr := strings.Join(permissionNames, " AND ")
			
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				m.logPermissionCheck(r.Context(), "", permissionsStr, false, "user_not_found", time.Since(startTime))
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required", permissionsStr)
				return
			}

			hasAll, err := m.permissionService.HasPermissions(r.Context(), user.ID, permissionNames)
			if err != nil {
				m.logPermissionCheck(r.Context(), user.ID, permissionsStr, false, fmt.Sprintf("check_error: %v", err), time.Since(startTime))
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify permissions", permissionsStr)
				return
			}

			if !hasAll {
				m.logPermissionCheck(r.Context(), user.ID, permissionsStr, false, "missing_permissions", time.Since(startTime))
				m.respondWithError(w, http.StatusForbidden, 
					fmt.Sprintf("Permission denied: requires all of [%s]", strings.Join(permissionNames, ", ")), permissionsStr)
				return
			}

			m.logPermissionCheck(r.Context(), user.ID, permissionsStr, true, "granted", time.Since(startTime))
			next.ServeHTTP(w, r)
		})
	}
}

// CheckAnyPermission requires user to have ANY of the specified permissions (OR logic)
func (m *PermissionMiddleware) CheckAnyPermission(permissionNames ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			permissionsStr := strings.Join(permissionNames, " OR ")
			
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				m.logPermissionCheck(r.Context(), "", permissionsStr, false, "user_not_found", time.Since(startTime))
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required", permissionsStr)
				return
			}

			hasAny, err := m.permissionService.HasAnyPermission(r.Context(), user.ID, permissionNames)
			if err != nil {
				m.logPermissionCheck(r.Context(), user.ID, permissionsStr, false, fmt.Sprintf("check_error: %v", err), time.Since(startTime))
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify permissions", permissionsStr)
				return
			}

			if !hasAny {
				m.logPermissionCheck(r.Context(), user.ID, permissionsStr, false, "no_matching_permission", time.Since(startTime))
				m.respondWithError(w, http.StatusForbidden, 
					fmt.Sprintf("Permission denied: requires one of [%s]", strings.Join(permissionNames, ", ")), permissionsStr)
				return
			}

			m.logPermissionCheck(r.Context(), user.ID, permissionsStr, true, "granted", time.Since(startTime))
			next.ServeHTTP(w, r)
		})
	}
}

// CheckResourceAction checks if user can perform an action on a resource type
// Uses the resource:action permission naming convention
func (m *PermissionMiddleware) CheckResourceAction(resource entity.ResourceType, action entity.ActionType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			permissionName := fmt.Sprintf("%s:%s", resource, action)
			
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				m.logPermissionCheck(r.Context(), "", permissionName, false, "user_not_found", time.Since(startTime))
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required", permissionName)
				return
			}

			canPerform, err := m.permissionService.CanUserPerformAction(r.Context(), user.ID, string(resource), string(action))
			if err != nil {
				m.logPermissionCheck(r.Context(), user.ID, permissionName, false, fmt.Sprintf("check_error: %v", err), time.Since(startTime))
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify permission", permissionName)
				return
			}

			if !canPerform {
				m.logPermissionCheck(r.Context(), user.ID, permissionName, false, "action_denied", time.Since(startTime))
				m.respondWithError(w, http.StatusForbidden, 
					fmt.Sprintf("Permission denied: cannot '%s' on '%s'", action, resource), permissionName)
				return
			}

			m.logPermissionCheck(r.Context(), user.ID, permissionName, true, "granted", time.Since(startTime))
			next.ServeHTTP(w, r)
		})
	}
}

// InjectPermissions adds user permissions to request context for later use
// Useful for handlers that need to check multiple permissions
func (m *PermissionMiddleware) InjectPermissions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := GetUserFromContext(r.Context())
		if err != nil || user == nil {
			// No user, skip injection
			next.ServeHTTP(w, r)
			return
		}

		permissions, err := m.permissionService.GetUserPermissions(r.Context(), user.ID)
		if err != nil {
			log.Printf("Permission injection failed for user %s: %v", user.ID, err)
			// Don't block request
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), PermissionsContextKey, permissions)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ResourceIDExtractor is a function type for extracting resource IDs from requests
type ResourceIDExtractor func(r *http.Request) string

// checkResourceOwnership verifies if a user owns a specific resource
// This is a placeholder - implement actual ownership logic based on your domain
func (m *PermissionMiddleware) checkResourceOwnership(ctx context.Context, userID, resourceType, resourceID string) (bool, error) {
	// TODO: Implement actual ownership checks based on resource type
	// This should query the appropriate repository to verify ownership
	
	// Example implementation:
	// switch resourceType {
	// case "ebook":
	//     return m.ebookRepo.IsOwner(ctx, resourceID, userID)
	// case "article":
	//     return m.articleRepo.IsOwner(ctx, resourceID, userID)
	// default:
	//     return false, fmt.Errorf("unknown resource type: %s", resourceType)
	// }
	
	log.Printf("Ownership check not implemented for resource type: %s", resourceType)
	return false, fmt.Errorf("ownership check not implemented for resource type: %s", resourceType)
}

// logPermissionCheck logs permission check results for audit trail
func (m *PermissionMiddleware) logPermissionCheck(ctx context.Context, userID, permission string, granted bool, reason string, duration time.Duration) {
	if m.auditLog == nil {
		return
	}

	m.auditLog.Log(PermissionAuditEntry{
		Timestamp:  time.Now(),
		UserID:     userID,
		Permission: permission,
		Granted:    granted,
		Reason:     reason,
		Duration:   duration,
	})

	// Also log to standard logger if debug mode is enabled
	if m.debugMode {
		status := "DENIED"
		if granted {
			status = "GRANTED"
		}
		log.Printf("[PERMISSION CHECK] user=%s permission=%s status=%s reason=%s duration=%v",
			userID, permission, status, reason, duration)
	}
}

// respondWithError sends a standardized error response
func (m *PermissionMiddleware) respondWithError(w http.ResponseWriter, statusCode int, message string, permission string) {
	response := PermissionErrorResponse{
		Error:      message,
		Status:     statusCode,
		Permission: permission,
		Timestamp:  time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// PermissionErrorResponse represents a permission error response
type PermissionErrorResponse struct {
	Error      string    `json:"error"`
	Status     int       `json:"status"`
	Permission string    `json:"permission,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// GetPermissionsFromContext retrieves injected permissions from context
func GetPermissionsFromContext(ctx context.Context) ([]*entity.Permission, bool) {
	permissions, ok := ctx.Value(PermissionsContextKey).([]*entity.Permission)
	return permissions, ok
}

// IsResourceOwner checks if the resource_owner flag is set in context
func IsResourceOwner(ctx context.Context) bool {
	owner, ok := ctx.Value("resource_owner").(bool)
	return ok && owner
}

// GetResourceID retrieves the resource ID from context
func GetResourceID(ctx context.Context) (string, bool) {
	resourceID, ok := ctx.Value("resource_id").(string)
	return resourceID, ok
}
