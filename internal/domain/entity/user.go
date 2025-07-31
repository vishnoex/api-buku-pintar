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

type Provider string

const (
	EmailPassword Provider = "email-password"
	Google		  Provider = "google"
)

// User represents a user in the system
type User struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Role      UserRole   `json:"role"`
	Provider  Provider   `json:"provider"`
	Avatar    string     `json:"avatar"`
	Status    UserStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
} 