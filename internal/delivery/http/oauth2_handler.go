package http

import (
	"buku-pintar/internal/constant"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"buku-pintar/internal/usecase"
	"buku-pintar/pkg/oauth2"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type OAuth2Handler struct {
	oauth2Service *oauth2.OAuth2Service
	userUsecase   usecase.UserUsecase
	tokenService  service.TokenService
}

func NewOAuth2Handler(
	oauth2Service *oauth2.OAuth2Service,
	userUsecase usecase.UserUsecase,
	tokenService service.TokenService,
) *OAuth2Handler {
	return &OAuth2Handler{
		oauth2Service: oauth2Service,
		userUsecase:   userUsecase,
		tokenService:  tokenService,
	}
}

// convertProviderToEntity converts oauth2.Provider to entity.OAuthProvider
func convertProviderToEntity(provider oauth2.Provider) entity.OAuthProvider {
	switch provider {
	case oauth2.ProviderGoogle:
		return entity.ProviderGoogle
	case oauth2.ProviderGitHub:
		return entity.ProviderGithub // Note: entity uses "Github" not "GitHub"
	case oauth2.ProviderFacebook:
		return entity.ProviderFacebook
	default:
		return entity.ProviderGoogle // Default fallback
	}
}

// OAuth2LoginRequest represents the OAuth2 login request
type OAuth2LoginRequest struct {
	Provider string `json:"provider"`
	State    string `json:"state,omitempty"`
}

// OAuth2CallbackRequest represents the OAuth2 callback request
type OAuth2CallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// OAuth2Response represents the OAuth2 response
type OAuth2Response struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

// OAuth2TokenResponse represents the OAuth2 token response
type OAuth2TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	User         *entity.User `json:"user"`
}

// generateState generates a random state parameter for OAuth2
func generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Login initiates OAuth2 login flow
func (h *OAuth2Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, constant.ERR_METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	var req OAuth2LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate provider
	provider := oauth2.Provider(req.Provider)
	if provider != oauth2.ProviderGoogle && 
	   provider != oauth2.ProviderGitHub && 
	   provider != oauth2.ProviderFacebook {
		http.Error(w, "Invalid OAuth2 provider", http.StatusBadRequest)
		return
	}

	// Generate state parameter
	state, err := generateState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	// Get authorization URL
	authURL, err := h.oauth2Service.GetAuthURL(provider, state)
	if err != nil {
		http.Error(w, "Failed to generate auth URL", http.StatusInternalServerError)
		return
	}

	response := OAuth2Response{
		AuthURL: authURL,
		State:   state,
	}

	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, constant.ERR_ENCODING_RESP, http.StatusInternalServerError)
		return
	}
}

// Callback handles OAuth2 callback and user authentication
func (h *OAuth2Handler) Callback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, constant.ERR_METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	var req OAuth2CallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract provider from query parameter or header
	providerStr := r.URL.Query().Get("provider")
	if providerStr == "" {
		providerStr = r.Header.Get("X-OAuth-Provider")
	}
	if providerStr == "" {
		http.Error(w, "Provider not specified", http.StatusBadRequest)
		return
	}

	provider := oauth2.Provider(providerStr)
	if provider != oauth2.ProviderGoogle && 
	   provider != oauth2.ProviderGitHub && 
	   provider != oauth2.ProviderFacebook {
		http.Error(w, "Invalid OAuth2 provider", http.StatusBadRequest)
		return
	}

	// Exchange authorization code for access token
	token, err := h.oauth2Service.ExchangeCode(r.Context(), provider, req.Code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	// Get user information from OAuth2 provider
	userInfo, err := h.oauth2Service.GetUserInfo(r.Context(), provider, token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Check if user exists in our system
	existingUser, err := h.userUsecase.GetUserByEmail(r.Context(), userInfo.Email)
	if err != nil {
		http.Error(w, "Failed to check existing user", http.StatusInternalServerError)
		return
	}

	var user *entity.User
	if existingUser == nil {
		// Create new user
		avatar := userInfo.Picture
		user = &entity.User{
			ID:       userInfo.ID,
			Email:    userInfo.Email,
			Name:     userInfo.Name,
			Avatar:   &avatar,
			Role:     entity.RoleReader,
			Status:   entity.StatusActive,
		}

		// Register user with OAuth2
		if err := h.userUsecase.RegisterWithOAuth2(r.Context(), user, provider); err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
	} else {
		user = existingUser
		// Update user information if needed
		if user.Avatar == nil || *user.Avatar != userInfo.Picture || user.Name != userInfo.Name {
			avatar := userInfo.Picture
			user.Avatar = &avatar
			user.Name = userInfo.Name
			if err := h.userUsecase.UpdateUser(r.Context(), user); err != nil {
				http.Error(w, "Failed to update user", http.StatusInternalServerError)
				return
			}
		}
	}

	// Store OAuth2 token in database (encrypted)
	entityProvider := convertProviderToEntity(provider)
	if err := h.tokenService.StoreOAuthToken(r.Context(), user.ID, entityProvider, token); err != nil {
		// Log error but don't fail the request
		// Token storage is important but shouldn't break authentication flow
		log.Printf("Warning: failed to store OAuth2 token for user %s: %v", user.ID, err)
	} else {
		log.Printf("Successfully stored OAuth2 token for user %s with provider %s", user.ID, entityProvider)
	}

	// Generate JWT token or session for the user
	// For now, we'll return the user info with the OAuth2 access token
	// TODO: Replace with JWT token generation for better security
	response := OAuth2TokenResponse{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresIn:   int64(token.Expiry.Sub(token.Expiry).Seconds()),
		User:        user,
	}

	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, constant.ERR_ENCODING_RESP, http.StatusInternalServerError)
		return
	}
}

