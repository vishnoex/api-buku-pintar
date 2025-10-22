# TokenService Interface Documentation

## Overview

The `TokenService` interface provides a comprehensive API for managing OAuth2 tokens, including storage, encryption, refresh, validation, and blacklisting operations. It serves as the business logic layer for all token-related operations in the authentication system.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    TokenService                          │
│  (Business Logic & Orchestration)                        │
└────────────┬────────────────────────────────────────────┘
             │
             ├─── Token Encryption (pkg/crypto)
             │    └─── AES-256-GCM encryption/decryption
             │
             ├─── OAuthTokenRepository (MySQL)
             │    └─── Persistent storage
             │
             ├─── OAuthTokenRedisRepository (Cache)
             │    └─── Fast access with TTL
             │
             ├─── TokenBlacklistRepository (MySQL)
             │    └─── Blacklist storage
             │
             └─── TokenBlacklistRedisRepository (Cache)
                  └─── Fast blacklist checks
```

## Core Features

### 1. OAuth Token Storage (Encrypted)
- **Automatic Encryption**: Tokens are automatically encrypted before storage
- **Provider Support**: Google, Facebook, Github, Apple
- **Refresh Token Management**: Handles both access and refresh tokens
- **Expiry Tracking**: Monitors token expiration

### 2. Token Refresh
- **Automatic Refresh**: Detects when tokens need refreshing
- **Provider Integration**: Works with OAuth2 providers
- **Smart Caching**: Updates both database and cache

### 3. Token Blacklisting (JWT)
- **Fast Validation**: Hash-based blacklist checking
- **Reason Tracking**: Tracks why tokens were blacklisted
- **Automatic Cleanup**: Removes expired blacklist entries

### 4. Security Operations
- **User Access Revocation**: Blacklist all tokens for a user
- **Provider-Specific Revocation**: Revoke specific provider access
- **Bulk Operations**: Efficient batch blacklisting

## Method Categories

### OAuth Token Storage (8 methods)

#### `StoreOAuthToken`
Stores an OAuth2 token with automatic encryption.

```go
func (s *TokenService) StoreOAuthToken(
    ctx context.Context,
    userID string,
    provider entity.OAuthProvider,
    token *oauth2.Token,
) error
```

**Usage:**
```go
// After OAuth2 callback
err := tokenService.StoreOAuthToken(ctx, userID, entity.ProviderGoogle, oauthToken)
```

**Features:**
- Encrypts access and refresh tokens
- Stores expiry information
- Updates cache automatically

#### `GetDecryptedOAuthToken`
Retrieves and automatically decrypts tokens.

```go
func (s *TokenService) GetDecryptedOAuthToken(
    ctx context.Context,
    userID string,
    provider entity.OAuthProvider,
) (accessToken, refreshToken string, err error)
```

**Usage:**
```go
access, refresh, err := tokenService.GetDecryptedOAuthToken(ctx, userID, entity.ProviderGoogle)
```

### Token Validation (3 methods)

#### `IsTokenValid`
Checks if a token exists and hasn't expired.

```go
func (s *TokenService) IsTokenValid(ctx context.Context, tokenID string) (bool, error)
```

#### `IsTokenExpired`
Checks if a token has expired.

```go
func (s *TokenService) IsTokenExpired(ctx context.Context, token *entity.OAuthToken) bool
```

#### `NeedsRefresh`
Determines if a token should be refreshed (expires within 5 minutes).

```go
func (s *TokenService) NeedsRefresh(ctx context.Context, token *entity.OAuthToken) bool
```

### Token Refresh (3 methods)

#### `RefreshOAuthToken`
Forces a token refresh from the OAuth provider.

```go
func (s *TokenService) RefreshOAuthToken(
    ctx context.Context,
    userID string,
    provider entity.OAuthProvider,
) (*oauth2.Token, error)
```

#### `RefreshTokenIfNeeded`
Refreshes token only if it's expiring soon.

```go
func (s *TokenService) RefreshTokenIfNeeded(
    ctx context.Context,
    userID string,
    provider entity.OAuthProvider,
) (*oauth2.Token, error)
```

**Smart Logic:**
- Checks if token needs refresh
- Returns existing token if still valid
- Refreshes and updates if expiring soon

### Token Blacklist (6 methods)

#### `BlacklistToken`
Blacklists a JWT token (hashes automatically).

```go
func (s *TokenService) BlacklistToken(
    ctx context.Context,
    token string,
    userID *string,
    reason entity.BlacklistReason,
    expiresAt time.Time,
) error
```

**Usage:**
```go
// User logout
expiresAt := time.Now().Add(24 * time.Hour)
err := tokenService.BlacklistToken(
    ctx,
    jwtToken,
    &userID,
    entity.ReasonLogout,
    expiresAt,
)
```

#### `IsTokenBlacklisted`
Fast check if a JWT token is blacklisted.

```go
func (s *TokenService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
```

**Usage in Middleware:**
```go
// In authentication middleware
blacklisted, err := tokenService.IsTokenBlacklisted(ctx, jwtToken)
if blacklisted {
    return errors.New("token has been revoked")
}
```

### Security Operations (2 methods)

#### `RevokeAllUserAccess`
Complete user token revocation.

```go
func (s *TokenService) RevokeAllUserAccess(
    ctx context.Context,
    userID string,
    reason entity.BlacklistReason,
) error
```

**Features:**
- Deletes all OAuth tokens for the user
- Blacklists all JWT tokens
- Updates all caches
- Used for: Account deletion, security breach, suspension

#### `RevokeUserProviderAccess`
Revokes access for a specific OAuth provider.

```go
func (s *TokenService) RevokeUserProviderAccess(
    ctx context.Context,
    userID string,
    provider entity.OAuthProvider,
    reason entity.BlacklistReason,
) error
```

### Cleanup Operations (3 methods)

#### `CleanupExpiredTokens`
Removes expired OAuth tokens from the database.

```go
func (s *TokenService) CleanupExpiredTokens(ctx context.Context) (int64, error)
```

**Usage (Cron Job):**
```go
// Run daily
count, err := tokenService.CleanupExpiredTokens(ctx)
log.Printf("Cleaned up %d expired tokens", count)
```

#### `CleanupExpiredBlacklistEntries`
Removes expired blacklist entries.

```go
func (s *TokenService) CleanupExpiredBlacklistEntries(ctx context.Context) (int64, error)
```

## Helper Types

### TokenRefreshResult
Result of a token refresh operation.

```go
type TokenRefreshResult struct {
    AccessToken  string    // New access token
    RefreshToken string    // New refresh token (if updated)
    TokenType    string    // Token type (usually "Bearer")
    Expiry       time.Time // Expiration time
    WasRefreshed bool      // True if token was actually refreshed
}
```

### TokenValidationResult
Comprehensive token validation result.

```go
type TokenValidationResult struct {
    IsValid       bool                // Token exists and not expired
    IsExpired     bool                // Token has expired
    IsBlacklisted bool                // Token is in blacklist
    NeedsRefresh  bool                // Should be refreshed soon
    Token         *entity.OAuthToken  // The actual token entity
}
```

## Implementation Example

### Service Implementation Structure

```go
type tokenServiceImpl struct {
    oauthTokenRepo      repository.OAuthTokenRepository
    oauthTokenRedis     repository.OAuthTokenRedisRepository
    blacklistRepo       repository.TokenBlacklistRepository
    blacklistRedis      repository.TokenBlacklistRedisRepository
    encryptor           *crypto.TokenEncryptor
    oauth2Config        map[entity.OAuthProvider]*oauth2.Config
}

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
```

## Use Cases

### 1. OAuth2 Callback Handler

```go
func HandleOAuth2Callback(ctx context.Context, code string, provider entity.OAuthProvider) error {
    // Exchange code for token
    oauth2Token, err := exchangeCodeForToken(code, provider)
    if err != nil {
        return err
    }
    
    // Get or create user
    user, err := getUserFromOAuth2(oauth2Token, provider)
    if err != nil {
        return err
    }
    
    // Store encrypted token
    return tokenService.StoreOAuthToken(ctx, user.ID, provider, oauth2Token)
}
```

### 2. API Call with Token Refresh

```go
func MakeAPICall(ctx context.Context, userID string, provider entity.OAuthProvider) error {
    // Get token and refresh if needed
    newToken, err := tokenService.RefreshTokenIfNeeded(ctx, userID, provider)
    if err != nil {
        return err
    }
    
    // Use the token
    return callExternalAPI(newToken.AccessToken)
}
```

### 3. User Logout

```go
func HandleLogout(ctx context.Context, userID string, jwtToken string) error {
    // Blacklist the JWT token
    expiresAt := time.Now().Add(24 * time.Hour) // JWT expiry
    err := tokenService.BlacklistToken(
        ctx,
        jwtToken,
        &userID,
        entity.ReasonLogout,
        expiresAt,
    )
    if err != nil {
        return err
    }
    
    // Delete all OAuth tokens
    return tokenService.DeleteOAuthTokensByUserID(ctx, userID)
}
```

### 4. Authentication Middleware

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractJWTToken(r)
        
        // Check if token is blacklisted
        blacklisted, err := tokenService.IsTokenBlacklisted(r.Context(), token)
        if err != nil {
            http.Error(w, "Internal error", http.StatusInternalServerError)
            return
        }
        
        if blacklisted {
            http.Error(w, "Token has been revoked", http.StatusUnauthorized)
            return
        }
        
        // Validate JWT and continue
        next.ServeHTTP(w, r)
    })
}
```

### 5. Account Security Event

```go
func HandlePasswordChange(ctx context.Context, userID string) error {
    // Revoke all user access (invalidate all tokens)
    return tokenService.RevokeAllUserAccess(
        ctx,
        userID,
        entity.ReasonPasswordChange,
    )
}
```

### 6. Scheduled Cleanup Job

```go
func RunDailyCleanup(ctx context.Context) {
    // Clean up expired OAuth tokens
    tokenCount, err := tokenService.CleanupExpiredTokens(ctx)
    log.Printf("Cleaned up %d expired OAuth tokens", tokenCount)
    
    // Clean up expired blacklist entries
    blacklistCount, err := tokenService.CleanupExpiredBlacklistEntries(ctx)
    log.Printf("Cleaned up %d expired blacklist entries", blacklistCount)
}
```

## Performance Considerations

### Caching Strategy
- **OAuth Tokens**: Cached with TTL matching token expiry
- **Blacklist Checks**: Cached for fast middleware validation
- **User Tokens**: List cached for 15 minutes

### Optimization Tips
1. Use `RefreshTokenIfNeeded` instead of `RefreshOAuthToken` to avoid unnecessary refreshes
2. Batch blacklist operations when possible
3. Run cleanup jobs during off-peak hours
4. Monitor cache hit rates for blacklist checks

## Error Handling

All methods return errors that should be properly handled:

```go
token, err := tokenService.GetDecryptedOAuthToken(ctx, userID, provider)
if err != nil {
    switch {
    case errors.Is(err, repository.ErrNotFound):
        // Token not found, may need to re-authenticate
    case errors.Is(err, crypto.ErrDecryptionFailed):
        // Encryption key may have changed
    default:
        // Other errors
    }
}
```

## Testing

Example test structure:

```go
func TestStoreOAuthToken(t *testing.T) {
    // Setup mocks
    mockRepo := &MockOAuthTokenRepository{}
    mockRedis := &MockOAuthTokenRedisRepository{}
    mockEncryptor := &MockTokenEncryptor{}
    
    service := NewTokenService(mockRepo, mockRedis, ..., mockEncryptor, ...)
    
    // Test storage
    token := &oauth2.Token{
        AccessToken:  "test-access-token",
        RefreshToken: "test-refresh-token",
        Expiry:       time.Now().Add(1 * time.Hour),
    }
    
    err := service.StoreOAuthToken(ctx, userID, entity.ProviderGoogle, token)
    assert.NoError(t, err)
    
    // Verify encryption was called
    assert.True(t, mockEncryptor.EncryptCalled)
}
```

## Security Notes

1. **Never Log Tokens**: Tokens should never appear in logs
2. **Encryption Required**: All tokens must be encrypted at rest
3. **Hash Blacklist**: Store token hashes, not raw JWT tokens
4. **Expiry Tracking**: Always set appropriate expiry times
5. **Secure Deletion**: Ensure tokens are completely removed on deletion
6. **Audit Trail**: Log all token operations for security auditing

## Related Documentation

- [Token Encryption](../../pkg/crypto/README.md)
- [OAuthToken Entity](../entity/oauth_token.go)
- [TokenBlacklist Entity](../entity/token_blacklist.go)
- [Token Repositories](../repository/)
