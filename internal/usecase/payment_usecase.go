package usecase

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// PaymentUsecase defines the interface for payment use cases
type PaymentUsecase interface {
	InitiatePayment(ctx context.Context, userID string, amount int64, currency, description string) (*entity.Payment, error)
	GetPayment(ctx context.Context, paymentID string) (*entity.Payment, error)
	HandleXenditCallback(ctx context.Context, callbackData map[string]interface{}) error
} 