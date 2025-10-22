package service

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/internal/domain/service"
	"buku-pintar/pkg/crypto"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

var (
	// ErrTokenNotFound is returned when a token is not found
	ErrTokenNotFound = errors.New("token not found")
	
	// ErrNoRefreshToken is returned when trying to refresh a token without a refresh token
	ErrNoRefreshToken = errors.New("no refresh token available")
	
	// ErrTokenEncryptionFailed is returned when token encryption fails
	ErrTokenEncryptionFailed = errors.New("token encryption failed")
	
	// ErrTokenDecryptionFailed is returned when token decryption fails
	ErrTokenDecryptionFailed = errors.New("token decryption failed")
	
	// ErrInvalidProvider is returned when an invalid OAuth provider is specified
	ErrInvalidProvider = errors.New("invalid OAuth provider")
)

// tokenServiceImpl implements the TokenService interface
type tokenServiceImpl struct {
	oauthTokenRepo      repository.OAuthTokenRepository
	oauthTokenRedis     repository.OAuthTokenRedisRepository
	blacklistRepo       repository.TokenBlacklistRepository
	blacklistRedis      repository.TokenBlacklistRedisRepository
	encryptor           *crypto.TokenEncryptor
	oauth2Config        map[entity.OAuthProvider]*oauth2.Config
}

// NewTokenService creates a new TokenService instance
func NewTokenService(
	oauthRepo repository.OAuthTokenRepository,
	oauthRedis repository.OAuthTokenRedisRepository,
	blacklistRepo repository.TokenBlacklistRepository,
	blacklistRedis repository.TokenBlacklistRedisRepository,
	encryptor *crypto.TokenEncryptor,
	oauth2Config map[entity.OAuthProvider]*oauth2.Config,
) service.TokenService {
	return &tokenServiceImpl{
		oauthTokenRepo:  oauthRepo,
		oauthTokenRedis: oauthRedis,
		blacklistRepo:   blacklistRepo,
		blacklistRedis:  blacklistRedis,
		encryptor:       encryptor,
		oauth2Config:    oauth2Config,
	}
}

// ============================================================================
// OAuth Token Storage Operations
// ============================================================================

// StoreOAuthToken stores an OAuth2 token with automatic encryption
func (s *tokenServiceImpl) StoreOAuthToken(
	ctx context.Context,
	userID string,
	provider entity.OAuthProvider,
	token *oauth2.Token,
) error {
	// Validate inputs
	if userID == "" {
		return errors.New("user ID is required")
	}
	if token == nil {
		return errors.New("token is required")
	}
	
	// Encrypt access token
	encryptedAccessToken, err := s.encryptor.Encrypt(token.AccessToken)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrTokenEncryptionFailed, err)
	}
	
	// Encrypt refresh token if present
	var encryptedRefreshToken *string
	if token.RefreshToken != "" {
		encrypted, err := s.encryptor.Encrypt(token.RefreshToken)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrTokenEncryptionFailed, err)
		}
		encryptedRefreshToken = &encrypted
	}
	
	// Determine token type
	tokenType := entity.TokenTypeBearer
	if token.TokenType != "" {
		tokenType = entity.TokenType(token.TokenType)
	}
	
	// Check if token already exists for this user and provider
	existingToken, err := s.oauthTokenRepo.GetByUserIDAndProvider(ctx, userID, provider)
	if err == nil && existingToken != nil {
		// Update existing token
		existingToken.AccessToken = encryptedAccessToken
		existingToken.RefreshToken = encryptedRefreshToken
		existingToken.TokenType = &tokenType
		existingToken.ExpiresAt = &token.Expiry
		existingToken.UpdatedAt = time.Now()
		
		// Update in database
		if err := s.oauthTokenRepo.Update(ctx, existingToken); err != nil {
			return fmt.Errorf("failed to update token: %w", err)
		}
		
		// Invalidate cache
		s.invalidateOAuthTokenCache(ctx, userID, provider, existingToken.ID)
		
		return nil
	}
	
	// Create new token entity
	oauthToken := &entity.OAuthToken{
		ID:           uuid.New().String(),
		UserID:       userID,
		Provider:     provider,
		AccessToken:  encryptedAccessToken,
		RefreshToken: encryptedRefreshToken,
		TokenType:    &tokenType,
		ExpiresAt:    &token.Expiry,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	
	// Store in database
	if err := s.oauthTokenRepo.Create(ctx, oauthToken); err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}
	
	// Cache the token (Redis will calculate appropriate TTL)
	if err := s.oauthTokenRedis.SetTokenByID(ctx, oauthToken); err != nil {
		// Log error but don't fail the operation
		// Cache is a performance optimization, not critical
		fmt.Printf("Warning: failed to cache token: %v\n", err)
	}
	
	// Also cache by user+provider for fast lookup
	_ = s.oauthTokenRedis.SetTokenByUserIDAndProvider(ctx, oauthToken)
	
	return nil
}

