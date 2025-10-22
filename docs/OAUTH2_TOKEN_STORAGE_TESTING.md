# OAuth2 Token Storage - Testing Guide

## Overview

This document provides comprehensive testing instructions for the OAuth2 token storage integration completed on Day 5 (OAuth2 Token Management Day 2).

## What Was Implemented

### 1. OAuth2Handler Updates
- ✅ Added `TokenService` dependency
- ✅ Created `convertProviderToEntity()` helper function
- ✅ Updated `Callback()` method to store OAuth2 tokens
- ✅ Updated `HandleOAuth2Redirect()` method to store tokens
- ✅ Added logging for token storage operations

### 2. Main.go Dependency Injection
- ✅ Initialized `TokenEncryptor` with encryption key from config
- ✅ Initialized `OAuthTokenRepository` (MySQL)
- ✅ Initialized `OAuthTokenRedisRepository` (Redis cache)
- ✅ Created OAuth2 configs map for token refresh
- ✅ Initialized `TokenService` with all dependencies
- ✅ Updated `OAuth2Handler` initialization with `TokenService`

### 3. Provider Type Conversion
- ✅ Maps `oauth2.Provider` → `entity.OAuthProvider`
- ✅ Handles: Google, GitHub, Facebook

## Test Scenarios

### Manual Testing

#### Test 1: Google OAuth Login (New User)

**Steps:**
1. Start the server: `go run cmd/api/main.go`
2. Make POST request to `/oauth2/login`:
```json
{
  "provider": "google"
}
```
3. Open the returned `auth_url` in browser
4. Complete Google OAuth consent
5. After callback, check database

**Expected Results:**
- User created in `users` table
- OAuth token stored in `oauth_tokens` table
- `access_token` is encrypted (not plain text)
- `refresh_token` is encrypted (if present)
- Provider is "google"
- Console shows: `Successfully stored OAuth2 token for user {userID} with provider google`

**Database Verification:**
```sql
SELECT id, user_id, provider, LENGTH(access_token), expires_at 
FROM oauth_tokens 
ORDER BY created_at DESC 
LIMIT 5;
```

#### Test 2: GitHub OAuth Login (Existing User)

**Steps:**
1. User already exists in database
2. Login with GitHub OAuth
3. Check if token is updated/created

**Expected Results:**
- User info updated if needed
- OAuth token for GitHub stored
- Previous Google token (if exists) remains unchanged
- Console shows successful token storage

#### Test 3: Token Encryption Verification

**Steps:**
1. Login with any OAuth provider
2. Query database for the token
3. Verify token is not readable plain text

**Database Check:**
```sql
SELECT access_token FROM oauth_tokens WHERE user_id = 'USER_ID';
```

**Expected:**
- Token should be base64-encoded gibberish
- Should start with random characters (nonce)
- Should NOT be a valid-looking OAuth token

#### Test 4: Token Refresh Capability

**Steps:**
1. Get user's OAuth token from database
2. Call TokenService to refresh:
```go
newToken, err := tokenService.RefreshTokenIfNeeded(ctx, userID, entity.ProviderGoogle)
```

**Expected:**
- If token expires soon, new token retrieved from Google
- Database updated with new encrypted token
- Old token replaced

#### Test 5: Multiple Provider Login

**Steps:**
1. User logs in with Google
2. Same user logs in with GitHub
3. Check database

**Expected:**
- Two separate records in `oauth_tokens`
- One for Google, one for GitHub
- Both associated with same `user_id`

### Automated Testing

#### Unit Test: Provider Conversion

