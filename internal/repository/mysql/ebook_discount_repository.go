package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"time"
)

type ebookDiscountRepository struct {
	db *sql.DB
}

// NewEbookDiscountRepository creates a new instance of EbookDiscountRepository
func NewEbookDiscountRepository(db *sql.DB) repository.EbookDiscountRepository {
	return &ebookDiscountRepository{
		db: db,
	}
}

func (r *ebookDiscountRepository) Create(ctx context.Context, discount *entity.EbookDiscount) error {
	query := `
		INSERT INTO ebook_discounts (id, ebook_id, discount_price, started_at, ended_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		discount.ID,
		discount.EbookID,
		discount.DiscountPrice,
		discount.StartedAt,
		discount.EndedAt,
		now,
		now,
	)
	
	return err
}

func (r *ebookDiscountRepository) GetByID(ctx context.Context, id string) (*entity.EbookDiscount, error) {
	query := `
		SELECT id, ebook_id, discount_price, started_at, ended_at, created_at, updated_at
		FROM ebook_discounts
		WHERE id = ?
	`
	
	discount := &entity.EbookDiscount{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&discount.ID,
		&discount.EbookID,
		&discount.DiscountPrice,
		&discount.StartedAt,
		&discount.EndedAt,
		&discount.CreatedAt,
		&discount.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return discount, nil
}

func (r *ebookDiscountRepository) Update(ctx context.Context, discount *entity.EbookDiscount) error {
	query := `
		UPDATE ebook_discounts
		SET ebook_id = ?, discount_price = ?, started_at = ?, ended_at = ?, updated_at = ?
		WHERE id = ?
	`
	
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		discount.EbookID,
		discount.DiscountPrice,
		discount.StartedAt,
		discount.EndedAt,
		now,
		discount.ID,
	)
	
	return err
}

func (r *ebookDiscountRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM ebook_discounts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ebookDiscountRepository) List(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error) {
	query := `
		SELECT id, ebook_id, discount_price, started_at, ended_at, created_at, updated_at
		FROM ebook_discounts
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var discounts []*entity.EbookDiscount
	for rows.Next() {
		discount := &entity.EbookDiscount{}
		err := rows.Scan(
			&discount.ID,
			&discount.EbookID,
			&discount.DiscountPrice,
			&discount.StartedAt,
			&discount.EndedAt,
			&discount.CreatedAt,
			&discount.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, discount)
	}
	
	return discounts, nil
}

func (r *ebookDiscountRepository) GetByEbookID(ctx context.Context, ebookID string) ([]*entity.EbookDiscount, error) {
	query := `
		SELECT id, ebook_id, discount_price, started_at, ended_at, created_at, updated_at
		FROM ebook_discounts
		WHERE ebook_id = ?
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, ebookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var discounts []*entity.EbookDiscount
	for rows.Next() {
		discount := &entity.EbookDiscount{}
		err := rows.Scan(
			&discount.ID,
			&discount.EbookID,
			&discount.DiscountPrice,
			&discount.StartedAt,
			&discount.EndedAt,
			&discount.CreatedAt,
			&discount.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, discount)
	}
	
	return discounts, nil
}

func (r *ebookDiscountRepository) GetActiveDiscounts(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error) {
	query := `
		SELECT id, ebook_id, discount_price, started_at, ended_at, created_at, updated_at
		FROM ebook_discounts
		WHERE started_at <= ? AND ended_at >= ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	now := time.Now()
	rows, err := r.db.QueryContext(ctx, query, now, now, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var discounts []*entity.EbookDiscount
	for rows.Next() {
		discount := &entity.EbookDiscount{}
		err := rows.Scan(
			&discount.ID,
			&discount.EbookID,
			&discount.DiscountPrice,
			&discount.StartedAt,
			&discount.EndedAt,
			&discount.CreatedAt,
			&discount.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, discount)
	}
	
	return discounts, nil
}

func (r *ebookDiscountRepository) GetActiveDiscountByEbookID(ctx context.Context, ebookID string) (*entity.EbookDiscount, error) {
	query := `
		SELECT id, ebook_id, discount_price, started_at, ended_at, created_at, updated_at
		FROM ebook_discounts
		WHERE ebook_id = ? AND started_at <= ? AND ended_at >= ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	now := time.Now()
	discount := &entity.EbookDiscount{}
	err := r.db.QueryRowContext(ctx, query, ebookID, now, now).Scan(
		&discount.ID,
		&discount.EbookID,
		&discount.DiscountPrice,
		&discount.StartedAt,
		&discount.EndedAt,
		&discount.CreatedAt,
		&discount.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return discount, nil
}

func (r *ebookDiscountRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM ebook_discounts`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *ebookDiscountRepository) CountByEbookID(ctx context.Context, ebookID string) (int64, error) {
	query := `SELECT COUNT(*) FROM ebook_discounts WHERE ebook_id = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, ebookID).Scan(&count)
	return count, err
}

func (r *ebookDiscountRepository) CountActiveDiscounts(ctx context.Context) (int64, error) {
	query := `
		SELECT COUNT(*) FROM ebook_discounts
		WHERE started_at <= ? AND ended_at >= ?
	`
	now := time.Now()
	var count int64
	err := r.db.QueryRowContext(ctx, query, now, now).Scan(&count)
	return count, err
}
