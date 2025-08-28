package service

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/internal/domain/service"
	"context"
	"log"
	"time"
)

type ebookService struct {
	ebookRepo      repository.EbookRepository
	ebookRedisRepo repository.EbookRedisRepository
	cacheTTL       time.Duration
}

// NewEbookService creates a new instance of EbookService
func NewEbookService(
	ebookRepo repository.EbookRepository,
	ebookRedisRepo repository.EbookRedisRepository,
) service.EbookService {
	return &ebookService{
		ebookRepo:      ebookRepo,
		ebookRedisRepo: ebookRedisRepo,
		cacheTTL:       15 * time.Minute, // 15 minutes cache TTL
	}
}

func (s *ebookService) CreateEbook(ctx context.Context, ebook *entity.Ebook) error {
	// Create ebook in database
	err := s.ebookRepo.Create(ctx, ebook)
	if err != nil {
		return err
	}

	// Invalidate cache after creation
	s.invalidateCache(ctx)
	return nil
}

func (s *ebookService) GetEbookByID(ctx context.Context, id string) (*entity.Ebook, error) {
	// Try to get from cache first
	cachedEbook, err := s.ebookRedisRepo.GetEbookByID(ctx, id)
	if err == nil && cachedEbook != nil {
		log.Println("Ebook retrieved from cache by ID")
		return cachedEbook, nil
	}

	// If not in cache, get from database
	ebook, err := s.ebookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if ebook != nil {
		err = s.ebookRedisRepo.SetEbookByID(ctx, ebook)
		if err != nil {
			log.Printf("Failed to cache ebook by ID: %v", err)
		}
	}

	return ebook, nil
}

func (s *ebookService) GetEbookBySlug(ctx context.Context, slug string) (*entity.Ebook, error) {
	// Try to get from cache first
	cachedEbook, err := s.ebookRedisRepo.GetEbookBySlug(ctx, slug)
	if err == nil && cachedEbook != nil {
		log.Println("Ebook retrieved from cache by slug")
		return cachedEbook, nil
	}

	// If not in cache, get from database
	ebook, err := s.ebookRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if ebook != nil {
		err = s.ebookRedisRepo.SetEbookBySlug(ctx, ebook)
		if err != nil {
			log.Printf("Failed to cache ebook by slug: %v", err)
		}
	}

	return ebook, nil
}

func (s *ebookService) UpdateEbook(ctx context.Context, ebook *entity.Ebook) error {
	// Update ebook in database
	err := s.ebookRepo.Update(ctx, ebook)
	if err != nil {
		return err
	}

	// Invalidate cache after update
	s.invalidateCache(ctx)
	return nil
}

func (s *ebookService) DeleteEbook(ctx context.Context, id string) error {
	// Delete ebook from database
	err := s.ebookRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache after deletion
	s.invalidateCache(ctx)
	return nil
}

func (s *ebookService) GetEbookList(ctx context.Context, limit, offset int) ([]*entity.EbookList, error) {
	// Try to get from cache first
	cachedEbooks, err := s.ebookRedisRepo.GetEbookList(ctx, limit, offset)
	if err == nil && cachedEbooks != nil && len(cachedEbooks) > 0 {
		log.Println("Ebook list retrieved from cache")
		return cachedEbooks, nil
	}

	// If not in cache, get from database
	ebooks, err := s.ebookRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if len(ebooks) > 0 {
		err = s.ebookRedisRepo.SetEbookList(ctx, ebooks, limit, offset)
		if err != nil {
			log.Printf("Failed to cache ebook list: %v", err)
		}
	}

	return ebooks, nil
}

