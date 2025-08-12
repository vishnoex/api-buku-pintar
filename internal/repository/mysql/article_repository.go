package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"time"
)

type articleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) repository.ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(ctx context.Context, article *entity.Article) error {
	query := `INSERT INTO articles (id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	article.CreatedAt = now
	article.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		article.ID,
		article.AuthorID,
		article.Title,
		article.Content,
		article.Slug,
		article.Excerpt,
		article.CoverImage,
		article.CategoryID,
		article.ContentStatusID,
		article.ReadingTime,
		article.PublishedAt,
		article.CreatedAt,
		article.UpdatedAt,
	)
	return err
}

func (r *articleRepository) GetByID(ctx context.Context, id string) (*entity.Article, error) {
	query := `SELECT id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at 
		FROM articles WHERE id = ?`

	article := &entity.Article{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&article.ID,
		&article.AuthorID,
		&article.Title,
		&article.Content,
		&article.Slug,
		&article.Excerpt,
		&article.CoverImage,
		&article.CategoryID,
		&article.ContentStatusID,
		&article.ReadingTime,
		&article.PublishedAt,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return article, nil
}

func (r *articleRepository) GetBySlug(ctx context.Context, slug string) (*entity.Article, error) {
	query := `SELECT id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at 
		FROM articles WHERE slug = ?`

	article := &entity.Article{}
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&article.ID,
		&article.AuthorID,
		&article.Title,
		&article.Content,
		&article.Slug,
		&article.Excerpt,
		&article.CoverImage,
		&article.CategoryID,
		&article.ContentStatusID,
		&article.ReadingTime,
		&article.PublishedAt,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return article, nil
}

func (r *articleRepository) Update(ctx context.Context, article *entity.Article) error {
	query := `UPDATE articles 
		SET author_id = ?, title = ?, content = ?, slug = ?, excerpt = ?, cover_image = ?, category_id = ?, content_status_id = ?, reading_time = ?, published_at = ?, updated_at = ?
		WHERE id = ?`

	article.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		article.AuthorID,
		article.Title,
		article.Content,
		article.Slug,
		article.Excerpt,
		article.CoverImage,
		article.CategoryID,
		article.ContentStatusID,
		article.ReadingTime,
		article.PublishedAt,
		article.UpdatedAt,
		article.ID,
	)
	return err
}

func (r *articleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM articles WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *articleRepository) List(ctx context.Context, limit, offset int) ([]*entity.Article, error) {
	query := `SELECT id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at 
		FROM articles ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*entity.Article
	for rows.Next() {
		article := &entity.Article{}
		err = rows.Scan(
			&article.ID,
			&article.AuthorID,
			&article.Title,
			&article.Content,
			&article.Slug,
			&article.Excerpt,
			&article.CoverImage,
			&article.CategoryID,
			&article.ContentStatusID,
			&article.ReadingTime,
			&article.PublishedAt,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *articleRepository) ListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Article, error) {
	query := `SELECT id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at 
		FROM articles WHERE author_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*entity.Article
	for rows.Next() {
		article := &entity.Article{}
		err = rows.Scan(
			&article.ID,
			&article.AuthorID,
			&article.Title,
			&article.Content,
			&article.Slug,
			&article.Excerpt,
			&article.CoverImage,
			&article.CategoryID,
			&article.ContentStatusID,
			&article.ReadingTime,
			&article.PublishedAt,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *articleRepository) ListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Article, error) {
	query := `SELECT id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at 
		FROM articles WHERE category_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*entity.Article
	for rows.Next() {
		article := &entity.Article{}
		err = rows.Scan(
			&article.ID,
			&article.AuthorID,
			&article.Title,
			&article.Content,
			&article.Slug,
			&article.Excerpt,
			&article.CoverImage,
			&article.CategoryID,
			&article.ContentStatusID,
			&article.ReadingTime,
			&article.PublishedAt,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *articleRepository) ListPublished(ctx context.Context, limit, offset int) ([]*entity.Article, error) {
	query := `SELECT id, author_id, title, content, slug, excerpt, cover_image, category_id, content_status_id, reading_time, published_at, created_at, updated_at 
		FROM articles WHERE published_at IS NOT NULL ORDER BY published_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*entity.Article
	for rows.Next() {
		article := &entity.Article{}
		err = rows.Scan(
			&article.ID,
			&article.AuthorID,
			&article.Title,
			&article.Content,
			&article.Slug,
			&article.Excerpt,
			&article.CoverImage,
			&article.CategoryID,
			&article.ContentStatusID,
			&article.ReadingTime,
			&article.PublishedAt,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *articleRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM articles`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *articleRepository) CountByAuthor(ctx context.Context, authorID string) (int64, error) {
	query := `SELECT COUNT(*) FROM articles WHERE author_id = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, authorID).Scan(&count)
	return count, err
}

func (r *articleRepository) CountByCategory(ctx context.Context, categoryID string) (int64, error) {
	query := `SELECT COUNT(*) FROM articles WHERE category_id = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&count)
	return count, err
}

func (r *articleRepository) CountPublished(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM articles WHERE published_at IS NOT NULL`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
