package http

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/usecase"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"buku-pintar/internal/constant"
)

type EbookHandler struct {
	ebookUsecase usecase.EbookUsecase
}

func NewEbookHandler(ebookUsecase usecase.EbookUsecase) *EbookHandler {
	return &EbookHandler{
		ebookUsecase: ebookUsecase,
	}
}

func (h *EbookHandler) ListEbooks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get ebooks from usecase
	ebooks, err := h.ebookUsecase.ListEbooks(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]any{
		"status": constant.STATUS_SUCCESS,
		"data":   ebooks,
	}

	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// List handles the HTTP GET request to retrieve a list of ebooks with pagination
func (h *EbookHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get ebooks from usecase
	ebooks, err := h.ebookUsecase.ListEbooks(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]any{
		"status": constant.STATUS_SUCCESS,
		"data":   ebooks,
	}

	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *EbookHandler) GetEbookByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, constant.INVALID_URL_PATH, http.StatusBadRequest)
		return
	}
	ebookID := pathParts[2]

	if ebookID == "" {
		http.Error(w, constant.EBOOK_ID_REQUIRED, http.StatusBadRequest)
		return
	}

	// Get ebook from usecase
	ebook, err := h.ebookUsecase.GetEbookByID(r.Context(), ebookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ebook == nil {
		http.Error(w, constant.EBOOK_NOT_FOUND, http.StatusNotFound)
		return
	}

	// Return response
	response := map[string]any{
		"status": constant.STATUS_SUCCESS,
		"data":   ebook,
	}

	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *EbookHandler) GetEbookBySlug(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, constant.INVALID_URL_PATH, http.StatusBadRequest)
		return
	}
	slug := pathParts[2]

	if slug == "" {
		http.Error(w, "slug is required", http.StatusBadRequest)
		return
	}

	// Get ebook from usecase
	ebook, err := h.ebookUsecase.GetEbookBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ebook == nil {
		http.Error(w, constant.EBOOK_NOT_FOUND, http.StatusNotFound)
		return
	}

	// Return response
	response := map[string]any{
		"status": constant.STATUS_SUCCESS,
		"data":   ebook,
	}

	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *EbookHandler) CreateEbook(w http.ResponseWriter, r *http.Request) {
	var ebook entity.Ebook

	if err := json.NewDecoder(r.Body).Decode(&ebook); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Create ebook using usecase
	err := h.ebookUsecase.CreateEbook(r.Context(), &ebook)
	if err != nil {
		// Check if it's a validation error
		if err.Error() == "title is required" ||
			err.Error() == "author_id is required" ||
			err.Error() == "category_id is required" ||
			err.Error() == "slug is required" ||
			err.Error() == "ebook with this slug already exists" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]any{
		"status":  constant.STATUS_SUCCESS,
		"message": "ebook created successfully",
		"data":    ebook,
	}

	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *EbookHandler) UpdateEbook(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, constant.INVALID_URL_PATH, http.StatusBadRequest)
		return
	}
	ebookID := pathParts[2]

	if ebookID == "" {
		http.Error(w, constant.EBOOK_ID_REQUIRED, http.StatusBadRequest)
		return
	}

	var ebook entity.Ebook
	if err := json.NewDecoder(r.Body).Decode(&ebook); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ebook.ID = ebookID // Ensure the ID is set from URL

	// Update ebook using usecase
	err := h.ebookUsecase.UpdateEbook(r.Context(), &ebook)
	if err != nil {
		// Check if it's a validation error
		if err.Error() == "id is required" ||
			err.Error() == "ebook not found" ||
			err.Error() == "ebook with this slug already exists" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]any{
		"status":  constant.STATUS_SUCCESS,
		"message": "ebook updated successfully",
		"data":    ebook,
	}

	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *EbookHandler) DeleteEbook(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	ebookID := pathParts[2]

	if ebookID == "" {
		http.Error(w, constant.EBOOK_ID_REQUIRED, http.StatusBadRequest)
		return
	}

	// Delete ebook using usecase
	err := h.ebookUsecase.DeleteEbook(r.Context(), ebookID)
	if err != nil {
		// Check if it's a validation error
		if err.Error() == "id is required" ||
			err.Error() == "ebook not found" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]any{
		"status":  constant.STATUS_SUCCESS,
		"message": "ebook deleted successfully",
	}

	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *EbookHandler) ListEbooksByCategory(w http.ResponseWriter, r *http.Request) {
	// Extract category ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, constant.INVALID_URL_PATH, http.StatusBadRequest)
		return
	}
	categoryID := pathParts[3]

	if categoryID == "" {
		http.Error(w, constant.CATEGORY_ID_REQUIRED, http.StatusBadRequest)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get ebooks by category from usecase
	ebooks, err := h.ebookUsecase.ListEbooksByCategory(r.Context(), categoryID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]any{
		"status": constant.STATUS_SUCCESS,
		"data":   ebooks,
	}

	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
