package usecase

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type paymentUsecase struct {
	paymentService service.PaymentService
}

func NewPaymentUsecase(paymentService service.PaymentService) PaymentUsecase {
	return &paymentUsecase{
		paymentService: paymentService,
	}
}

func (u *paymentUsecase) InitiatePayment(ctx context.Context, userID string, amount int64, currency, description string) (*entity.Payment, error) {
	payment := &entity.Payment{
		ID:          uuid.New().String(),
		UserID:      userID,
		Amount:      amount,
		Currency:    currency,
		Description: description,
		Status:      entity.PaymentStatusPending,
	}

	err := u.paymentService.InitiatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (u *paymentUsecase) GetPayment(ctx context.Context, paymentID string) (*entity.Payment, error) {
	return u.paymentService.GetPaymentByID(ctx, paymentID)
}

func (u *paymentUsecase) HandleXenditCallback(ctx context.Context, callbackData map[string]interface{}) error {
	externalID, ok := callbackData["external_id"].(string)
	if !ok {
		return errors.New("invalid callback data: missing external_id")
	}

	status, ok := callbackData["status"].(string)
	if !ok {
		return errors.New("invalid callback data: missing status")
	}

	payment, err := u.paymentService.GetPaymentByXenditReference(ctx, externalID)
	if err != nil {
		return fmt.Errorf("failed to get payment by xendit reference: %w", err)
	}
	if payment == nil {
		return fmt.Errorf("payment with xendit reference %s not found", externalID)
	}

	var paymentStatus entity.PaymentStatus
	switch status {
	case "PAID":
		paymentStatus = entity.PaymentStatusPaid
	case "EXPIRED":
		paymentStatus = entity.PaymentStatusExpired
	default:
		paymentStatus = entity.PaymentStatusFailed
	}

	return u.paymentService.UpdatePaymentStatus(ctx, payment.ID, paymentStatus)
} 