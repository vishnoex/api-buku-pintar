package entity

import "time"

type PaymentProviderStatus string

const (
	PaymentProviderStatusActive   PaymentProviderStatus = "active"
	PaymentProviderStatusInactive PaymentProviderStatus = "inactive"
)

// PaymentProvider represents a payment provider in the system
// Clean Architecture: Entity layer, no dependencies on infrastructure
type PaymentProvider struct {
	ID        string                `db:"id" json:"id"`
	Name      string                `db:"name" json:"name"`
	Status    PaymentProviderStatus `db:"status" json:"status"`
	CreatedAt time.Time             `db:"created_at" json:"created_at"`
	UpdatedAt time.Time             `db:"updated_at" json:"updated_at"`
}
