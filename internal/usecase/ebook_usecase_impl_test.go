package usecase

import (
	"buku-pintar/internal/domain/entity"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

// MockEbookRepository is a simple mock implementation for testing
type MockEbookRepository struct {
	ebooks     []*entity.Ebook
	ebook      *entity.Ebook
	err        error
	count      int64
	createFunc func(ctx context.Context, ebook *entity.Ebook) error
	listFunc   func(ctx context.Context, limit, offset int) ([]*entity.Ebook, error)
	getByIDFunc func(ctx context.Context, id string) (*entity.Ebook, error)
	getBySlugFunc func(ctx context.Context, slug string) (*entity.Ebook, error)
	// For testing specific scenarios
	existingBySlug *entity.Ebook
}

func (m *MockEbookRepository) Create(ctx context.Context, ebook *entity.Ebook) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, ebook)
	}
	return m.err
}

func (m *MockEbookRepository) GetByID(ctx context.Context, id string) (*entity.Ebook, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return m.ebook, m.err
}

func (m *MockEbookRepository) GetBySlug(ctx context.Context, slug string) (*entity.Ebook, error) {
	if m.getBySlugFunc != nil {
		return m.getBySlugFunc(ctx, slug)
	}
	if m.existingBySlug != nil {
		return m.existingBySlug, m.err
	}
	return m.ebook, m.err
}

func (m *MockEbookRepository) Update(ctx context.Context, ebook *entity.Ebook) error {
	return m.err
}

func (m *MockEbookRepository) Delete(ctx context.Context, id string) error {
	return m.err
}

func (m *MockEbookRepository) List(ctx context.Context, limit, offset int) ([]*entity.Ebook, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, limit, offset)
	}
	return m.ebooks, m.err
}

func (m *MockEbookRepository) ListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Ebook, error) {
	return m.ebooks, m.err
}

func (m *MockEbookRepository) ListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Ebook, error) {
	return m.ebooks, m.err
}

func (m *MockEbookRepository) Count(ctx context.Context) (int64, error) {
	return m.count, m.err
}

func (m *MockEbookRepository) CountByCategory(ctx context.Context, categoryID string) (int64, error) {
	return m.count, m.err
}

func (m *MockEbookRepository) CountByAuthor(ctx context.Context, authorID string) (int64, error) {
	return m.count, m.err
}

func TestEbookUsecase_ListEbooks(t *testing.T) {
	tests := []struct {
		name           string
		limit          int
		offset         int
		mockEbooks     []*entity.Ebook
		mockError      error
		expectedLength int
		expectedError  bool
	}{
		{
			name:   "should return list of ebooks successfully",
			limit:  10,
			offset: 0,
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
			expectedLength: 2,
			expectedError:  false,
		},
		{
			name:           "should return empty list when no ebooks exist",
			limit:          10,
			offset:         0,
			mockEbooks:     []*entity.Ebook{},
			mockError:      nil,
			expectedLength: 0,
			expectedError:  false,
		},
		{
			name:           "should return error when repository fails",
			limit:          10,
			offset:         0,
			mockEbooks:     nil,
			mockError:      errors.New("database error"),
			expectedLength: 0,
			expectedError:  true,
		},
		{
			name:           "should use default limit when limit is 0",
			limit:          0,
			offset:         0,
			mockEbooks:     []*entity.Ebook{},
			mockError:      nil,
			expectedLength: 0,
			expectedError:  false,
		},
		{
			name:           "should use default limit when limit is negative",
			limit:          -5,
			offset:         0,
			mockEbooks:     []*entity.Ebook{},
			mockError:      nil,
			expectedLength: 0,
			expectedError:  false,
		},
		{
			name:           "should use 0 offset when offset is negative",
			limit:          10,
			offset:         -5,
			mockEbooks:     []*entity.Ebook{},
			mockError:      nil,
			expectedLength: 0,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockEbookRepository{
				ebooks: tt.mockEbooks,
				err:    tt.mockError,
			}
			usecase := NewEbookUsecase(mockRepo)
			ctx := context.Background()

			// Act
			result, err := usecase.ListEbooks(ctx, tt.limit, tt.offset)

			// Assert
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if len(result) != tt.expectedLength {
					t.Errorf("expected %d ebooks, got %d", tt.expectedLength, len(result))
				}
			}
		})
	}
}