```go
func TestConvertProviderToEntity(t *testing.T) {
    tests := []struct {
        name     string
        provider oauth2.Provider
        want     entity.OAuthProvider
    }{
        {
            name:     "Google",
            provider: oauth2.ProviderGoogle,
            want:     entity.ProviderGoogle,
        },
        {
            name:     "GitHub",
            provider: oauth2.ProviderGitHub,
            want:     entity.ProviderGithub,
        },
        {
            name:     "Facebook",
            provider: oauth2.ProviderFacebook,
            want:     entity.ProviderFacebook,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := convertProviderToEntity(tt.provider)
            if got != tt.want {
                t.Errorf("convertProviderToEntity() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

#### Integration Test: OAuth Callback with Token Storage

```go
func TestOAuth2Handler_CallbackStoresToken(t *testing.T) {
    // Setup mocks
    mockOAuth2Service := &MockOAuth2Service{}
    mockUserUsecase := &MockUserUsecase{}
    mockTokenService := &MockTokenService{}
    
    handler := NewOAuth2Handler(mockOAuth2Service, mockUserUsecase, mockTokenService)
    
    // Mock OAuth2 token
    mockToken := &oauth2lib.Token{
        AccessToken:  "test-access-token",
        RefreshToken: "test-refresh-token",
        Expiry:       time.Now().Add(1 * time.Hour),
        TokenType:    "Bearer",
    }
    
    mockOAuth2Service.On("ExchangeCode", mock.Anything, mock.Anything, mock.Anything).Return(mockToken, nil)
    mockOAuth2Service.On("GetUserInfo", mock.Anything, mock.Anything, mock.Anything).Return(&oauth2.UserInfo{
        ID:      "123",
        Email:   "test@example.com",
        Name:    "Test User",
        Picture: "https://example.com/pic.jpg",
    }, nil)
    
    mockUserUsecase.On("GetUserByEmail", mock.Anything, "test@example.com").Return(nil, nil)
    mockUserUsecase.On("RegisterWithOAuth2", mock.Anything, mock.Anything, mock.Anything).Return(nil)
    
    // This is the critical assertion
    mockTokenService.On("StoreOAuthToken", mock.Anything, mock.Anything, entity.ProviderGoogle, mockToken).Return(nil)
    
    // Execute
    req := httptest.NewRequest("POST", "/oauth2/callback?provider=google", strings.NewReader(`{"code":"test-code","state":"test-state"}`))
    w := httptest.NewRecorder()
    
    handler.Callback(w, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    mockTokenService.AssertCalled(t, "StoreOAuthToken", mock.Anything, mock.Anything, entity.ProviderGoogle, mockToken)
}
```

## Verification Checklist

### Pre-Deployment Checks

- [ ] **Configuration**
  - [ ] `security.token_encryption_key` is set in config.json
  - [ ] Key is at least 32 characters (for AES-256)
  - [ ] OAuth2 providers configured (client ID, secret, redirect URL)

- [ ] **Database**
  - [ ] `oauth_tokens` table exists (migration 000022)
  - [ ] Table has proper indexes
  - [ ] `token_blacklist` table exists (migration 000023)

- [ ] **Code Compilation**
  - [ ] No compilation errors
  - [ ] All imports resolved
  - [ ] Dependencies injected correctly

- [ ] **Logging**
  - [ ] Success logs appear: "Successfully stored OAuth2 token..."
  - [ ] Warning logs appear on failure (but don't crash app)
  - [ ] No error logs for normal operations

### Post-Deployment Checks

- [ ] **Functionality**
  - [ ] Users can login with Google
  - [ ] Users can login with GitHub
  - [ ] Users can login with Facebook
  - [ ] Tokens stored in database after login
  - [ ] Tokens are encrypted (not plain text)

- [ ] **Performance**
  - [ ] Login response time < 2 seconds
  - [ ] Database queries optimized (use indexes)
  - [ ] Redis caching working
  - [ ] No memory leaks

- [ ] **Security**
  - [ ] OAuth2 tokens never appear in logs
  - [ ] Tokens encrypted in database
  - [ ] Encryption key secured (not in version control)
  - [ ] HTTPS enforced for OAuth callbacks

## Common Issues & Solutions

### Issue 1: Token Not Stored
**Symptoms:** Console shows "Warning: failed to store OAuth2 token"

**Solutions:**
1. Check database connection
2. Verify `oauth_tokens` table exists
3. Check user_id is valid
4. Verify encryption key is set

### Issue 2: Encryption Error
**Symptoms:** "Failed to initialize token encryptor"

**Solutions:**
1. Set `security.token_encryption_key` in config.json
2. Ensure key is valid (any string works, but use 32+ chars)
3. Check config file is loaded correctly

### Issue 3: Provider Mismatch
**Symptoms:** Token stored with wrong provider

**Solutions:**
1. Verify `convertProviderToEntity()` mapping
2. Check provider parameter in callback request
3. Ensure provider string matches entity constants

### Issue 4: Token Not Found After Login
**Symptoms:** Token stored but can't retrieve

**Solutions:**
1. Check user_id matches
2. Verify provider string matches
3. Check Redis cache invalidation
4. Query database directly to verify storage

## Database Queries for Debugging

### Check Token Storage
```sql
-- Count tokens per user
SELECT user_id, provider, COUNT(*) as token_count
FROM oauth_tokens
GROUP BY user_id, provider;

-- Recent tokens
SELECT id, user_id, provider, created_at, expires_at
FROM oauth_tokens
ORDER BY created_at DESC
LIMIT 10;

-- Tokens for specific user
SELECT provider, created_at, expires_at
FROM oauth_tokens
WHERE user_id = 'USER_ID';
```

### Check Token Encryption
```sql
-- This should return gibberish, not readable tokens
SELECT LEFT(access_token, 50) as token_preview
FROM oauth_tokens
LIMIT 5;
```

### Check Token Expiry
```sql
-- Tokens expiring soon (within 5 minutes)
SELECT user_id, provider, expires_at
FROM oauth_tokens
WHERE expires_at <= DATE_ADD(NOW(), INTERVAL 5 MINUTE)
AND expires_at > NOW();

-- Expired tokens
SELECT COUNT(*) as expired_count
FROM oauth_tokens
WHERE expires_at < NOW();
```

## Performance Monitoring

### Metrics to Track

1. **Token Storage Success Rate**
   - Target: >99.9%
   - Alert if: <95%

2. **Token Storage Latency**
   - Target: <100ms
   - Alert if: >500ms

3. **Encryption Performance**
   - Target: >1M ops/sec
   - Current: 1.4M encrypt, 2.9M decrypt

4. **Database Query Time**
   - Insert: <50ms
   - Select: <10ms (with cache)
   - Select: <50ms (without cache)

## Security Audit Checklist

- [ ] Tokens encrypted with AES-256-GCM
- [ ] Encryption key stored securely (environment variable)
- [ ] Tokens never logged in plain text
- [ ] Database connection secured (SSL/TLS)
- [ ] Redis connection secured (AUTH enabled)
- [ ] OAuth2 callbacks use HTTPS
- [ ] CSRF protection on OAuth callback
- [ ] State parameter validated
- [ ] Token expiry enforced
- [ ] Refresh tokens rotated

## Next Steps

After successful testing:

1. **Complete TokenBlacklist Repositories**
   - Implement MySQL repository
   - Implement Redis repository
   - Update TokenService initialization

2. **Implement Token Refresh Endpoint**
   - Add `/api/tokens/refresh` endpoint
   - Check token expiry
   - Call OAuth2 provider for refresh

3. **Add Token Validation**
   - Check token exists in database
   - Verify not expired
   - Verify not blacklisted

4. **Integration Testing**
   - Test complete OAuth flow
   - Test token refresh flow
   - Test logout flow with blacklist

5. **Load Testing**
   - Concurrent OAuth logins
   - Token storage throughput
   - Cache performance

## Related Documentation

- [OAuth2 Callback Update Plan](./OAUTH2_CALLBACK_UPDATE.md)
- [TokenService Documentation](../domain/service/TOKEN_SERVICE.md)
- [Token Encryption Guide](../../pkg/crypto/README.md)
- [Sprint Progress](../../SPRINT.md)