// GetProviders returns available OAuth2 providers
func (h *OAuth2Handler) GetProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, constant.ERR_METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	providers := []string{}
	if _, exists := h.oauth2Service.GetProvider(oauth2.ProviderGoogle); exists {
		providers = append(providers, "google")
	}
	if _, exists := h.oauth2Service.GetProvider(oauth2.ProviderGitHub); exists {
		providers = append(providers, "github")
	}
	if _, exists := h.oauth2Service.GetProvider(oauth2.ProviderFacebook); exists {
		providers = append(providers, "facebook")
	}

	response := map[string]interface{}{
		"providers": providers,
	}

	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, constant.ERR_ENCODING_RESP, http.StatusInternalServerError)
		return
	}
}

// HandleOAuth2Redirect handles the redirect from OAuth2 provider
func (h *OAuth2Handler) HandleOAuth2Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, constant.ERR_METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	// Extract provider from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid OAuth2 redirect path", http.StatusBadRequest)
		return
	}

	providerStr := pathParts[2] // /oauth2/{provider}/redirect
	provider := oauth2.Provider(providerStr)

	// Get authorization code and state from query parameters
	code := r.URL.Query().Get("code")
	_ = r.URL.Query().Get("state") // Ignore state for now

	if code == "" {
		http.Error(w, "Authorization code not provided", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	token, err := h.oauth2Service.ExchangeCode(r.Context(), provider, code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	// Get user information
	userInfo, err := h.oauth2Service.GetUserInfo(r.Context(), provider, token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Check if user exists
	existingUser, err := h.userUsecase.GetUserByEmail(r.Context(), userInfo.Email)
	if err != nil {
		http.Error(w, "Failed to check existing user", http.StatusInternalServerError)
		return
	}

	var user *entity.User
	if existingUser == nil {
		// Create new user
		avatar := userInfo.Picture
		user = &entity.User{
			ID:       userInfo.ID,
			Email:    userInfo.Email,
			Name:     userInfo.Name,
			Avatar:   &avatar,
			Role:     entity.RoleReader,
			Status:   entity.StatusActive,
		}

		if err := h.userUsecase.RegisterWithOAuth2(r.Context(), user, provider); err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
	} else {
		user = existingUser
		// Update user information if needed
		if user.Avatar == nil || *user.Avatar != userInfo.Picture || user.Name != userInfo.Name {
			avatar := userInfo.Picture
			user.Avatar = &avatar
			user.Name = userInfo.Name
			if err := h.userUsecase.UpdateUser(r.Context(), user); err != nil {
				http.Error(w, "Failed to update user", http.StatusInternalServerError)
				return
			}
		}
	}

	// Store OAuth2 token in database (encrypted)
	entityProvider := convertProviderToEntity(provider)
	if err := h.tokenService.StoreOAuthToken(r.Context(), user.ID, entityProvider, token); err != nil {
		// Log error but don't fail the redirect
		log.Printf("Warning: failed to store OAuth2 token for user %s: %v", user.ID, err)
	} else {
		log.Printf("Successfully stored OAuth2 token for user %s with provider %s", user.ID, entityProvider)
	}

	// Redirect to frontend with user info and token
	redirectURL := fmt.Sprintf("/oauth2/success?user_id=%s&email=%s&name=%s&provider=%s&access_token=%s",
		url.QueryEscape(user.ID),
		url.QueryEscape(user.Email),
		url.QueryEscape(user.Name),
		url.QueryEscape(string(provider)),
		url.QueryEscape(token.AccessToken),
	)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
