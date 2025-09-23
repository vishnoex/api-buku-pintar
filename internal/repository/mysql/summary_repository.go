package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SummaryRepositoryImpl struct {
	db *sql.DB
}

func NewSummaryRepositoryImpl(db *sql.DB) repository.SummaryRepository {
	return &SummaryRepositoryImpl{
		db: db,
	}
}

func (r *SummaryRepositoryImpl) CreateSummary(ctx context.Context, summary *entity.EbookSummary) error {
	// Generate UUID if not provided
	if summary.ID == "" {
		summary.ID = uuid.New().String()
	}

	query := `
		INSERT INTO ebook_summaries (id, ebook_id, description, url, audio_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	summary.CreatedAt = now
	summary.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		summary.ID,
		summary.EbookID,
		summary.Description,
		summary.URL,
		summary.AudioURL,
		summary.CreatedAt,
		summary.UpdatedAt,
	)

	return err
}

func (r *SummaryRepositoryImpl) GetSummaryByID(ctx context.Context, id string) (*entity.EbookSummary, error) {
	query := `
		SELECT id, ebook_id, description, url, audio_url, created_at, updated_at
		FROM ebook_summaries
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var summary entity.EbookSummary
	err := row.Scan(
		&summary.ID,
		&summary.EbookID,
		&summary.Description,
		&summary.URL,
		&summary.AudioURL,
		&summary.CreatedAt,
		&summary.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &summary, nil
}

func (r *SummaryRepositoryImpl) UpdateSummary(ctx context.Context, summary *entity.EbookSummary) error {
	query := `
		UPDATE ebook_summaries
		SET ebook_id = ?, description = ?, url = ?, audio_url = ?, updated_at = ?
		WHERE id = ?
	`

	summary.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		summary.EbookID,
		summary.Description,
		summary.URL,
		summary.AudioURL,
		summary.UpdatedAt,
		summary.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("summary not found")
	}

	return nil
}

func (r *SummaryRepositoryImpl) DeleteSummary(ctx context.Context, id string) error {
	query := `DELETE FROM ebook_summaries WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("summary not found")
	}

	return nil
}

func (r *SummaryRepositoryImpl) ListSummaries(ctx context.Context, limit, offset int) ([]*entity.EbookSummaryList, error) {
	query := `
		SELECT 
			es.id, 
			es.ebook_id, 
			e.title as ebook_title,
			e.slug,
			es.description, 
			es.url, 
			es.audio_url,
			e.duration,
			es.created_at, 
			es.updated_at
		FROM ebook_summaries es
		JOIN ebooks e ON es.ebook_id = e.id
		ORDER BY es.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []*entity.EbookSummaryList
	for rows.Next() {
		var summary entity.EbookSummaryList
		err := rows.Scan(
			&summary.ID,
			&summary.EbookID,
			&summary.EbookTitle,
			&summary.Slug,
			&summary.Description,
			&summary.URL,
			&summary.AudioURL,
			&summary.Duration,
			&summary.CreatedAt,
			&summary.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, &summary)
	}

	return summaries, nil
}

func (r *SummaryRepositoryImpl) GetSummariesByEbookID(ctx context.Context, ebookID string, limit, offset int) ([]*entity.EbookSummary, error) {
	query := `
		SELECT id, ebook_id, description, url, audio_url, created_at, updated_at
		FROM ebook_summaries
		WHERE ebook_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, ebookID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []*entity.EbookSummary
	for rows.Next() {
		var summary entity.EbookSummary
		err := rows.Scan(
			&summary.ID,
			&summary.EbookID,
			&summary.Description,
			&summary.URL,
			&summary.AudioURL,
			&summary.CreatedAt,
			&summary.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, &summary)
	}

	return summaries, nil
}

func (r *SummaryRepositoryImpl) CountSummaries(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM ebook_summaries`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *SummaryRepositoryImpl) CountSummariesByEbookID(ctx context.Context, ebookID string) (int64, error) {
	query := `SELECT COUNT(*) FROM ebook_summaries WHERE ebook_id = ?`

	var count int64
	err := r.db.QueryRowContext(ctx, query, ebookID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
