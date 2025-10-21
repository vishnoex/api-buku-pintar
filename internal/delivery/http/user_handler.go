package http

import (
	"buku-pintar/internal/constant"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/usecase"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

type RegisterRequest struct {
	User    entity.User `json:"user"`
	IDToken string      `json:"id_token"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, constant.ERR_METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	// Only OAuth2 registration is supported
	http.Error(w, "Direct registration not supported. Please use OAuth2 authentication.", http.StatusBadRequest)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, constant.ERR_METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	user, err := h.userUsecase.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, constant.ERR_ENCODING_RESP, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, constant.ERR_METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userUsecase.UpdateUser(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, constant.ERR_ENCODING_RESP, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, constant.ERR_METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if err := h.userUsecase.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
} 