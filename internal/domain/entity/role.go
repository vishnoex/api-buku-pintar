package entity

import "time"

// Role represents a role in the system for RBAC (Role-Based Access Control)
// Clean Architecture: Entity layer, no dependencies on infrastructure
type Role struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// RoleType represents predefined role types in the system
type RoleType string

const (
	RoleTypeAdmin   RoleType = "admin"
	RoleTypeEditor  RoleType = "editor"
	RoleTypeReader  RoleType = "reader"
	RoleTypePremium RoleType = "premium"
)
