package http

import (
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

func (m *MockEbookUsecase) ListEbooks(ctx context.Context, limit, offset int) ([]*entity.Ebook, error) {
	return m.ebooks, m.err
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
					ID:          uuid.New().String(),
					Title:       "Test Ebook 1",
					AuthorID:    "author-1",
					CategoryID:  "category-1",
					Price:       1000,
					Language:    "en",
					Format:      entity.FormatPDF,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					ID:          uuid.New().String(),
					Title:       "Test Ebook 2",
					AuthorID:    "author-2",
					CategoryID:  "category-2",
					Price:       2000,
					Language:    "en",
					Format:      entity.FormatEPUB,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
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

				if data, ok := response["data"].([]interface{}); ok {
					if len(data) != tt.expectedCount {
						t.Errorf("expected %d ebooks, got %d", tt.expectedCount, len(data))
					}
				} else {
					t.Errorf("expected 'data' field in response")
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
				ID:          "ebook-1",
				Title:       "Test Ebook",
				AuthorID:    "author-1",
				CategoryID:  "category-1",
				Price:       1000,
				Language:    "en",
				Format:      entity.FormatPDF,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
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

				if data, ok := response["data"].(map[string]interface{}); ok {
					if id, ok := data["id"].(string); !ok || id != tt.ebookID {
						t.Errorf("expected ebook ID %s, got %s", tt.ebookID, id)
					}
				} else {
					t.Errorf("expected 'data' field in response")
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
		})
	}
}
