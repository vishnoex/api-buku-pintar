package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"time"
)

type userLoginRepository struct {
	db *sql.DB
}

func NewUserLoginRepository(db *sql.DB) repository.UserLoginRepository {
	return &userLoginRepository{db: db}
}

func (r *userLoginRepository) Create(ctx context.Context, userLogin *entity.UserLogin) error {
	query := `INSERT INTO user_logins (id, user_id, login_provider_id, status, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	userLogin.CreatedAt = now
	userLogin.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		userLogin.ID,
		userLogin.UserID,
		userLogin.LoginProviderID,
		userLogin.Status,
		userLogin.CreatedAt,
		userLogin.UpdatedAt,
	)
	return err
}

func (r *userLoginRepository) GetByID(ctx context.Context, id string) (*entity.UserLogin, error) {
	query := `SELECT id, user_id, login_provider_id, status, created_at, updated_at 
		FROM user_logins WHERE id = ?`
	
	userLogin := &entity.UserLogin{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&userLogin.ID,
		&userLogin.UserID,
		&userLogin.LoginProviderID,
		&userLogin.Status,
		&userLogin.CreatedAt,
		&userLogin.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return userLogin, nil
}

func (r *userLoginRepository) GetByUserAndProvider(ctx context.Context, userID, loginProviderID string) (*entity.UserLogin, error) {
	query := `SELECT id, user_id, login_provider_id, status, created_at, updated_at 
		FROM user_logins WHERE user_id = ? AND login_provider_id = ?`
	
	userLogin := &entity.UserLogin{}
	err := r.db.QueryRowContext(ctx, query, userID, loginProviderID).Scan(
		&userLogin.ID,
		&userLogin.UserID,
		&userLogin.LoginProviderID,
		&userLogin.Status,
		&userLogin.CreatedAt,
		&userLogin.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return userLogin, nil
}

func (r *userLoginRepository) Update(ctx context.Context, userLogin *entity.UserLogin) error {
	query := `UPDATE user_logins 
		SET user_id = ?, login_provider_id = ?, status = ?, updated_at = ?
		WHERE id = ?`
	
	userLogin.UpdatedAt = time.Now()
	
	_, err := r.db.ExecContext(ctx, query,
		userLogin.UserID,
		userLogin.LoginProviderID,
		userLogin.Status,
		userLogin.UpdatedAt,
		userLogin.ID,
	)
	return err
}

func (r *userLoginRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM user_logins WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userLoginRepository) ListByUser(ctx context.Context, userID string) ([]*entity.UserLogin, error) {
	query := `SELECT id, user_id, login_provider_id, status, created_at, updated_at 
		FROM user_logins WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userLogins []*entity.UserLogin
	for rows.Next() {
		userLogin := &entity.UserLogin{}
		err := rows.Scan(
			&userLogin.ID,
			&userLogin.UserID,
			&userLogin.LoginProviderID,
			&userLogin.Status,
			&userLogin.CreatedAt,
			&userLogin.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		userLogins = append(userLogins, userLogin)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userLogins, nil
}

func (r *userLoginRepository) ListActiveByUser(ctx context.Context, userID string) ([]*entity.UserLogin, error) {
	query := `SELECT id, user_id, login_provider_id, status, created_at, updated_at 
		FROM user_logins WHERE user_id = ? AND status = 'active' ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userLogins []*entity.UserLogin
	for rows.Next() {
		userLogin := &entity.UserLogin{}
		err := rows.Scan(
			&userLogin.ID,
			&userLogin.UserID,
			&userLogin.LoginProviderID,
			&userLogin.Status,
			&userLogin.CreatedAt,
			&userLogin.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		userLogins = append(userLogins, userLogin)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userLogins, nil
}

func (r *userLoginRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM user_logins WHERE user_id = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *userLoginRepository) CountActiveByUser(ctx context.Context, userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM user_logins WHERE user_id = ? AND status = 'active'`
	var count int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}
