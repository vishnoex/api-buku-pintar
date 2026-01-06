package http

import (
	"buku-pintar/internal/constant"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"encoding/json"
	"fmt"
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

// TokenValidationRequest represents the request body for token validation
type TokenValidationRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	Provider string `json:"provider" validate:"required,oneof=google facebook github apple"`
}

// TokenValidationResponse represents the response for token validation
type TokenValidationResponse struct {
	Success       bool      `json:"success"`
	Message       string    `json:"message"`
	IsValid       bool      `json:"is_valid"`
	IsExpired     bool      `json:"is_expired"`
	IsBlacklisted bool      `json:"is_blacklisted"`
	NeedsRefresh  bool      `json:"needs_refresh"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	ExpiresIn     int64     `json:"expires_in,omitempty"` // seconds until expiration
}

// LogoutRequest represents the request body for logout
type LogoutRequest struct {
	UserID              string `json:"user_id" validate:"required"`
	Provider            string `json:"provider,omitempty"`                                  // Optional: specific provider to logout from
	RevokeAllProviders  bool   `json:"revoke_all_providers,omitempty"`                      // If true, logout from all providers
	BlacklistJWT        bool   `json:"blacklist_jwt,omitempty"`                             // If true, also blacklist the JWT token
	JWTToken            string `json:"jwt_token,omitempty"`                                 // JWT token to blacklist (if blacklist_jwt is true)
}

// LogoutResponse represents the response for logout
type LogoutResponse struct {
	Success            bool      `json:"success"`
	Message            string    `json:"message"`
	RevokedProviders   []string  `json:"revoked_providers,omitempty"`   // List of providers that were logged out
	TokensRevoked      int       `json:"tokens_revoked"`                // Number of OAuth tokens revoked
	JWTBlacklisted     bool      `json:"jwt_blacklisted,omitempty"`     // Whether JWT was blacklisted
	LoggedOutAt        time.Time `json:"logged_out_at"`                 // Timestamp of logout
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

// ValidateToken handles POST /api/tokens/validate
// Validates an OAuth2 token against the database, checking expiration and blacklist status
func (h *TokenHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, constant.ERR_METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req TokenValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[TokenHandler] Failed to decode validation request: %v", err)
		respondWithValidationError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.UserID == "" {
		respondWithValidationError(w, "user_id is required", http.StatusBadRequest)
		return
	}
	if req.Provider == "" {
		respondWithValidationError(w, "provider is required", http.StatusBadRequest)
		return
	}

	// Convert provider string to entity.OAuthProvider
	provider, err := parseProvider(req.Provider)
	if err != nil {
		respondWithValidationError(w, "Invalid provider. Must be one of: google, facebook, github, apple", http.StatusBadRequest)
		return
	}

	// Get token from database
	ctx := r.Context()
	token, err := h.tokenService.GetOAuthToken(ctx, req.UserID, provider)
	if err != nil {
		log.Printf("[TokenHandler] Failed to get token for user %s with provider %s: %v", 
			req.UserID, req.Provider, err)
		respondWithValidationError(w, "Failed to retrieve token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if token exists
	if token == nil {
		log.Printf("[TokenHandler] No token found for user %s with provider %s", req.UserID, req.Provider)
		response := TokenValidationResponse{
			Success:       false,
			Message:       "No token found for this user and provider",
			IsValid:       false,
			IsExpired:     false,
			IsBlacklisted: false,
			NeedsRefresh:  false,
		}
		sendValidationResponse(w, response, http.StatusNotFound)
		return
	}

	// Check if token is expired
	isExpired := token.IsExpired()
	needsRefresh := token.NeedsRefresh()

	// Get decrypted access token to check blacklist
	accessToken, _, err := h.tokenService.GetDecryptedOAuthToken(ctx, req.UserID, provider)
	if err != nil {
		log.Printf("[TokenHandler] Failed to decrypt token: %v", err)
		respondWithValidationError(w, "Failed to decrypt token", http.StatusInternalServerError)
		return
	}

	// Check if token is blacklisted (for JWT tokens)
	isBlacklisted := false
	if accessToken != "" {
		isBlacklisted, err = h.tokenService.IsTokenBlacklisted(ctx, accessToken)
		if err != nil {
			log.Printf("[TokenHandler] Failed to check blacklist status: %v", err)
			// Continue with validation - blacklist check failure shouldn't block
		}
	}

	// Calculate seconds until expiration
	var expiresIn int64 = 0
	var expiresAt *time.Time = nil
	if token.ExpiresAt != nil {
		expiresAt = token.ExpiresAt
		if !isExpired {
			expiresIn = int64(time.Until(*token.ExpiresAt).Seconds())
			if expiresIn < 0 {
				expiresIn = 0
			}
		}
	}

	// Determine if token is valid
	isValid := !isExpired && !isBlacklisted && token != nil

	// Build response
	response := TokenValidationResponse{
		Success:       true,
		Message:       getValidationMessage(isValid, isExpired, isBlacklisted, needsRefresh),
		IsValid:       isValid,
		IsExpired:     isExpired,
		IsBlacklisted: isBlacklisted,
		NeedsRefresh:  needsRefresh,
		ExpiresAt:     expiresAt,
		ExpiresIn:     expiresIn,
	}

	// Send response with appropriate status code
	statusCode := http.StatusOK
	if !isValid {
		statusCode = http.StatusUnauthorized
	}

	sendValidationResponse(w, response, statusCode)

	log.Printf("[TokenHandler] Token validation complete for user %s with provider %s (valid: %v, expired: %v, blacklisted: %v)", 
		req.UserID, req.Provider, isValid, isExpired, isBlacklisted)
}

// Logout handles POST /api/tokens/logout
// Logs out a user by invalidating their OAuth2 tokens and optionally blacklisting JWT tokens
func (h *TokenHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, constant.ERR_METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[TokenHandler] Failed to decode logout request: %v", err)
		respondWithLogoutError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.UserID == "" {
		respondWithLogoutError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var revokedProviders []string
	var tokensRevoked int
	var jwtBlacklisted bool

	// Determine logout strategy
	if req.RevokeAllProviders {
		// Logout from all providers
		log.Printf("[TokenHandler] Revoking all access for user %s", req.UserID)
		
		// Get all tokens for the user first to track which providers were revoked
		tokens, err := h.tokenService.GetOAuthTokensByUserID(ctx, req.UserID)
		if err != nil {
			log.Printf("[TokenHandler] Failed to get user tokens: %v", err)
			respondWithLogoutError(w, "Failed to retrieve user tokens", http.StatusInternalServerError)
			return
		}

		// Track unique providers
		providerMap := make(map[entity.OAuthProvider]bool)
		for _, token := range tokens {
			providerMap[token.Provider] = true
		}

		// Revoke all user access
		err = h.tokenService.RevokeAllUserAccess(ctx, req.UserID, entity.ReasonLogout)
		if err != nil {
			log.Printf("[TokenHandler] Failed to revoke all user access: %v", err)
			respondWithLogoutError(w, "Failed to logout from all providers: "+err.Error(), http.StatusInternalServerError)
			return
		}

		tokensRevoked = len(tokens)
		for provider := range providerMap {
			revokedProviders = append(revokedProviders, string(provider))
		}

		log.Printf("[TokenHandler] Successfully revoked all access for user %s (%d tokens from %d providers)", 
			req.UserID, tokensRevoked, len(revokedProviders))

	} else if req.Provider != "" {
		// Logout from specific provider
		provider, err := parseProvider(req.Provider)
		if err != nil {
			respondWithLogoutError(w, "Invalid provider. Must be one of: google, facebook, github, apple", http.StatusBadRequest)
			return
		}

		log.Printf("[TokenHandler] Revoking access for user %s from provider %s", req.UserID, req.Provider)

		// Revoke user-provider access
		err = h.tokenService.RevokeUserProviderAccess(ctx, req.UserID, provider, entity.ReasonLogout)
		if err != nil {
			log.Printf("[TokenHandler] Failed to revoke provider access: %v", err)
			respondWithLogoutError(w, "Failed to logout from provider: "+err.Error(), http.StatusInternalServerError)
			return
		}

		tokensRevoked = 1
		revokedProviders = append(revokedProviders, req.Provider)

		log.Printf("[TokenHandler] Successfully revoked access for user %s from provider %s", req.UserID, req.Provider)

	} else {
		// No provider specified and not revoking all - invalid request
		respondWithLogoutError(w, "Either 'provider' or 'revoke_all_providers' must be specified", http.StatusBadRequest)
		return
	}

	// Optionally blacklist JWT token
	if req.BlacklistJWT && req.JWTToken != "" {
		log.Printf("[TokenHandler] Blacklisting JWT token for user %s", req.UserID)

		// Calculate JWT expiration (typically 24 hours from now, adjust as needed)
		jwtExpiry := time.Now().Add(24 * time.Hour)

		err := h.tokenService.BlacklistToken(ctx, req.JWTToken, &req.UserID, entity.ReasonLogout, jwtExpiry)
		if err != nil {
			log.Printf("[TokenHandler] Failed to blacklist JWT token: %v", err)
			// Don't fail the entire logout if JWT blacklist fails - OAuth tokens are already revoked
			// Just log the error and continue
		} else {
			jwtBlacklisted = true
			log.Printf("[TokenHandler] Successfully blacklisted JWT token for user %s", req.UserID)
		}
	}

	// Build success response
	message := buildLogoutMessage(tokensRevoked, len(revokedProviders), jwtBlacklisted)
	response := LogoutResponse{
		Success:          true,
		Message:          message,
		RevokedProviders: revokedProviders,
		TokensRevoked:    tokensRevoked,
		JWTBlacklisted:   jwtBlacklisted,
		LoggedOutAt:      time.Now(),
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[TokenHandler] Failed to encode logout response: %v", err)
	}

	log.Printf("[TokenHandler] Logout completed for user %s (tokens_revoked: %d, providers: %v, jwt_blacklisted: %v)", 
		req.UserID, tokensRevoked, revokedProviders, jwtBlacklisted)
}

// buildLogoutMessage creates an appropriate logout message based on what was revoked
func buildLogoutMessage(tokensRevoked, providersRevoked int, jwtBlacklisted bool) string {
	if tokensRevoked == 0 {
		return "No active tokens found to revoke"
	}

	message := fmt.Sprintf("Successfully logged out (%d token(s) from %d provider(s) revoked)", 
		tokensRevoked, providersRevoked)

	if jwtBlacklisted {
		message += ", JWT token blacklisted"
	}

	return message
}

// respondWithLogoutError sends a JSON error response for logout
func respondWithLogoutError(w http.ResponseWriter, message string, statusCode int) {
	response := LogoutResponse{
		Success:       false,
		Message:       message,
		TokensRevoked: 0,
		LoggedOutAt:   time.Now(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[TokenHandler] Failed to encode error response: %v", err)
	}
}

// getValidationMessage returns an appropriate message based on token status
func getValidationMessage(isValid, isExpired, isBlacklisted, needsRefresh bool) string {
	if isBlacklisted {
		return "Token has been revoked and is no longer valid"
	}
	if isExpired {
		return "Token has expired"
	}
	if !isValid {
		return "Token is not valid"
	}
	if needsRefresh {
		return "Token is valid but should be refreshed soon"
	}
	return "Token is valid"
}

// sendValidationResponse sends a JSON validation response
func sendValidationResponse(w http.ResponseWriter, response TokenValidationResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[TokenHandler] Failed to encode validation response: %v", err)
	}
}

// respondWithValidationError sends a JSON error response for validation
func respondWithValidationError(w http.ResponseWriter, message string, statusCode int) {
	response := TokenValidationResponse{
		Success: false,
		Message: message,
		IsValid: false,
	}
	sendValidationResponse(w, response, statusCode)
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
