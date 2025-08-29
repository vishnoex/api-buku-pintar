package service

import (
	"buku-pintar/internal/domain/entity"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

// MockEbookRepository is a mock implementation for testing
type MockEbookRepository struct {
	ebooks    []*entity.Ebook
	ebookList []*entity.EbookList
	ebook     *entity.Ebook
	err       error
	count     int64
}

func (m *MockEbookRepository) Create(ctx context.Context, ebook *entity.Ebook) error {
	return m.err
}

func (m *MockEbookRepository) GetByID(ctx context.Context, id string) (*entity.Ebook, error) {
	return m.ebook, m.err
}

func (m *MockEbookRepository) GetBySlug(ctx context.Context, slug string) (*entity.Ebook, error) {
	return m.ebook, m.err
}

func (m *MockEbookRepository) Update(ctx context.Context, ebook *entity.Ebook) error {
	return m.err
}

func (m *MockEbookRepository) Delete(ctx context.Context, id string) error {
	return m.err
}

func (m *MockEbookRepository) List(ctx context.Context, limit, offset int) ([]*entity.EbookList, error) {
	return m.ebookList, m.err
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

// MockEbookRedisRepository is a mock implementation for testing
type MockEbookRedisRepository struct {
	ebooks    []*entity.Ebook
	ebookList []*entity.EbookList
	ebook     *entity.Ebook
	err       error
	count     int64
	cacheHit  bool
}

func (m *MockEbookRedisRepository) GetEbookList(ctx context.Context, limit, offset int) ([]*entity.EbookList, error) {
	if m.cacheHit {
		return m.ebookList, nil
	}
	return nil, m.err
}

func (m *MockEbookRedisRepository) SetEbookList(ctx context.Context, ebooks []*entity.EbookList, limit, offset int) error {
	return m.err
}

func (m *MockEbookRedisRepository) GetEbookTotal(ctx context.Context) (int64, error) {
	if m.cacheHit {
		return m.count, nil
	}
	return 0, m.err
}

func (m *MockEbookRedisRepository) SetEbookTotal(ctx context.Context, count int64) error {
	return m.err
}

func (m *MockEbookRedisRepository) GetEbookByID(ctx context.Context, id string) (*entity.Ebook, error) {
	if m.cacheHit {
		return m.ebook, nil
	}
	return nil, m.err
}

func (m *MockEbookRedisRepository) SetEbookByID(ctx context.Context, ebook *entity.Ebook) error {
	return m.err
}

func (m *MockEbookRedisRepository) GetEbookBySlug(ctx context.Context, slug string) (*entity.Ebook, error) {
	if m.cacheHit {
		return m.ebook, nil
	}
	return nil, m.err
}

func (m *MockEbookRedisRepository) SetEbookBySlug(ctx context.Context, ebook *entity.Ebook) error {
	return m.err
}

func (m *MockEbookRedisRepository) GetEbookListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Ebook, error) {
	if m.cacheHit {
		return m.ebooks, nil
	}
	return nil, m.err
}

func (m *MockEbookRedisRepository) SetEbookListByCategory(ctx context.Context, ebooks []*entity.Ebook, categoryID string, limit, offset int) error {
	return m.err
}

func (m *MockEbookRedisRepository) GetEbookCountByCategory(ctx context.Context, categoryID string) (int64, error) {
	if m.cacheHit {
		return m.count, nil
	}
	return 0, m.err
}

func (m *MockEbookRedisRepository) SetEbookCountByCategory(ctx context.Context, categoryID string, count int64) error {
	return m.err
}

func (m *MockEbookRedisRepository) GetEbookListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Ebook, error) {
	if m.cacheHit {
		return m.ebooks, nil
	}
	return nil, m.err
}

func (m *MockEbookRedisRepository) SetEbookListByAuthor(ctx context.Context, ebooks []*entity.Ebook, authorID string, limit, offset int) error {
	return m.err
}

func (m *MockEbookRedisRepository) GetEbookCountByAuthor(ctx context.Context, authorID string) (int64, error) {
	if m.cacheHit {
		return m.count, nil
	}
	return 0, m.err
}

func (m *MockEbookRedisRepository) SetEbookCountByAuthor(ctx context.Context, authorID string, count int64) error {
	return m.err
}

func (m *MockEbookRedisRepository) InvalidateEbookCache(ctx context.Context) error {
	return m.err
}

func TestEbookService_GetEbookByID_CacheHit(t *testing.T) {
	// Arrange
	expectedEbook := &entity.Ebook{
		ID:         uuid.New().String(),
		Title:      "Test Ebook",
		AuthorID:   "author-1",
		CategoryID: "category-1",
		Price:      1000,
		Language:   "en",
		Format:     entity.FormatPDF,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mockRepo := &MockEbookRepository{}
	mockRedisRepo := &MockEbookRedisRepository{
		ebook:    expectedEbook,
		cacheHit: true,
	}

	service := NewEbookService(mockRepo, mockRedisRepo)
	ctx := context.Background()

	// Act
	result, err := service.GetEbookByID(ctx, expectedEbook.ID)

	// Assert
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	if result == nil {
		t.Error("expected ebook but got nil")
		return
	}
	if result.ID != expectedEbook.ID {
		t.Errorf("expected ebook ID %s, got %s", expectedEbook.ID, result.ID)
	}
}

func TestEbookService_GetEbookByID_CacheMiss(t *testing.T) {
	// Arrange
	expectedEbook := &entity.Ebook{
		ID:         uuid.New().String(),
		Title:      "Test Ebook",
		AuthorID:   "author-1",
		CategoryID: "category-1",
		Price:      1000,
		Language:   "en",
		Format:     entity.FormatPDF,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mockRepo := &MockEbookRepository{
		ebook: expectedEbook,
	}
	mockRedisRepo := &MockEbookRedisRepository{
		cacheHit: false,
	}

	service := NewEbookService(mockRepo, mockRedisRepo)
	ctx := context.Background()

	// Act
	result, err := service.GetEbookByID(ctx, expectedEbook.ID)

	// Assert
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	if result == nil {
		t.Error("expected ebook but got nil")
		return
	}
	if result.ID != expectedEbook.ID {
		t.Errorf("expected ebook ID %s, got %s", expectedEbook.ID, result.ID)
	}
}

func TestEbookService_GetEbookList_CacheHit(t *testing.T) {
	// Arrange
	expectedEbooks := []*entity.EbookList{
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
	}

	mockRepo := &MockEbookRepository{}
	mockRedisRepo := &MockEbookRedisRepository{
		ebookList: expectedEbooks,
		cacheHit:  true,
	}

	service := NewEbookService(mockRepo, mockRedisRepo)
	ctx := context.Background()

	// Act
	result, err := service.GetEbookList(ctx, 10, 0)

	// Assert
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	if len(result) != len(expectedEbooks) {
		t.Errorf("expected %d ebooks, got %d", len(expectedEbooks), len(result))
	}
}

func TestEbookService_GetEbookCount_CacheHit(t *testing.T) {
	// Arrange
	expectedCount := int64(100)

	mockRepo := &MockEbookRepository{}
	mockRedisRepo := &MockEbookRedisRepository{
		count:    expectedCount,
		cacheHit: true,
	}

	service := NewEbookService(mockRepo, mockRedisRepo)
	ctx := context.Background()

	// Act
	result, err := service.GetEbookCount(ctx)

	// Assert
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	if result != expectedCount {
		t.Errorf("expected count %d, got %d", expectedCount, result)
	}
}
