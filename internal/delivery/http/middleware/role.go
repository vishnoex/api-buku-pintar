package middleware

import (
	"buku-pintar/internal/domain/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// RoleMiddleware handles role-based and permission-based access control
type RoleMiddleware struct {
	roleService       service.RoleService
	permissionService service.PermissionService
}

// NewRoleMiddleware creates a new instance of RoleMiddleware
func NewRoleMiddleware(
	roleService service.RoleService,
	permissionService service.PermissionService,
) *RoleMiddleware {
	return &RoleMiddleware{
		roleService:       roleService,
		permissionService: permissionService,
	}
}

// RequireRole creates a middleware that requires the user to have a specific role
// Usage: router.Use(roleMiddleware.RequireRole("admin"))
func (m *RoleMiddleware) RequireRole(roleName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context (set by auth middleware)
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				log.Printf("Role middleware: user not found in context")
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			// Check if user has a role assigned
			if user.RoleID == nil {
				log.Printf("Role middleware: user %s has no role assigned", user.ID)
				m.respondWithError(w, http.StatusForbidden, "No role assigned to user")
				return
			}

			// Get user's role
			role, err := m.roleService.GetRoleByID(r.Context(), *user.RoleID)
			if err != nil {
				log.Printf("Role middleware: failed to get role: %v", err)
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify role")
				return
			}

			if role == nil {
				log.Printf("Role middleware: role not found for user %s", user.ID)
				m.respondWithError(w, http.StatusForbidden, "Invalid role")
				return
			}

			// Check if user's role matches required role
			if !strings.EqualFold(role.Name, roleName) {
				log.Printf("Role middleware: user %s role '%s' does not match required '%s'", 
					user.ID, role.Name, roleName)
				m.respondWithError(w, http.StatusForbidden, 
					fmt.Sprintf("Requires '%s' role", roleName))
				return
			}

			// Role check passed, proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole creates a middleware that requires the user to have any of the specified roles
// Usage: router.Use(roleMiddleware.RequireAnyRole("admin", "editor"))
func (m *RoleMiddleware) RequireAnyRole(roleNames ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				log.Printf("Role middleware: user not found in context")
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			// Check if user has a role assigned
			if user.RoleID == nil {
				log.Printf("Role middleware: user %s has no role assigned", user.ID)
				m.respondWithError(w, http.StatusForbidden, "No role assigned to user")
				return
			}

			// Get user's role
			role, err := m.roleService.GetRoleByID(r.Context(), *user.RoleID)
			if err != nil {
				log.Printf("Role middleware: failed to get role: %v", err)
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify role")
				return
			}

			if role == nil {
				log.Printf("Role middleware: role not found for user %s", user.ID)
				m.respondWithError(w, http.StatusForbidden, "Invalid role")
				return
			}

			// Check if user's role matches any of the required roles
			hasRole := false
			for _, requiredRole := range roleNames {
				if strings.EqualFold(role.Name, requiredRole) {
					hasRole = true
					break
				}
			}

			if !hasRole {
				log.Printf("Role middleware: user %s role '%s' does not match any required roles: %v", 
					user.ID, role.Name, roleNames)
				m.respondWithError(w, http.StatusForbidden, 
					fmt.Sprintf("Requires one of: %s", strings.Join(roleNames, ", ")))
				return
			}

			// Role check passed, proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission creates a middleware that requires the user to have a specific permission
// Usage: router.Use(roleMiddleware.RequirePermission("ebook:create"))
func (m *RoleMiddleware) RequirePermission(permissionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				log.Printf("Permission middleware: user not found in context")
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			// Check if user has the required permission
			hasPermission, err := m.permissionService.HasPermission(r.Context(), user.ID, permissionName)
			if err != nil {
				log.Printf("Permission middleware: failed to check permission: %v", err)
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify permission")
				return
			}

			if !hasPermission {
				log.Printf("Permission middleware: user %s does not have permission '%s'", 
					user.ID, permissionName)
				m.respondWithError(w, http.StatusForbidden, 
					fmt.Sprintf("Permission denied: requires '%s'", permissionName))
				return
			}

			// Permission check passed, proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAllPermissions creates a middleware that requires the user to have ALL specified permissions
// Usage: router.Use(roleMiddleware.RequireAllPermissions("ebook:create", "ebook:publish"))
func (m *RoleMiddleware) RequireAllPermissions(permissionNames ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				log.Printf("Permission middleware: user not found in context")
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			// Check if user has all required permissions
			hasAllPermissions, err := m.permissionService.HasPermissions(r.Context(), user.ID, permissionNames)
			if err != nil {
				log.Printf("Permission middleware: failed to check permissions: %v", err)
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify permissions")
				return
			}

			if !hasAllPermissions {
				log.Printf("Permission middleware: user %s does not have all required permissions: %v", 
					user.ID, permissionNames)
				m.respondWithError(w, http.StatusForbidden, 
					fmt.Sprintf("Permission denied: requires all of %s", strings.Join(permissionNames, ", ")))
				return
			}

			// Permission check passed, proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyPermission creates a middleware that requires the user to have ANY of the specified permissions
// Usage: router.Use(roleMiddleware.RequireAnyPermission("ebook:read", "ebook:read_free"))
func (m *RoleMiddleware) RequireAnyPermission(permissionNames ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				log.Printf("Permission middleware: user not found in context")
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			// Check if user has any of the required permissions
			hasAnyPermission, err := m.permissionService.HasAnyPermission(r.Context(), user.ID, permissionNames)
			if err != nil {
				log.Printf("Permission middleware: failed to check permissions: %v", err)
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify permissions")
				return
			}

			if !hasAnyPermission {
				log.Printf("Permission middleware: user %s does not have any required permissions: %v", 
					user.ID, permissionNames)
				m.respondWithError(w, http.StatusForbidden, 
					fmt.Sprintf("Permission denied: requires one of %s", strings.Join(permissionNames, ", ")))
				return
			}

			// Permission check passed, proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// RequireResourceAction creates a middleware that requires the user to have permission for a resource-action pair
// Usage: router.Use(roleMiddleware.RequireResourceAction("ebook", "create"))
func (m *RoleMiddleware) RequireResourceAction(resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				log.Printf("Permission middleware: user not found in context")
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			// Check if user can perform action on resource
			canPerform, err := m.permissionService.CanUserPerformAction(r.Context(), user.ID, resource, action)
			if err != nil {
				log.Printf("Permission middleware: failed to check resource action: %v", err)
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify permission")
				return
			}

			if !canPerform {
				log.Printf("Permission middleware: user %s cannot perform '%s' on '%s'", 
					user.ID, action, resource)
				m.respondWithError(w, http.StatusForbidden, 
					fmt.Sprintf("Permission denied: cannot '%s' on '%s'", action, resource))
				return
			}

			// Permission check passed, proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// CheckPermissionFunc is a custom function type for checking permissions
// This allows for dynamic permission checks based on request context
type CheckPermissionFunc func(ctx context.Context, userID string) (bool, error)

// RequireCustomPermission creates a middleware with a custom permission check function
// Usage: router.Use(roleMiddleware.RequireCustomPermission(func(ctx, userID) { ... }))
func (m *RoleMiddleware) RequireCustomPermission(checkFunc CheckPermissionFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			user, err := GetUserFromContext(r.Context())
			if err != nil || user == nil {
				log.Printf("Permission middleware: user not found in context")
				m.respondWithError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			// Run custom permission check
			hasPermission, err := checkFunc(r.Context(), user.ID)
			if err != nil {
				log.Printf("Permission middleware: custom check failed: %v", err)
				m.respondWithError(w, http.StatusInternalServerError, "Failed to verify permission")
				return
			}

			if !hasPermission {
				log.Printf("Permission middleware: user %s failed custom permission check", user.ID)
				m.respondWithError(w, http.StatusForbidden, "Permission denied")
				return
			}

			// Permission check passed, proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// respondWithError sends an error response with proper JSON format
func (m *RoleMiddleware) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, `{"error":"%s","status":%d}`, message, statusCode)
}

// InjectRoleIntoContext adds the user's role information to the request context
// This can be useful for handlers that need role information without querying the database
func (m *RoleMiddleware) InjectRoleIntoContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user from context
		user, err := GetUserFromContext(r.Context())
		if err != nil || user == nil {
			// No user in context, skip role injection
			next.ServeHTTP(w, r)
			return
		}

		// Check if user has a role assigned
		if user.RoleID == nil {
			// No role assigned, skip injection
			next.ServeHTTP(w, r)
			return
		}

		// Get user's role
		role, err := m.roleService.GetRoleByID(r.Context(), *user.RoleID)
		if err != nil {
			log.Printf("Role injection: failed to get role: %v", err)
			// Don't block request on role retrieval failure
			next.ServeHTTP(w, r)
			return
		}

		// Inject role into context
		ctx := context.WithValue(r.Context(), RoleContextKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// InjectPermissionsIntoContext adds the user's permissions to the request context
// This can be useful for handlers that need to check multiple permissions
func (m *RoleMiddleware) InjectPermissionsIntoContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user from context
		user, err := GetUserFromContext(r.Context())
		if err != nil || user == nil {
			// No user in context, skip permission injection
			next.ServeHTTP(w, r)
			return
		}

		// Get user's permissions
		permissions, err := m.permissionService.GetUserPermissions(r.Context(), user.ID)
		if err != nil {
			log.Printf("Permission injection: failed to get permissions: %v", err)
			// Don't block request on permission retrieval failure
			next.ServeHTTP(w, r)
			return
		}

		// Inject permissions into context
		ctx := context.WithValue(r.Context(), PermissionsContextKey, permissions)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Context keys for storing role and permissions
const (
	RoleContextKey        ContextKey = "role"
	PermissionsContextKey ContextKey = "permissions"
)
