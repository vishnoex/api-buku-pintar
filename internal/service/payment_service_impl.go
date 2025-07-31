package service

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/internal/domain/service"
	"context"

	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
)

type paymentService struct {
	paymentRepo repository.PaymentRepository
	xenditKey   string
}

func NewPaymentService(paymentRepo repository.PaymentRepository, xenditKey string) service.PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		xenditKey:   xenditKey,
	}
}

func (s *paymentService) InitiatePayment(ctx context.Context, payment *entity.Payment) error {
	xendit.Opt.SecretKey = s.xenditKey

	// Create invoice on Xendit
	data := invoice.CreateParams{
		ExternalID:  payment.ID,
		Amount:      float64(payment.Amount),
		Description: payment.Description,
		Currency:    payment.Currency,
	}

	resp, err := invoice.Create(&data)
	if err != nil {
		return err
	}

	// Save Xendit reference and update status
	payment.XenditReference = resp.ID
	payment.Status = entity.PaymentStatusPending

	// Create payment record in our database
	return s.paymentRepo.Create(ctx, payment)
}

func (s *paymentService) GetPaymentByID(ctx context.Context, id string) (*entity.Payment, error) {
	return s.paymentRepo.GetByID(ctx, id)
}

func (s *paymentService) GetPaymentByXenditReference(ctx context.Context, ref string) (*entity.Payment, error) {
	return s.paymentRepo.GetByXenditReference(ctx, ref)
}

func (s *paymentService) UpdatePaymentStatus(ctx context.Context, id string, status entity.PaymentStatus) error {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if payment == nil {
		return nil // Or return an error
	}

	payment.Status = status
	return s.paymentRepo.Update(ctx, payment)
}

func (s *paymentService) ListPaymentsByUserID(ctx context.Context, userID string) ([]*entity.Payment, error) {
	return s.paymentRepo.ListByUserID(ctx, userID)
} 