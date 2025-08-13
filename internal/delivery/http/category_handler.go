package http

import (
	"buku-pintar/internal/constant"
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/usecase"
	"net/http"
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
