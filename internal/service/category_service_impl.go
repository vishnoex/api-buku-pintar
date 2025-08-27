package service

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/internal/domain/service"
	"context"
	"log"
	"time"
)

type categoryService struct {
	categoryRepo     repository.CategoryRepository
	categoryRedisRepo repository.CategoryRedisRepository
	cacheTTL         time.Duration
}

// NewCategoryService creates a new instance of CategoryService
func NewCategoryService(
	categoryRepo repository.CategoryRepository,
	categoryRedisRepo repository.CategoryRedisRepository,
) service.CategoryService {
	return &categoryService{
		categoryRepo:      categoryRepo,
		categoryRedisRepo: categoryRedisRepo,
		cacheTTL:          10 * time.Minute, // 10 minutes cache TTL
	}
}

func (s *categoryService) GetCategoryList(ctx context.Context, limit, offset int) ([]*entity.Category, error) {
	// Try to get from cache first
	cachedCategories, err := s.categoryRedisRepo.GetCategoryList(ctx, limit, offset)
	if err == nil && cachedCategories != nil && len(cachedCategories) > 0 {
		log.Println("Category list retrieved from cache")
		return cachedCategories, nil
	}

	// If not in cache, get from database
	categories, err := s.categoryRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if len(categories) > 0 {
		err = s.categoryRedisRepo.SetCategoryList(ctx, categories, limit, offset)
		if err != nil {
			log.Printf("Failed to cache category list: %v", err)
		}
	}

	return categories, nil
}

func (s *categoryService) GetActiveCategoryList(ctx context.Context, limit, offset int) ([]*entity.Category, error) {
	// Try to get from cache first
	cachedCategories, err := s.categoryRedisRepo.GetActiveCategoryList(ctx, limit, offset)
	if err == nil && cachedCategories != nil && len(cachedCategories) > 0 {
		log.Println("Active category list retrieved from cache")
		return cachedCategories, nil
	}

	// If not in cache, get from database
	categories, err := s.categoryRepo.ListActive(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if len(categories) > 0 {
		err = s.categoryRedisRepo.SetActiveCategoryList(ctx, categories, limit, offset)
		if err != nil {
			log.Printf("Failed to cache active category list: %v", err)
		}
	}

	return categories, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id string) (*entity.Category, error) {
	// Try to get from cache first
	cachedCategory, err := s.categoryRedisRepo.GetCategoryByID(ctx, id)
	if err == nil && cachedCategory != nil {
		log.Println("Category retrieved from cache")
		return cachedCategory, nil
	}

	// If not in cache, get from database
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if category != nil {
		err = s.categoryRedisRepo.SetCategoryByID(ctx, category)
		if err != nil {
			log.Printf("Failed to cache category: %v", err)
		}
	}

	return category, nil
}

func (s *categoryService) GetCategoryByName(ctx context.Context, name string) (*entity.Category, error) {
	// For name-based queries, we'll get directly from database
	// as caching by name might not be as beneficial
	return s.categoryRepo.GetByName(ctx, name)
}

func (s *categoryService) CreateCategory(ctx context.Context, category *entity.Category) error {
	// Create category in database
	err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return err
	}

	// Invalidate cache after creation
	s.invalidateCache(ctx)
	return nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, category *entity.Category) error {
	// Update category in database
	err := s.categoryRepo.Update(ctx, category)
	if err != nil {
		return err
	}

	// Invalidate cache after update
	s.invalidateCache(ctx)
	return nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id string) error {
	// Delete category from database
	err := s.categoryRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache after deletion
	s.invalidateCache(ctx)
	return nil
}

func (s *categoryService) GetCategoryCount(ctx context.Context) (int64, error) {
	// Try to get count from cache first
	cachedCount, err := s.categoryRedisRepo.GetCategoryTotal(ctx)
	if err == nil && cachedCount > 0 {
		log.Println("Category count retrieved from cache")
		return cachedCount, nil
	}

	// If not in cache, get from database
	count, err := s.categoryRepo.Count(ctx)
	if err != nil {
		return 0, err
	}

	// Cache the count
	err = s.categoryRedisRepo.SetCategoryTotal(ctx, count)
	if err != nil {
		log.Printf("Failed to cache category count: %v", err)
	}

	return count, nil
}

func (s *categoryService) GetActiveCategoryCount(ctx context.Context) (int64, error) {
	// Try to get count from cache first
	cachedCount, err := s.categoryRedisRepo.GetActiveCategoryTotal(ctx)
	if err == nil && cachedCount > 0 {
		log.Println("Active category count retrieved from cache")
		return cachedCount, nil
	}

	// If not in cache, get from database
	count, err := s.categoryRepo.CountActive(ctx)
	if err != nil {
		return 0, err
	}

	// Cache the count
	err = s.categoryRedisRepo.SetActiveCategoryTotal(ctx, count)
	if err != nil {
		log.Printf("Failed to cache active category count: %v", err)
	}

	return count, nil
}

func (s *categoryService) GetCategoriesByParent(ctx context.Context, parentID string, limit, offset int) ([]*entity.Category, error) {
	// For parent-based queries, we'll get directly from database
	// as caching by parent might not be as beneficial
	return s.categoryRepo.ListByParent(ctx, parentID, limit, offset)
}

func (s *categoryService) GetCategoryCountByParent(ctx context.Context, parentID string) (int64, error) {
	// For parent-based count queries, we'll get directly from database
	return s.categoryRepo.CountByParent(ctx, parentID)
}

// invalidateCache clears all category-related cache
func (s *categoryService) invalidateCache(ctx context.Context) {
	err := s.categoryRedisRepo.InvalidateCategoryCache(ctx)
	if err != nil {
		log.Printf("Failed to invalidate category cache: %v", err)
	} else {
		log.Println("Category cache invalidated")
	}
}
