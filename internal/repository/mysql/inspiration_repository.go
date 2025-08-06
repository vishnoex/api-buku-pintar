package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"time"
)

type inspirationRepository struct {
	db *sql.DB
}

func NewInspirationRepository(db *sql.DB) repository.InspirationRepository {
	return &inspirationRepository{db: db}
}

func (r *inspirationRepository) Create(ctx context.Context, inspiration *entity.Inspiration) error {
	query := `INSERT INTO inspirations (id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	inspiration.CreatedAt = now
	inspiration.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		inspiration.ID,
		inspiration.AuthorID,
		inspiration.Title,
		inspiration.Content,
		inspiration.Slug,
		inspiration.Excerpt,
		inspiration.CoverImage,
		inspiration.CategoryID,
		inspiration.ContentStatusID,
		inspiration.ReadingTime,
		inspiration.PublishedAt,
		inspiration.CreatedAt,
		inspiration.UpdatedAt,
	)
	return err
}

func (r *inspirationRepository) GetByID(ctx context.Context, id string) (*entity.Inspiration, error) {
	query := `SELECT id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at 
		FROM inspirations WHERE id = ?`
	
	inspiration := &entity.Inspiration{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&inspiration.ID,
		&inspiration.AuthorID,
		&inspiration.Title,
		&inspiration.Content,
		&inspiration.Slug,
		&inspiration.Excerpt,
		&inspiration.CoverImage,
		&inspiration.CategoryID,
		&inspiration.ContentStatusID,
		&inspiration.ReadingTime,
		&inspiration.PublishedAt,
		&inspiration.CreatedAt,
		&inspiration.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return inspiration, nil
}

func (r *inspirationRepository) GetBySlug(ctx context.Context, slug string) (*entity.Inspiration, error) {
	query := `SELECT id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at 
		FROM inspirations WHERE slug = ?`
	
	inspiration := &entity.Inspiration{}
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&inspiration.ID,
		&inspiration.AuthorID,
		&inspiration.Title,
		&inspiration.Content,
		&inspiration.Slug,
		&inspiration.Excerpt,
		&inspiration.CoverImage,
		&inspiration.CategoryID,
		&inspiration.ContentStatusID,
		&inspiration.ReadingTime,
		&inspiration.PublishedAt,
		&inspiration.CreatedAt,
		&inspiration.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return inspiration, nil
}

func (r *inspirationRepository) Update(ctx context.Context, inspiration *entity.Inspiration) error {
	query := `UPDATE inspirations 
		SET author_id = ?, title = ?, content = ?, slug = ?, excerpt = ?, cover_image = ?, category_id = ?, content_status_id = ?, reading_time = ?, published_at = ?, updated_at = ?
		WHERE id = ?`
	
	inspiration.UpdatedAt = time.Now()
	
	_, err := r.db.ExecContext(ctx, query,
		inspiration.AuthorID,
		inspiration.Title,
		inspiration.Content,
		inspiration.Slug,
		inspiration.Excerpt,
		inspiration.CoverImage,
		inspiration.CategoryID,
		inspiration.ContentStatusID,
		inspiration.ReadingTime,
		inspiration.PublishedAt,
		inspiration.UpdatedAt,
		inspiration.ID,
	)
	return err
}

func (r *inspirationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM inspirations WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *inspirationRepository) List(ctx context.Context, limit, offset int) ([]*entity.Inspiration, error) {
	query := `SELECT id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at 
		FROM inspirations ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inspirations []*entity.Inspiration
	for rows.Next() {
		inspiration := &entity.Inspiration{}
		err := rows.Scan(
			&inspiration.ID,
			&inspiration.AuthorID,
			&inspiration.Title,
			&inspiration.Content,
			&inspiration.Slug,
			&inspiration.Excerpt,
			&inspiration.CoverImage,
			&inspiration.CategoryID,
			&inspiration.ContentStatusID,
			&inspiration.ReadingTime,
			&inspiration.PublishedAt,
			&inspiration.CreatedAt,
			&inspiration.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		inspirations = append(inspirations, inspiration)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inspirations, nil
}

func (r *inspirationRepository) ListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Inspiration, error) {
	query := `SELECT id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at 
		FROM inspirations WHERE author_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inspirations []*entity.Inspiration
	for rows.Next() {
		inspiration := &entity.Inspiration{}
		err := rows.Scan(
			&inspiration.ID,
			&inspiration.AuthorID,
			&inspiration.Title,
			&inspiration.Content,
			&inspiration.Slug,
			&inspiration.Excerpt,
			&inspiration.CoverImage,
			&inspiration.CategoryID,
			&inspiration.ContentStatusID,
			&inspiration.ReadingTime,
			&inspiration.PublishedAt,
			&inspiration.CreatedAt,
			&inspiration.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		inspirations = append(inspirations, inspiration)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inspirations, nil
}

func (r *inspirationRepository) ListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Inspiration, error) {
	query := `SELECT id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at 
		FROM inspirations WHERE category_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inspirations []*entity.Inspiration
	for rows.Next() {
		inspiration := &entity.Inspiration{}
		err := rows.Scan(
			&inspiration.ID,
			&inspiration.AuthorID,
			&inspiration.Title,
			&inspiration.Content,
			&inspiration.Slug,
			&inspiration.Excerpt,
			&inspiration.CoverImage,
			&inspiration.CategoryID,
			&inspiration.ContentStatusID,
			&inspiration.ReadingTime,
			&inspiration.PublishedAt,
			&inspiration.CreatedAt,
			&inspiration.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		inspirations = append(inspirations, inspiration)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inspirations, nil
}

func (r *inspirationRepository) ListPublished(ctx context.Context, limit, offset int) ([]*entity.Inspiration, error) {
	query := `SELECT id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at 
		FROM inspirations WHERE published_at IS NOT NULL ORDER BY published_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inspirations []*entity.Inspiration
	for rows.Next() {
		inspiration := &entity.Inspiration{}
		err := rows.Scan(
			&inspiration.ID,
			&inspiration.AuthorID,
			&inspiration.Title,
			&inspiration.Content,
			&inspiration.Slug,
			&inspiration.Excerpt,
			&inspiration.CoverImage,
			&inspiration.CategoryID,
			&inspiration.ContentStatusID,
			&inspiration.ReadingTime,
			&inspiration.PublishedAt,
			&inspiration.CreatedAt,
			&inspiration.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		inspirations = append(inspirations, inspiration)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inspirations, nil
}

func (r *inspirationRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM inspirations`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *inspirationRepository) CountByAuthor(ctx context.Context, authorID string) (int64, error) {
	query := `SELECT COUNT(*) FROM inspirations WHERE author_id = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, authorID).Scan(&count)
	return count, err
}

func (r *inspirationRepository) CountByCategory(ctx context.Context, categoryID string) (int64, error) {
	query := `SELECT COUNT(*) FROM inspirations WHERE category_id = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&count)
	return count, err
}

func (r *inspirationRepository) CountPublished(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM inspirations WHERE published_at IS NOT NULL`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
