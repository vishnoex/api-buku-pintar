package middleware

import (
	"buku-pintar/internal/domain/entity"
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

// ClaimKey represents a type for Firebase token claim keys
type ClaimKey string

const (
	EmailClaim ClaimKey = "email"
	NameClaim  ClaimKey = "name"
)

type AuthMiddleware struct {
	firebaseAuth *auth.Client
}

func NewAuthMiddleware(firebaseAuth *auth.Client) *AuthMiddleware {
	return &AuthMiddleware{
		firebaseAuth: firebaseAuth,
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
		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		if idToken == authHeader {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Verify the ID token
		token, err := m.firebaseAuth.VerifyIDToken(r.Context(), idToken)
		if err != nil {
			http.Error(w, "Invalid ID token", http.StatusUnauthorized)
			return
		}

		// Create user context
		user := &entity.User{
			ID:     token.UID, // Use Firebase UID as the user ID
			Email:  token.Claims[string(EmailClaim)].(string),
			Name:   token.Claims[string(NameClaim)].(string),
			Role:   entity.RoleReader, // Default role
			Status: entity.StatusActive, // Default status
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext retrieves the user from the context
func GetUserFromContext(ctx context.Context) (*entity.User, error) {
	user, ok := ctx.Value("user").(*entity.User)
	if !ok {
		return nil, nil
	}
	return user, nil
} 