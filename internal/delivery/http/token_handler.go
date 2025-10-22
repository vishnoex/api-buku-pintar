package http

import (
	"buku-pintar/internal/constant"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// TokenHandler handles token-related HTTP requests
type TokenHandler struct {
	tokenService service.TokenService
}

// NewTokenHandler creates a new TokenHandler instance
func NewTokenHandler(tokenService service.TokenService) *TokenHandler {
	return &TokenHandler{
		tokenService: tokenService,
	}
}

// TokenRefreshRequest represents the request body for token refresh
type TokenRefreshRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	Provider string `json:"provider" validate:"required,oneof=google facebook github apple"`
}

// TokenRefreshResponse represents the response for token refresh
type TokenRefreshResponse struct {
	Success       bool                `json:"success"`
	Message       string              `json:"message"`
	Token         *TokenInfo          `json:"token,omitempty"`
	RefreshedFrom string              `json:"refreshed_from,omitempty"` // "cache" or "provider"
	ExpiresIn     int64               `json:"expires_in,omitempty"`     // seconds until expiration
}

// TokenInfo contains token details in the response
type TokenInfo struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scope        string    `json:"scope,omitempty"`
	HasRefresh   bool      `json:"has_refresh_token"`
}

// RefreshToken handles POST /api/tokens/refresh
// Refreshes an OAuth2 token if needed, or returns the existing valid token
func (h *TokenHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, constant.ERR_METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req TokenRefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[TokenHandler] Failed to decode refresh request: %v", err)
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.UserID == "" {
		respondWithError(w, "user_id is required", http.StatusBadRequest)
		return
	}
	if req.Provider == "" {
		respondWithError(w, "provider is required", http.StatusBadRequest)
		return
	}

	// Convert provider string to entity.OAuthProvider
	provider, err := parseProvider(req.Provider)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call TokenService to refresh token if needed
	ctx := r.Context()
	token, err := h.tokenService.RefreshTokenIfNeeded(ctx, req.UserID, provider)
	if err != nil {
		log.Printf("[TokenHandler] Failed to refresh token for user %s with provider %s: %v", 
			req.UserID, req.Provider, err)
		respondWithError(w, "Failed to refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if token was found
	if token == nil {
		respondWithError(w, "No token found for this user and provider", http.StatusNotFound)
		return
	}

	// Calculate seconds until expiration
	expiresIn := int64(0)
	if !token.Expiry.IsZero() {
		expiresIn = int64(time.Until(token.Expiry).Seconds())
		if expiresIn < 0 {
			expiresIn = 0
		}
	}

	// Determine if token was refreshed (if expiry is very recent, it was likely refreshed)
	wasRefreshed := time.Since(token.Expiry) < 1*time.Minute

	// Build response
	response := TokenRefreshResponse{
		Success:       true,
		Message:       getMessage(wasRefreshed),
		RefreshedFrom: getSource(wasRefreshed),
		ExpiresIn:     expiresIn,
		Token: &TokenInfo{
			AccessToken:  token.AccessToken,
			TokenType:    token.TokenType,
			ExpiresAt:    token.Expiry,
			Scope:        "", // oauth2.Token doesn't have Scope field directly
			HasRefresh:   token.RefreshToken != "",
		},
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[TokenHandler] Failed to encode response: %v", err)
	}

	log.Printf("[TokenHandler] Successfully handled token refresh for user %s with provider %s (was_refreshed: %v)", 
		req.UserID, req.Provider, wasRefreshed)
}

// parseProvider converts a string provider name to entity.OAuthProvider
func parseProvider(provider string) (entity.OAuthProvider, error) {
	switch provider {
	case "google":
		return entity.ProviderGoogle, nil
	case "facebook":
		return entity.ProviderFacebook, nil
	case "github":
		return entity.ProviderGithub, nil
	case "apple":
		return entity.ProviderApple, nil
	default:
		return "", http.ErrNotSupported // Using standard error
	}
}

// getMessage returns an appropriate message based on whether token was refreshed
func getMessage(wasRefreshed bool) string {
	if wasRefreshed {
		return "Token successfully refreshed from OAuth provider"
	}
	return "Token is still valid, no refresh needed"
}

// getSource returns where the token came from
func getSource(wasRefreshed bool) string {
	if wasRefreshed {
		return "provider"
	}
	return "cache"
}

// respondWithError sends a JSON error response
func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	response := TokenRefreshResponse{
		Success: false,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
