package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
	"time"
)

// OAuthTokenRepository defines the interface for OAuth token data operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type OAuthTokenRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, token *entity.OAuthToken) error
	GetByID(ctx context.Context, id string) (*entity.OAuthToken, error)
	Update(ctx context.Context, token *entity.OAuthToken) error
	Delete(ctx context.Context, id string) error
	
	// Query operations
	GetByUserIDAndProvider(ctx context.Context, userID string, provider entity.OAuthProvider) (*entity.OAuthToken, error)
	GetByUserID(ctx context.Context, userID string) ([]*entity.OAuthToken, error)
	GetByProvider(ctx context.Context, provider entity.OAuthProvider, limit, offset int) ([]*entity.OAuthToken, error)
	
	// Token validation and refresh operations
	IsTokenValid(ctx context.Context, tokenID string) (bool, error)
	GetExpiredTokens(ctx context.Context, limit int) ([]*entity.OAuthToken, error)
	GetTokensExpiringBefore(ctx context.Context, expiryTime time.Time, limit int) ([]*entity.OAuthToken, error)
	
	// Bulk operations
	DeleteByUserID(ctx context.Context, userID string) error
	DeleteExpiredTokens(ctx context.Context) (int64, error) // Returns number of deleted tokens
	
	// Count operations
	Count(ctx context.Context) (int64, error)
	CountByUserID(ctx context.Context, userID string) (int64, error)
	CountByProvider(ctx context.Context, provider entity.OAuthProvider) (int64, error)
}

// OAuthTokenRedisRepository defines the interface for OAuth token Redis operations
// Used for caching to improve token validation performance
type OAuthTokenRedisRepository interface {
	// Token caching
	GetTokenByID(ctx context.Context, id string) (*entity.OAuthToken, error)
	SetTokenByID(ctx context.Context, token *entity.OAuthToken) error
	
	// User-provider token caching (most common query)
	GetTokenByUserIDAndProvider(ctx context.Context, userID string, provider entity.OAuthProvider) (*entity.OAuthToken, error)
	SetTokenByUserIDAndProvider(ctx context.Context, token *entity.OAuthToken) error
	
	// User tokens list caching
	GetTokensByUserID(ctx context.Context, userID string) ([]*entity.OAuthToken, error)
	SetTokensByUserID(ctx context.Context, userID string, tokens []*entity.OAuthToken) error
	
	// Cache invalidation
	InvalidateTokenCache(ctx context.Context) error
	InvalidateTokenByID(ctx context.Context, tokenID string) error
	InvalidateTokensByUserID(ctx context.Context, userID string) error
	InvalidateTokenByUserIDAndProvider(ctx context.Context, userID string, provider entity.OAuthProvider) error
}
