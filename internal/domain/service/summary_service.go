package service

import (
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/entity"
	"context"
)

// SummaryService defines the interface for summary business logic
type SummaryService interface {
	CreateSummary(ctx context.Context, summary *entity.EbookSummary) error
	GetSummaryByID(ctx context.Context, id string) (*entity.EbookSummary, error)
	UpdateSummary(ctx context.Context, summary *entity.EbookSummary) error
	DeleteSummary(ctx context.Context, id string) error
	ListSummaries(ctx context.Context, limit, offset int) ([]*response.EbookSummaryResponse, error)
	GetSummariesByEbookID(ctx context.Context, ebookID string, limit, offset int) ([]*entity.EbookSummary, error)
	CountSummaries(ctx context.Context) (int64, error)
	CountSummariesByEbookID(ctx context.Context, ebookID string) (int64, error)
}
