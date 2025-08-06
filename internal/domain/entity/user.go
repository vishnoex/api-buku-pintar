package entity

import "time"

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleReader   UserRole = "reader"
	RoleEditor   UserRole = "editor"
)

// UserStatus represents the status of a user in the system
type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusSuspended UserStatus = "suspended"
)

// User represents a user in the system
// Clean Architecture: Entity layer, no dependencies on infrastructure
type User struct {
	ID                string     `db:"id" json:"id"`
	Name              string     `db:"name" json:"name"`
	Email             string     `db:"email" json:"email"`
	Password          *string    `db:"password" json:"password"`
	Role              UserRole   `db:"role" json:"role"`
	Avatar            *string    `db:"avatar" json:"avatar"`
	Status            UserStatus `db:"status" json:"status"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
} 