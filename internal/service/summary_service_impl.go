package service

import (
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/internal/domain/service"
	"context"
	"fmt"
)

type SummaryServiceImpl struct {
	summaryRepository repository.SummaryRepository
	summaryRedisRepository repository.SummaryRedisRepository
}

func NewSummaryServiceImpl(
	summaryRepository repository.SummaryRepository,
	summaryRedisRepository repository.SummaryRedisRepository,
) service.SummaryService {
	return &SummaryServiceImpl{
		summaryRepository: summaryRepository,
		summaryRedisRepository: summaryRedisRepository,
	}
}

func (s *SummaryServiceImpl) CreateSummary(ctx context.Context, summary *entity.EbookSummary) error {
	// Create summary in database
	err := s.summaryRepository.CreateSummary(ctx, summary)
	if err != nil {
		return err
	}

	// Clear related caches
	_ = s.summaryRedisRepository.ClearCache(ctx)

	return nil
}

func (s *SummaryServiceImpl) GetSummaryByID(ctx context.Context, id string) (*entity.EbookSummary, error) {
	// Try to get from cache first
	summary, err := s.summaryRedisRepository.GetSummaryByID(ctx, id)
	if err == nil && summary != nil {
		return summary, nil
	}

	// If not in cache, get from database
	summary, err = s.summaryRepository.GetSummaryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if summary != nil {
		_ = s.summaryRedisRepository.SetSummaryByID(ctx, id, summary)
	}

	return summary, nil
}

func (s *SummaryServiceImpl) UpdateSummary(ctx context.Context, summary *entity.EbookSummary) error {
	// Update in database
	err := s.summaryRepository.UpdateSummary(ctx, summary)
	if err != nil {
		return err
	}

	// Clear related caches
	_ = s.summaryRedisRepository.ClearCache(ctx)

	return nil
}

func (s *SummaryServiceImpl) DeleteSummary(ctx context.Context, id string) error {
	// Delete from database
	err := s.summaryRepository.DeleteSummary(ctx, id)
	if err != nil {
		return err
	}

	// Clear related caches
	_ = s.summaryRedisRepository.ClearCache(ctx)

	return nil
}

func (s *SummaryServiceImpl) ListSummaries(ctx context.Context, limit, offset int) ([]*response.EbookSummaryResponse, error) {
	// Try to get from cache first
	summaries, err := s.summaryRedisRepository.GetSummariesList(ctx, limit, offset)
	if err == nil && summaries != nil {
		// Convert entity to response
		return s.convertToSummaryResponses(summaries), nil
	}

	// If not in cache, get from database
	summaries, err = s.summaryRepository.ListSummaries(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if summaries != nil {
		_ = s.summaryRedisRepository.SetSummariesList(ctx, limit, offset, summaries)
	}

	// Convert entity to response
	return s.convertToSummaryResponses(summaries), nil
}

func (s *SummaryServiceImpl) GetSummariesByEbookID(ctx context.Context, ebookID string, limit, offset int) ([]*entity.EbookSummary, error) {
	// Try to get from cache first
	summaries, err := s.summaryRedisRepository.GetSummariesByEbookID(ctx, ebookID, limit, offset)
	if err == nil && summaries != nil {
		return summaries, nil
	}

	// If not in cache, get from database
	summaries, err = s.summaryRepository.GetSummariesByEbookID(ctx, ebookID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if summaries != nil {
		_ = s.summaryRedisRepository.SetSummariesByEbookID(ctx, ebookID, limit, offset, summaries)
	}

	return summaries, nil
}

func (s *SummaryServiceImpl) CountSummaries(ctx context.Context) (int64, error) {
	// Try to get from cache first
	count, err := s.summaryRedisRepository.GetSummariesCount(ctx)
	if err == nil {
		return count, nil
	}

	// If not in cache, get from database
	count, err = s.summaryRepository.CountSummaries(ctx)
	if err != nil {
		return 0, err
	}

	// Cache the result
	_ = s.summaryRedisRepository.SetSummariesCount(ctx, count)

	return count, nil
}

func (s *SummaryServiceImpl) CountSummariesByEbookID(ctx context.Context, ebookID string) (int64, error) {
	// Try to get from cache first
	count, err := s.summaryRedisRepository.GetSummariesCountByEbookID(ctx, ebookID)
	if err == nil {
		return count, nil
	}

	// If not in cache, get from database
	count, err = s.summaryRepository.CountSummariesByEbookID(ctx, ebookID)
	if err != nil {
		return 0, err
	}

	// Cache the result
	_ = s.summaryRedisRepository.SetSummariesCountByEbookID(ctx, ebookID, count)

	return count, nil
}

// convertToSummaryResponses converts entity summaries to response summaries
func (s *SummaryServiceImpl) convertToSummaryResponses(summaries []*entity.EbookSummaryList) []*response.EbookSummaryResponse {
	var responses []*response.EbookSummaryResponse
	for _, summary := range summaries {
		responses = append(responses, &response.EbookSummaryResponse{
			ID:          summary.ID,
			EbookID:     summary.EbookID,
			EbookTitle:  summary.EbookTitle,
			Slug:        summary.Slug,
			Description: summary.Description,
			URL:         summary.URL,
			AudioURL:    summary.AudioURL,
			Duration:    fmt.Sprintf("%d minutes", summary.Duration),
			// Note: EbookTitle, Slug, and Duration would need to be fetched from the ebook entity
			// For now, we'll leave them empty as they require additional database queries
		})
	}
	return responses
}
