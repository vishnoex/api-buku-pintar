package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"time"
)

type oauthTokenRepository struct {
	db *sql.DB
}

// NewOAuthTokenRepository creates a new instance of OAuthTokenRepository
func NewOAuthTokenRepository(db *sql.DB) repository.OAuthTokenRepository {
	return &oauthTokenRepository{db: db}
}

// Create creates a new OAuth token
func (r *oauthTokenRepository) Create(ctx context.Context, token *entity.OAuthToken) error {
	query := `INSERT INTO oauth_tokens (id, user_id, provider, access_token, refresh_token, token_type, expires_at, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	token.CreatedAt = now
	token.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.Provider,
		token.AccessToken,
		token.RefreshToken,
		token.TokenType,
		token.ExpiresAt,
		token.CreatedAt,
		token.UpdatedAt,
	)
	return err
}

// GetByID retrieves an OAuth token by its ID
func (r *oauthTokenRepository) GetByID(ctx context.Context, id string) (*entity.OAuthToken, error) {
	query := `SELECT id, user_id, provider, access_token, refresh_token, token_type, expires_at, created_at, updated_at 
		FROM oauth_tokens WHERE id = ?`

	token := &entity.OAuthToken{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&token.ID,
		&token.UserID,
		&token.Provider,
		&token.AccessToken,
		&token.RefreshToken,
		&token.TokenType,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return token, nil
}

// Update updates an existing OAuth token
func (r *oauthTokenRepository) Update(ctx context.Context, token *entity.OAuthToken) error {
	query := `UPDATE oauth_tokens 
		SET access_token = ?, refresh_token = ?, token_type = ?, expires_at = ?, updated_at = ? 
		WHERE id = ?`

	token.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		token.AccessToken,
		token.RefreshToken,
		token.TokenType,
		token.ExpiresAt,
		token.UpdatedAt,
		token.ID,
	)
	return err
}

// Delete deletes an OAuth token by ID
func (r *oauthTokenRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM oauth_tokens WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByUserIDAndProvider retrieves an OAuth token by user ID and provider
func (r *oauthTokenRepository) GetByUserIDAndProvider(ctx context.Context, userID string, provider entity.OAuthProvider) (*entity.OAuthToken, error) {
	query := `SELECT id, user_id, provider, access_token, refresh_token, token_type, expires_at, created_at, updated_at 
		FROM oauth_tokens 
		WHERE user_id = ? AND provider = ?
		ORDER BY created_at DESC
		LIMIT 1`

	token := &entity.OAuthToken{}
	err := r.db.QueryRowContext(ctx, query, userID, provider).Scan(
		&token.ID,
		&token.UserID,
		&token.Provider,
		&token.AccessToken,
		&token.RefreshToken,
		&token.TokenType,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return token, nil
}

// GetByUserID retrieves all OAuth tokens for a user
func (r *oauthTokenRepository) GetByUserID(ctx context.Context, userID string) ([]*entity.OAuthToken, error) {
	query := `SELECT id, user_id, provider, access_token, refresh_token, token_type, expires_at, created_at, updated_at 
		FROM oauth_tokens 
		WHERE user_id = ?
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*entity.OAuthToken
	for rows.Next() {
		token := &entity.OAuthToken{}
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Provider,
			&token.AccessToken,
			&token.RefreshToken,
			&token.TokenType,
			&token.ExpiresAt,
			&token.CreatedAt,
			&token.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

// GetByProvider retrieves OAuth tokens by provider with pagination
func (r *oauthTokenRepository) GetByProvider(ctx context.Context, provider entity.OAuthProvider, limit, offset int) ([]*entity.OAuthToken, error) {
	query := `SELECT id, user_id, provider, access_token, refresh_token, token_type, expires_at, created_at, updated_at 
		FROM oauth_tokens 
		WHERE provider = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, provider, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*entity.OAuthToken
	for rows.Next() {
		token := &entity.OAuthToken{}
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Provider,
			&token.AccessToken,
			&token.RefreshToken,
			&token.TokenType,
			&token.ExpiresAt,
			&token.CreatedAt,
			&token.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

// IsTokenValid checks if a token exists and is not expired
func (r *oauthTokenRepository) IsTokenValid(ctx context.Context, tokenID string) (bool, error) {
	query := `SELECT COUNT(*) 
		FROM oauth_tokens 
		WHERE id = ? AND (expires_at IS NULL OR expires_at > NOW())`

	var count int
	err := r.db.QueryRowContext(ctx, query, tokenID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetExpiredTokens retrieves expired tokens for cleanup
func (r *oauthTokenRepository) GetExpiredTokens(ctx context.Context, limit int) ([]*entity.OAuthToken, error) {
	query := `SELECT id, user_id, provider, access_token, refresh_token, token_type, expires_at, created_at, updated_at 
		FROM oauth_tokens 
		WHERE expires_at IS NOT NULL AND expires_at <= NOW()
		ORDER BY expires_at ASC
		LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*entity.OAuthToken
	for rows.Next() {
		token := &entity.OAuthToken{}
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Provider,
			&token.AccessToken,
			&token.RefreshToken,
			&token.TokenType,
			&token.ExpiresAt,
			&token.CreatedAt,
			&token.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

// GetTokensExpiringBefore retrieves tokens expiring before a specific time
func (r *oauthTokenRepository) GetTokensExpiringBefore(ctx context.Context, expiryTime time.Time, limit int) ([]*entity.OAuthToken, error) {
	query := `SELECT id, user_id, provider, access_token, refresh_token, token_type, expires_at, created_at, updated_at 
		FROM oauth_tokens 
		WHERE expires_at IS NOT NULL AND expires_at <= ?
		ORDER BY expires_at ASC
		LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, expiryTime, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*entity.OAuthToken
	for rows.Next() {
		token := &entity.OAuthToken{}
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Provider,
			&token.AccessToken,
			&token.RefreshToken,
			&token.TokenType,
			&token.ExpiresAt,
			&token.CreatedAt,
			&token.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

// DeleteByUserID deletes all OAuth tokens for a user
func (r *oauthTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM oauth_tokens WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// DeleteExpiredTokens deletes all expired tokens
func (r *oauthTokenRepository) DeleteExpiredTokens(ctx context.Context) (int64, error) {
	query := `DELETE FROM oauth_tokens 
		WHERE expires_at IS NOT NULL AND expires_at <= NOW()`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// Count returns the total number of OAuth tokens
func (r *oauthTokenRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM oauth_tokens`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountByUserID returns the number of OAuth tokens for a user
func (r *oauthTokenRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM oauth_tokens WHERE user_id = ?`

	var count int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountByProvider returns the number of OAuth tokens for a provider
func (r *oauthTokenRepository) CountByProvider(ctx context.Context, provider entity.OAuthProvider) (int64, error) {
	query := `SELECT COUNT(*) FROM oauth_tokens WHERE provider = ?`

	var count int64
	err := r.db.QueryRowContext(ctx, query, provider).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
