package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"time"
)

type bannerRepository struct {
	db *sql.DB
}

func NewBannerRepository(db *sql.DB) repository.BannerRepository {
	return &bannerRepository{db: db}
}

func (r *bannerRepository) Create(ctx context.Context, banner *entity.Banner) error {
	query := `INSERT INTO banners (id, title, image_url, link, cta_label, background_color, is_active, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	banner.CreatedAt = now
	banner.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		banner.ID,
		banner.Title,
		banner.ImageURL,
		banner.Link,
		banner.CTALabel,
		banner.BackgroundColor,
		banner.IsActive,
		banner.CreatedAt,
		banner.UpdatedAt,
	)
	return err
}

func (r *bannerRepository) GetByID(ctx context.Context, id string) (*entity.Banner, error) {
	query := `SELECT id, title, image_url, link, cta_label, background_color, is_active, created_at, updated_at 
		FROM banners WHERE id = ?`
	
	banner := &entity.Banner{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&banner.ID,
		&banner.Title,
		&banner.ImageURL,
		&banner.Link,
		&banner.CTALabel,
		&banner.BackgroundColor,
		&banner.IsActive,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return banner, nil
}

func (r *bannerRepository) Update(ctx context.Context, banner *entity.Banner) error {
	query := `UPDATE banners 
		SET title = ?, image_url = ?, link = ?, cta_label = ?, background_color = ?, is_active = ?, updated_at = ?
		WHERE id = ?`
	
	banner.UpdatedAt = time.Now()
	
	_, err := r.db.ExecContext(ctx, query,
		banner.Title,
		banner.ImageURL,
		banner.Link,
		banner.CTALabel,
		banner.BackgroundColor,
		banner.IsActive,
		banner.UpdatedAt,
		banner.ID,
	)
	return err
}

func (r *bannerRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM banners WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *bannerRepository) List(ctx context.Context, limit, offset int) ([]*entity.Banner, error) {
	query := `SELECT id, title, image_url, link, cta_label, background_color, is_active, created_at, updated_at 
		FROM banners ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []*entity.Banner
	for rows.Next() {
		banner := &entity.Banner{}
		err := rows.Scan(
			&banner.ID,
			&banner.Title,
			&banner.ImageURL,
			&banner.Link,
			&banner.CTALabel,
			&banner.BackgroundColor,
			&banner.IsActive,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

func (r *bannerRepository) ListActive(ctx context.Context, limit, offset int) ([]*entity.Banner, error) {
	query := `SELECT id, title, image_url, link, cta_label, background_color, is_active, created_at, updated_at 
		FROM banners WHERE is_active = true ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []*entity.Banner
	for rows.Next() {
		banner := &entity.Banner{}
		err := rows.Scan(
			&banner.ID,
			&banner.Title,
			&banner.ImageURL,
			&banner.Link,
			&banner.CTALabel,
			&banner.BackgroundColor,
			&banner.IsActive,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

func (r *bannerRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM banners`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *bannerRepository) CountActive(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM banners WHERE is_active = true`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
