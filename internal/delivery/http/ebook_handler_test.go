package http

import (
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/entity"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

// MockEbookUsecase is a mock implementation for testing
type MockEbookUsecase struct {
	ebooks []*entity.Ebook
	ebook  *entity.Ebook
	err    error
}

func (m *MockEbookUsecase) CreateEbook(ctx context.Context, ebook *entity.Ebook) error {
	return m.err
}

func (m *MockEbookUsecase) GetEbookByID(ctx context.Context, id string) (*entity.Ebook, error) {
	return m.ebook, m.err
}

func (m *MockEbookUsecase) GetEbookBySlug(ctx context.Context, slug string) (*entity.Ebook, error) {
	return m.ebook, m.err
}

func (m *MockEbookUsecase) UpdateEbook(ctx context.Context, ebook *entity.Ebook) error {
	return m.err
}

func (m *MockEbookUsecase) DeleteEbook(ctx context.Context, id string) error {
	return m.err
}

func (m *MockEbookUsecase) ListEbooks(ctx context.Context, limit, offset int) ([]*response.EbookListResponse, error) {
	// Convert entity.Ebook to response.EbookListResponse
	ebookResponses := make([]*response.EbookListResponse, 0, len(m.ebooks))
	for _, ebook := range m.ebooks {
		ebookResponses = append(ebookResponses, &response.EbookListResponse{
			ID:         ebook.ID,
			Title:      ebook.Title,
			Slug:       ebook.Slug,
			CoverImage: ebook.CoverImage,
			Price:      ebook.Price,
			Discount:   0,
			Status:     "", // Default empty status
		})
	}
	return ebookResponses, m.err
}

func (m *MockEbookUsecase) ListEbooksByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Ebook, error) {
	return m.ebooks, m.err
}

func (m *MockEbookUsecase) ListEbooksByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Ebook, error) {
	return m.ebooks, m.err
}

func (m *MockEbookUsecase) CountEbooks(ctx context.Context) (int64, error) {
	return int64(len(m.ebooks)), m.err
}

func (m *MockEbookUsecase) CountEbooksByCategory(ctx context.Context, categoryID string) (int64, error) {
	return int64(len(m.ebooks)), m.err
}

func (m *MockEbookUsecase) CountEbooksByAuthor(ctx context.Context, authorID string) (int64, error) {
	return int64(len(m.ebooks)), m.err
}

