package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"time"
)

type contentStatusRepository struct {
	db *sql.DB
}

func NewContentStatusRepository(db *sql.DB) repository.ContentStatusRepository {
	return &contentStatusRepository{db: db}
}

func (r *contentStatusRepository) Create(ctx context.Context, contentStatus *entity.ContentStatus) error {
	query := `INSERT INTO content_statuses (id, name, created_at, updated_at) 
		VALUES (?, ?, ?, ?)`

	now := time.Now()
	contentStatus.CreatedAt = now
	contentStatus.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		contentStatus.ID,
		contentStatus.Name,
		contentStatus.CreatedAt,
		contentStatus.UpdatedAt,
	)
	return err
}

func (r *contentStatusRepository) GetByID(ctx context.Context, id string) (*entity.ContentStatus, error) {
	query := `SELECT id, name, created_at, updated_at 
		FROM content_statuses WHERE id = ?`

	contentStatus := &entity.ContentStatus{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&contentStatus.ID,
		&contentStatus.Name,
		&contentStatus.CreatedAt,
		&contentStatus.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return contentStatus, nil
}

func (r *contentStatusRepository) GetByName(ctx context.Context, name string) (*entity.ContentStatus, error) {
	query := `SELECT id, name, created_at, updated_at 
		FROM content_statuses WHERE name = ?`

	contentStatus := &entity.ContentStatus{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&contentStatus.ID,
		&contentStatus.Name,
		&contentStatus.CreatedAt,
		&contentStatus.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return contentStatus, nil
}

func (r *contentStatusRepository) Update(ctx context.Context, contentStatus *entity.ContentStatus) error {
	query := `UPDATE content_statuses 
		SET name = ?, updated_at = ?
		WHERE id = ?`

	contentStatus.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		contentStatus.Name,
		contentStatus.UpdatedAt,
		contentStatus.ID,
	)
	return err
}

func (r *contentStatusRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM content_statuses WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *contentStatusRepository) List(ctx context.Context, limit, offset int) ([]*entity.ContentStatus, error) {
	query := `SELECT id, name, created_at, updated_at 
		FROM content_statuses ORDER BY name ASC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contentStatuses []*entity.ContentStatus
	for rows.Next() {
		contentStatus := &entity.ContentStatus{}
		err = rows.Scan(
			&contentStatus.ID,
			&contentStatus.Name,
			&contentStatus.CreatedAt,
			&contentStatus.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		contentStatuses = append(contentStatuses, contentStatus)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contentStatuses, nil
}

func (r *contentStatusRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM content_statuses`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