func TestEbookUsecase_ListEbooksByCategory(t *testing.T) {
	tests := []struct {
		name           string
		categoryID     string
		limit          int
		offset         int
		mockEbooks     []*entity.Ebook
		mockError      error
		expectedLength int
		expectedError  bool
	}{
		{
			name:       "should return ebooks by category successfully",
			categoryID: "category-1",
			limit:      10,
			offset:     0,
			mockEbooks: []*entity.Ebook{
				{
					ID:          uuid.New().String(),
					Title:       "Category Ebook 1",
					AuthorID:    "author-1",
					CategoryID:  "category-1",
					Price:       1000,
					Language:    "en",
					Format:      entity.FormatPDF,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			mockError:      nil,
			expectedLength: 1,
			expectedError:  false,
		},
		{
			name:           "should return empty list when no ebooks in category",
			categoryID:     "category-2",
			limit:          10,
			offset:         0,
			mockEbooks:     []*entity.Ebook{},
			mockError:      nil,
			expectedLength: 0,
			expectedError:  false,
		},
		{
			name:           "should return error when repository fails",
			categoryID:     "category-1",
			limit:          10,
			offset:         0,
			mockEbooks:     nil,
			mockError:      errors.New("database error"),
			expectedLength: 0,
			expectedError:  true,
		},
		{
			name:           "should return error when category_id is empty",
			categoryID:     "",
			limit:          10,
			offset:         0,
			mockEbooks:     nil,
			mockError:      nil,
			expectedLength: 0,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockEbookRepository{
				ebooks: tt.mockEbooks,
				err:    tt.mockError,
			}
			usecase := NewEbookUsecase(mockRepo)
			ctx := context.Background()

			// Act
			result, err := usecase.ListEbooksByCategory(ctx, tt.categoryID, tt.limit, tt.offset)

			// Assert
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if len(result) != tt.expectedLength {
					t.Errorf("expected %d ebooks, got %d", tt.expectedLength, len(result))
				}
			}
		})
	}
}

func TestEbookUsecase_GetEbookByID(t *testing.T) {
	tests := []struct {
		name          string
		ebookID       string
		mockEbook     *entity.Ebook
		mockError     error
		expectedFound bool
		expectedError bool
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
			mockError:     nil,
			expectedFound: true,
			expectedError: false,
		},
		{
			name:          "should return nil when ebook not found",
			ebookID:       "non-existent",
			mockEbook:     nil,
			mockError:     nil,
			expectedFound: false,
			expectedError: false,
		},
		{
			name:          "should return error when repository fails",
			ebookID:       "ebook-1",
			mockEbook:     nil,
			mockError:     errors.New("database error"),
			expectedFound: false,
			expectedError: true,
		},
		{
			name:          "should return error when ID is empty",
			ebookID:       "",
			mockEbook:     nil,
			mockError:     nil,
			expectedFound: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockEbookRepository{
				ebook: tt.mockEbook,
				err:   tt.mockError,
			}
			usecase := NewEbookUsecase(mockRepo)
			ctx := context.Background()

			// Act
			result, err := usecase.GetEbookByID(ctx, tt.ebookID)

			// Assert
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if tt.expectedFound {
					if result == nil {
						t.Errorf("expected ebook but got nil")
					} else if result.ID != tt.ebookID {
						t.Errorf("expected ebook ID %s, got %s", tt.ebookID, result.ID)
					}
				} else {
					if result != nil {
						t.Errorf("expected nil but got ebook")
					}
				}
			}
		})
	}
}

func TestEbookUsecase_GetEbookBySlug(t *testing.T) {
	tests := []struct {
		name          string
		slug          string
		mockEbook     *entity.Ebook
		mockError     error
		expectedFound bool
		expectedError bool
	}{
		{
			name:  "should return ebook by slug successfully",
			slug:  "test-ebook",
			mockEbook: &entity.Ebook{
				ID:          "ebook-1",
				Title:       "Test Ebook",
				Slug:        "test-ebook",
				AuthorID:    "author-1",
				CategoryID:  "category-1",
				Price:       1000,
				Language:    "en",
				Format:      entity.FormatPDF,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			mockError:     nil,
			expectedFound: true,
			expectedError: false,
		},
		{
			name:          "should return nil when ebook not found",
			slug:          "non-existent",
			mockEbook:     nil,
			mockError:     nil,
			expectedFound: false,
			expectedError: false,
		},
		{
			name:          "should return error when repository fails",
			slug:          "test-ebook",
			mockEbook:     nil,
			mockError:     errors.New("database error"),
			expectedFound: false,
			expectedError: true,
		},
		{
			name:          "should return error when slug is empty",
			slug:          "",
			mockEbook:     nil,
			mockError:     nil,
			expectedFound: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockEbookRepository{
				ebook: tt.mockEbook,
				err:   tt.mockError,
			}
			usecase := NewEbookUsecase(mockRepo)
			ctx := context.Background()

			// Act
			result, err := usecase.GetEbookBySlug(ctx, tt.slug)

			// Assert
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if tt.expectedFound {
					if result == nil {
						t.Errorf("expected ebook but got nil")
					} else if result.Slug != tt.slug {
						t.Errorf("expected ebook slug %s, got %s", tt.slug, result.Slug)
					}
				} else {
					if result != nil {
						t.Errorf("expected nil but got ebook")
					}
				}
			}
		})
	}
}

