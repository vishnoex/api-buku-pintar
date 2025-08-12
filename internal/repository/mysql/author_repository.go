package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"time"
)

type authorRepository struct {
	db *sql.DB
}

func NewAuthorRepository(db *sql.DB) repository.AuthorRepository {
	return &authorRepository{db: db}
}

func (r *authorRepository) Create(ctx context.Context, author *entity.Author) error {
	query := `INSERT INTO authors (id, name, avatar, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	author.CreatedAt = now
	author.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		author.ID,
		author.Name,
		author.Avatar,
		author.CreatedAt,
		author.UpdatedAt,
	)
	return err
}

func (r *authorRepository) GetByID(ctx context.Context, id string) (*entity.Author, error) {
	query := `SELECT id, name, avatar, created_at, updated_at 
		FROM authors WHERE id = ?`

	author := &entity.Author{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&author.ID,
		&author.Name,
		&author.Avatar,
		&author.CreatedAt,
		&author.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return author, nil
}

func (r *authorRepository) GetByName(ctx context.Context, name string) (*entity.Author, error) {
	query := `SELECT id, name, avatar, created_at, updated_at 
		FROM authors WHERE name = ?`

	author := &entity.Author{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&author.ID,
		&author.Name,
		&author.Avatar,
		&author.CreatedAt,
		&author.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return author, nil
}

func (r *authorRepository) Update(ctx context.Context, author *entity.Author) error {
	query := `UPDATE authors 
		SET name = ?, avatar = ?, updated_at = ?
		WHERE id = ?`

	author.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		author.Name,
		author.Avatar,
		author.UpdatedAt,
		author.ID,
	)
	return err
}

func (r *authorRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM authors WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *authorRepository) List(ctx context.Context, limit, offset int) ([]*entity.Author, error) {
	query := `SELECT id, name, avatar, created_at, updated_at 
		FROM authors ORDER BY name ASC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []*entity.Author
	for rows.Next() {
		author := &entity.Author{}
		err = rows.Scan(
			&author.ID,
			&author.Name,
			&author.Avatar,
			&author.CreatedAt,
			&author.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (r *authorRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM authors`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
