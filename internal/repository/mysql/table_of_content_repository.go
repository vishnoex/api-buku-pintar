package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"time"
)

type tableOfContentRepository struct {
	db *sql.DB
}

func NewTableOfContentRepository(db *sql.DB) repository.TableOfContentRepository {
	return &tableOfContentRepository{db: db}
}

func (r *tableOfContentRepository) Create(ctx context.Context, toc *entity.TableOfContent) error {
	query := `INSERT INTO table_of_contents (id, ebook_id, title, page_number, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	toc.CreatedAt = now
	toc.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		toc.ID,
		toc.EbookID,
		toc.Title,
		toc.PageNumber,
		toc.CreatedAt,
		toc.UpdatedAt,
	)
	return err
}

func (r *tableOfContentRepository) GetByID(ctx context.Context, id string) (*entity.TableOfContent, error) {
	query := `SELECT id, ebook_id, title, page_number, created_at, updated_at 
		FROM table_of_contents WHERE id = ?`
	
	toc := &entity.TableOfContent{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&toc.ID,
		&toc.EbookID,
		&toc.Title,
		&toc.PageNumber,
		&toc.CreatedAt,
		&toc.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return toc, nil
}

func (r *tableOfContentRepository) Update(ctx context.Context, toc *entity.TableOfContent) error {
	query := `UPDATE table_of_contents 
		SET ebook_id = ?, title = ?, page_number = ?, updated_at = ?
		WHERE id = ?`
	
	toc.UpdatedAt = time.Now()
	
	_, err := r.db.ExecContext(ctx, query,
		toc.EbookID,
		toc.Title,
		toc.PageNumber,
		toc.UpdatedAt,
		toc.ID,
	)
	return err
}

func (r *tableOfContentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM table_of_contents WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *tableOfContentRepository) ListByEbook(ctx context.Context, ebookID string) ([]*entity.TableOfContent, error) {
	query := `SELECT id, ebook_id, title, page_number, created_at, updated_at 
		FROM table_of_contents WHERE ebook_id = ? ORDER BY page_number ASC`

	rows, err := r.db.QueryContext(ctx, query, ebookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tableOfContents []*entity.TableOfContent
	for rows.Next() {
		toc := &entity.TableOfContent{}
		err := rows.Scan(
			&toc.ID,
			&toc.EbookID,
			&toc.Title,
			&toc.PageNumber,
			&toc.CreatedAt,
			&toc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tableOfContents = append(tableOfContents, toc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tableOfContents, nil
}

func (r *tableOfContentRepository) DeleteByEbook(ctx context.Context, ebookID string) error {
	query := `DELETE FROM table_of_contents WHERE ebook_id = ?`
	_, err := r.db.ExecContext(ctx, query, ebookID)
	return err
}

func (r *tableOfContentRepository) CountByEbook(ctx context.Context, ebookID string) (int64, error) {
	query := `SELECT COUNT(*) FROM table_of_contents WHERE ebook_id = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, ebookID).Scan(&count)
	return count, err
}
