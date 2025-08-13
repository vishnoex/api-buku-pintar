package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"time"
)

type seoMetaRepository struct {
	db *sql.DB
}

func NewSeoMetaRepository(db *sql.DB) repository.SeoMetaRepository {
	return &seoMetaRepository{db: db}
}

func (r *seoMetaRepository) Create(ctx context.Context, seoMeta *entity.SeoMeta) error {
	query := `INSERT INTO seo_metadatas (id, title, description, keywords, entity, entity_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	seoMeta.CreatedAt = now
	seoMeta.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		seoMeta.ID,
		seoMeta.Title,
		seoMeta.Description,
		seoMeta.Keywords,
		seoMeta.Entity,
		seoMeta.EntityID,
		seoMeta.CreatedAt,
		seoMeta.UpdatedAt,
	)
	return err
}

func (r *seoMetaRepository) GetByID(ctx context.Context, id string) (*entity.SeoMeta, error) {
	query := `SELECT id, title, description, keywords, entity, entity_id, created_at, updated_at 
		FROM seo_metadatas WHERE id = ?`

	seoMeta := &entity.SeoMeta{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&seoMeta.ID,
		&seoMeta.Title,
		&seoMeta.Description,
		&seoMeta.Keywords,
		&seoMeta.Entity,
		&seoMeta.EntityID,
		&seoMeta.CreatedAt,
		&seoMeta.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return seoMeta, nil
}

func (r *seoMetaRepository) GetByEntity(ctx context.Context, entityType entity.SeoEntityType, entityID string) (*entity.SeoMeta, error) {
	query := `SELECT id, title, description, keywords, entity, entity_id, created_at, updated_at 
		FROM seo_metadatas WHERE entity = ? AND entity_id = ?`

	seoMeta := &entity.SeoMeta{}
	err := r.db.QueryRowContext(ctx, query, entityType, entityID).Scan(
		&seoMeta.ID,
		&seoMeta.Title,
		&seoMeta.Description,
		&seoMeta.Keywords,
		&seoMeta.Entity,
		&seoMeta.EntityID,
		&seoMeta.CreatedAt,
		&seoMeta.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return seoMeta, nil
}

func (r *seoMetaRepository) Update(ctx context.Context, seoMeta *entity.SeoMeta) error {
	query := `UPDATE seo_metadatas 
		SET title = ?, description = ?, keywords = ?, entity = ?, entity_id = ?, updated_at = ?
		WHERE id = ?`

	seoMeta.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		seoMeta.Title,
		seoMeta.Description,
		seoMeta.Keywords,
		seoMeta.Entity,
		seoMeta.EntityID,
		seoMeta.UpdatedAt,
		seoMeta.ID,
	)
	return err
}

func (r *seoMetaRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM seo_metadatas WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *seoMetaRepository) ListByEntity(ctx context.Context, entityType entity.SeoEntityType, limit, offset int) ([]*entity.SeoMeta, error) {
	query := `SELECT id, title, description, keywords, entity, entity_id, created_at, updated_at 
		FROM seo_metadatas WHERE entity = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, entityType, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seoMetas []*entity.SeoMeta
	for rows.Next() {
		seoMeta := &entity.SeoMeta{}
		err = rows.Scan(
			&seoMeta.ID,
			&seoMeta.Title,
			&seoMeta.Description,
			&seoMeta.Keywords,
			&seoMeta.Entity,
			&seoMeta.EntityID,
			&seoMeta.CreatedAt,
			&seoMeta.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		seoMetas = append(seoMetas, seoMeta)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return seoMetas, nil
}

func (r *seoMetaRepository) CountByEntity(ctx context.Context, entityType entity.SeoEntityType) (int64, error) {
	query := `SELECT COUNT(*) FROM seo_metadatas WHERE entity = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, entityType).Scan(&count)
	return count, err
}