// GetOAuthToken retrieves an OAuth token (returns encrypted version)
func (s *tokenServiceImpl) GetOAuthToken(
	ctx context.Context,
	userID string,
	provider entity.OAuthProvider,
) (*entity.OAuthToken, error) {
	// Try cache first
	token, err := s.oauthTokenRedis.GetTokenByUserIDAndProvider(ctx, userID, provider)
	if err == nil && token != nil {
		return token, nil
	}
	
	// Fall back to database
	token, err = s.oauthTokenRepo.GetByUserIDAndProvider(ctx, userID, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	if token == nil {
		return nil, ErrTokenNotFound
	}
	
	// Cache for next time
	_ = s.oauthTokenRedis.SetTokenByUserIDAndProvider(ctx, token)
	
	return token, nil
}

// GetDecryptedOAuthToken retrieves and decrypts an OAuth token
func (s *tokenServiceImpl) GetDecryptedOAuthToken(
	ctx context.Context,
	userID string,
	provider entity.OAuthProvider,
) (accessToken, refreshToken string, err error) {
	// Get encrypted token
	token, err := s.GetOAuthToken(ctx, userID, provider)
	if err != nil {
		return "", "", err
	}
	
	// Decrypt access token
	decryptedAccessToken, err := s.encryptor.Decrypt(token.AccessToken)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", ErrTokenDecryptionFailed, err)
	}
	
	// Decrypt refresh token if present
	var decryptedRefreshToken string
	if token.RefreshToken != nil && *token.RefreshToken != "" {
		decryptedRefreshToken, err = s.encryptor.Decrypt(*token.RefreshToken)
		if err != nil {
			return "", "", fmt.Errorf("%w: %v", ErrTokenDecryptionFailed, err)
		}
	}
	
	return decryptedAccessToken, decryptedRefreshToken, nil
}

// UpdateOAuthToken updates an existing OAuth token
func (s *tokenServiceImpl) UpdateOAuthToken(ctx context.Context, token *entity.OAuthToken) error {
	if token == nil || token.ID == "" {
		return errors.New("invalid token")
	}
	
	token.UpdatedAt = time.Now()
	
	if err := s.oauthTokenRepo.Update(ctx, token); err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}
	
	// Invalidate cache
	s.invalidateOAuthTokenCache(ctx, token.UserID, token.Provider, token.ID)
	
	return nil
}

// DeleteOAuthToken deletes an OAuth token by ID
func (s *tokenServiceImpl) DeleteOAuthToken(ctx context.Context, tokenID string) error {
	// Get token to know user and provider for cache invalidation
	token, err := s.oauthTokenRepo.GetByID(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	
	if err := s.oauthTokenRepo.Delete(ctx, tokenID); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	
	if token != nil {
		s.invalidateOAuthTokenCache(ctx, token.UserID, token.Provider, tokenID)
	}
	
	return nil
}

// DeleteOAuthTokensByUserID deletes all OAuth tokens for a user
func (s *tokenServiceImpl) DeleteOAuthTokensByUserID(ctx context.Context, userID string) error {
	if err := s.oauthTokenRepo.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete tokens: %w", err)
	}
	
	// Invalidate all user token caches
	_ = s.oauthTokenRedis.InvalidateTokensByUserID(ctx, userID)
	
	return nil
}

// ============================================================================
// OAuth Token Retrieval Operations
// ============================================================================