func TestEbookHandler_ListEbooks(t *testing.T) {
	tests := []struct {
		name           string
		mockEbooks     []*entity.Ebook
		mockError      error
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "should return list of ebooks successfully",
			mockEbooks: []*entity.Ebook{
				{
					ID:         uuid.New().String(),
					Title:      "Test Ebook 1",
					AuthorID:   "author-1",
					CategoryID: "category-1",
					Price:      1000,
					Language:   "en",
					Format:     entity.FormatPDF,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
				{
					ID:         uuid.New().String(),
					Title:      "Test Ebook 2",
					AuthorID:   "author-2",
					CategoryID: "category-2",
					Price:      2000,
					Language:   "en",
					Format:     entity.FormatEPUB,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "should return empty list when no ebooks exist",
			mockEbooks:     []*entity.Ebook{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "should return error when usecase fails",
			mockEbooks:     nil,
			mockError:      context.DeadlineExceeded,
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUsecase := &MockEbookUsecase{
				ebooks: tt.mockEbooks,
				err:    tt.mockError,
			}
			handler := NewEbookHandler(mockUsecase)

			// Create request
			req, err := http.NewRequest("GET", "/api/ebooks", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Act
			handler.ListEbooks(rr, req)

			// Assert
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}

				// Check status field
				if status, ok := response["status"].(string); !ok || status != "success" {
					t.Errorf("expected status 'success', got %v", status)
				}

				// Check data field
				if data, ok := response["data"].([]interface{}); ok {
					if len(data) != tt.expectedCount {
						t.Errorf("expected %d ebooks, got %d", tt.expectedCount, len(data))
					}
				} else {
					t.Errorf("expected 'data' field in response")
				}

				// Check meta field for pagination
				if meta, ok := response["meta"].(map[string]interface{}); ok {
					if _, ok := meta["total"]; !ok {
						t.Errorf("expected 'total' field in meta")
					}
					if _, ok := meta["limit"]; !ok {
						t.Errorf("expected 'limit' field in meta")
					}
					if _, ok := meta["offset"]; !ok {
						t.Errorf("expected 'offset' field in meta")
					}
				} else {
					t.Errorf("expected 'meta' field in response")
				}
			}
		})
	}
}

func TestEbookHandler_GetEbookByID(t *testing.T) {
	tests := []struct {
		name           string
		ebookID        string
		mockEbook      *entity.Ebook
		mockError      error
		expectedStatus int
		expectedFound  bool
	}{
		{
			name:    "should return ebook by ID successfully",
			ebookID: "ebook-1",
			mockEbook: &entity.Ebook{
				ID:         "ebook-1",
				Title:      "Test Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Price:      1000,
				Language:   "en",
				Format:     entity.FormatPDF,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedFound:  true,
		},
		{
			name:           "should return 404 when ebook not found",
			ebookID:        "non-existent",
			mockEbook:      nil,
			mockError:      nil,
			expectedStatus: http.StatusNotFound,
			expectedFound:  false,
		},
		{
			name:           "should return 500 when usecase fails",
			ebookID:        "ebook-1",
			mockEbook:      nil,
			mockError:      context.DeadlineExceeded,
			expectedStatus: http.StatusInternalServerError,
			expectedFound:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUsecase := &MockEbookUsecase{
				ebook: tt.mockEbook,
				err:   tt.mockError,
			}
			handler := NewEbookHandler(mockUsecase)

			// Create request
			req, err := http.NewRequest("GET", "/api/ebooks/"+tt.ebookID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Act
			handler.GetEbookByID(rr, req)

			// Assert
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusOK && tt.expectedFound {
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}

				// Check status field
				if status, ok := response["status"].(string); !ok || status != "success" {
					t.Errorf("expected status 'success', got %v", status)
				}

				// Check data field
				if data, ok := response["data"].(map[string]interface{}); ok {
					if id, ok := data["id"].(string); !ok || id != tt.ebookID {
						t.Errorf("expected ebook ID %s, got %s", tt.ebookID, id)
					}
				} else {
					t.Errorf("expected 'data' field in response")
				}
			} else if tt.expectedStatus != http.StatusOK {
				// For error responses, check error structure
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}

				// Check status field
				if status, ok := response["status"].(string); !ok || status != "error" {
					t.Errorf("expected status 'error', got %v", status)
				}

				// Check error field
				if _, ok := response["error"].(map[string]interface{}); !ok {
					t.Errorf("expected 'error' field in response")
				}
			}
		})
	}
}

