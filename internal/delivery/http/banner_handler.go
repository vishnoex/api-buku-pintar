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

type BannerHandler struct {
	usecase usecase.BannerUsecase
}

func NewBannerHandler(banner usecase.BannerUsecase) *BannerHandler {
	return &BannerHandler{
		usecase: banner,
	}
}

func (h *BannerHandler) ListBanner(w http.ResponseWriter, r *http.Request) {
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

	banners, err := h.usecase.ListBanner(ctx, limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	total, err := h.usecase.CountBanner(ctx)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	response.WritePaginated(w, banners, total, limit, offset)
}

func (h *BannerHandler) GetBannerByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract ID from URL path (/banners/view/{id})
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, constant.ERR_BANNER_ID_REQUIRED)
		return
	}
	id := pathParts[len(pathParts)-1]

	if id == "" {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, constant.ERR_BANNER_ID_REQUIRED)
		return
	}

	banner, err := h.usecase.GetBannerByID(ctx, id)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	if banner == nil {
		response.WriteError(w, http.StatusNotFound, constant.ERR_CODE_NOT_FOUND, "banner not found")
		return
	}

	response.WriteSuccess(w, http.StatusOK, banner, "banner retrieved successfully")
}

func (h *BannerHandler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var banner entity.Banner
	if err := json.NewDecoder(r.Body).Decode(&banner); err != nil {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, "invalid request body")
		return
	}

	// Generate ID if not provided
	if banner.ID == "" {
		banner.ID = uuid.New().String()
	}

	err := h.usecase.CreateBanner(ctx, &banner)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusCreated, map[string]string{"id": banner.ID, "message": "banner created successfully"}, "banner created successfully")
}

func (h *BannerHandler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract ID from URL path (/banners/edit/{id})
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, constant.ERR_BANNER_ID_REQUIRED)
		return
	}
	id := pathParts[len(pathParts)-1]

	if id == "" {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, constant.ERR_BANNER_ID_REQUIRED)
		return
	}

	var banner entity.Banner
	if err := json.NewDecoder(r.Body).Decode(&banner); err != nil {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, "invalid request body")
		return
	}

	banner.ID = id

	err := h.usecase.UpdateBanner(ctx, &banner)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, map[string]string{"message": "banner updated successfully"}, "banner updated successfully")
}

func (h *BannerHandler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract ID from URL path (/banners/delete/{id})
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, constant.ERR_BANNER_ID_REQUIRED)
		return
	}
	id := pathParts[len(pathParts)-1]

	if id == "" {
		response.WriteError(w, http.StatusBadRequest, constant.ERR_CODE_BAD_REQUEST, constant.ERR_BANNER_ID_REQUIRED)
		return
	}

	err := h.usecase.DeleteBanner(ctx, id)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, map[string]string{"message": "banner deleted successfully"}, "banner deleted successfully")
}

func (h *BannerHandler) ListActiveBanner(w http.ResponseWriter, r *http.Request) {
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

	banners, err := h.usecase.ListActiveBanner(ctx, limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	total, err := h.usecase.CountActiveBanner(ctx)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, constant.ERR_CODE_SERVER_ERROR, err.Error())
		return
	}

	response.WritePaginated(w, banners, total, limit, offset)
}