// GetOAuthTokenByID retrieves an OAuth token by its ID
func (s *tokenServiceImpl) GetOAuthTokenByID(ctx context.Context, tokenID string) (*entity.OAuthToken, error) {
	token, err := s.oauthTokenRepo.GetByID(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	if token == nil {
		return nil, ErrTokenNotFound
	}
	return token, nil
}

// GetOAuthTokensByUserID retrieves all OAuth tokens for a user
func (s *tokenServiceImpl) GetOAuthTokensByUserID(ctx context.Context, userID string) ([]*entity.OAuthToken, error) {
	// Try cache first
	tokens, err := s.oauthTokenRedis.GetTokensByUserID(ctx, userID)
	if err == nil && tokens != nil {
		return tokens, nil
	}
	
	// Fall back to database
	tokens, err = s.oauthTokenRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens: %w", err)
	}
	
	// Cache for next time
	_ = s.oauthTokenRedis.SetTokensByUserID(ctx, userID, tokens)
	
	return tokens, nil
}

// GetOAuthTokensByProvider retrieves all OAuth tokens for a provider with pagination
func (s *tokenServiceImpl) GetOAuthTokensByProvider(
	ctx context.Context,
	provider entity.OAuthProvider,
	limit, offset int,
) ([]*entity.OAuthToken, error) {
	tokens, err := s.oauthTokenRepo.GetByProvider(ctx, provider, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens: %w", err)
	}
	return tokens, nil
}

// ============================================================================
// Token Validation Operations
// ============================================================================

// IsTokenValid checks if a token exists and hasn't expired
func (s *tokenServiceImpl) IsTokenValid(ctx context.Context, tokenID string) (bool, error) {
	token, err := s.GetOAuthTokenByID(ctx, tokenID)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return false, nil
		}
		return false, err
	}
	
	return !token.IsExpired(), nil
}

// IsTokenExpired checks if a token has expired
func (s *tokenServiceImpl) IsTokenExpired(ctx context.Context, token *entity.OAuthToken) bool {
	if token == nil {
		return true
	}
	return token.IsExpired()
}

// NeedsRefresh checks if a token needs to be refreshed
func (s *tokenServiceImpl) NeedsRefresh(ctx context.Context, token *entity.OAuthToken) bool {
	if token == nil {
		return false
	}
	return token.NeedsRefresh()
}

// ============================================================================
// Token Refresh Operations
// ============================================================================

// RefreshOAuthToken refreshes an OAuth token using the refresh token
func (s *tokenServiceImpl) RefreshOAuthToken(
	ctx context.Context,
	userID string,
	provider entity.OAuthProvider,
) (*oauth2.Token, error) {
	// Get current token
	token, err := s.GetOAuthToken(ctx, userID, provider)
	if err != nil {
		return nil, err
	}
	
	// Check if refresh token exists
	if !token.HasRefreshToken() {
		return nil, ErrNoRefreshToken
	}
	
	// Decrypt refresh token
	decryptedRefreshToken, err := s.encryptor.Decrypt(*token.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenDecryptionFailed, err)
	}
	
	// Get OAuth2 config for provider
	config, ok := s.oauth2Config[provider]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrInvalidProvider, provider)
	}
	
	// Create token source and refresh
	tokenSource := config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: decryptedRefreshToken,
	})
	
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	
	// Store the new token
	if err := s.StoreOAuthToken(ctx, userID, provider, newToken); err != nil {
		return nil, fmt.Errorf("failed to store refreshed token: %w", err)
	}
	
	return newToken, nil
}

// RefreshTokenIfNeeded refreshes a token only if it needs refreshing
func (s *tokenServiceImpl) RefreshTokenIfNeeded(
	ctx context.Context,
	userID string,
	provider entity.OAuthProvider,
) (*oauth2.Token, error) {
	// Get current token
	token, err := s.GetOAuthToken(ctx, userID, provider)
	if err != nil {
		return nil, err
	}
	
	// Check if refresh is needed
	if !token.NeedsRefresh() {
		// Token is still valid, decrypt and return
		accessToken, refreshToken, err := s.GetDecryptedOAuthToken(ctx, userID, provider)
		if err != nil {
			return nil, err
		}
		
		return &oauth2.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    string(*token.TokenType),
			Expiry:       *token.ExpiresAt,
		}, nil
	}
	
	// Token needs refresh
	return s.RefreshOAuthToken(ctx, userID, provider)
}

