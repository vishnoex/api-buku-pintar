package service

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// PaymentService defines the interface for payment business operations
type PaymentService interface {
	InitiatePayment(ctx context.Context, payment *entity.Payment) error
	GetPaymentByID(ctx context.Context, id string) (*entity.Payment, error)
	GetPaymentByXenditReference(ctx context.Context, ref string) (*entity.Payment, error)
	UpdatePaymentStatus(ctx context.Context, id string, status entity.PaymentStatus) error
	ListPaymentsByUserID(ctx context.Context, userID string) ([]*entity.Payment, error)
}
