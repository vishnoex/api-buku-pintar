package entity

import "time"

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusExpired   PaymentStatus = "expired"
)

// Payment represents a payment transaction in the system
// Clean Architecture: Entity layer, no dependencies on infrastructure
// XenditReference can store invoice ID or payment reference from Xendit
// UserID links to the user making the payment
// Amount is in smallest currency unit (e.g., cents)
type Payment struct {
	ID              string        `db:"id" json:"id"`
	UserID          string        `db:"user_id" json:"user_id"`
	Amount          int64         `db:"amount" json:"amount"`
	Currency        string        `db:"currency" json:"currency"`
	Status          PaymentStatus `db:"status" json:"status"`
	XenditReference string        `db:"xendit_reference" json:"xendit_reference"`
	Description     string        `db:"description" json:"description"`
	CreatedAt       time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at" json:"updated_at"`
}