// HandleTokenRefresh handles updating a token after refresh
func (s *tokenServiceImpl) HandleTokenRefresh(
	ctx context.Context,
	token *entity.OAuthToken,
	newToken *oauth2.Token,
) error {
	if token == nil || newToken == nil {
		return errors.New("invalid parameters")
	}
	
	return s.StoreOAuthToken(ctx, token.UserID, token.Provider, newToken)
}

// ============================================================================
// Token Expiration Handling
// ============================================================================

// GetExpiredTokens retrieves expired tokens
func (s *tokenServiceImpl) GetExpiredTokens(ctx context.Context, limit int) ([]*entity.OAuthToken, error) {
	tokens, err := s.oauthTokenRepo.GetExpiredTokens(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired tokens: %w", err)
	}
	return tokens, nil
}

// GetTokensExpiringBefore retrieves tokens expiring before a specific time
func (s *tokenServiceImpl) GetTokensExpiringBefore(
	ctx context.Context,
	expiryTime time.Time,
	limit int,
) ([]*entity.OAuthToken, error) {
	tokens, err := s.oauthTokenRepo.GetTokensExpiringBefore(ctx, expiryTime, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring tokens: %w", err)
	}
	return tokens, nil
}

// CleanupExpiredTokens removes expired tokens from the database
func (s *tokenServiceImpl) CleanupExpiredTokens(ctx context.Context) (int64, error) {
	count, err := s.oauthTokenRepo.DeleteExpiredTokens(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}
	
	// Note: Individual token caches will expire naturally based on their TTL
	// No need to manually invalidate
	
	return count, nil
}

// ============================================================================
// Token Counting Operations
// ============================================================================

// CountOAuthTokens counts all OAuth tokens
func (s *tokenServiceImpl) CountOAuthTokens(ctx context.Context) (int64, error) {
	count, err := s.oauthTokenRepo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count tokens: %w", err)
	}
	return count, nil
}

// CountOAuthTokensByUserID counts OAuth tokens for a specific user
func (s *tokenServiceImpl) CountOAuthTokensByUserID(ctx context.Context, userID string) (int64, error) {
	count, err := s.oauthTokenRepo.CountByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count user tokens: %w", err)
	}
	return count, nil
}

// CountOAuthTokensByProvider counts OAuth tokens for a specific provider
func (s *tokenServiceImpl) CountOAuthTokensByProvider(ctx context.Context, provider entity.OAuthProvider) (int64, error) {
	count, err := s.oauthTokenRepo.CountByProvider(ctx, provider)
	if err != nil {
		return 0, fmt.Errorf("failed to count provider tokens: %w", err)
	}
	return count, nil
}

// ============================================================================
// Token Blacklist Operations (for JWT tokens)
// ============================================================================

// BlacklistToken adds a JWT token to the blacklist (hashes automatically)
func (s *tokenServiceImpl) BlacklistToken(
	ctx context.Context,
	token string,
	userID *string,
	reason entity.BlacklistReason,
	expiresAt time.Time,
) error {
	// Hash the token
	tokenHash := crypto.HashToken(token)
	
	return s.BlacklistTokenWithHash(ctx, tokenHash, userID, reason, expiresAt)
}

// BlacklistTokenWithHash adds a token hash to the blacklist
func (s *tokenServiceImpl) BlacklistTokenWithHash(
	ctx context.Context,
	tokenHash string,
	userID *string,
	reason entity.BlacklistReason,
	expiresAt time.Time,
) error {
	// Check if already blacklisted
	exists, err := s.blacklistRepo.IsTokenBlacklisted(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to check blacklist: %w", err)
	}
	if exists {
		return nil // Already blacklisted
	}
	
	// Create blacklist entry
	reasonPtr := &reason
	blacklist := &entity.TokenBlacklist{
		ID:            uuid.New().String(),
		TokenHash:     tokenHash,
		UserID:        userID,
		Reason:        reasonPtr,
		ExpiresAt:     expiresAt,
		BlacklistedAt: time.Now(),
	}
	
	// Store in database
	if err := s.blacklistRepo.Create(ctx, blacklist); err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}
	
	// Cache the blacklist status
	if err := s.blacklistRedis.SetTokenBlacklisted(ctx, tokenHash, expiresAt); err != nil {
		fmt.Printf("Warning: failed to cache blacklist: %v\n", err)
	}
	
	return nil
}

