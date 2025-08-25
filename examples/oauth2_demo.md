# OAuth2 Integration Demo

This document demonstrates how to use the OAuth2 integration in the Buku Pintar API.

## Prerequisites

1. Configure OAuth2 providers in your `config.json`
2. Set up OAuth2 applications in Google, GitHub, and Facebook developer consoles
3. Ensure the API is running on the configured port

## OAuth2 Flow Examples

### 1. Google OAuth2 Login

#### Step 1: Initiate Login
```bash
curl -X POST http://localhost:8080/oauth2/login \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "google"
  }'
```

**Response:**
```json
{
  "auth_url": "https://accounts.google.com/oauth2/authorize?client_id=...&redirect_uri=...&scope=...&state=...",
  "state": "random_state_string"
}
```

#### Step 2: User Authorization
The user is redirected to Google's authorization page where they grant permissions.

#### Step 3: Handle Callback
Google redirects back to your application with an authorization code:
```
GET /oauth2/google/redirect?code=AUTHORIZATION_CODE&state=STATE
```

### 2. GitHub OAuth2 Login

#### Step 1: Initiate Login
```bash
curl -X POST http://localhost:8080/oauth2/login \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "github"
  }'
```

#### Step 2: Handle Callback
GitHub redirects back to your application:
```
GET /oauth2/github/redirect?code=AUTHORIZATION_CODE&state=STATE
```

### 3. Facebook OAuth2 Login

#### Step 1: Initiate Login
```bash
curl -X POST http://localhost:8080/oauth2/login \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "facebook"
  }'
```

#### Step 2: Handle Callback
Facebook redirects back to your application:
```
GET /oauth2/facebook/redirect?code=AUTHORIZATION_CODE&state=STATE
```

## Get Available Providers

```bash
curl -X GET http://localhost:8080/oauth2/providers
```

**Response:**
```json
{
  "providers": ["google", "github", "facebook"]
}
```

## Using OAuth2 Tokens

After successful OAuth2 authentication, you can use the returned access token to access protected endpoints:

```bash
curl -X GET http://localhost:8080/users \
  -H "Authorization: Bearer YOUR_OAUTH2_ACCESS_TOKEN"
```

## Frontend Integration Example

### React Component Example

```jsx
import React, { useState } from 'react';

const OAuth2Login = () => {
  const [loading, setLoading] = useState(false);

  const handleOAuth2Login = async (provider) => {
    setLoading(true);
    try {
      const response = await fetch('/oauth2/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ provider }),
      });
      
      const data = await response.json();
      
      // Redirect user to OAuth2 provider
      window.location.href = data.auth_url;
    } catch (error) {
      console.error('OAuth2 login failed:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2>Login with OAuth2</h2>
      <button 
        onClick={() => handleOAuth2Login('google')}
        disabled={loading}
      >
        Login with Google
      </button>
      
      <button 
        onClick={() => handleOAuth2Login('github')}
        disabled={loading}
      >
        Login with GitHub
      </button>
      
      <button 
        onClick={() => handleOAuth2Login('facebook')}
        disabled={loading}
      >
        Login with Facebook
      </button>
    </div>
  );
};

export default OAuth2Login;
```

### Handle OAuth2 Callback

```jsx
import React, { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';

const OAuth2Callback = () => {
  const location = useLocation();
  const [user, setUser] = useState(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const code = params.get('code');
    const state = params.get('state');

    if (code) {
      // Exchange code for token
      exchangeCodeForToken(code, state);
    }
  }, [location]);

  const exchangeCodeForToken = async (code, state) => {
    try {
      const response = await fetch('/oauth2/callback', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ code, state }),
      });

      if (response.ok) {
        const data = await response.json();
        setUser(data.user);
        
        // Store token in localStorage or secure storage
        localStorage.setItem('access_token', data.access_token);
        
        // Redirect to dashboard or home page
        window.location.href = '/dashboard';
      } else {
        setError('Authentication failed');
      }
    } catch (error) {
      setError('Authentication error: ' + error.message);
    }
  };

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (user) {
    return <div>Welcome, {user.name}!</div>;
  }

  return <div>Processing authentication...</div>;
};

export default OAuth2Callback;
```

## Error Handling

The OAuth2 integration includes comprehensive error handling:

- **Invalid Provider**: Returns 400 Bad Request for unsupported OAuth2 providers
- **Missing Configuration**: Returns 500 Internal Server Error if OAuth2 provider is not configured
- **Invalid Authorization Code**: Returns 500 Internal Server Error if code exchange fails
- **User Info Retrieval Failure**: Returns 500 Internal Server Error if user information cannot be retrieved

## Security Considerations

1. **State Parameter**: Always validate the state parameter to prevent CSRF attacks
2. **HTTPS**: Use HTTPS in production to secure OAuth2 communications
3. **Token Storage**: Store OAuth2 tokens securely (e.g., in HTTP-only cookies or secure storage)
4. **Scope Limitation**: Request only necessary OAuth2 scopes for your application
5. **Token Validation**: Implement proper token validation and refresh mechanisms

## Testing

You can test the OAuth2 integration using the provided test suite:

```bash
# Run all tests
go test ./...

# Run OAuth2 tests specifically
go test ./pkg/oauth2/...

# Run with verbose output
go test -v ./pkg/oauth2/...
```

## Troubleshooting

### Common Issues

1. **Provider Not Configured**: Ensure OAuth2 provider credentials are properly set in `config.json`
2. **Redirect URI Mismatch**: Verify redirect URIs match between your app and OAuth2 provider settings
3. **Invalid Client ID/Secret**: Double-check OAuth2 application credentials
4. **Scope Issues**: Ensure requested OAuth2 scopes are enabled in provider settings

### Debug Mode

Enable debug logging by setting the environment variable:
```bash
export DEBUG=true
go run cmd/api/main.go
```

This will provide detailed information about OAuth2 flows and token exchanges.
