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
	ID              string        `json:"id"`
	UserID          string        `json:"user_id"`
	Amount          int64         `json:"amount"`
	Currency        string        `json:"currency"`
	Status          PaymentStatus `json:"status"`
	XenditReference string        `json:"xendit_reference"`
	Description     string        `json:"description"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
} 