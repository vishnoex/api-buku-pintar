package http

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/usecase"
	"encoding/json"
	"net/http"
	"strings"

	"buku-pintar/internal/constant"
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/helper"
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
	limit, offset := helper.HandlePagination(r)

	// Get ebooks from usecase
	ebooks, err := h.ebookUsecase.ListEbooks(r.Context(), limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Get total count for pagination
	total, err := h.ebookUsecase.CountEbooks(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Return response
	response.WritePaginated(w, ebooks, total, limit, offset)
}

// GetEbookByID handles the HTTP GET request to retrieve an ebook by its ID
func (h *EbookHandler) GetEbookByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		response.WriteError(w, http.StatusBadRequest, "invalid_url_path", constant.INVALID_URL_PATH)
		return
	}
	ebookID := pathParts[2]

	if ebookID == "" {
		response.WriteError(w, http.StatusBadRequest, "ebook_id_required", constant.EBOOK_ID_REQUIRED)
		return
	}

	// Get ebook from usecase
	ebook, err := h.ebookUsecase.GetEbookByID(r.Context(), ebookID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	if ebook == nil {
		response.WriteError(w, http.StatusNotFound, "ebook_not_found", constant.EBOOK_NOT_FOUND)
		return
	}

	// Return response
	response.WriteSuccess(w, http.StatusOK, ebook, "")
}

func (h *EbookHandler) GetEbookBySlug(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		response.WriteError(w, http.StatusBadRequest, "invalid_url_path", constant.INVALID_URL_PATH)
		return
	}
	slug := pathParts[3]

	if slug == "" {
		response.WriteError(w, http.StatusBadRequest, "slug_required", "slug is required")
		return
	}

	// Get ebook from usecase
	ebook, err := h.ebookUsecase.GetEbookBySlug(r.Context(), slug)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	if ebook == nil {
		response.WriteError(w, http.StatusNotFound, "ebook_not_found", constant.EBOOK_NOT_FOUND)
		return
	}

	// Return response
	response.WriteSuccess(w, http.StatusOK, ebook, "")
}

func (h *EbookHandler) CreateEbook(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var ebook entity.Ebook
	err := json.NewDecoder(r.Body).Decode(&ebook)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid_request_body", err.Error())
		return
	}

	// Validate ebook data
	if ebook.Title == "" {
		response.WriteError(w, http.StatusBadRequest, "ebook_title_required", "ebook title is required")
		return
	}

	// Create ebook
	err = h.ebookUsecase.CreateEbook(r.Context(), &ebook)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Return response
	response.WriteSuccess(w, http.StatusCreated, &ebook, "Ebook created successfully")
}

func (h *EbookHandler) UpdateEbook(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		response.WriteError(w, http.StatusBadRequest, "invalid_url_path", constant.INVALID_URL_PATH)
		return
	}
	ebookID := pathParts[2]

	if ebookID == "" {
		response.WriteError(w, http.StatusBadRequest, "ebook_id_required", constant.EBOOK_ID_REQUIRED)
		return
	}

	// Parse request body
	var ebook entity.Ebook
	if err := json.NewDecoder(r.Body).Decode(&ebook); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid_request_body", err.Error())
		return
	}

	// Set the ID from the URL
	ebook.ID = ebookID

	// Update ebook
	err := h.ebookUsecase.UpdateEbook(r.Context(), &ebook)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Return response
	response.WriteSuccess(w, http.StatusOK, &ebook, "Ebook updated successfully")
}

func (h *EbookHandler) DeleteEbook(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		response.WriteError(w, http.StatusBadRequest, "invalid_url_path", constant.INVALID_URL_PATH)
		return
	}
	ebookID := pathParts[2]

	if ebookID == "" {
		response.WriteError(w, http.StatusBadRequest, "ebook_id_required", constant.EBOOK_ID_REQUIRED)
		return
	}

	// Delete ebook using usecase
	err := h.ebookUsecase.DeleteEbook(r.Context(), ebookID)
	if err != nil {
		// Check if it's a validation error
		if err.Error() == "id is required" ||
			err.Error() == "ebook not found" {
			response.WriteError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Return response
	response.WriteSuccess(w, http.StatusOK, nil, "Ebook deleted successfully")
}

func (h *EbookHandler) ListEbooksByCategory(w http.ResponseWriter, r *http.Request) {
	// Extract category ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.WriteError(w, http.StatusBadRequest, "invalid_url_path", constant.INVALID_URL_PATH)
		return
	}
	categoryID := pathParts[3]

	if categoryID == "" {
		response.WriteError(w, http.StatusBadRequest, "category_id_required", constant.CATEGORY_ID_REQUIRED)
		return
	}
	limit, offset := helper.HandlePagination(r)

	// Get ebooks from usecase
	ebooks, err := h.ebookUsecase.ListEbooksByCategory(r.Context(), categoryID, limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Get total count for pagination
	total, err := h.ebookUsecase.CountEbooksByCategory(r.Context(), categoryID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Return response
	response.WritePaginated(w, ebooks, total, limit, offset)
}
