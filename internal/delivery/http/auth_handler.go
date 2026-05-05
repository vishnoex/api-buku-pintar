package http

import (
	"buku-pintar/internal/constant"
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"buku-pintar/internal/usecase"
	"buku-pintar/pkg/supabase"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strings"
)

// AuthHandler handles Supabase-backed signup and verification flows.
type AuthHandler struct {
	supabaseAuth *supabase.Authenticator
	userUsecase  usecase.UserUsecase
	roleService  service.RoleService
}

func NewAuthHandler(
	supabaseAuth *supabase.Authenticator,
	userUsecase usecase.UserUsecase,
	roleService service.RoleService,
) *AuthHandler {
	return &AuthHandler{
		supabaseAuth: supabaseAuth,
		userUsecase:  userUsecase,
		roleService:  roleService,
	}
}

type AuthRegisterRequest struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

type AuthRegisterResponse struct {
	Message  string   `json:"message"`
	Email    string   `json:"email"`
	Verified bool     `json:"verified"`
	Roles    []string `json:"roles"`
}

type AuthVerifyEmailRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthVerifyEmailResponse struct {
	Verified bool                `json:"verified"`
	User     *entity.User        `json:"user"`
	Roles    []*entity.Role      `json:"roles"`
	Session  AuthSessionResponse `json:"session"`
}

type AuthSessionResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed", constant.ERR_METHOD_NOT_ALLOWED)
		return
	}

	var req AuthRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, "invalid request body")
		return
	}

	roles, err := h.validateRegisterRequest(r, &req)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, err.Error())
		return
	}

	if _, err := h.supabaseAuth.SignUp(r.Context(), req.Email, req.Password, roles); err != nil {
		response.WriteError(w, http.StatusBadGateway, "bad_gateway", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusCreated, AuthRegisterResponse{
		Message:  "Registration success. Please check your email to verify your account",
		Email:    req.Email,
		Verified: false,
		Roles:    roles,
	}, "registration successful")
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if req, ok := verifyEmailRequestFromQuery(r); ok {
			h.verifyEmailWithTokens(w, r, req)
			return
		}
		writeVerifyEmailCallbackPage(w, r)
		return
	case http.MethodPost:
		req, ok := decodeVerifyEmailRequest(w, r)
		if !ok {
			return
		}
		h.verifyEmailWithTokens(w, r, req)
		return
	default:
		response.WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed", constant.ERR_METHOD_NOT_ALLOWED)
		return
	}
}

func (h *AuthHandler) verifyEmailWithTokens(w http.ResponseWriter, r *http.Request, req AuthVerifyEmailRequest) {
	supabaseUser, err := h.supabaseAuth.GetUser(r.Context(), req.AccessToken)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized", err.Error())
		return
	}
	if supabaseUser.EmailConfirmedAt == nil {
		response.WriteError(w, http.StatusForbidden, "forbidden", "email not verified")
		return
	}

	roles, err := h.rolesFromNames(r, supabaseUser.Roles())
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, err.Error())
		return
	}
	primaryRole := roles[0]

	user, err := h.userUsecase.GetUserByID(r.Context(), supabaseUser.ID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}
	if user == nil {
		user = &entity.User{
			ID:     supabaseUser.ID,
			Name:   displayNameFromMetadata(supabaseUser.UserMetadata, supabaseUser.Email, supabaseUser.ID),
			Email:  supabaseUser.Email,
			RoleID: &primaryRole.ID,
			Role:   legacyRole(primaryRole.Name),
			Status: entity.StatusActive,
			Avatar: avatarFromMetadata(supabaseUser.UserMetadata),
		}
		if err := h.userUsecase.CreateUser(r.Context(), user); err != nil {
			response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
			return
		}
	} else {
		user.Name = displayNameFromMetadata(supabaseUser.UserMetadata, supabaseUser.Email, supabaseUser.ID)
		user.Email = supabaseUser.Email
		user.RoleID = &primaryRole.ID
		user.Role = legacyRole(primaryRole.Name)
		user.Status = entity.StatusActive
		if avatar := avatarFromMetadata(supabaseUser.UserMetadata); avatar != nil {
			user.Avatar = avatar
		}
		if err := h.userUsecase.UpdateUser(r.Context(), user); err != nil {
			response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, AuthVerifyEmailResponse{
		Verified: true,
		User:     user,
		Roles:    roles,
		Session:  AuthSessionResponse(req),
	}, "email verified successfully")
}

func decodeVerifyEmailRequest(w http.ResponseWriter, r *http.Request) (AuthVerifyEmailRequest, bool) {
	var req AuthVerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, "invalid request body")
		return AuthVerifyEmailRequest{}, false
	}
	if strings.TrimSpace(req.AccessToken) == "" || strings.TrimSpace(req.RefreshToken) == "" {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, "access_token and refresh_token are required")
		return AuthVerifyEmailRequest{}, false
	}
	return req, true
}

