package response

import (
	"buku-pintar/internal/constant"
	"encoding/json"
	"net/http"
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

// WriteError writes an error response to the ResponseWriter
func WriteError(w http.ResponseWriter, statusCode int, errorCode, errorMessage string) {
	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.Header().Set(constant.ACCESS_CONTROL_ALLOW_ORIGIN, "*")
	w.Header().Set(constant.ACCESS_CONTROL_ALLOW_HEADER, "Content-Type")
	w.WriteHeader(statusCode)
	resp := NewErrorResponse(errorCode, errorMessage)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// WriteSuccess writes a success response to the ResponseWriter
func WriteSuccess(w http.ResponseWriter, statusCode int, data any, message string) {
	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.Header().Set(constant.ACCESS_CONTROL_ALLOW_ORIGIN, "*")
	w.Header().Set(constant.ACCESS_CONTROL_ALLOW_HEADER, "Content-Type")
	w.WriteHeader(statusCode)
	resp := NewSuccessResponse(data, message)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// WritePaginated writes a paginated response to the ResponseWriter
func WritePaginated(w http.ResponseWriter, data any, total int64, limit, offset int) {
	w.Header().Set(constant.CONTENT_TYPE, constant.APPLICATION_JSON)
	w.Header().Set(constant.ACCESS_CONTROL_ALLOW_ORIGIN, "*")
	w.Header().Set(constant.ACCESS_CONTROL_ALLOW_METHOD, "GET")
	w.Header().Set(constant.ACCESS_CONTROL_ALLOW_HEADER, "Content-Type")
	w.WriteHeader(http.StatusOK)
	resp := NewPaginatedResponse(data, total, limit, offset)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
