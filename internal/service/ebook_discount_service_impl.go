package service

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/internal/domain/service"
	"context"
	"log"
	"time"
)

type ebookDiscountService struct {
	discountRepo      repository.EbookDiscountRepository
	discountRedisRepo repository.EbookDiscountRedisRepository
	cacheTTL          time.Duration
}

// NewEbookDiscountService creates a new instance of EbookDiscountService
func NewEbookDiscountService(
	discountRepo repository.EbookDiscountRepository,
	discountRedisRepo repository.EbookDiscountRedisRepository,
) service.EbookDiscountService {
	return &ebookDiscountService{
		discountRepo:      discountRepo,
		discountRedisRepo: discountRedisRepo,
		cacheTTL:          10 * time.Minute, // 10 minutes cache TTL
	}
}

func (s *ebookDiscountService) CreateDiscount(ctx context.Context, discount *entity.EbookDiscount) error {
	// Create discount in database
	err := s.discountRepo.Create(ctx, discount)
	if err != nil {
		return err
	}

	// Invalidate cache after creation
	s.invalidateCache(ctx)
	return nil
}

func (s *ebookDiscountService) GetDiscountByID(ctx context.Context, id string) (*entity.EbookDiscount, error) {
	// Try to get from cache first
	cachedDiscount, err := s.discountRedisRepo.GetDiscountByID(ctx, id)
	if err == nil && cachedDiscount != nil {
		log.Println("Discount retrieved from cache by ID")
		return cachedDiscount, nil
	}

	// If not in cache, get from database
	discount, err := s.discountRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if discount != nil {
		err = s.discountRedisRepo.SetDiscountByID(ctx, discount)
		if err != nil {
			log.Printf("Failed to cache discount by ID: %v", err)
		}
	}

	return discount, nil
}

func (s *ebookDiscountService) UpdateDiscount(ctx context.Context, discount *entity.EbookDiscount) error {
	// Update discount in database
	err := s.discountRepo.Update(ctx, discount)
	if err != nil {
		return err
	}

	// Invalidate cache after update
	s.invalidateCache(ctx)
	return nil
}

func (s *ebookDiscountService) DeleteDiscount(ctx context.Context, id string) error {
	// Delete discount from database
	err := s.discountRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache after deletion
	s.invalidateCache(ctx)
	return nil
}

func (s *ebookDiscountService) GetDiscountList(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error) {
	// Try to get from cache first
	cachedDiscounts, err := s.discountRedisRepo.GetDiscountList(ctx, limit, offset)
	if err == nil && cachedDiscounts != nil && len(cachedDiscounts) > 0 {
		log.Println("Discount list retrieved from cache")
		return cachedDiscounts, nil
	}

	// If not in cache, get from database
	discounts, err := s.discountRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if len(discounts) > 0 {
		err = s.discountRedisRepo.SetDiscountList(ctx, discounts, limit, offset)
		if err != nil {
			log.Printf("Failed to cache discount list: %v", err)
		}
	}

	return discounts, nil
}

