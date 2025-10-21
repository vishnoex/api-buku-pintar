package entity

// RolePermission represents the many-to-many relationship between roles and permissions
// This is a junction table entity for RBAC (Role-Based Access Control)
// Clean Architecture: Entity layer, no dependencies on infrastructure
type RolePermission struct {
	RoleID       string `db:"role_id" json:"role_id"`
	PermissionID string `db:"permission_id" json:"permission_id"`
}

// RoleWithPermissions represents a role with its associated permissions
// This is a helper struct for queries that join roles and permissions
type RoleWithPermissions struct {
	Role        Role         `json:"role"`
	Permissions []Permission `json:"permissions"`
}

// PermissionWithRoles represents a permission with its associated roles
// This is a helper struct for queries that join permissions and roles
type PermissionWithRoles struct {
	Permission Permission `json:"permission"`
	Roles      []Role     `json:"roles"`
}
