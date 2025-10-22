package service

import (
	"buku-pintar/internal/domain/entity"
	"context"
	"time"

	"golang.org/x/oauth2"
)

// TokenService defines the interface for OAuth2 token management operations
// Clean Architecture: Service layer, orchestrates token storage, encryption, refresh, and blacklisting
type TokenService interface {
	// OAuth Token Storage Operations (with automatic encryption)
	StoreOAuthToken(ctx context.Context, userID string, provider entity.OAuthProvider, token *oauth2.Token) error
	GetOAuthToken(ctx context.Context, userID string, provider entity.OAuthProvider) (*entity.OAuthToken, error)
	GetDecryptedOAuthToken(ctx context.Context, userID string, provider entity.OAuthProvider) (accessToken, refreshToken string, err error)
	UpdateOAuthToken(ctx context.Context, token *entity.OAuthToken) error
	DeleteOAuthToken(ctx context.Context, tokenID string) error
	DeleteOAuthTokensByUserID(ctx context.Context, userID string) error
	
	// OAuth Token Retrieval Operations
	GetOAuthTokenByID(ctx context.Context, tokenID string) (*entity.OAuthToken, error)
	GetOAuthTokensByUserID(ctx context.Context, userID string) ([]*entity.OAuthToken, error)
	GetOAuthTokensByProvider(ctx context.Context, provider entity.OAuthProvider, limit, offset int) ([]*entity.OAuthToken, error)
	
	// Token Validation Operations
	IsTokenValid(ctx context.Context, tokenID string) (bool, error)
	IsTokenExpired(ctx context.Context, token *entity.OAuthToken) bool
	NeedsRefresh(ctx context.Context, token *entity.OAuthToken) bool
	
	// Token Refresh Operations
	RefreshOAuthToken(ctx context.Context, userID string, provider entity.OAuthProvider) (*oauth2.Token, error)
	RefreshTokenIfNeeded(ctx context.Context, userID string, provider entity.OAuthProvider) (*oauth2.Token, error)
	HandleTokenRefresh(ctx context.Context, token *entity.OAuthToken, newToken *oauth2.Token) error
	
	// Token Expiration Handling
	GetExpiredTokens(ctx context.Context, limit int) ([]*entity.OAuthToken, error)
	GetTokensExpiringBefore(ctx context.Context, expiryTime time.Time, limit int) ([]*entity.OAuthToken, error)
	CleanupExpiredTokens(ctx context.Context) (int64, error)
	
	// Token Counting Operations
	CountOAuthTokens(ctx context.Context) (int64, error)
	CountOAuthTokensByUserID(ctx context.Context, userID string) (int64, error)
	CountOAuthTokensByProvider(ctx context.Context, provider entity.OAuthProvider) (int64, error)
	
	// Token Blacklist Operations (for JWT tokens)
	BlacklistToken(ctx context.Context, token string, userID *string, reason entity.BlacklistReason, expiresAt time.Time) error
	BlacklistTokenWithHash(ctx context.Context, tokenHash string, userID *string, reason entity.BlacklistReason, expiresAt time.Time) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	IsTokenHashBlacklisted(ctx context.Context, tokenHash string) (bool, error)
	
	// User Token Blacklist Operations
	BlacklistAllUserTokens(ctx context.Context, userID string, reason entity.BlacklistReason, expiresAt time.Time) error
	GetBlacklistedTokensByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.TokenBlacklist, error)
	
	// Blacklist Query Operations
	GetBlacklistByID(ctx context.Context, id string) (*entity.TokenBlacklist, error)
	GetBlacklistByTokenHash(ctx context.Context, tokenHash string) (*entity.TokenBlacklist, error)
	GetBlacklistByReason(ctx context.Context, reason entity.BlacklistReason, limit, offset int) ([]*entity.TokenBlacklist, error)
	
	// Blacklist Cleanup Operations
	GetExpiredBlacklistEntries(ctx context.Context, limit int) ([]*entity.TokenBlacklist, error)
	CleanupExpiredBlacklistEntries(ctx context.Context) (int64, error)
	DeleteBlacklistEntriesExpiringBefore(ctx context.Context, expiryTime time.Time) (int64, error)
	DeleteBlacklistByUserID(ctx context.Context, userID string) error
	
	// Blacklist Counting Operations
	CountBlacklistedTokens(ctx context.Context) (int64, error)
	CountBlacklistedTokensByUserID(ctx context.Context, userID string) (int64, error)
	CountBlacklistedTokensByReason(ctx context.Context, reason entity.BlacklistReason) (int64, error)
	CountExpiredBlacklistEntries(ctx context.Context) (int64, error)
	
	// Bulk Operations
	BlacklistMultipleTokens(ctx context.Context, tokens []string, userID *string, reason entity.BlacklistReason, expiresAt time.Time) error
	
	// Security Operations
	RevokeAllUserAccess(ctx context.Context, userID string, reason entity.BlacklistReason) error
	RevokeUserProviderAccess(ctx context.Context, userID string, provider entity.OAuthProvider, reason entity.BlacklistReason) error
	
	// Token Hash Utility
	HashToken(token string) string
}

// TokenRefreshResult represents the result of a token refresh operation
type TokenRefreshResult struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	Expiry       time.Time
	WasRefreshed bool
}

// TokenValidationResult represents the result of token validation
type TokenValidationResult struct {
	IsValid      bool
	IsExpired    bool
	IsBlacklisted bool
	NeedsRefresh bool
	Token        *entity.OAuthToken
}
