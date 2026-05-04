package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"buku-pintar/pkg/config"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

const defaultAudience = "authenticated"

// Claims contains the Supabase access token fields used by the API.
type Claims struct {
	Email       string                 `json:"email"`
	Phone       string                 `json:"phone"`
	Role        string                 `json:"role"`
	AppMetadata map[string]interface{} `json:"app_metadata"`
	UserMeta    map[string]interface{} `json:"user_metadata"`
	jwt.RegisteredClaims
}

// Authenticator verifies Supabase Auth JWT access tokens.
type Authenticator struct {
	projectURL      string
	anonKey         string
	jwtSecret       []byte
	jwksURL         string
	issuer          string
	audience        string
	emailRedirectTo string
	client          *http.Client
	mu              sync.RWMutex
	jwks            *keyfunc.JWKS
}

// SignUpUser contains the Supabase auth user fields returned during signup.
type SignUpUser struct {
	ID               string     `json:"id"`
	Email            string     `json:"email"`
	EmailConfirmedAt *time.Time `json:"email_confirmed_at"`
}

// SignUpResponse contains the subset of Supabase signup response used by the API.
type SignUpResponse struct {
	User    SignUpUser  `json:"user"`
	Session interface{} `json:"session"`
}

// User contains the Supabase auth user fields used after email verification.
type User struct {
	ID               string                 `json:"id"`
	Email            string                 `json:"email"`
	EmailConfirmedAt *time.Time             `json:"email_confirmed_at"`
	UserMetadata     map[string]interface{} `json:"user_metadata"`
	AppMetadata      map[string]interface{} `json:"app_metadata"`
}

// Roles returns normalized role names from user metadata.
func (u *User) Roles() []string {
	return metadataRoles(u.UserMetadata)
}

