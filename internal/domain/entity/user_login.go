package entity

import "time"

// UserLoginStatus represents the status of a user login
type UserLoginStatus string

const (
	UserLoginStatusActive   UserLoginStatus = "active"
	UserLoginStatusInactive UserLoginStatus = "inactive"
)

// UserLogin represents a user login in the system
// Clean Architecture: Entity layer, no dependencies on infrastructure
type UserLogin struct {
	ID              string          `db:"id" json:"id"`
	UserID          string          `db:"user_id" json:"user_id"`
	LoginProviderID string          `db:"login_provider_id" json:"login_provider_id"`
	Status          UserLoginStatus `db:"status" json:"status"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
}
