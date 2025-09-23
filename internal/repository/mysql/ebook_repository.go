package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"log"
	"time"
)

type ebookRepository struct {
	db *sql.DB
}

func NewEbookRepository(db *sql.DB) repository.EbookRepository {
	return &ebookRepository{db: db}
}

func (r *ebookRepository) Create(ctx context.Context, ebook *entity.Ebook) error {
	query := `INSERT INTO ebooks (id, author_id, title, synopsis, slug, cover_image, category_id, content_status_id, price, language, duration, filesize, format, page_count, preview_page, url, published_at, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	ebook.CreatedAt = now
	ebook.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		ebook.ID,
		ebook.AuthorID,
		ebook.Title,
		ebook.Synopsis,
		ebook.Slug,
		ebook.CoverImage,
		ebook.CategoryID,
		ebook.ContentStatusID,
		ebook.Price,
		ebook.Language,
		ebook.Duration,
		ebook.Filesize,
		ebook.Format,
		ebook.PageCount,
		ebook.PreviewPage,
		ebook.URL,
		ebook.PublishedAt,
		ebook.CreatedAt,
		ebook.UpdatedAt,
	)
	return err
}

func (r *ebookRepository) GetByID(ctx context.Context, id string) (*entity.Ebook, error) {
	query := `SELECT id, author_id, title, synopsis, slug, cover_image, category_id, content_status_id, price, language, duration, filesize, format, page_count, preview_page, url, published_at, created_at, updated_at 
		FROM ebooks WHERE id = ?`

	ebook := &entity.Ebook{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ebook.ID,
		&ebook.AuthorID,
		&ebook.Title,
		&ebook.Synopsis,
		&ebook.Slug,
		&ebook.CoverImage,
		&ebook.CategoryID,
		&ebook.ContentStatusID,
		&ebook.Price,
		&ebook.Language,
		&ebook.Duration,
		&ebook.Filesize,
		&ebook.Format,
		&ebook.PageCount,
		&ebook.PreviewPage,
		&ebook.URL,
		&ebook.PublishedAt,
		&ebook.CreatedAt,
		&ebook.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return ebook, nil
}

func (r *ebookRepository) GetBySlug(ctx context.Context, slug string) (*entity.Ebook, error) {
	query := `SELECT id, author_id, title, synopsis, slug, cover_image, category_id, content_status_id, price, language, duration, filesize, format, page_count, preview_page, url, published_at, created_at, updated_at 
		FROM ebooks WHERE slug = ?`

	ebook := &entity.Ebook{}
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&ebook.ID,
		&ebook.AuthorID,
		&ebook.Title,
		&ebook.Synopsis,
		&ebook.Slug,
		&ebook.CoverImage,
		&ebook.CategoryID,
		&ebook.ContentStatusID,
		&ebook.Price,
		&ebook.Language,
		&ebook.Duration,
		&ebook.Filesize,
		&ebook.Format,
		&ebook.PageCount,
		&ebook.PreviewPage,
		&ebook.URL,
		&ebook.PublishedAt,
		&ebook.CreatedAt,
		&ebook.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return ebook, nil
}

func (r *ebookRepository) Update(ctx context.Context, ebook *entity.Ebook) error {
	query := `UPDATE ebooks 
		SET author_id = ?, title = ?, synopsis = ?, slug = ?, cover_image = ?, category_id = ?, content_status_id = ?, price = ?, language = ?, duration = ?, filesize = ?, format = ?, page_count = ?, preview_page = ?, url = ?, published_at = ?, updated_at = ?
		WHERE id = ?`

	ebook.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		ebook.AuthorID,
		ebook.Title,
		ebook.Synopsis,
		ebook.Slug,
		ebook.CoverImage,
		ebook.CategoryID,
		ebook.ContentStatusID,
		ebook.Price,
		ebook.Language,
		ebook.Duration,
		ebook.Filesize,
		ebook.Format,
		ebook.PageCount,
		ebook.PreviewPage,
		ebook.URL,
		ebook.PublishedAt,
		ebook.UpdatedAt,
		ebook.ID,
	)
	return err
}

func (r *ebookRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM ebooks WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ebookRepository) List(ctx context.Context, limit, offset int) ([]*entity.EbookList, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	query := `SELECT
				e.id, e.title, e.slug, e.cover_image, e.price, ed.discount_price AS discount
			FROM ebooks e
			LEFT JOIN content_statuses cs ON cs.id = e.content_status_id
			LEFT JOIN
				ebook_discounts ed
					ON ed.ebook_id = e.id AND ed.started_at <= ? AND ed.ended_at >= ? 
			WHERE e.published_at IS NOT NULL AND e.published_at <= ? AND cs.name = "published"
			ORDER BY e.published_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, now, now, now, limit, offset)
	log.Println(now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ebooks []*entity.EbookList
	for rows.Next() {
		ebook := &entity.EbookList{}
		err = rows.Scan(
			&ebook.ID,
			&ebook.Title,
			&ebook.Slug,
			&ebook.CoverImage,
			&ebook.Price,
			&ebook.Discount,
		)
		if err != nil {
			return nil, err
		}
		ebooks = append(ebooks, ebook)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ebooks, nil
}

func (r *ebookRepository) ListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Ebook, error) {
	query := `SELECT id, author_id, title, synopsis, slug, cover_image, category_id, content_status_id, price, language, duration, filesize, format, page_count, preview_page, url, published_at, created_at, updated_at 
		FROM ebooks WHERE category_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ebooks []*entity.Ebook
	for rows.Next() {
		ebook := &entity.Ebook{}
		err = rows.Scan(
			&ebook.ID,
			&ebook.AuthorID,
			&ebook.Title,
			&ebook.Synopsis,
			&ebook.Slug,
			&ebook.CoverImage,
			&ebook.CategoryID,
			&ebook.ContentStatusID,
			&ebook.Price,
			&ebook.Language,
			&ebook.Duration,
			&ebook.Filesize,
			&ebook.Format,
			&ebook.PageCount,
			&ebook.PreviewPage,
			&ebook.URL,
			&ebook.PublishedAt,
			&ebook.CreatedAt,
			&ebook.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ebooks = append(ebooks, ebook)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ebooks, nil
}

func (r *ebookRepository) ListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Ebook, error) {
	query := `SELECT id, author_id, title, synopsis, slug, cover_image, category_id, content_status_id, price, language, duration, filesize, format, page_count, preview_page, url, published_at, created_at, updated_at 
		FROM ebooks WHERE author_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ebooks []*entity.Ebook
	for rows.Next() {
		ebook := &entity.Ebook{}
		err = rows.Scan(
			&ebook.ID,
			&ebook.AuthorID,
			&ebook.Title,
			&ebook.Synopsis,
			&ebook.Slug,
			&ebook.CoverImage,
			&ebook.CategoryID,
			&ebook.ContentStatusID,
			&ebook.Price,
			&ebook.Language,
			&ebook.Duration,
			&ebook.Filesize,
			&ebook.Format,
			&ebook.PageCount,
			&ebook.PreviewPage,
			&ebook.URL,
			&ebook.PublishedAt,
			&ebook.CreatedAt,
			&ebook.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ebooks = append(ebooks, ebook)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ebooks, nil
}

func (r *ebookRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM ebooks`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *ebookRepository) CountByCategory(ctx context.Context, categoryID string) (int64, error) {
	query := `SELECT COUNT(*) FROM ebooks WHERE category_id = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&count)
	return count, err
}

func (r *ebookRepository) CountByAuthor(ctx context.Context, authorID string) (int64, error) {
	query := `SELECT COUNT(*) FROM ebooks WHERE author_id = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, authorID).Scan(&count)
	return count, err
}
