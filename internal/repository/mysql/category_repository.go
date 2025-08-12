package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"time"
)

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) repository.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) error {
	query := `INSERT INTO categories (id, name, description, icon_link, parent_id, order_number, is_active, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		category.ID,
		category.Name,
		category.Description,
		category.IconLink,
		category.ParentID,
		category.OrderNumber,
		category.IsActive,
		category.CreatedAt,
		category.UpdatedAt,
	)
	return err
}

func (r *categoryRepository) GetByID(ctx context.Context, id string) (*entity.Category, error) {
	query := `SELECT id, name, description, icon_link, parent_id, order_number, is_active, created_at, updated_at 
		FROM categories WHERE id = ?`

	category := &entity.Category{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.IconLink,
		&category.ParentID,
		&category.OrderNumber,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) GetByName(ctx context.Context, name string) (*entity.Category, error) {
	query := `SELECT id, name, description, icon_link, parent_id, order_number, is_active, created_at, updated_at 
		FROM categories WHERE name = ?`

	category := &entity.Category{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.IconLink,
		&category.ParentID,
		&category.OrderNumber,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) error {
	query := `UPDATE categories 
		SET name = ?, description = ?, icon_link = ?, parent_id = ?, order_number = ?, is_active = ?, updated_at = ?
		WHERE id = ?`

	category.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		category.Name,
		category.Description,
		category.IconLink,
		category.ParentID,
		category.OrderNumber,
		category.IsActive,
		category.UpdatedAt,
		category.ID,
	)
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM categories WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *categoryRepository) List(ctx context.Context, limit, offset int) ([]*entity.Category, error) {
	query := `SELECT id, name, description, icon_link, parent_id, order_number, is_active, created_at, updated_at 
		FROM categories ORDER BY order_number ASC, created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		category := &entity.Category{}
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.IconLink,
			&category.ParentID,
			&category.OrderNumber,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) ListActive(ctx context.Context, limit, offset int) ([]*entity.Category, error) {
	query := `SELECT id, name, description, icon_link, parent_id, order_number, is_active, created_at, updated_at 
		FROM categories WHERE is_active = true ORDER BY order_number ASC, created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		category := &entity.Category{}
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.IconLink,
			&category.ParentID,
			&category.OrderNumber,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) ListByParent(ctx context.Context, parentID string, limit, offset int) ([]*entity.Category, error) {
	query := `SELECT id, name, description, icon_link, parent_id, order_number, is_active, created_at, updated_at 
		FROM categories WHERE parent_id = ? ORDER BY order_number ASC, created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, parentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		category := &entity.Category{}
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.IconLink,
			&category.ParentID,
			&category.OrderNumber,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM categories`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *categoryRepository) CountActive(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM categories WHERE is_active = true`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *categoryRepository) CountByParent(ctx context.Context, parentID string) (int64, error) {
	query := `SELECT COUNT(*) FROM categories WHERE parent_id = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, parentID).Scan(&count)
	return count, err
}
