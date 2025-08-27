package service

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/internal/domain/service"
	"context"
	"log"
	"time"
)

type bannerService struct {
	bannerRepo     repository.BannerRepository
	bannerRedisRepo repository.BannerRedisRepository
	cacheTTL       time.Duration
}

// NewBannerService creates a new instance of BannerService
func NewBannerService(
	bannerRepo repository.BannerRepository,
	bannerRedisRepo repository.BannerRedisRepository,
) service.BannerService {
	return &bannerService{
		bannerRepo:      bannerRepo,
		bannerRedisRepo: bannerRedisRepo,
		cacheTTL:        5 * time.Minute, // 5 minutes cache TTL
	}
}

func (s *bannerService) GetBannerList(ctx context.Context, limit, offset int) ([]*entity.Banner, error) {
	// Try to get from cache first
	cachedBanners, err := s.bannerRedisRepo.GetBannerList(ctx, limit, offset)
	if err == nil && cachedBanners != nil && len(cachedBanners) > 0 {
		log.Println("Banner list retrieved from cache")
		return cachedBanners, nil
	}

	// If not in cache, get from database
	banners, err := s.bannerRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if len(banners) > 0 {
		err = s.bannerRedisRepo.SetBannerList(ctx, banners, limit, offset)
		if err != nil {
			log.Printf("Failed to cache banner list: %v", err)
		}
	}

	return banners, nil
}

func (s *bannerService) GetBannerByID(ctx context.Context, id string) (*entity.Banner, error) {
	// For individual banner, we'll get directly from database
	// as caching individual banners might not be as beneficial
	return s.bannerRepo.GetByID(ctx, id)
}

func (s *bannerService) CreateBanner(ctx context.Context, banner *entity.Banner) error {
	// Create banner in database
	err := s.bannerRepo.Create(ctx, banner)
	if err != nil {
		return err
	}

	// Invalidate cache after creation
	s.invalidateCache(ctx)
	return nil
}

func (s *bannerService) UpdateBanner(ctx context.Context, banner *entity.Banner) error {
	// Update banner in database
	err := s.bannerRepo.Update(ctx, banner)
	if err != nil {
		return err
	}

	// Invalidate cache after update
	s.invalidateCache(ctx)
	return nil
}

func (s *bannerService) DeleteBanner(ctx context.Context, id string) error {
	// Delete banner from database
	err := s.bannerRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache after deletion
	s.invalidateCache(ctx)
	return nil
}

func (s *bannerService) GetBannerCount(ctx context.Context) (int64, error) {
	// Try to get count from cache first
	cachedCount, err := s.bannerRedisRepo.GetBannerTotal(ctx)
	if err == nil && cachedCount > 0 {
		log.Println("Banner count retrieved from cache")
		return cachedCount, nil
	}

	// If not in cache, get from database
	count, err := s.bannerRepo.Count(ctx)
	if err != nil {
		return 0, err
	}

	// Cache the count
	err = s.bannerRedisRepo.SetBannerTotal(ctx, count)
	if err != nil {
		log.Printf("Failed to cache banner count: %v", err)
	}

	return count, nil
}

func (s *bannerService) GetActiveBannerList(ctx context.Context, limit, offset int) ([]*entity.Banner, error) {
	// For active banners, we'll get directly from database
	// as they might change frequently and caching might not be beneficial
	return s.bannerRepo.ListActive(ctx, limit, offset)
}

func (s *bannerService) GetActiveBannerCount(ctx context.Context) (int64, error) {
	// For active banner count, we'll get directly from database
	return s.bannerRepo.CountActive(ctx)
}

// invalidateCache clears all banner-related cache
func (s *bannerService) invalidateCache(ctx context.Context) {
	// Clear banner count cache
	err := s.bannerRedisRepo.SetBannerTotal(ctx, 0)
	if err != nil {
		log.Printf("Failed to invalidate banner count cache: %v", err)
	}

	// Note: In a more sophisticated implementation, you might want to
	// implement cache invalidation for specific keys or use cache tags
	log.Println("Banner cache invalidated")
}
