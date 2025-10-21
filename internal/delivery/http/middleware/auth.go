package middleware

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/pkg/oauth2"
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ContextKey represents a type for context keys
type ContextKey string

const (
	UserContextKey ContextKey = "user"
)

type AuthMiddleware struct {
	oauth2Service *oauth2.OAuth2Service
}

func NewAuthMiddleware(oauth2Service *oauth2.OAuth2Service) *AuthMiddleware {
	return &AuthMiddleware{
		oauth2Service: oauth2Service,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Check if the header has the Bearer prefix
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Try OAuth2 token validation
		user, err := m.authenticateWithOAuth2(r.Context(), token)
		if err == nil && user != nil {
			// OAuth2 authentication successful
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Authentication failed
		http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
	})
}

// authenticateWithOAuth2 attempts to authenticate using OAuth2 token
func (m *AuthMiddleware) authenticateWithOAuth2(ctx context.Context, accessToken string) (*entity.User, error) {
	// For OAuth2, we would typically validate the token with the provider
	// and retrieve user information. For now, we'll return nil as this
	// would require additional implementation to validate OAuth2 tokens.
	// In a production environment, you might want to:
	// 1. Store OAuth2 tokens in a secure way
	// 2. Validate tokens with the respective providers
	// 3. Implement token refresh logic
	
	return nil, fmt.Errorf("OAuth2 token validation not implemented")
}

// GetUserFromContext retrieves the user from the context
func GetUserFromContext(ctx context.Context) (*entity.User, error) {
	user, ok := ctx.Value(UserContextKey).(*entity.User)
	if !ok {
		return nil, nil
	}
	return user, nil
} 