package usecase

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

// MockEbookService is a simple mock implementation for testing
type MockEbookService struct {
	ebooks        []*entity.Ebook
	ebookList     []*entity.EbookList
	ebook         *entity.Ebook
	ebookDetail   *entity.EbookDetail
	err           error
	count         int64
	createFunc    func(ctx context.Context, ebook *entity.Ebook) error
	listFunc      func(ctx context.Context, limit, offset int) ([]*entity.EbookList, error)
	getByIDFunc   func(ctx context.Context, id string) (*entity.Ebook, error)
	getBySlugFunc func(ctx context.Context, slug string) (*entity.EbookDetail, error)
	// For testing specific scenarios
	existingBySlug *entity.EbookDetail
}

type MockEbookDiscountService struct {
	discountList []*entity.EbookDiscount
	discount     *entity.EbookDiscount
	err          error
	count        int64
	createFunc   func(ctx context.Context, discount *entity.EbookDiscount) error
	listFunc     func(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error)
	getByIDFunc  func(ctx context.Context, id string) (*entity.EbookDiscount, error)
	getByEbookIDFunc func(ctx context.Context, ebookID string) ([]*entity.EbookDiscount, error)
	getActiveDiscountsFunc func(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error)
	getActiveDiscountByEbookIDFunc func(ctx context.Context, ebookID string) (*entity.EbookDiscount, error)
	countFunc func(ctx context.Context) (int64, error)
	countByEbookIDFunc func(ctx context.Context, ebookID string) (int64, error)
	countActiveDiscountsFunc func(ctx context.Context) (int64, error)
}

// Implement all EbookDiscountService interface methods
func (m *MockEbookDiscountService) CreateDiscount(ctx context.Context, discount *entity.EbookDiscount) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, discount)
	}
	return m.err
}

func (m *MockEbookDiscountService) GetDiscountByID(ctx context.Context, id string) (*entity.EbookDiscount, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return m.discount, m.err
}

func (m *MockEbookDiscountService) UpdateDiscount(ctx context.Context, discount *entity.EbookDiscount) error {
	return m.err
}

func (m *MockEbookDiscountService) DeleteDiscount(ctx context.Context, id string) error {
	return m.err
}

func (m *MockEbookDiscountService) GetDiscountList(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, limit, offset)
	}
	return m.discountList, m.err
}

func (m *MockEbookDiscountService) GetDiscountsByEbookID(ctx context.Context, ebookID string) ([]*entity.EbookDiscount, error) {
	if m.getByEbookIDFunc != nil {
		return m.getByEbookIDFunc(ctx, ebookID)
	}
	return m.discountList, m.err
}

func (m *MockEbookDiscountService) GetActiveDiscounts(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error) {
	if m.getActiveDiscountsFunc != nil {
		return m.getActiveDiscountsFunc(ctx, limit, offset)
	}
	return m.discountList, m.err
}

func (m *MockEbookDiscountService) GetActiveDiscountByEbookID(ctx context.Context, ebookID string) (*entity.EbookDiscount, error) {
	if m.getActiveDiscountByEbookIDFunc != nil {
		return m.getActiveDiscountByEbookIDFunc(ctx, ebookID)
	}
	return m.discount, m.err
}

func (m *MockEbookDiscountService) GetDiscountCount(ctx context.Context) (int64, error) {
	if m.countFunc != nil {
		return m.countFunc(ctx)
	}
	return m.count, m.err
}

func (m *MockEbookDiscountService) GetDiscountCountByEbookID(ctx context.Context, ebookID string) (int64, error) {
	if m.countByEbookIDFunc != nil {
		return m.countByEbookIDFunc(ctx, ebookID)
	}
	return m.count, m.err
}

func (m *MockEbookDiscountService) GetActiveDiscountCount(ctx context.Context) (int64, error) {
	if m.countActiveDiscountsFunc != nil {
		return m.countActiveDiscountsFunc(ctx)
	}
	return m.count, m.err
}

func (m *MockEbookService) CreateEbook(ctx context.Context, ebook *entity.Ebook) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, ebook)
	}
	return m.err
}

