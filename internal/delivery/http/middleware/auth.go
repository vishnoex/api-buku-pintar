package middleware

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/pkg/supabase"
	"context"
	"errors"
	"net/http"
	"strings"
)

// ContextKey represents a type for context keys
type ContextKey string

const (
	UserContextKey ContextKey = "user"
)

type AuthMiddleware struct {
	supabaseAuth *supabase.Authenticator
	userRepo     repository.UserRepository
}

func NewAuthMiddleware(supabaseAuth *supabase.Authenticator, userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		supabaseAuth: supabaseAuth,
		userRepo:     userRepo,
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

		user, err := m.authenticateWithSupabase(r.Context(), token)
		if err != nil {
			http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) authenticateWithSupabase(ctx context.Context, accessToken string) (*entity.User, error) {
	claims, err := m.supabaseAuth.VerifyToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	user, err := m.userRepo.GetByID(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	return nil, errors.New("verified local user not found")
}

// GetUserFromContext retrieves the user from the context
func GetUserFromContext(ctx context.Context) (*entity.User, error) {
	user, ok := ctx.Value(UserContextKey).(*entity.User)
	if !ok {
		return nil, nil
	}
	return user, nil
}