func TestEbookUsecase_CreateEbook(t *testing.T) {
	tests := []struct {
		name          string
		ebook         *entity.Ebook
		mockError     error
		mockExisting  *entity.Ebook
		expectedError bool
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
			mockError:     nil,
			mockExisting:  nil,
			expectedError: false,
		},
		{
			name: "should return error when title is empty",
			ebook: &entity.Ebook{
				ID:         uuid.New().String(),
				Title:      "",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "test-ebook",
			},
			mockError:     nil,
			mockExisting:  nil,
			expectedError: true,
		},
		{
			name: "should return error when author_id is empty",
			ebook: &entity.Ebook{
				ID:         uuid.New().String(),
				Title:      "Test Ebook",
				AuthorID:   "",
				CategoryID: "category-1",
				Slug:       "test-ebook",
			},
			mockError:     nil,
			mockExisting:  nil,
			expectedError: true,
		},
		{
			name: "should return error when category_id is empty",
			ebook: &entity.Ebook{
				ID:         uuid.New().String(),
				Title:      "Test Ebook",
				AuthorID:   "author-1",
				CategoryID: "",
				Slug:       "test-ebook",
			},
			mockError:     nil,
			mockExisting:  nil,
			expectedError: true,
		},
		{
			name: "should return error when slug is empty",
			ebook: &entity.Ebook{
				ID:         uuid.New().String(),
				Title:      "Test Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "",
			},
			mockError:     nil,
			mockExisting:  nil,
			expectedError: true,
		},
		{
			name: "should return error when slug already exists",
			ebook: &entity.Ebook{
				ID:         uuid.New().String(),
				Title:      "Test Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "existing-ebook",
			},
			mockError: nil,
			mockExisting: &entity.Ebook{
				ID:   "existing-id",
				Slug: "existing-ebook",
			},
			expectedError: true,
		},
		{
			name: "should return error when repository fails",
			ebook: &entity.Ebook{
				ID:         uuid.New().String(),
				Title:      "Test Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "test-ebook",
			},
			mockError:     errors.New("database error"),
			mockExisting:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockEbookRepository{
				ebook: tt.mockExisting,
				err:   tt.mockError,
			}
			usecase := NewEbookUsecase(mockRepo)
			ctx := context.Background()

			// Act
			err := usecase.CreateEbook(ctx, tt.ebook)

			// Assert
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestEbookUsecase_UpdateEbook(t *testing.T) {
	tests := []struct {
		name          string
		ebook         *entity.Ebook
		mockExisting  *entity.Ebook
		mockError     error
		expectedError bool
	}{
		{
			name: "should update ebook successfully",
			ebook: &entity.Ebook{
				ID:         "ebook-1",
				Title:      "Updated Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "updated-ebook",
			},
			mockExisting: &entity.Ebook{
				ID:   "ebook-1",
				Slug: "old-slug",
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "should return error when ID is empty",
			ebook: &entity.Ebook{
				ID:         "",
				Title:      "Updated Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "updated-ebook",
			},
			mockExisting:  nil,
			mockError:     nil,
			expectedError: true,
		},
		{
			name: "should return error when ebook not found",
			ebook: &entity.Ebook{
				ID:         "non-existent",
				Title:      "Updated Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "updated-ebook",
			},
			mockExisting:  nil,
			mockError:     nil,
			expectedError: true,
		},
		{
			name: "should return error when new slug already exists",
			ebook: &entity.Ebook{
				ID:         "ebook-1",
				Title:      "Updated Ebook",
				AuthorID:   "author-1",
				CategoryID: "category-1",
				Slug:       "existing-slug",
			},
			mockExisting: &entity.Ebook{
				ID:   "ebook-1",
				Slug: "old-slug",
			},
			mockError: nil,
			// Mock that another ebook with the new slug exists
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockEbookRepository{
				ebook: tt.mockExisting,
				err:   tt.mockError,
			}
					// For update test, we need to handle the slug check differently
		if tt.name == "should update ebook successfully" {
			mockRepo.getBySlugFunc = func(ctx context.Context, slug string) (*entity.Ebook, error) {
				// Return nil for the new slug, indicating no existing ebook with that slug
				return nil, nil
			}
		}
			usecase := NewEbookUsecase(mockRepo)
			ctx := context.Background()

			// Act
			err := usecase.UpdateEbook(ctx, tt.ebook)

			// Assert
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestEbookUsecase_DeleteEbook(t *testing.T) {
	tests := []struct {
		name          string
		ebookID       string
		mockExisting  *entity.Ebook
		mockError     error
		expectedError bool
	}{
		{
			name:    "should delete ebook successfully",
			ebookID: "ebook-1",
			mockExisting: &entity.Ebook{
				ID: "ebook-1",
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "should return error when ID is empty",
			ebookID:       "",
			mockExisting:  nil,
			mockError:     nil,
			expectedError: true,
		},
		{
			name:          "should return error when ebook not found",
			ebookID:       "non-existent",
			mockExisting:  nil,
			mockError:     nil,
			expectedError: true,
		},
		{
			name:          "should return error when repository fails",
			ebookID:       "ebook-1",
			mockExisting: &entity.Ebook{
				ID: "ebook-1",
			},
			mockError:     errors.New("database error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockEbookRepository{
				ebook: tt.mockExisting,
				err:   tt.mockError,
			}
			usecase := NewEbookUsecase(mockRepo)
			ctx := context.Background()

			// Act
			err := usecase.DeleteEbook(ctx, tt.ebookID)

			// Assert
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}