func (m *MockEbookService) GetEbookByID(ctx context.Context, id string) (*entity.Ebook, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return m.ebook, m.err
}

func (m *MockEbookService) GetEbookBySlug(ctx context.Context, slug string) (*entity.EbookDetail, error) {
	if m.getBySlugFunc != nil {
		return m.getBySlugFunc(ctx, slug)
	}
	if m.existingBySlug != nil {
		return m.existingBySlug, m.err
	}
	return m.ebookDetail, m.err
}

func (m *MockEbookService) UpdateEbook(ctx context.Context, ebook *entity.Ebook) error {
	return m.err
}

func (m *MockEbookService) DeleteEbook(ctx context.Context, id string) error {
	return m.err
}

func (m *MockEbookService) GetEbookList(ctx context.Context, limit, offset int) ([]*entity.EbookList, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, limit, offset)
	}
	return m.ebookList, m.err
}

func (m *MockEbookService) GetEbookListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Ebook, error) {
	return m.ebooks, m.err
}

func (m *MockEbookService) GetEbookListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Ebook, error) {
	return m.ebooks, m.err
}

func (m *MockEbookService) GetEbookCount(ctx context.Context) (int64, error) {
	return m.count, m.err
}

func (m *MockEbookService) GetEbookCountByCategory(ctx context.Context, categoryID string) (int64, error) {
	return m.count, m.err
}

func (m *MockEbookService) GetEbookCountByAuthor(ctx context.Context, authorID string) (int64, error) {
	return m.count, m.err
}

// Ensure MockEbookService implements service.EbookService interface
var _ service.EbookService = (*MockEbookService)(nil)

// Ensure MockEbookDiscountService implements service.EbookDiscountService interface
var _ service.EbookDiscountService = (*MockEbookDiscountService)(nil)

