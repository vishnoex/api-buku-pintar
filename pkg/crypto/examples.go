package crypto

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/pkg/config"
	"context"
	"time"

	"github.com/google/uuid"
)

// This file contains example usage patterns for the token encryption utilities
// These examples show how to integrate encryption with OAuth token management

// Example 1: Initialize encryptor from config
func ExampleInitializeEncryptor(cfg *config.Config) (*TokenEncryptor, error) {
	// Create encryptor from configuration
	return NewTokenEncryptorFromString(cfg.Security.TokenEncryptionKey)
}

// Example 2: Encrypt OAuth tokens before storing
func ExampleStoreOAuthToken(encryptor *TokenEncryptor, userID string, provider string, accessToken, refreshToken string, expiresAt time.Time) (*entity.OAuthToken, error) {
	// Encrypt access token
	encryptedAccess, err := encryptor.EncryptToken(accessToken)
	if err != nil {
		return nil, err
	}

	// Encrypt refresh token if present
	var encryptedRefresh *string
	if refreshToken != "" {
		encrypted, err := encryptor.EncryptRefreshToken(refreshToken)
		if err != nil {
			return nil, err
		}
		encryptedRefresh = &encrypted
	}

	// Create entity with encrypted tokens
	token := &entity.OAuthToken{
		ID:           uuid.New().String(),
		UserID:       userID,
		Provider:     entity.OAuthProvider(provider),
		AccessToken:  encryptedAccess,  // Stored encrypted
		RefreshToken: encryptedRefresh, // Stored encrypted
		ExpiresAt:    &expiresAt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return token, nil
}

// Example 3: Decrypt OAuth tokens after retrieval
func ExampleRetrieveOAuthToken(encryptor *TokenEncryptor, encryptedToken *entity.OAuthToken) (accessToken string, refreshToken string, err error) {
	// Decrypt access token
	accessToken, err = encryptor.DecryptToken(encryptedToken.AccessToken)
	if err != nil {
		return "", "", err
	}

	// Decrypt refresh token if present
	if encryptedToken.RefreshToken != nil && *encryptedToken.RefreshToken != "" {
		refreshToken, err = encryptor.DecryptRefreshToken(*encryptedToken.RefreshToken)
		if err != nil {
			return "", "", err
		}
	}

	return accessToken, refreshToken, nil
}

// Example 4: Hash JWT token for blacklisting
func ExampleBlacklistJWTToken(jwtToken string, userID string, reason entity.BlacklistReason, expiresAt time.Time) (*entity.TokenBlacklist, error) {
	// Hash the JWT token (never store raw JWT)
	tokenHash := HashToken(jwtToken)

	// Create blacklist entry
	blacklist := &entity.TokenBlacklist{
		ID:            uuid.New().String(),
		TokenHash:     tokenHash,
		UserID:        &userID,
		BlacklistedAt: time.Now(),
		ExpiresAt:     expiresAt,
		Reason:        &reason,
	}

	return blacklist, nil
}

// Example 5: Check if token is blacklisted
func ExampleCheckTokenBlacklist(jwtToken string, blacklistedHashes map[string]bool) bool {
	// Hash the token
	tokenHash := HashToken(jwtToken)

	// Check if hash exists in blacklist
	return blacklistedHashes[tokenHash]
}

// Example 6: Service layer integration
type OAuthTokenService struct {
	encryptor *TokenEncryptor
	// repository interfaces would be here
}

func NewOAuthTokenService(cfg *config.Config) (*OAuthTokenService, error) {
	encryptor, err := NewTokenEncryptorFromString(cfg.Security.TokenEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &OAuthTokenService{
		encryptor: encryptor,
	}, nil
}

// Example method: Store token with automatic encryption
func (s *OAuthTokenService) StoreToken(ctx context.Context, userID string, provider entity.OAuthProvider, accessToken, refreshToken string, expiresAt time.Time) error {
	// Encrypt tokens
	encryptedAccess, err := s.encryptor.EncryptToken(accessToken)
	if err != nil {
		return err
	}

	var encryptedRefresh *string
	if refreshToken != "" {
		encrypted, err := s.encryptor.EncryptRefreshToken(refreshToken)
		if err != nil {
			return err
		}
		encryptedRefresh = &encrypted
	}

	token := &entity.OAuthToken{
		ID:           uuid.New().String(),
		UserID:       userID,
		Provider:     provider,
		AccessToken:  encryptedAccess,
		RefreshToken: encryptedRefresh,
		ExpiresAt:    &expiresAt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Store in repository (would call repository.Create here)
	_ = token

	return nil
}

// Example method: Get decrypted token
func (s *OAuthTokenService) GetDecryptedToken(ctx context.Context, userID string, provider entity.OAuthProvider) (string, string, error) {
	// Retrieve from repository (would call repository.GetByUserIDAndProvider here)
	// For this example, assume we got an encrypted token
	var encryptedToken *entity.OAuthToken

	// Decrypt access token
	accessToken, err := s.encryptor.DecryptToken(encryptedToken.AccessToken)
	if err != nil {
		return "", "", err
	}

	// Decrypt refresh token if present
	refreshToken := ""
	if encryptedToken.RefreshToken != nil && *encryptedToken.RefreshToken != "" {
		refreshToken, err = s.encryptor.DecryptRefreshToken(*encryptedToken.RefreshToken)
		if err != nil {
			return "", "", err
		}
	}

	return accessToken, refreshToken, nil
}

// Example 7: Middleware integration for token blacklist checking
func ExampleAuthMiddleware(jwtToken string) (bool, error) {
	// Hash the JWT token
	tokenHash := HashToken(jwtToken)

	// Check in Redis/Database if this hash is blacklisted
	// This is fast because we're checking hash, not full token
	// isBlacklisted := repository.IsTokenBlacklisted(ctx, tokenHash)

	_ = tokenHash
	return false, nil // Would return actual blacklist status
}

// Example 8: Batch encrypt multiple tokens
func ExampleBatchEncryptTokens(encryptor *TokenEncryptor, tokens []string) ([]string, error) {
	encrypted := make([]string, len(tokens))

	for i, token := range tokens {
		enc, err := encryptor.Encrypt(token)
		if err != nil {
			return nil, err
		}
		encrypted[i] = enc
	}

	return encrypted, nil
}

// Example 9: Generate and save encryption key
func ExampleGenerateEncryptionKey() (string, error) {
	// Generate a new random encryption key
	keyBase64, err := GenerateRandomKeyBase64()
	if err != nil {
		return "", err
	}

	// This key should be saved to your config file or environment variables
	// Example output: "Rk1WQjJGNzVITktMOFBRUlNUVVZXWFlaQUJDREVGR0g="
	return keyBase64, nil
}

// Example 10: Handle key rotation (advanced)
func ExampleKeyRotation(oldKey, newKey string, encryptedToken string) (string, error) {
	// Create encryptor with old key
	oldEncryptor, err := NewTokenEncryptorFromString(oldKey)
	if err != nil {
		return "", err
	}

	// Create encryptor with new key
	newEncryptor, err := NewTokenEncryptorFromString(newKey)
	if err != nil {
		return "", err
	}

	// Decrypt with old key
	plaintext, err := oldEncryptor.Decrypt(encryptedToken)
	if err != nil {
		return "", err
	}

	// Re-encrypt with new key
	newEncrypted, err := newEncryptor.Encrypt(plaintext)
	if err != nil {
		return "", err
	}

	return newEncrypted, nil
}