// IsTokenBlacklisted checks if a JWT token is blacklisted
func (s *tokenServiceImpl) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	// Hash the token
	tokenHash := crypto.HashToken(token)
	
	return s.IsTokenHashBlacklisted(ctx, tokenHash)
}

// IsTokenHashBlacklisted checks if a token hash is blacklisted
func (s *tokenServiceImpl) IsTokenHashBlacklisted(ctx context.Context, tokenHash string) (bool, error) {
	// Try cache first for fast lookup
	blacklisted, err := s.blacklistRedis.IsTokenBlacklisted(ctx, tokenHash)
	if err == nil {
		return blacklisted, nil
	}
	
	// Fall back to database
	blacklisted, err = s.blacklistRepo.IsTokenBlacklisted(ctx, tokenHash)
	if err != nil {
		return false, fmt.Errorf("failed to check blacklist: %w", err)
	}
	
	return blacklisted, nil
}

// ============================================================================
// User Token Blacklist Operations
// ============================================================================

// BlacklistAllUserTokens blacklists all tokens for a specific user
func (s *tokenServiceImpl) BlacklistAllUserTokens(
	ctx context.Context,
	userID string,
	reason entity.BlacklistReason,
	expiresAt time.Time,
) error {
	// This is a placeholder - in real implementation, you would:
	// 1. Get all active JWT tokens for the user (if you track them)
	// 2. Blacklist each one
	// For now, we'll just mark the operation as complete
	
	// Could also delete all OAuth tokens as part of user lockout
	return s.DeleteOAuthTokensByUserID(ctx, userID)
}

// GetBlacklistedTokensByUserID retrieves blacklisted tokens for a user
func (s *tokenServiceImpl) GetBlacklistedTokensByUserID(
	ctx context.Context,
	userID string,
	limit, offset int,
) ([]*entity.TokenBlacklist, error) {
	tokens, err := s.blacklistRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get blacklisted tokens: %w", err)
	}
	return tokens, nil
}

// ============================================================================
// Blacklist Query Operations
// ============================================================================

// GetBlacklistByID retrieves a blacklist entry by ID
func (s *tokenServiceImpl) GetBlacklistByID(ctx context.Context, id string) (*entity.TokenBlacklist, error) {
	entry, err := s.blacklistRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get blacklist entry: %w", err)
	}
	return entry, nil
}

// GetBlacklistByTokenHash retrieves a blacklist entry by token hash
func (s *tokenServiceImpl) GetBlacklistByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*entity.TokenBlacklist, error) {
	entry, err := s.blacklistRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get blacklist entry: %w", err)
	}
	return entry, nil
}

// GetBlacklistByReason retrieves blacklist entries by reason
func (s *tokenServiceImpl) GetBlacklistByReason(
	ctx context.Context,
	reason entity.BlacklistReason,
	limit, offset int,
) ([]*entity.TokenBlacklist, error) {
	entries, err := s.blacklistRepo.GetByReason(ctx, reason, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get blacklist entries: %w", err)
	}
	return entries, nil
}

// ============================================================================
// Blacklist Cleanup Operations
// ============================================================================

// GetExpiredBlacklistEntries retrieves expired blacklist entries
func (s *tokenServiceImpl) GetExpiredBlacklistEntries(
	ctx context.Context,
	limit int,
) ([]*entity.TokenBlacklist, error) {
	entries, err := s.blacklistRepo.GetExpiredEntries(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired entries: %w", err)
	}
	return entries, nil
}

// CleanupExpiredBlacklistEntries removes expired entries from blacklist
func (s *tokenServiceImpl) CleanupExpiredBlacklistEntries(ctx context.Context) (int64, error) {
	count, err := s.blacklistRepo.DeleteExpiredEntries(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup blacklist: %w", err)
	}
	return count, nil
}

// DeleteBlacklistEntriesExpiringBefore deletes entries expiring before a time
func (s *tokenServiceImpl) DeleteBlacklistEntriesExpiringBefore(
	ctx context.Context,
	expiryTime time.Time,
) (int64, error) {
	count, err := s.blacklistRepo.DeleteEntriesExpiringBefore(ctx, expiryTime)
	if err != nil {
		return 0, fmt.Errorf("failed to delete entries: %w", err)
	}
	return count, nil
}

