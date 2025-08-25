package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// Provider represents an OAuth2 provider
type Provider string

const (
	ProviderGoogle   Provider = "google"
	ProviderGitHub   Provider = "github"
	ProviderFacebook Provider = "facebook"
)

// UserInfo represents user information from OAuth2 provider
type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	Provider  Provider `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
}

// OAuth2Service handles OAuth2 authentication flows
type OAuth2Service struct {
	providers map[Provider]*oauth2.Config
}

// NewOAuth2Service creates a new OAuth2 service instance
func NewOAuth2Service() *OAuth2Service {
	return &OAuth2Service{
		providers: make(map[Provider]*oauth2.Config),
	}
}

// AddGoogleProvider adds Google OAuth2 provider
func (s *OAuth2Service) AddGoogleProvider(clientID, clientSecret, redirectURL string) {
	s.providers[ProviderGoogle] = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

// AddGitHubProvider adds GitHub OAuth2 provider
func (s *OAuth2Service) AddGitHubProvider(clientID, clientSecret, redirectURL string) {
	s.providers[ProviderGitHub] = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"user:email",
			"read:user",
		},
		Endpoint: github.Endpoint,
	}
}

// AddFacebookProvider adds Facebook OAuth2 provider
func (s *OAuth2Service) AddFacebookProvider(clientID, clientSecret, redirectURL string) {
	s.providers[ProviderFacebook] = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"email",
			"public_profile",
		},
		Endpoint: facebook.Endpoint,
	}
}

// GetAuthURL generates the authorization URL for the specified provider
func (s *OAuth2Service) GetAuthURL(provider Provider, state string) (string, error) {
	config, exists := s.providers[provider]
	if !exists {
		return "", fmt.Errorf("provider %s not configured", provider)
	}

	return config.AuthCodeURL(state), nil
}

// ExchangeCode exchanges authorization code for access token
func (s *OAuth2Service) ExchangeCode(ctx context.Context, provider Provider, code string) (*oauth2.Token, error) {
	config, exists := s.providers[provider]
	if !exists {
		return nil, fmt.Errorf("provider %s not configured", provider)
	}

	return config.Exchange(ctx, code)
}

// GetUserInfo retrieves user information from the OAuth2 provider
func (s *OAuth2Service) GetUserInfo(ctx context.Context, provider Provider, token *oauth2.Token) (*UserInfo, error) {
	client := s.providers[provider].Client(ctx, token)
	
	var userInfo *UserInfo
	var err error

	switch provider {
	case ProviderGoogle:
		userInfo, err = s.getGoogleUserInfo(client)
	case ProviderGitHub:
		userInfo, err = s.getGitHubUserInfo(client)
	case ProviderFacebook:
		userInfo, err = s.getFacebookUserInfo(client)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	if err != nil {
		return nil, err
	}

	userInfo.Provider = provider
	userInfo.CreatedAt = time.Now()
	return userInfo, nil
}

// getGoogleUserInfo retrieves user information from Google
func (s *OAuth2Service) getGoogleUserInfo(client *http.Client) (*UserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:      googleUser.ID,
		Email:   googleUser.Email,
		Name:    googleUser.Name,
		Picture: googleUser.Picture,
	}, nil
}

// getGitHubUserInfo retrieves user information from GitHub
func (s *OAuth2Service) getGitHubUserInfo(client *http.Client) (*UserInfo, error) {
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:      fmt.Sprintf("%d", githubUser.ID),
		Email:   githubUser.Email,
		Name:    githubUser.Name,
		Picture: githubUser.AvatarURL,
	}, nil
}

// getFacebookUserInfo retrieves user information from Facebook
func (s *OAuth2Service) getFacebookUserInfo(client *http.Client) (*UserInfo, error) {
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var facebookUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&facebookUser); err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:      facebookUser.ID,
		Email:   facebookUser.Email,
		Name:    facebookUser.Name,
		Picture: facebookUser.Picture.Data.URL,
	}, nil
}

// GetProvider returns the OAuth2 config for a specific provider
func (s *OAuth2Service) GetProvider(provider Provider) (*oauth2.Config, bool) {
	config, exists := s.providers[provider]
	return config, exists
}
