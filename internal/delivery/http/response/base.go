package response

import (
	"buku-pintar/internal/constant"
)

// Response represents the standard API response structure
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

// Meta contains metadata for paginated responses
type Meta struct {
	Total       int64 `json:"total"`
	Limit       int   `json:"limit"`
	Offset      int   `json:"offset"`
	CurrentPage int   `json:"current_page,omitempty"`
	TotalPages  int   `json:"total_pages,omitempty"`
}

// Error contains error details
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data any, message string) *Response {
	return &Response{
		Status:  constant.STATUS_SUCCESS,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message string) *Response {
	return &Response{
		Status: constant.STATUS_ERROR,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data any, total int64, limit, offset int) *Response {
	return &Response{
		Status: constant.STATUS_SUCCESS,
		Data:   data,
		Meta: &Meta{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	}
}
