package entity

import "time"

// OAuthProvider represents the OAuth2 provider type
type OAuthProvider string

const (
	ProviderGoogle   OAuthProvider = "google"
	ProviderFacebook OAuthProvider = "facebook"
	ProviderGithub   OAuthProvider = "github"
	ProviderApple    OAuthProvider = "apple"
)

// TokenType represents the type of OAuth2 token
type TokenType string

const (
	TokenTypeBearer TokenType = "Bearer"
	TokenTypeMAC    TokenType = "MAC"
)

// OAuthToken represents an OAuth2 token stored in the system
// This stores access tokens and refresh tokens from OAuth2 providers
// Clean Architecture: Entity layer, no dependencies on infrastructure
type OAuthToken struct {
	ID           string        `db:"id" json:"id"`
	UserID       string        `db:"user_id" json:"user_id"`
	Provider     OAuthProvider `db:"provider" json:"provider"`
	AccessToken  string        `db:"access_token" json:"access_token"`   // Encrypted in storage
	RefreshToken *string       `db:"refresh_token" json:"refresh_token"` // Encrypted in storage, optional
	TokenType    *TokenType    `db:"token_type" json:"token_type"`       // Usually "Bearer"
	ExpiresAt    *time.Time    `db:"expires_at" json:"expires_at"`       // When the access token expires
	CreatedAt    time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at" json:"updated_at"`
}

// IsExpired checks if the OAuth token has expired
func (t *OAuthToken) IsExpired() bool {
	if t.ExpiresAt == nil {
		return false // No expiration set
	}
	return time.Now().After(*t.ExpiresAt)
}

// NeedsRefresh checks if the token should be refreshed
// Returns true if the token will expire within the next 5 minutes
func (t *OAuthToken) NeedsRefresh() bool {
	if t.ExpiresAt == nil {
		return false // No expiration set
	}
	// Refresh if token expires in less than 5 minutes
	refreshThreshold := time.Now().Add(5 * time.Minute)
	return t.ExpiresAt.Before(refreshThreshold)
}

// HasRefreshToken checks if the token has a refresh token available
func (t *OAuthToken) HasRefreshToken() bool {
	return t.RefreshToken != nil && *t.RefreshToken != ""
}
