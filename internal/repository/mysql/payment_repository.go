package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"time"
)

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) repository.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	query := `INSERT INTO payments (id, user_id, amount, currency, status, xendit_reference, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	payment.CreatedAt = now
	payment.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		payment.ID,
		payment.UserID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.XenditReference,
		payment.Description,
		payment.CreatedAt,
		payment.UpdatedAt,
	)
	return err
}

func (r *paymentRepository) GetByID(ctx context.Context, id string) (*entity.Payment, error) {
	query := `SELECT id, user_id, amount, currency, status, xendit_reference, description, created_at, updated_at
		FROM payments WHERE id = ?`

	payment := &entity.Payment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.XenditReference,
		&payment.Description,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return payment, nil
}

func (r *paymentRepository) GetByXenditReference(ctx context.Context, ref string) (*entity.Payment, error) {
	query := `SELECT id, user_id, amount, currency, status, xendit_reference, description, created_at, updated_at
		FROM payments WHERE xendit_reference = ?`

	payment := &entity.Payment{}
	err := r.db.QueryRowContext(ctx, query, ref).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.XenditReference,
		&payment.Description,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return payment, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *entity.Payment) error {
	query := `UPDATE payments
		SET user_id = ?, amount = ?, currency = ?, status = ?, xendit_reference = ?, description = ?, updated_at = ?
		WHERE id = ?`

	payment.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		payment.UserID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.XenditReference,
		payment.Description,
		payment.UpdatedAt,
		payment.ID,
	)
	return err
}

func (r *paymentRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.Payment, error) {
	query := `SELECT id, user_id, amount, currency, status, xendit_reference, description, created_at, updated_at
		FROM payments WHERE user_id = ?`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*entity.Payment
	for rows.Next() {
		payment := &entity.Payment{}
		err = rows.Scan(
			&payment.ID,
			&payment.UserID,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.XenditReference,
			&payment.Description,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
