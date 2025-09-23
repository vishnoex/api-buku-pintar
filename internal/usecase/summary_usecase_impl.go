package usecase

import (
	"buku-pintar/internal/constant"
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"context"
)

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type SummaryUsecaseImpl struct {
	summaryService service.SummaryService
}

func NewSummaryUsecaseImpl(summaryService service.SummaryService) SummaryUsecase {
	return &SummaryUsecaseImpl{
		summaryService: summaryService,
	}
}

func (u *SummaryUsecaseImpl) CreateSummary(ctx context.Context, summary *entity.EbookSummary) error {
	// Validate required fields
	if summary.EbookID == "" {
		return &ValidationError{Message: constant.EBOOK_ID_REQUIRED_VALIDATION}
	}
	if summary.Description == "" {
		return &ValidationError{Message: constant.DESCRIPTION_REQUIRED}
	}

	// Create summary through service
	return u.summaryService.CreateSummary(ctx, summary)
}

func (u *SummaryUsecaseImpl) GetSummaryByID(ctx context.Context, id string) (*entity.EbookSummary, error) {
	if id == "" {
		return nil, &ValidationError{Message: constant.ERR_ID_REQUIRED}
	}

	return u.summaryService.GetSummaryByID(ctx, id)
}

func (u *SummaryUsecaseImpl) UpdateSummary(ctx context.Context, summary *entity.EbookSummary) error {
	if summary.ID == "" {
		return &ValidationError{Message: constant.ERR_ID_REQUIRED}
	}

	return u.summaryService.UpdateSummary(ctx, summary)
}

func (u *SummaryUsecaseImpl) DeleteSummary(ctx context.Context, id string) error {
	if id == "" {
		return &ValidationError{Message: constant.ERR_ID_REQUIRED}
	}

	return u.summaryService.DeleteSummary(ctx, id)
}

func (u *SummaryUsecaseImpl) ListSummaries(ctx context.Context, limit, offset int) ([]*response.EbookSummaryResponse, error) {
	return u.summaryService.ListSummaries(ctx, limit, offset)
}

func (u *SummaryUsecaseImpl) GetSummariesByEbookID(ctx context.Context, ebookID string, limit, offset int) ([]*entity.EbookSummary, error) {
	if ebookID == "" {
		return nil, &ValidationError{Message: constant.EBOOK_ID_REQUIRED_VALIDATION}
	}

	return u.summaryService.GetSummariesByEbookID(ctx, ebookID, limit, offset)
}

func (u *SummaryUsecaseImpl) CountSummaries(ctx context.Context) (int64, error) {
	return u.summaryService.CountSummaries(ctx)
}

func (u *SummaryUsecaseImpl) CountSummariesByEbookID(ctx context.Context, ebookID string) (int64, error) {
	if ebookID == "" {
		return 0, &ValidationError{Message: constant.EBOOK_ID_REQUIRED_VALIDATION}
	}

	return u.summaryService.CountSummariesByEbookID(ctx, ebookID)
}
