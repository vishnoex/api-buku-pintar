package http

import (
	"buku-pintar/internal/delivery/http/response"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFoundHandler(t *testing.T) {
	// Create a test request to a non-existent route
	req, err := http.NewRequest("GET", "/non-existent-route", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the NotFoundHandler directly
	NotFoundHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	// Check the response body
	var responseBody response.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
		t.Fatal("Failed to unmarshal response body:", err)
	}

	// Check the response structure
	if responseBody.Status != "error" {
		t.Errorf("Expected status 'error', got '%s'", responseBody.Status)
	}

	if responseBody.Error == nil {
		t.Error("Expected error field to be present")
		return
	}

	if responseBody.Error.Code != "page_not_found" {
		t.Errorf("Expected error code 'page_not_found', got '%s'", responseBody.Error.Code)
	}

	if responseBody.Error.Message != "The requested page was not found" {
		t.Errorf("Expected error message 'The requested page was not found', got '%s'", responseBody.Error.Message)
	}
}

func TestNotFoundHandlerWithDifferentMethods(t *testing.T) {
	testCases := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "GET request",
			method: "GET",
			path:   "/non-existent-route",
		},
		{
			name:   "POST request",
			method: "POST",
			path:   "/non-existent-route",
		},
		{
			name:   "PUT request",
			method: "PUT",
			path:   "/non-existent-route",
		},
		{
			name:   "DELETE request",
			method: "DELETE",
			path:   "/non-existent-route",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			NotFoundHandler(rr, req)

			// Check the status code
			if status := rr.Code; status != http.StatusNotFound {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
			}

			// Check the response body
			var responseBody response.Response
			if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
				t.Fatal("Failed to unmarshal response body:", err)
			}

			// Check the response structure
			if responseBody.Status != "error" {
				t.Errorf("Expected status 'error', got '%s'", responseBody.Status)
			}

			if responseBody.Error == nil {
				t.Error("Expected error field to be present")
				return
			}

			if responseBody.Error.Code != "page_not_found" {
				t.Errorf("Expected error code 'page_not_found', got '%s'", responseBody.Error.Code)
			}

			if responseBody.Error.Message != "The requested page was not found" {
				t.Errorf("Expected error message 'The requested page was not found', got '%s'", responseBody.Error.Message)
			}
		})
	}
}
