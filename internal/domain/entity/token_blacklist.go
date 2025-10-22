package entity

import "time"

// BlacklistReason represents the reason why a token was blacklisted
type BlacklistReason string

const (
	ReasonLogout          BlacklistReason = "logout"            // User logged out
	ReasonPasswordChange  BlacklistReason = "password_change"   // User changed password
	ReasonSecurityBreach  BlacklistReason = "security_breach"   // Security incident detected
	ReasonAdminRevoked    BlacklistReason = "admin_revoked"     // Admin revoked the token
	ReasonAccountDeleted  BlacklistReason = "account_deleted"   // User account was deleted
	ReasonAccountSuspended BlacklistReason = "account_suspended" // User account was suspended
	ReasonTokenExpired    BlacklistReason = "token_expired"     // Token expired (for cleanup)
)

// TokenBlacklist represents a blacklisted JWT token
// Blacklisted tokens are invalid and cannot be used for authentication
// Clean Architecture: Entity layer, no dependencies on infrastructure
type TokenBlacklist struct {
	ID            string           `db:"id" json:"id"`
	TokenHash     string           `db:"token_hash" json:"token_hash"` // SHA256 hash of the token
	UserID        *string          `db:"user_id" json:"user_id"`       // Optional, for audit purposes
	BlacklistedAt time.Time        `db:"blacklisted_at" json:"blacklisted_at"`
	ExpiresAt     time.Time        `db:"expires_at" json:"expires_at"` // When the blacklist entry can be removed
	Reason        *BlacklistReason `db:"reason" json:"reason"`         // Why the token was blacklisted
}

// IsExpired checks if the blacklist entry has expired and can be removed
func (tb *TokenBlacklist) IsExpired() bool {
	return time.Now().After(tb.ExpiresAt)
}

// CanBeCleanedUp checks if this blacklist entry can be safely removed
// Returns true if the token has expired
func (tb *TokenBlacklist) CanBeCleanedUp() bool {
	return tb.IsExpired()
}

// ReasonString returns the reason as a string, or "unknown" if not set
func (tb *TokenBlacklist) ReasonString() string {
	if tb.Reason == nil {
		return "unknown"
	}
	return string(*tb.Reason)
}
