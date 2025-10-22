package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
	"time"
)

// TokenBlacklistRepository defines the interface for token blacklist data operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type TokenBlacklistRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, blacklist *entity.TokenBlacklist) error
	GetByID(ctx context.Context, id string) (*entity.TokenBlacklist, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*entity.TokenBlacklist, error)
	Delete(ctx context.Context, id string) error
	
	// Query operations
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.TokenBlacklist, error)
	GetByReason(ctx context.Context, reason entity.BlacklistReason, limit, offset int) ([]*entity.TokenBlacklist, error)
	List(ctx context.Context, limit, offset int) ([]*entity.TokenBlacklist, error)
	
	// Validation operations
	IsTokenBlacklisted(ctx context.Context, tokenHash string) (bool, error)
	
	// Cleanup operations
	GetExpiredEntries(ctx context.Context, limit int) ([]*entity.TokenBlacklist, error)
	DeleteExpiredEntries(ctx context.Context) (int64, error) // Returns number of deleted entries
	DeleteByUserID(ctx context.Context, userID string) error
	DeleteEntriesExpiringBefore(ctx context.Context, expiryTime time.Time) (int64, error)
	
	// Bulk operations
	CreateBatch(ctx context.Context, blacklist []*entity.TokenBlacklist) error
	BlacklistUserTokens(ctx context.Context, userID string, reason entity.BlacklistReason, expiresAt time.Time) error
	
	// Count operations
	Count(ctx context.Context) (int64, error)
	CountByUserID(ctx context.Context, userID string) (int64, error)
	CountByReason(ctx context.Context, reason entity.BlacklistReason) (int64, error)
	CountExpiredEntries(ctx context.Context) (int64, error)
}

// TokenBlacklistRedisRepository defines the interface for token blacklist Redis operations
// Used for caching to improve token validation performance (critical for auth middleware)
type TokenBlacklistRedisRepository interface {
	// Token hash checking (most critical operation - must be fast)
	IsTokenBlacklisted(ctx context.Context, tokenHash string) (bool, error)
	SetTokenBlacklisted(ctx context.Context, tokenHash string, expiresAt time.Time) error
	
	// Blacklist entry caching
	GetBlacklistByTokenHash(ctx context.Context, tokenHash string) (*entity.TokenBlacklist, error)
	SetBlacklistByTokenHash(ctx context.Context, blacklist *entity.TokenBlacklist) error
	
	// User blacklist caching
	GetBlacklistByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.TokenBlacklist, error)
	SetBlacklistByUserID(ctx context.Context, userID string, blacklist []*entity.TokenBlacklist) error
	
	// Cache invalidation
	InvalidateBlacklistCache(ctx context.Context) error
	InvalidateBlacklistByTokenHash(ctx context.Context, tokenHash string) error
	InvalidateBlacklistByUserID(ctx context.Context, userID string) error
	
	// Bulk operations for user logout/security events
	BlacklistUserTokens(ctx context.Context, userID string, tokenHashes []string, expiresAt time.Time) error
}
