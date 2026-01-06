package middleware

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"context"
	"log"
	"net/http"
	"strings"
)

// TokenValidationMiddleware provides database-backed token validation
// Checks if OAuth2 tokens exist in database, are not expired, and not blacklisted
type TokenValidationMiddleware struct {
	tokenService service.TokenService
	userService  interface{} // Will be used for user lookups if needed
}

// NewTokenValidationMiddleware creates a new token validation middleware
func NewTokenValidationMiddleware(tokenService service.TokenService) *TokenValidationMiddleware {
	return &TokenValidationMiddleware{
		tokenService: tokenService,
	}
}

// ValidateToken is middleware that validates OAuth2 tokens against the database
// It checks:
// 1. Token exists in database
// 2. Token is not expired
// 3. Token is not blacklisted (for JWT tokens)
func (m *TokenValidationMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("[TokenValidation] Missing Authorization header")
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Extract Bearer token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			log.Printf("[TokenValidation] Invalid authorization header format")
			http.Error(w, "Invalid authorization header format. Use: Bearer <token>", http.StatusUnauthorized)
			return
		}

		// Check if token is blacklisted (for JWT tokens)
		isBlacklisted, err := m.tokenService.IsTokenBlacklisted(r.Context(), token)
		if err != nil {
			log.Printf("[TokenValidation] Error checking token blacklist: %v", err)
			// Continue - blacklist check failure shouldn't block valid tokens
		} else if isBlacklisted {
			log.Printf("[TokenValidation] Token is blacklisted")
			http.Error(w, "Token has been revoked", http.StatusUnauthorized)
			return
		}

		// For OAuth2 tokens, we need user_id and provider from context
		// This middleware should be used after basic auth middleware that sets user context
		// If we have user context, we can validate the OAuth2 token in database
		userID := r.Context().Value(UserIDContextKey)
		if userID != nil {
			userIDStr, ok := userID.(string)
			if ok && userIDStr != "" {
				// We have user ID, now we need to validate their OAuth tokens
				// Get all tokens for this user and check if any are still valid
				tokens, err := m.tokenService.GetOAuthTokensByUserID(r.Context(), userIDStr)
				if err != nil {
					log.Printf("[TokenValidation] Error fetching user tokens: %v", err)
				} else if len(tokens) == 0 {
					log.Printf("[TokenValidation] No tokens found for user %s", userIDStr)
					http.Error(w, "No valid authentication tokens found", http.StatusUnauthorized)
					return
				} else {
					// Check if any token is valid (not expired)
					hasValidToken := false
					for _, t := range tokens {
						if !t.IsExpired() {
							hasValidToken = true
							break
						}
					}
					
					if !hasValidToken {
						log.Printf("[TokenValidation] All tokens expired for user %s", userIDStr)
						http.Error(w, "All authentication tokens have expired", http.StatusUnauthorized)
						return
					}
				}
			}
		}

		// Token validation passed
		log.Printf("[TokenValidation] Token validation successful")
		next.ServeHTTP(w, r)
	})
}

// ValidateOAuthToken validates a specific OAuth2 token for a user and provider
// This is more strict validation that requires user_id and provider
func (m *TokenValidationMiddleware) ValidateOAuthToken(userID string, provider entity.OAuthProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get the token from database
		token, err := m.tokenService.GetOAuthToken(ctx, userID, provider)
		if err != nil {
			log.Printf("[TokenValidation] Error fetching OAuth token for user %s, provider %s: %v", 
				userID, provider, err)
			http.Error(w, "Failed to validate token", http.StatusInternalServerError)
			return
		}

		if token == nil {
			log.Printf("[TokenValidation] No token found for user %s, provider %s", userID, provider)
			http.Error(w, "No authentication token found", http.StatusUnauthorized)
			return
		}

		// Check if token is expired
		if token.IsExpired() {
			log.Printf("[TokenValidation] Token expired for user %s, provider %s", userID, provider)
			http.Error(w, "Authentication token has expired", http.StatusUnauthorized)
			return
		}

		// Check if token needs refresh (optional warning)
		if token.NeedsRefresh() {
			log.Printf("[TokenValidation] Warning: Token for user %s, provider %s needs refresh", 
				userID, provider)
			// Add header to indicate token should be refreshed soon
			w.Header().Set("X-Token-Refresh-Needed", "true")
		}

		// Token is valid, add token info to context
		ctx = context.WithValue(ctx, TokenContextKey, token)
		r = r.WithContext(ctx)
		
		log.Printf("[TokenValidation] OAuth token validation successful for user %s, provider %s", 
			userID, provider)
	})
}

// QuickBlacklistCheck performs a fast blacklist check without full token validation
// Useful for high-performance endpoints where full validation isn't needed
func (m *TokenValidationMiddleware) QuickBlacklistCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			next.ServeHTTP(w, r)
			return
		}

		// Quick blacklist check using Redis cache
		isBlacklisted, err := m.tokenService.IsTokenBlacklisted(r.Context(), token)
		if err != nil {
			log.Printf("[TokenValidation] Blacklist check error: %v", err)
			// On error, allow request to continue
			next.ServeHTTP(w, r)
			return
		}

		if isBlacklisted {
			log.Printf("[TokenValidation] Token is blacklisted")
			http.Error(w, "Token has been revoked", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Context keys for token validation
const (
	TokenContextKey   ContextKey = "oauth_token"
	UserIDContextKey  ContextKey = "user_id"
	ProviderContextKey ContextKey = "provider"
)

// GetTokenFromContext retrieves the OAuth token from the request context
func GetTokenFromContext(ctx context.Context) (*entity.OAuthToken, bool) {
	token, ok := ctx.Value(TokenContextKey).(*entity.OAuthToken)
	return token, ok
}

// TokenValidationResult represents the result of token validation
type TokenValidationResult struct {
	IsValid       bool
	IsExpired     bool
	IsBlacklisted bool
	NeedsRefresh  bool
	ErrorMessage  string
	Token         *entity.OAuthToken
}

// ValidateTokenComprehensive performs comprehensive token validation and returns detailed results
// This is useful for endpoints that need to know why validation failed
func (m *TokenValidationMiddleware) ValidateTokenComprehensive(
	ctx context.Context,
	userID string,
	provider entity.OAuthProvider,
	accessToken string,
) *TokenValidationResult {
	result := &TokenValidationResult{
		IsValid: false,
	}

	// Check if JWT token is blacklisted
	isBlacklisted, err := m.tokenService.IsTokenBlacklisted(ctx, accessToken)
	if err != nil {
		result.ErrorMessage = "Failed to check token blacklist"
		return result
	}
	result.IsBlacklisted = isBlacklisted
	if isBlacklisted {
		result.ErrorMessage = "Token has been revoked"
		return result
	}

	// Get OAuth token from database
	token, err := m.tokenService.GetOAuthToken(ctx, userID, provider)
	if err != nil {
		result.ErrorMessage = "Failed to retrieve token from database"
		return result
	}
	if token == nil {
		result.ErrorMessage = "Token not found in database"
		return result
	}
	result.Token = token

	// Check if token is expired
	result.IsExpired = token.IsExpired()
	if result.IsExpired {
		result.ErrorMessage = "Token has expired"
		return result
	}

	// Check if token needs refresh
	result.NeedsRefresh = token.NeedsRefresh()

	// Token is valid
	result.IsValid = true
	return result
}