func (s *ebookDiscountService) GetDiscountsByEbookID(ctx context.Context, ebookID string) ([]*entity.EbookDiscount, error) {
	// Try to get from cache first
	cachedDiscounts, err := s.discountRedisRepo.GetDiscountsByEbookID(ctx, ebookID)
	if err == nil && cachedDiscounts != nil && len(cachedDiscounts) > 0 {
		log.Println("Discounts by ebook ID retrieved from cache")
		return cachedDiscounts, nil
	}

	// If not in cache, get from database
	discounts, err := s.discountRepo.GetByEbookID(ctx, ebookID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if len(discounts) > 0 {
		err = s.discountRedisRepo.SetDiscountsByEbookID(ctx, ebookID, discounts)
		if err != nil {
			log.Printf("Failed to cache discounts by ebook ID: %v", err)
		}
	}

	return discounts, nil
}

func (s *ebookDiscountService) GetActiveDiscounts(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error) {
	// Try to get from cache first
	cachedDiscounts, err := s.discountRedisRepo.GetActiveDiscounts(ctx, limit, offset)
	if err == nil && cachedDiscounts != nil && len(cachedDiscounts) > 0 {
		log.Println("Active discounts retrieved from cache")
		return cachedDiscounts, nil
	}

	// If not in cache, get from database
	discounts, err := s.discountRepo.GetActiveDiscounts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if len(discounts) > 0 {
		err = s.discountRedisRepo.SetActiveDiscounts(ctx, discounts, limit, offset)
		if err != nil {
			log.Printf("Failed to cache active discounts: %v", err)
		}
	}

	return discounts, nil
}

func (s *ebookDiscountService) GetActiveDiscountByEbookID(ctx context.Context, ebookID string) (*entity.EbookDiscount, error) {
	// Try to get from cache first
	cachedDiscount, err := s.discountRedisRepo.GetActiveDiscountByEbookID(ctx, ebookID)
	if err == nil && cachedDiscount != nil {
		log.Println("Active discount by ebook ID retrieved from cache")
		return cachedDiscount, nil
	}

	// If not in cache, get from database
	discount, err := s.discountRepo.GetActiveDiscountByEbookID(ctx, ebookID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if discount != nil {
		err = s.discountRedisRepo.SetActiveDiscountByEbookID(ctx, ebookID, discount)
		if err != nil {
			log.Printf("Failed to cache active discount by ebook ID: %v", err)
		}
	}

	return discount, nil
}

func (s *ebookDiscountService) GetDiscountCount(ctx context.Context) (int64, error) {
	// Try to get count from cache first
	cachedCount, err := s.discountRedisRepo.GetDiscountCount(ctx)
	if err == nil && cachedCount > 0 {
		log.Println("Discount count retrieved from cache")
		return cachedCount, nil
	}

	// If not in cache, get from database
	count, err := s.discountRepo.Count(ctx)
	if err != nil {
		return 0, err
	}

	// Cache the count
	err = s.discountRedisRepo.SetDiscountCount(ctx, count)
	if err != nil {
		log.Printf("Failed to cache discount count: %v", err)
	}

	return count, nil
}

func (s *ebookDiscountService) GetDiscountCountByEbookID(ctx context.Context, ebookID string) (int64, error) {
	// Try to get count from cache first
	cachedCount, err := s.discountRedisRepo.GetDiscountCountByEbookID(ctx, ebookID)
	if err == nil && cachedCount > 0 {
		log.Println("Discount count by ebook ID retrieved from cache")
		return cachedCount, nil
	}

	// If not in cache, get from database
	count, err := s.discountRepo.CountByEbookID(ctx, ebookID)
	if err != nil {
		return 0, err
	}

	// Cache the count
	err = s.discountRedisRepo.SetDiscountCountByEbookID(ctx, ebookID, count)
	if err != nil {
		log.Printf("Failed to cache discount count by ebook ID: %v", err)
	}

	return count, nil
}

func (s *ebookDiscountService) GetActiveDiscountCount(ctx context.Context) (int64, error) {
	// Try to get count from cache first
	cachedCount, err := s.discountRedisRepo.GetActiveDiscountCount(ctx)
	if err == nil && cachedCount > 0 {
		log.Println("Active discount count retrieved from cache")
		return cachedCount, nil
	}

	// If not in cache, get from database
	count, err := s.discountRepo.CountActiveDiscounts(ctx)
	if err != nil {
		return 0, err
	}

	// Cache the count
	err = s.discountRedisRepo.SetActiveDiscountCount(ctx, count)
	if err != nil {
		log.Printf("Failed to cache active discount count: %v", err)
	}

	return count, nil
}

// invalidateCache clears all discount-related cache
func (s *ebookDiscountService) invalidateCache(ctx context.Context) {
	err := s.discountRedisRepo.InvalidateDiscountCache(ctx)
	if err != nil {
		log.Printf("Failed to invalidate discount cache: %v", err)
	} else {
		log.Println("Discount cache invalidated")
	}
}