func verifyEmailRequestFromQuery(r *http.Request) (AuthVerifyEmailRequest, bool) {
	query := r.URL.Query()
	req := AuthVerifyEmailRequest{
		AccessToken:  strings.TrimSpace(query.Get("access_token")),
		RefreshToken: strings.TrimSpace(query.Get("refresh_token")),
	}
	return req, req.AccessToken != "" && req.RefreshToken != ""
}

func writeVerifyEmailCallbackPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constant.CONTENT_TYPE, "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := verifyEmailCallbackTemplate.Execute(w, map[string]string{
		"Endpoint": r.URL.Path,
	}); err != nil {
		http.Error(w, constant.ERR_ENCODING_RESP, http.StatusInternalServerError)
	}
}

var verifyEmailCallbackTemplate = template.Must(template.New("verify-email-callback").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Email Verification</title>
  <style>
    body { font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; margin: 0; min-height: 100vh; display: grid; place-items: center; background: #f8fafc; color: #0f172a; }
    main { width: min(92vw, 34rem); padding: 2rem; border-radius: 1rem; background: #fff; box-shadow: 0 20px 45px rgba(15, 23, 42, .12); }
    h1 { margin: 0 0 .75rem; font-size: 1.5rem; }
    p { margin: 0; line-height: 1.6; color: #475569; }
    pre { white-space: pre-wrap; overflow-wrap: anywhere; margin-top: 1rem; padding: 1rem; border-radius: .75rem; background: #f1f5f9; color: #334155; }
  </style>
</head>
<body>
  <main>
    <h1 id="title">Verifying email...</h1>
    <p id="message">Please wait while your email verification is completed.</p>
    <pre id="details" hidden></pre>
  </main>
  <script>
    (async function () {
      const title = document.getElementById('title');
      const message = document.getElementById('message');
      const details = document.getElementById('details');
      const params = new URLSearchParams(window.location.hash.slice(1) || window.location.search.slice(1));
      const accessToken = params.get('access_token');
      const refreshToken = params.get('refresh_token');

      if (!accessToken || !refreshToken) {
        title.textContent = 'Verification link is invalid';
        message.textContent = 'The verification callback did not include the required tokens.';
        return;
      }

      try {
        const response = await fetch('{{ .Endpoint }}', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ access_token: accessToken, refresh_token: refreshToken }),
        });
        const payload = await response.json();
        if (!response.ok) {
          throw payload;
        }
        title.textContent = 'Email verified successfully';
        message.textContent = 'Your account is active. You can close this page and continue signing in.';
      } catch (error) {
        title.textContent = 'Email verification failed';
        message.textContent = 'We could not complete the verification request.';
        details.hidden = false;
        details.textContent = JSON.stringify(error, null, 2);
      } finally {
        if (window.history.replaceState) {
          window.history.replaceState(null, document.title, window.location.pathname + window.location.search);
        }
      }
    })();
  </script>
</body>
</html>`))

func (h *AuthHandler) validateRegisterRequest(r *http.Request, req *AuthRegisterRequest) ([]string, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" {
		return nil, errors.New("email is required")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	roles := normalizeRoles(req.Roles)
	if len(roles) == 0 {
		return nil, errors.New("at least one role is required")
	}
	if len(roles) > 1 {
		return nil, errors.New("only one role is supported per user")
	}
	if _, err := h.roleService.GetRoleByName(r.Context(), roles[0]); err != nil {
		return nil, errors.New("role not found")
	}
	return roles, nil
}

func (h *AuthHandler) rolesFromNames(r *http.Request, roleNames []string) ([]*entity.Role, error) {
	roles := normalizeRoles(roleNames)
	if len(roles) == 0 {
		roles = []string{string(entity.RoleTypeReader)}
	}
	if len(roles) > 1 {
		return nil, errors.New("only one role is supported per user")
	}

	role, err := h.roleService.GetRoleByName(r.Context(), roles[0])
	if err != nil {
		return nil, errors.New("role not found")
	}
	return []*entity.Role{role}, nil
}

func normalizeRoles(roles []string) []string {
	seen := make(map[string]struct{}, len(roles))
	normalized := make([]string, 0, len(roles))
	for _, role := range roles {
		role = strings.TrimSpace(strings.ToLower(role))
		if role == "" {
			continue
		}
		if _, exists := seen[role]; exists {
			continue
		}
		seen[role] = struct{}{}
		normalized = append(normalized, role)
	}
	return normalized
}

func displayNameFromMetadata(metadata map[string]interface{}, email, id string) string {
	for _, key := range []string{"name", "full_name", "user_name"} {
		if value, ok := metadata[key].(string); ok && value != "" {
			return value
		}
	}
	if email != "" {
		return email
	}
	return id
}

func avatarFromMetadata(metadata map[string]interface{}) *string {
	for _, key := range []string{"avatar_url", "picture"} {
		if value, ok := metadata[key].(string); ok && value != "" {
			return &value
		}
	}
	return nil
}

func legacyRole(roleName string) entity.UserRole {
	switch entity.RoleType(strings.ToLower(roleName)) {
	case entity.RoleTypeAdmin:
		return entity.RoleAdmin
	case entity.RoleTypeEditor:
		return entity.RoleEditor
	default:
		return entity.RoleReader
	}
}
