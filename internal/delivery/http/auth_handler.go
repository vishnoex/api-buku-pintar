package http

import (
	"buku-pintar/internal/constant"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"buku-pintar/internal/usecase"
	"buku-pintar/pkg/supabase"
	"encoding/json"
	"errors"
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
		http.Error(w, constant.ERR_METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	var req AuthRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	roles, err := h.validateRegisterRequest(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := h.supabaseAuth.SignUp(r.Context(), req.Email, req.Password, roles); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(AuthRegisterResponse{
		Message:  "Registration success. Please check your email to verify your account",
		Email:    req.Email,
		Verified: false,
		Roles:    roles,
	})
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, constant.ERR_METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	var req AuthVerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.AccessToken) == "" || strings.TrimSpace(req.RefreshToken) == "" {
		http.Error(w, "access_token and refresh_token are required", http.StatusBadRequest)
		return
	}

	supabaseUser, err := h.supabaseAuth.GetUser(r.Context(), req.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if supabaseUser.EmailConfirmedAt == nil {
		http.Error(w, "Email not verified", http.StatusForbidden)
		return
	}

	roles, err := h.rolesFromNames(r, supabaseUser.Roles())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	primaryRole := roles[0]

	user, err := h.userUsecase.GetUserByID(r.Context(), supabaseUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(AuthVerifyEmailResponse{
		Verified: true,
		User:     user,
		Roles:    roles,
		Session: AuthSessionResponse{
			AccessToken:  req.AccessToken,
			RefreshToken: req.RefreshToken,
		},
	})
}

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