// NewAuthenticator creates a Supabase JWT authenticator from app config.
func NewAuthenticator(cfg config.SupabaseConfig) (*Authenticator, error) {
	projectURL := strings.TrimRight(cfg.ProjectURL, "/")
	issuer := strings.TrimRight(cfg.Issuer, "/")
	if issuer == "" {
		issuer = projectURL + "/auth/v1"
	}

	jwksURL := cfg.JWKSURL
	if jwksURL == "" && projectURL != "" {
		jwksURL = projectURL + "/auth/v1/.well-known/jwks.json"
	}

	audience := cfg.Audience
	if audience == "" {
		audience = defaultAudience
	}

	if jwksURL == "" && cfg.JWTSecret == "" {
		return nil, errors.New("supabase project_url/jwks_url or jwt_secret is required")
	}

	return &Authenticator{
		projectURL:      projectURL,
		anonKey:         cfg.AnonKey,
		jwtSecret:       []byte(cfg.JWTSecret),
		jwksURL:         jwksURL,
		issuer:          issuer,
		audience:        audience,
		emailRedirectTo: cfg.EmailRedirectTo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// SignUp creates a Supabase Auth user and stores requested roles in user metadata.
func (a *Authenticator) SignUp(ctx context.Context, email, password string, roles []string) (*SignUpResponse, error) {
	if a.projectURL == "" {
		return nil, errors.New("supabase project_url is required")
	}
	if a.anonKey == "" {
		return nil, errors.New("supabase anon_key is required")
	}

	body := map[string]interface{}{
		"email":    email,
		"password": password,
		"data": map[string]interface{}{
			"roles": roles,
		},
	}
	if a.emailRedirectTo != "" {
		body["gotrue_meta_security"] = map[string]string{
			"captcha_token": "",
		}
		body["options"] = map[string]string{
			"email_redirect_to": a.emailRedirectTo,
		}
	}

	var response SignUpResponse
	if err := a.doJSON(ctx, http.MethodPost, a.projectURL+"/auth/v1/signup", a.anonKey, a.anonKey, body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetUser retrieves the Supabase Auth user for an access token.
func (a *Authenticator) GetUser(ctx context.Context, accessToken string) (*User, error) {
	if a.projectURL == "" {
		return nil, errors.New("supabase project_url is required")
	}
	if accessToken == "" {
		return nil, errors.New("access token is required")
	}

	var user User
	if err := a.doJSON(ctx, http.MethodGet, a.projectURL+"/auth/v1/user", accessToken, a.anonKey, nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// VerifyToken validates a Supabase JWT and returns its claims.
func (a *Authenticator) VerifyToken(ctx context.Context, rawToken string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
		switch token.Method.Alg() {
		case jwt.SigningMethodHS256.Alg():
			if len(a.jwtSecret) == 0 {
				return nil, errors.New("supabase jwt_secret is required for HS256 tokens")
			}
			return a.jwtSecret, nil
		case jwt.SigningMethodRS256.Alg(), jwt.SigningMethodES256.Alg():
			jwks, err := a.getJWKS(ctx)
			if err != nil {
				return nil, err
			}
			return jwks.Keyfunc(token)
		default:
			return nil, fmt.Errorf("unsupported signing method: %s", token.Method.Alg())
		}
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid supabase token")
	}
	if claims.Subject == "" {
		return nil, errors.New("supabase token subject is required")
	}
	if a.issuer != "" && !claims.VerifyIssuer(a.issuer, true) {
		return nil, errors.New("invalid supabase token issuer")
	}
	if a.audience != "" && !claims.VerifyAudience(a.audience, true) {
		return nil, errors.New("invalid supabase token audience")
	}

	return claims, nil
}

func (a *Authenticator) getJWKS(ctx context.Context) (*keyfunc.JWKS, error) {
	a.mu.RLock()
	if a.jwks != nil {
		jwks := a.jwks
		a.mu.RUnlock()
		return jwks, nil
	}
	a.mu.RUnlock()

	a.mu.Lock()
	defer a.mu.Unlock()
	if a.jwks != nil {
		return a.jwks, nil
	}
	if a.jwksURL == "" {
		return nil, errors.New("supabase jwks_url is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.jwksURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("failed to fetch supabase jwks: status %d", resp.StatusCode)
	}

	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	jwks, err := keyfunc.NewJSON(raw)
	if err != nil {
		return nil, err
	}
	a.jwks = jwks
	return jwks, nil
}

func (a *Authenticator) doJSON(ctx context.Context, method, url, bearerToken, apiKey string, body interface{}, out interface{}) error {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	if apiKey != "" {
		req.Header.Set("apikey", apiKey)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		message, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase auth request failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(message)))
	}

	if out == nil || resp.StatusCode == http.StatusNoContent {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// DisplayName returns the most useful name available in Supabase user metadata.
func (c *Claims) DisplayName() string {
	for _, key := range []string{"name", "full_name", "user_name"} {
		if value, ok := c.UserMeta[key].(string); ok && value != "" {
			return value
		}
	}
	if c.Email != "" {
		return c.Email
	}
	return c.Subject
}

// AvatarURL returns the avatar URL from Supabase user metadata when present.
func (c *Claims) AvatarURL() *string {
	for _, key := range []string{"avatar_url", "picture"} {
		if value, ok := c.UserMeta[key].(string); ok && value != "" {
			return &value
		}
	}
	return nil
}

// Roles returns normalized role names from JWT user metadata.
func (c *Claims) Roles() []string {
	return metadataRoles(c.UserMeta)
}

func metadataRoles(metadata map[string]interface{}) []string {
	if metadata == nil {
		return nil
	}
	value, ok := metadata["roles"]
	if !ok {
		return nil
	}
	switch typed := value.(type) {
	case []string:
		return typed
	case []interface{}:
		roles := make([]string, 0, len(typed))
		for _, item := range typed {
			if role, ok := item.(string); ok {
				roles = append(roles, role)
			}
		}
		return roles
	case string:
		return []string{typed}
	default:
		return nil
	}
}