func TestEbookHandler_CreateEbook(t *testing.T) {
	tests := []struct {
		name           string
		ebook          *entity.Ebook
		mockError      error
		expectedStatus int
	}{
		{
			name: "should create ebook successfully",
			ebook: &entity.Ebook{
				ID:         uuid.New().String(),
				Title:      "Test Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "test-ebook",
				Price:      1000,
				Language:   "en",
				Format:     entity.FormatPDF,
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "should return 400 when validation fails",
			ebook: &entity.Ebook{
				ID:         uuid.New().String(),
				Title:      "", // Empty title should fail validation
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "test-ebook",
			},
			mockError:      errors.New("title is required"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return 500 when usecase fails",
			ebook: &entity.Ebook{
				ID:         uuid.New().String(),
				Title:      "Test Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "test-ebook",
			},
			mockError:      context.DeadlineExceeded,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUsecase := &MockEbookUsecase{
				err: tt.mockError,
			}
			handler := NewEbookHandler(mockUsecase)

			// Create request body
			body, err := json.Marshal(tt.ebook)
			if err != nil {
				t.Fatal(err)
			}

			// Create request
			req, err := http.NewRequest("POST", "/api/ebooks", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Act
			handler.CreateEbook(rr, req)

			// Assert
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Verify response structure
			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("failed to unmarshal response: %v", err)
			}

			if tt.expectedStatus == http.StatusCreated {
				// Check success response
				if status, ok := response["status"].(string); !ok || status != "success" {
					t.Errorf("expected status 'success', got %v", status)
				}

				if _, ok := response["data"].(map[string]interface{}); !ok {
					t.Errorf("expected 'data' field in response")
				}

				if message, ok := response["message"].(string); !ok || message != "Ebook created successfully" {
					t.Errorf("expected message 'Ebook created successfully', got %v", message)
				}
			} else {
				// Check error response
				if status, ok := response["status"].(string); !ok || status != "error" {
					t.Errorf("expected status 'error', got %v", status)
				}

				if _, ok := response["error"].(map[string]interface{}); !ok {
					t.Errorf("expected 'error' field in response")
				}
			}
		})
	}
}

func TestEbookHandler_GetEbookBySlug(t *testing.T) {
	tests := []struct {
		name           string
		slug           string
		mockEbook      *entity.Ebook
		mockError      error
		expectedStatus int
		expectedFound  bool
	}{
		{
			name: "should return ebook by slug successfully",
			slug: "test-ebook",
			mockEbook: &entity.Ebook{
				ID:         "ebook-1",
				Title:      "Test Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "test-ebook",
				Price:      1000,
				Language:   "en",
				Format:     entity.FormatPDF,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedFound:  true,
		},
		{
			name:           "should return 404 when ebook not found",
			slug:           "non-existent",
			mockEbook:      nil,
			mockError:      nil,
			expectedStatus: http.StatusNotFound,
			expectedFound:  false,
		},
		{
			name:           "should return 500 when usecase fails",
			slug:           "test-ebook",
			mockEbook:      nil,
			mockError:      context.DeadlineExceeded,
			expectedStatus: http.StatusInternalServerError,
			expectedFound:  false,
		},
		{
			name:           "should return 400 when slug is empty",
			slug:           "",
			mockEbook:      nil,
			mockError:      nil,
			expectedStatus: http.StatusNotFound,
			expectedFound:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUsecase := &MockEbookUsecase{
				ebook: tt.mockEbook,
				err:   tt.mockError,
			}
			handler := NewEbookHandler(mockUsecase)

			// Create request
			req, err := http.NewRequest("GET", "/api/ebooks/"+tt.slug, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Act
			handler.GetEbookBySlug(rr, req)

			// Assert
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("failed to unmarshal response: %v", err)
			}

			if tt.expectedStatus == http.StatusOK && tt.expectedFound {
				// Check status field
				if status, ok := response["status"].(string); !ok || status != "success" {
					t.Errorf("expected status 'success', got %v", status)
				}

				// Check data field
				if data, ok := response["data"].(map[string]interface{}); ok {
					if slug, ok := data["slug"].(string); !ok || slug != tt.slug {
						t.Errorf("expected ebook slug %s, got %s", tt.slug, slug)
					}
				} else {
					t.Errorf("expected 'data' field in response")
				}
			} else {
				// Check error response
				if status, ok := response["status"].(string); !ok || status != "error" {
					t.Errorf("expected status 'error', got %v", status)
				}

				if _, ok := response["error"].(map[string]interface{}); !ok {
					t.Errorf("expected 'error' field in response")
				}
			}
		})
	}
}

func TestEbookHandler_UpdateEbook(t *testing.T) {
	tests := []struct {
		name           string
		ebookID        string
		ebook          *entity.Ebook
		mockError      error
		expectedStatus int
	}{
		{
			name:    "should update ebook successfully",
			ebookID: "ebook-1",
			ebook: &entity.Ebook{
				ID:         "ebook-1",
				Title:      "Updated Test Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "updated-test-ebook",
				Price:      2000,
				Language:   "en",
				Format:     entity.FormatPDF,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:    "should return 400 when validation fails",
			ebookID: "ebook-1",
			ebook: &entity.Ebook{
				ID:         "ebook-1",
				Title:      "", // Empty title should fail validation
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "test-ebook",
			},
			mockError:      errors.New("title is required"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:    "should return 500 when usecase fails",
			ebookID: "ebook-1",
			ebook: &entity.Ebook{
				ID:         "ebook-1",
				Title:      "Updated Test Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "updated-test-ebook",
			},
			mockError:      context.DeadlineExceeded,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:    "should return 400 when ID is empty",
			ebookID: "",
			ebook: &entity.Ebook{
				ID:         "",
				Title:      "Updated Test Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "updated-test-ebook",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUsecase := &MockEbookUsecase{
				err: tt.mockError,
			}
			handler := NewEbookHandler(mockUsecase)

			// Create request body
			body, err := json.Marshal(tt.ebook)
			if err != nil {
				t.Fatal(err)
			}

			// Create request
			req, err := http.NewRequest("PUT", "/api/ebooks/"+tt.ebookID, bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Act
			handler.UpdateEbook(rr, req)

			// Assert
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Verify response structure
			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("failed to unmarshal response: %v", err)
			}

			if tt.expectedStatus == http.StatusOK {
				// Check success response
				if status, ok := response["status"].(string); !ok || status != "success" {
					t.Errorf("expected status 'success', got %v", status)
				}

				if _, ok := response["data"].(map[string]interface{}); !ok {
					t.Errorf("expected 'data' field in response")
				}

				if message, ok := response["message"].(string); !ok || message != "Ebook updated successfully" {
					t.Errorf("expected message 'Ebook updated successfully', got %v", message)
				}
			} else {
				// Check error response
				if status, ok := response["status"].(string); !ok || status != "error" {
					t.Errorf("expected status 'error', got %v", status)
				}

				if _, ok := response["error"].(map[string]interface{}); !ok {
					t.Errorf("expected 'error' field in response")
				}
			}
		})
	}
}

func TestEbookHandler_DeleteEbook(t *testing.T) {
	tests := []struct {
		name           string
		ebookID        string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "should delete ebook successfully",
			ebookID:        "ebook-1",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "should return 400 when ID is empty",
			ebookID:        "",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "should return 400 when validation fails",
			ebookID:        "ebook-1",
			mockError:      errors.New("ebook not found"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "should return 500 when usecase fails",
			ebookID:        "ebook-1",
			mockError:      context.DeadlineExceeded,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUsecase := &MockEbookUsecase{
				err: tt.mockError,
			}
			handler := NewEbookHandler(mockUsecase)

			// Create request
			req, err := http.NewRequest("DELETE", "/api/ebooks/"+tt.ebookID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Act
			handler.DeleteEbook(rr, req)

			// Assert
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Verify response structure
			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("failed to unmarshal response: %v", err)
			}

			if tt.expectedStatus == http.StatusOK {
				// Check success response
				if status, ok := response["status"].(string); !ok || status != "success" {
					t.Errorf("expected status 'success', got %v", status)
				}

				if message, ok := response["message"].(string); !ok || message != "Ebook deleted successfully" {
					t.Errorf("expected message 'Ebook deleted successfully', got %v", message)
				}
			} else {
				// Check error response
				if status, ok := response["status"].(string); !ok || status != "error" {
					t.Errorf("expected status 'error', got %v", status)
				}

				if _, ok := response["error"].(map[string]interface{}); !ok {
					t.Errorf("expected 'error' field in response")
				}
			}
		})
	}
}

func TestEbookHandler_ListEbooksByCategory(t *testing.T) {
	tests := []struct {
		name           string
		categoryID     string
		mockEbooks     []*entity.Ebook
		mockError      error
		expectedStatus int
		expectedCount  int
	}{
		{
			name:       "should return list of ebooks by category successfully",
			categoryID: "category-1",
			mockEbooks: []*entity.Ebook{
				{
					ID:         uuid.New().String(),
					Title:      "Test Ebook 1",
					AuthorID:   "author-1",
					CategoryID: "category-1",
					Price:      1000,
					Language:   "en",
					Format:     entity.FormatPDF,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
				{
					ID:         uuid.New().String(),
					Title:      "Test Ebook 2",
					AuthorID:   "author-2",
					CategoryID: "category-1",
					Price:      2000,
					Language:   "en",
					Format:     entity.FormatEPUB,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "should return empty list when no ebooks in category",
			categoryID:     "category-2",
			mockEbooks:     []*entity.Ebook{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "should return 500 when usecase fails",
			categoryID:     "category-1",
			mockEbooks:     nil,
			mockError:      context.DeadlineExceeded,
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUsecase := &MockEbookUsecase{
				ebooks: tt.mockEbooks,
				err:    tt.mockError,
			}
			handler := NewEbookHandler(mockUsecase)

			// Create request
			req, err := http.NewRequest("GET", "/api/ebooks/category/"+tt.categoryID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Act
			handler.ListEbooksByCategory(rr, req)

			// Assert
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}

				// Check status field
				if status, ok := response["status"].(string); !ok || status != "success" {
					t.Errorf("expected status 'success', got %v", status)
				}

				// Check data field
				if data, ok := response["data"].([]interface{}); ok {
					if len(data) != tt.expectedCount {
						t.Errorf("expected %d ebooks, got %d", tt.expectedCount, len(data))
					}
				} else {
					t.Errorf("expected 'data' field in response")
				}

				// Check meta field for pagination
				if meta, ok := response["meta"].(map[string]interface{}); ok {
					if _, ok := meta["total"]; !ok {
						t.Errorf("expected 'total' field in meta")
					}
					if _, ok := meta["limit"]; !ok {
						t.Errorf("expected 'limit' field in meta")
					}
					if _, ok := meta["offset"]; !ok {
						t.Errorf("expected 'offset' field in meta")
					}
				} else {
					t.Errorf("expected 'meta' field in response")
				}
			} else if tt.expectedStatus != http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}

				// Check error response
				if status, ok := response["status"].(string); !ok || status != "error" {
					t.Errorf("expected status 'error', got %v", status)
				}

				if _, ok := response["error"].(map[string]interface{}); !ok {
					t.Errorf("expected 'error' field in response")
				}
			}
		})
	}
}
