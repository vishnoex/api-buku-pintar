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

type SummaryHandler struct {
	summaryUsecase usecase.SummaryUsecase
}

func NewSummaryHandler(summaryUsecase usecase.SummaryUsecase) *SummaryHandler {
	return &SummaryHandler{
		summaryUsecase: summaryUsecase,
	}
}

// ListSummaries handles GET /summaries - List all summaries with pagination
func (h *SummaryHandler) ListSummaries(w http.ResponseWriter, r *http.Request) {
	limit, offset := helper.HandlePagination(r)

	// Get summaries from usecase
	summaries, err := h.summaryUsecase.ListSummaries(r.Context(), limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Get total count for pagination
	total, err := h.summaryUsecase.CountSummaries(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Return response
	response.WritePaginated(w, summaries, total, limit, offset)
}

// GetSummaryByID handles GET /summaries/{id} - Get summary by ID
func (h *SummaryHandler) GetSummaryByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		response.WriteError(w, http.StatusBadRequest, "invalid_url_path", constant.INVALID_URL_PATH)
		return
	}
	summaryID := pathParts[2]

	if summaryID == "" {
		response.WriteError(w, http.StatusBadRequest, "summary_id_required", constant.SUMMARY_ID_REQUIRED)
		return
	}

	// Get summary from usecase
	summary, err := h.summaryUsecase.GetSummaryByID(r.Context(), summaryID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	if summary == nil {
		response.WriteError(w, http.StatusNotFound, "summary_not_found", "Summary not found")
		return
	}

	// Return response
	response.WriteSuccess(w, http.StatusOK, summary, "")
}

// GetSummariesByEbookID handles GET /summaries/ebook/{ebookID} - Get summaries by ebook ID
func (h *SummaryHandler) GetSummariesByEbookID(w http.ResponseWriter, r *http.Request) {
	// Extract ebook ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.WriteError(w, http.StatusBadRequest, "invalid_url_path", constant.INVALID_URL_PATH)
		return
	}
	ebookID := pathParts[3]

	if ebookID == "" {
		response.WriteError(w, http.StatusBadRequest, "ebook_id_required", constant.EBOOK_ID_REQUIRED)
		return
	}

	limit, offset := helper.HandlePagination(r)

	// Get summaries from usecase
	summaries, err := h.summaryUsecase.GetSummariesByEbookID(r.Context(), ebookID, limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Get total count for pagination
	total, err := h.summaryUsecase.CountSummariesByEbookID(r.Context(), ebookID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Return response
	response.WritePaginated(w, summaries, total, limit, offset)
}

// CreateSummary handles POST /summaries - Create new summary (protected)
func (h *SummaryHandler) CreateSummary(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var summary entity.EbookSummary
	err := json.NewDecoder(r.Body).Decode(&summary)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid_request_body", err.Error())
		return
	}

	// Validate summary data
	if summary.EbookID == "" {
		response.WriteError(w, http.StatusBadRequest, "ebook_id_required", constant.EBOOK_ID_REQUIRED)
		return
	}

	if summary.Description == "" {
		response.WriteError(w, http.StatusBadRequest, "description_required", constant.DESCRIPTION_REQUIRED)
		return
	}

	// Create summary
	err = h.summaryUsecase.CreateSummary(r.Context(), &summary)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Return response
	response.WriteSuccess(w, http.StatusCreated, &summary, "Summary created successfully")
}

// UpdateSummary handles PUT /summaries/{id} - Update summary (protected)
func (h *SummaryHandler) UpdateSummary(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		response.WriteError(w, http.StatusBadRequest, "invalid_url_path", constant.INVALID_URL_PATH)
		return
	}
	summaryID := pathParts[2]

	if summaryID == "" {
		response.WriteError(w, http.StatusBadRequest, "summary_id_required", constant.SUMMARY_ID_REQUIRED)
		return
	}

	// Parse request body
	var summary entity.EbookSummary
	if err := json.NewDecoder(r.Body).Decode(&summary); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid_request_body", err.Error())
		return
	}

	// Set the ID from the URL
	summary.ID = summaryID

	// Update summary
	err := h.summaryUsecase.UpdateSummary(r.Context(), &summary)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Return response
	response.WriteSuccess(w, http.StatusOK, &summary, "Summary updated successfully")
}

// DeleteSummary handles DELETE /summaries/{id} - Delete summary (protected)
func (h *SummaryHandler) DeleteSummary(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		response.WriteError(w, http.StatusBadRequest, "invalid_url_path", constant.INVALID_URL_PATH)
		return
	}
	summaryID := pathParts[2]

	if summaryID == "" {
		response.WriteError(w, http.StatusBadRequest, "summary_id_required", constant.SUMMARY_ID_REQUIRED)
		return
	}

	// Delete summary using usecase
	err := h.summaryUsecase.DeleteSummary(r.Context(), summaryID)
	if err != nil {
		// Check if it's a validation error
		if err.Error() == "id is required" ||
			err.Error() == "summary not found" {
			response.WriteError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	// Return response
	response.WriteSuccess(w, http.StatusOK, nil, "Summary deleted successfully")
}
