package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// PaymentRepository defines the interface for payment data operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) error
	GetByID(ctx context.Context, id string) (*entity.Payment, error)
	GetByXenditReference(ctx context.Context, ref string) (*entity.Payment, error)
	Update(ctx context.Context, payment *entity.Payment) error
	ListByUserID(ctx context.Context, userID string) ([]*entity.Payment, error)
} 