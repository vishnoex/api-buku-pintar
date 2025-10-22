package redis

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type oauthTokenRedisRepository struct {
	client *redis.Client
}

// NewOAuthTokenRedisRepository creates a new instance of OAuthTokenRedisRepository
func NewOAuthTokenRedisRepository(client *redis.Client) repository.OAuthTokenRedisRepository {
	return &oauthTokenRedisRepository{
		client: client,
	}
}

// GetTokenByID retrieves an OAuth token from cache by ID
func (r *oauthTokenRedisRepository) GetTokenByID(ctx context.Context, id string) (*entity.OAuthToken, error) {
	key := fmt.Sprintf("oauth_token:id:%s", id)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var token entity.OAuthToken
	err = json.Unmarshal([]byte(data), &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// SetTokenByID stores an OAuth token in cache by ID
func (r *oauthTokenRedisRepository) SetTokenByID(ctx context.Context, token *entity.OAuthToken) error {
	key := fmt.Sprintf("oauth_token:id:%s", token.ID)

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	// Calculate TTL based on token expiry
	ttl := r.calculateTTL(token.ExpiresAt)

	return r.client.Set(ctx, key, data, ttl).Err()
}

// GetTokenByUserIDAndProvider retrieves an OAuth token from cache by user ID and provider
func (r *oauthTokenRedisRepository) GetTokenByUserIDAndProvider(ctx context.Context, userID string, provider entity.OAuthProvider) (*entity.OAuthToken, error) {
	key := fmt.Sprintf("oauth_token:user:%s:provider:%s", userID, provider)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var token entity.OAuthToken
	err = json.Unmarshal([]byte(data), &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// SetTokenByUserIDAndProvider stores an OAuth token in cache by user ID and provider
func (r *oauthTokenRedisRepository) SetTokenByUserIDAndProvider(ctx context.Context, token *entity.OAuthToken) error {
	key := fmt.Sprintf("oauth_token:user:%s:provider:%s", token.UserID, token.Provider)

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	// Calculate TTL based on token expiry
	ttl := r.calculateTTL(token.ExpiresAt)

	return r.client.Set(ctx, key, data, ttl).Err()
}

// GetTokensByUserID retrieves all OAuth tokens for a user from cache
func (r *oauthTokenRedisRepository) GetTokensByUserID(ctx context.Context, userID string) ([]*entity.OAuthToken, error) {
	key := fmt.Sprintf("oauth_token:user:%s:list", userID)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var tokens []*entity.OAuthToken
	err = json.Unmarshal([]byte(data), &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// SetTokensByUserID stores all OAuth tokens for a user in cache
func (r *oauthTokenRedisRepository) SetTokensByUserID(ctx context.Context, userID string, tokens []*entity.OAuthToken) error {
	key := fmt.Sprintf("oauth_token:user:%s:list", userID)

	data, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	// Cache for 15 minutes (token lists change less frequently)
	return r.client.Set(ctx, key, data, 15*time.Minute).Err()
}

// InvalidateTokenCache invalidates all OAuth token caches
func (r *oauthTokenRedisRepository) InvalidateTokenCache(ctx context.Context) error {
	pattern := "oauth_token:*"

	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// InvalidateTokenByID invalidates cache for a specific token by ID
func (r *oauthTokenRedisRepository) InvalidateTokenByID(ctx context.Context, tokenID string) error {
	key := fmt.Sprintf("oauth_token:id:%s", tokenID)
	return r.client.Del(ctx, key).Err()
}

// InvalidateTokensByUserID invalidates all token caches for a specific user
func (r *oauthTokenRedisRepository) InvalidateTokensByUserID(ctx context.Context, userID string) error {
	// Invalidate the user's token list
	listKey := fmt.Sprintf("oauth_token:user:%s:list", userID)
	if err := r.client.Del(ctx, listKey).Err(); err != nil {
		return err
	}

	// Invalidate all user-provider combinations for this user
	pattern := fmt.Sprintf("oauth_token:user:%s:provider:*", userID)
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// InvalidateTokenByUserIDAndProvider invalidates cache for a specific user-provider combination
func (r *oauthTokenRedisRepository) InvalidateTokenByUserIDAndProvider(ctx context.Context, userID string, provider entity.OAuthProvider) error {
	// Invalidate the specific user-provider token
	providerKey := fmt.Sprintf("oauth_token:user:%s:provider:%s", userID, provider)
	if err := r.client.Del(ctx, providerKey).Err(); err != nil {
		return err
	}

	// Also invalidate the user's token list since it changed
	listKey := fmt.Sprintf("oauth_token:user:%s:list", userID)
	return r.client.Del(ctx, listKey).Err()
}

// calculateTTL calculates the cache TTL based on token expiration
// If token has no expiry, use default 1 hour
// If token expires soon, use the remaining time
// Otherwise, use 1 hour or remaining time, whichever is shorter
func (r *oauthTokenRedisRepository) calculateTTL(expiresAt *time.Time) time.Duration {
	const defaultTTL = 1 * time.Hour
	const maxTTL = 1 * time.Hour

	if expiresAt == nil {
		return defaultTTL
	}

	timeUntilExpiry := time.Until(*expiresAt)
	if timeUntilExpiry <= 0 {
		// Token already expired, use minimal TTL to avoid caching
		return 1 * time.Minute
	}

	// Use the shorter of maxTTL or time until expiry
	if timeUntilExpiry < maxTTL {
		return timeUntilExpiry
	}

	return maxTTL
}
