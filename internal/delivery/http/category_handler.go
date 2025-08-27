package http

import (
	"buku-pintar/internal/constant"
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/usecase"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type CategoryHandler struct {
	usecase usecase.CategoryUsecase
}

func NewCategoryHandler(usecase usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{usecase: usecase}
}

func (h *CategoryHandler) ListCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit := 10
	offset := 0

	// Parse query parameters for pagination
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	categories, err := h.usecase.ListCategory(ctx, limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	total, err := h.usecase.CountCategory(ctx)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	response.WritePaginated(w, categories, total, limit, offset)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract ID from URL path (/categories/view/{id})
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, constant.ERR_CATEGORY_ID_REQUIRED)
		return
	}
	id := pathParts[len(pathParts)-1]

	if id == "" {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, constant.ERR_CATEGORY_ID_REQUIRED)
		return
	}

	category, err := h.usecase.GetCategoryByID(ctx, id)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	if category == nil {
		response.WriteError(w, http.StatusNotFound, constant.ERR_CODE_NOT_FOUND, "category not found")
		return
	}

	response.WriteSuccess(w, http.StatusOK, category, "category retrieved successfully")
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var category entity.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, "invalid request body")
		return
	}

	// Generate ID if not provided
	if category.ID == "" {
		category.ID = uuid.New().String()
	}

	err := h.usecase.CreateCategory(ctx, &category)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusCreated, map[string]string{"id": category.ID, "message": "category created successfully"}, "category created successfully")
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract ID from URL path (/categories/edit/{id})
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, constant.ERR_CATEGORY_ID_REQUIRED)
		return
	}
	id := pathParts[len(pathParts)-1]

	if id == "" {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, constant.ERR_CATEGORY_ID_REQUIRED)
		return
	}

	var category entity.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, "invalid request body")
		return
	}

	category.ID = id

	err := h.usecase.UpdateCategory(ctx, &category)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, map[string]string{"message": "category updated successfully"}, "category updated successfully")
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract ID from URL path (/categories/delete/{id})
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, constant.ERR_CATEGORY_ID_REQUIRED)
		return
	}
	id := pathParts[len(pathParts)-1]

	if id == "" {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, constant.ERR_CATEGORY_ID_REQUIRED)
		return
	}

	err := h.usecase.DeleteCategory(ctx, id)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, map[string]string{"message": "category deleted successfully"}, "category deleted successfully")
}

func (h *CategoryHandler) ListAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit := 10
	offset := 0

	// Parse query parameters for pagination
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	categories, err := h.usecase.ListAllCategories(ctx, limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	total, err := h.usecase.CountAllCategories(ctx)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	response.WritePaginated(w, categories, total, limit, offset)
}

func (h *CategoryHandler) ListCategoriesByParent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit := 10
	offset := 0

	// Parse query parameters for pagination
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Extract parent ID from URL path (/categories/parent/{parentID})
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, "parent category ID is required")
		return
	}
	parentID := pathParts[len(pathParts)-1]

	if parentID == "" {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, "parent category ID is required")
		return
	}

	categories, err := h.usecase.ListCategoriesByParent(ctx, parentID, limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	total, err := h.usecase.CountCategoriesByParent(ctx, parentID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	response.WritePaginated(w, categories, total, limit, offset)
}
