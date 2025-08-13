package helper

import (
	"net/http"
	"strconv"
)

func HandlePagination(r *http.Request) (int, int) {
	limit := 10
	offset := 0

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

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

	// Maximum 100 items each page
	if limit > 100 {
		limit = 100
	}

	return limit, offset
}