func TestEbookUsecase_ListEbooks(t *testing.T) {
	tests := []struct {
		name           string
		limit          int
		offset         int
		mockEbookList  []*entity.EbookList
		mockDiscountList []*entity.EbookDiscount
		mockError      error
		expectedLength int
		expectedError  bool
	}{
		{
			name:   "should return list of ebooks successfully",
			limit:  10,
			offset: 0,
			mockEbookList: []*entity.EbookList{
				{
					ID:         uuid.New().String(),
					Title:      "Test Ebook 1",
					Slug:       "test-ebook-1",
					CoverImage: "cover1.jpg",
					Price:      1000,
				},
				{
					ID:         uuid.New().String(),
					Title:      "Test Ebook 2",
					Slug:       "test-ebook-2",
					CoverImage: "cover2.jpg",
					Price:      2000,
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
			mockEbookList:  []*entity.EbookList{},
			mockError:      nil,
			expectedLength: 0,
			expectedError:  false,
		},
		{
			name:           "should return error when repository fails",
			limit:          10,
			offset:         0,
			mockEbookList:  nil,
			mockError:      errors.New("database error"),
			expectedLength: 0,
			expectedError:  true,
		},
		{
			name:           "should use default limit when limit is 0",
			limit:          0,
			offset:         0,
			mockEbookList:  []*entity.EbookList{},
			mockError:      nil,
			expectedLength: 0,
			expectedError:  false,
		},
		{
			name:           "should use default limit when limit is negative",
			limit:          -5,
			offset:         0,
			mockEbookList:  []*entity.EbookList{},
			mockError:      nil,
			expectedLength: 0,
			expectedError:  false,
		},
		{
			name:           "should use 0 offset when offset is negative",
			limit:          10,
			offset:         -5,
			mockEbookList:  []*entity.EbookList{},
			mockError:      nil,
			expectedLength: 0,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockEbookService{
				ebookList: tt.mockEbookList,
				err:       tt.mockError,
			}
			mockRepoDiscount := &MockEbookDiscountService{
				discountList: tt.mockDiscountList,
				err:       tt.mockError,
			}
			usecase := NewEbookUsecase(mockRepo, mockRepoDiscount)
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
		mockDiscountList []*entity.EbookDiscount
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
					ID:         uuid.New().String(),
					Title:      "Category Ebook 1",
					AuthorID:   "author-1",
					CategoryID: "category-1",
					Price:      1000,
					Language:   "en",
					Format:     entity.FormatPDF,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
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
			mockRepo := &MockEbookService{
				ebooks: tt.mockEbooks,
				err:    tt.mockError,
			}
			mockRepoDiscount := &MockEbookDiscountService{
				discountList: tt.mockDiscountList,
				err:       tt.mockError,
			}
			usecase := NewEbookUsecase(mockRepo, mockRepoDiscount)
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
		mockDiscount  *entity.EbookDiscount
		mockError     error
		expectedFound bool
		expectedError bool
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
			mockRepo := &MockEbookService{
				ebook: tt.mockEbook,
				err:   tt.mockError,
			}
			mockRepoDiscount := &MockEbookDiscountService{
				discount: tt.mockDiscount,
				err:       tt.mockError,
			}
			usecase := NewEbookUsecase(mockRepo, mockRepoDiscount)
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
		mockEbook     *entity.EbookDetail
		mockDiscount  *entity.EbookDiscount
		mockError     error
		expectedFound bool
		expectedError bool
	}{
		{
			name: "should return ebook by slug successfully",
			slug: "test-ebook",
			mockEbook: &entity.EbookDetail{
				ID:         "ebook-1",
				Title:      "Test Ebook",
				Slug:       "test-ebook",
				AuthorID:   "author-1",
				Price:      1000,
				Language:   "en",
				Format:     entity.FormatPDF,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
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
			mockRepo := &MockEbookService{
				ebookDetail: tt.mockEbook,
				err:        tt.mockError,
			}
			mockRepoDiscount := &MockEbookDiscountService{
				discount: tt.mockDiscount,
				err:       tt.mockError,
			}
			usecase := NewEbookUsecase(mockRepo, mockRepoDiscount)
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
						t.Errorf("expected ebook response but got nil")
					} else if result.Slug != tt.slug {
						t.Errorf("expected ebook slug %s, got %s", tt.slug, result.Slug)
					}
				} else {
					if result != nil {
						t.Errorf("expected nil but got ebook response")
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
		mockDiscount  *entity.EbookDiscount
		mockError     error
		mockExisting  *entity.EbookDetail
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
			mockExisting: &entity.EbookDetail{
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
			mockRepo := &MockEbookService{
				ebookDetail: tt.mockExisting,
				err:        tt.mockError,
			}
			mockRepoDiscount := &MockEbookDiscountService{
				discount: tt.mockDiscount,
				err:       tt.mockError,
			}
			usecase := NewEbookUsecase(mockRepo, mockRepoDiscount)
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
		mockDiscount  *entity.EbookDiscount
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
			mockRepo := &MockEbookService{
				ebook: tt.mockExisting,
				err:   tt.mockError,
			}

			mockRepoDiscount := &MockEbookDiscountService{
				discount: tt.mockDiscount,
				err:       tt.mockError,
			}
			// For update test, we need to handle the slug check differently
			if tt.name == "should update ebook successfully" {
				mockRepo.getBySlugFunc = func(ctx context.Context, slug string) (*entity.EbookDetail, error) {
					// Return nil for the new slug, indicating no existing ebook with that slug
					return nil, nil
				}
			} else if tt.name == "should return error when new slug already exists" {
				mockRepo.getBySlugFunc = func(ctx context.Context, slug string) (*entity.EbookDetail, error) {
					// Return an existing ebook for the new slug, indicating conflict
					return &entity.EbookDetail{
						ID:   "other-ebook-id",
						Slug: slug,
					}, nil
				}
			}
			usecase := NewEbookUsecase(mockRepo, mockRepoDiscount)
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
		mockDiscount  *entity.EbookDiscount
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
			name:    "should return error when repository fails",
			ebookID: "ebook-1",
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
			mockRepo := &MockEbookService{
				ebook: tt.mockExisting,
				err:   tt.mockError,
			}
			mockRepoDiscount := &MockEbookDiscountService{
				discount: tt.mockDiscount,
				err:       tt.mockError,
			}
			usecase := NewEbookUsecase(mockRepo, mockRepoDiscount)
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