func (s *ebookService) GetEbookListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Ebook, error) {
	// Try to get from cache first
	cachedEbooks, err := s.ebookRedisRepo.GetEbookListByCategory(ctx, categoryID, limit, offset)
	if err == nil && cachedEbooks != nil && len(cachedEbooks) > 0 {
		log.Println("Ebook list by category retrieved from cache")
		return cachedEbooks, nil
	}

	// If not in cache, get from database
	ebooks, err := s.ebookRepo.ListByCategory(ctx, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if len(ebooks) > 0 {
		err = s.ebookRedisRepo.SetEbookListByCategory(ctx, ebooks, categoryID, limit, offset)
		if err != nil {
			log.Printf("Failed to cache ebook list by category: %v", err)
		}
	}

	return ebooks, nil
}

func (s *ebookService) GetEbookListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Ebook, error) {
	// Try to get from cache first
	cachedEbooks, err := s.ebookRedisRepo.GetEbookListByAuthor(ctx, authorID, limit, offset)
	if err == nil && cachedEbooks != nil && len(cachedEbooks) > 0 {
		log.Println("Ebook list by author retrieved from cache")
		return cachedEbooks, nil
	}

	// If not in cache, get from database
	ebooks, err := s.ebookRepo.ListByAuthor(ctx, authorID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if len(ebooks) > 0 {
		err = s.ebookRedisRepo.SetEbookListByAuthor(ctx, ebooks, authorID, limit, offset)
		if err != nil {
			log.Printf("Failed to cache ebook list by author: %v", err)
		}
	}

	return ebooks, nil
}

func (s *ebookService) GetEbookCount(ctx context.Context) (int64, error) {
	// Try to get count from cache first
	cachedCount, err := s.ebookRedisRepo.GetEbookTotal(ctx)
	if err == nil && cachedCount > 0 {
		log.Println("Ebook count retrieved from cache")
		return cachedCount, nil
	}

	// If not in cache, get from database
	count, err := s.ebookRepo.Count(ctx)
	if err != nil {
		return 0, err
	}

	// Cache the count
	err = s.ebookRedisRepo.SetEbookTotal(ctx, count)
	if err != nil {
		log.Printf("Failed to cache ebook count: %v", err)
	}

	return count, nil
}

func (s *ebookService) GetEbookCountByCategory(ctx context.Context, categoryID string) (int64, error) {
	// Try to get count from cache first
	cachedCount, err := s.ebookRedisRepo.GetEbookCountByCategory(ctx, categoryID)
	if err == nil && cachedCount > 0 {
		log.Println("Ebook count by category retrieved from cache")
		return cachedCount, nil
	}

	// If not in cache, get from database
	count, err := s.ebookRepo.CountByCategory(ctx, categoryID)
	if err != nil {
		return 0, err
	}

	// Cache the count
	err = s.ebookRedisRepo.SetEbookCountByCategory(ctx, categoryID, count)
	if err != nil {
		log.Printf("Failed to cache ebook count by category: %v", err)
	}

	return count, nil
}

func (s *ebookService) GetEbookCountByAuthor(ctx context.Context, authorID string) (int64, error) {
	// Try to get count from cache first
	cachedCount, err := s.ebookRedisRepo.GetEbookCountByAuthor(ctx, authorID)
	if err == nil && cachedCount > 0 {
		log.Println("Ebook count by author retrieved from cache")
		return cachedCount, nil
	}

	// If not in cache, get from database
	count, err := s.ebookRepo.CountByAuthor(ctx, authorID)
	if err != nil {
		return 0, err
	}

	// Cache the count
	err = s.ebookRedisRepo.SetEbookCountByAuthor(ctx, authorID, count)
	if err != nil {
		log.Printf("Failed to cache ebook count by author: %v", err)
	}

	return count, nil
}

// invalidateCache clears all ebook-related cache
func (s *ebookService) invalidateCache(ctx context.Context) {
	err := s.ebookRedisRepo.InvalidateEbookCache(ctx)
	if err != nil {
		log.Printf("Failed to invalidate ebook cache: %v", err)
	} else {
		log.Println("Ebook cache invalidated")
	}
}