// DeleteBlacklistByUserID deletes all blacklist entries for a user
func (s *tokenServiceImpl) DeleteBlacklistByUserID(ctx context.Context, userID string) error {
	if err := s.blacklistRepo.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user blacklist: %w", err)
	}
	
	// Invalidate cache
	_ = s.blacklistRedis.InvalidateBlacklistByUserID(ctx, userID)
	
	return nil
}

// ============================================================================
// Blacklist Counting Operations
// ============================================================================

// CountBlacklistedTokens counts all blacklisted tokens
func (s *tokenServiceImpl) CountBlacklistedTokens(ctx context.Context) (int64, error) {
	count, err := s.blacklistRepo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count blacklisted tokens: %w", err)
	}
	return count, nil
}

// CountBlacklistedTokensByUserID counts blacklisted tokens for a user
func (s *tokenServiceImpl) CountBlacklistedTokensByUserID(ctx context.Context, userID string) (int64, error) {
	count, err := s.blacklistRepo.CountByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count user blacklisted tokens: %w", err)
	}
	return count, nil
}

// CountBlacklistedTokensByReason counts blacklisted tokens by reason
func (s *tokenServiceImpl) CountBlacklistedTokensByReason(
	ctx context.Context,
	reason entity.BlacklistReason,
) (int64, error) {
	count, err := s.blacklistRepo.CountByReason(ctx, reason)
	if err != nil {
		return 0, fmt.Errorf("failed to count by reason: %w", err)
	}
	return count, nil
}

// CountExpiredBlacklistEntries counts expired blacklist entries
func (s *tokenServiceImpl) CountExpiredBlacklistEntries(ctx context.Context) (int64, error) {
	count, err := s.blacklistRepo.CountExpiredEntries(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count expired entries: %w", err)
	}
	return count, nil
}

// ============================================================================
// Bulk Operations
// ============================================================================

// BlacklistMultipleTokens blacklists multiple tokens at once
func (s *tokenServiceImpl) BlacklistMultipleTokens(
	ctx context.Context,
	tokens []string,
	userID *string,
	reason entity.BlacklistReason,
	expiresAt time.Time,
) error {
	for _, token := range tokens {
		if err := s.BlacklistToken(ctx, token, userID, reason, expiresAt); err != nil {
			return fmt.Errorf("failed to blacklist token: %w", err)
		}
	}
	return nil
}

// ============================================================================
// Security Operations
// ============================================================================

// RevokeAllUserAccess revokes all access for a user
func (s *tokenServiceImpl) RevokeAllUserAccess(
	ctx context.Context,
	userID string,
	reason entity.BlacklistReason,
) error {
	// Delete all OAuth tokens
	if err := s.DeleteOAuthTokensByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete OAuth tokens: %w", err)
	}
	
	// Blacklist user's tokens (if tracking JWT tokens per user)
	expiresAt := time.Now().Add(24 * time.Hour) // Typical JWT expiry
	if err := s.BlacklistAllUserTokens(ctx, userID, reason, expiresAt); err != nil {
		return fmt.Errorf("failed to blacklist tokens: %w", err)
	}
	
	return nil
}

// RevokeUserProviderAccess revokes access for a specific provider
func (s *tokenServiceImpl) RevokeUserProviderAccess(
	ctx context.Context,
	userID string,
	provider entity.OAuthProvider,
	reason entity.BlacklistReason,
) error {
	// Get the token
	token, err := s.GetOAuthToken(ctx, userID, provider)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return nil // No token to revoke
		}
		return err
	}
	
	// Delete the token
	return s.DeleteOAuthToken(ctx, token.ID)
}

// ============================================================================
// Token Hash Utility
// ============================================================================

// HashToken hashes a token using SHA-256
func (s *tokenServiceImpl) HashToken(token string) string {
	return crypto.HashToken(token)
}

// ============================================================================
// Helper Methods
// ============================================================================

// invalidateOAuthTokenCache invalidates all caches related to an OAuth token
func (s *tokenServiceImpl) invalidateOAuthTokenCache(
	ctx context.Context,
	userID string,
	provider entity.OAuthProvider,
	tokenID string,
) {
	_ = s.oauthTokenRedis.InvalidateTokenByID(ctx, tokenID)
	_ = s.oauthTokenRedis.InvalidateTokenByUserIDAndProvider(ctx, userID, provider)
	_ = s.oauthTokenRedis.InvalidateTokensByUserID(ctx, userID)
}
