# Token Refresh Endpoint - Quick Reference

## Endpoint
```
POST /tokens/refresh
```

## Authentication
```http
Authorization: Bearer <your_jwt_token>
```

## Request Body
```json
{
  "user_id": "your-user-id",
  "provider": "google"
}
```

**Valid Providers:** `google`, `facebook`, `github`, `apple`

## Success Response
```json
{
  "success": true,
  "message": "Token is still valid, no refresh needed",
  "refreshed_from": "cache",
  "expires_in": 3599,
  "token": {
    "access_token": "ya29.a0AfH6SMB...",
    "token_type": "Bearer",
    "expires_at": "2025-10-23T15:30:00Z",
    "scope": "",
    "has_refresh_token": true
  }
}
```

## cURL Example
```bash
curl -X POST http://localhost:8080/tokens/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "user_id": "YOUR_USER_ID",
    "provider": "google"
  }'
```

## JavaScript Example
```javascript
const response = await fetch('http://localhost:8080/tokens/refresh', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${jwtToken}`
  },
  body: JSON.stringify({
    user_id: userId,
    provider: 'google'
  })
});

const data = await response.json();
console.log('Access Token:', data.token.access_token);
console.log('Expires in:', data.expires_in, 'seconds');
```

## Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Whether the request was successful |
| `message` | string | Human-readable status message |
| `refreshed_from` | string | `"cache"` or `"provider"` |
| `expires_in` | number | Seconds until token expires |
| `token.access_token` | string | OAuth2 access token |
| `token.token_type` | string | Usually `"Bearer"` |
| `token.expires_at` | string | ISO 8601 timestamp |
| `token.has_refresh_token` | boolean | Whether refresh token exists |

## Smart Refresh Behavior

- **Token expires in > 5 minutes:** Returns existing token from cache (fast)
- **Token expires in < 5 minutes:** Refreshes from OAuth provider (slower)
- **Token already expired:** Attempts refresh immediately

## Common Errors

### 401 Unauthorized
```json
{"error": "Unauthorized"}
```
**Solution:** Include valid JWT token in Authorization header

### 400 Bad Request
```json
{
  "success": false,
  "message": "user_id is required"
}
```
**Solution:** Include both `user_id` and `provider` in request body

### 404 Not Found
```json
{
  "success": false,
  "message": "No token found for this user and provider"
}
```
**Solution:** User must first authenticate via `/oauth2/login`

### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Failed to refresh token: invalid_grant"
}
```
**Solution:** Refresh token expired, user must re-authenticate

## Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success - token returned or refreshed |
| 400 | Bad request - invalid parameters |
| 401 | Unauthorized - missing or invalid JWT |
| 404 | Not found - no token for user/provider |
| 405 | Method not allowed - must use POST |
| 500 | Server error - token refresh failed |

## Best Practices

1. **Check expiry before calling:** Only refresh when needed
2. **Handle errors gracefully:** Redirect to login on 404/500
3. **Cache tokens client-side:** Reduce unnecessary refresh calls
4. **Use HTTPS in production:** Protect tokens in transit
5. **Implement retry logic:** Handle temporary provider failures

## Testing Commands

### Test with valid token
```bash
JWT_TOKEN="your_token"
USER_ID="your_id"

curl -X POST http://localhost:8080/tokens/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{"user_id":"'$USER_ID'","provider":"google"}' | jq
```

### Test without authentication
```bash
curl -X POST http://localhost:8080/tokens/refresh \
  -H "Content-Type: application/json" \
  -d '{"user_id":"test","provider":"google"}' | jq
```

### Test with invalid provider
```bash
curl -X POST http://localhost:8080/tokens/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{"user_id":"'$USER_ID'","provider":"invalid"}' | jq
```

## Integration Flow

```
1. User logs in via /oauth2/login
   → OAuth2Handler stores encrypted tokens
   
2. User uses the application
   → JWT token remains valid
   
3. Access token approaches expiration (< 5 minutes)
   → Frontend calls /tokens/refresh
   
4. TokenHandler checks and refreshes if needed
   → Returns fresh access token
   
5. User continues using application
   → No re-authentication needed
```

## Documentation

- **Full API Docs:** `internal/delivery/http/TOKEN_REFRESH_ENDPOINT.md`
- **Implementation Summary:** `TOKEN_REFRESH_IMPLEMENTATION_SUMMARY.md`
- **TokenService API:** `internal/domain/service/TOKEN_SERVICE.md`
