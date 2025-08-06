package entity

import "time"

// LoginProviderStatus represents the status of a login provider
type LoginProviderStatus string

const (
	LoginProviderStatusActive   LoginProviderStatus = "active"
	LoginProviderStatusInactive LoginProviderStatus = "inactive"
)

// LoginProvider represents a login provider in the system
// Clean Architecture: Entity layer, no dependencies on infrastructure
type LoginProvider struct {
	ID        string               `db:"id" json:"id"`
	Name      string               `db:"name" json:"name"`
	Status    LoginProviderStatus  `db:"status" json:"status"`
	CreatedAt time.Time            `db:"created_at" json:"created_at"`
	UpdatedAt time.Time            `db:"updated_at" json:"updated_at"`
}
