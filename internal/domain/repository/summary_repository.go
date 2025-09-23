package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// SummaryRepository defines the interface for summary data access
type SummaryRepository interface {
	CreateSummary(ctx context.Context, summary *entity.EbookSummary) error
	GetSummaryByID(ctx context.Context, id string) (*entity.EbookSummary, error)
	UpdateSummary(ctx context.Context, summary *entity.EbookSummary) error
	DeleteSummary(ctx context.Context, id string) error
	ListSummaries(ctx context.Context, limit, offset int) ([]*entity.EbookSummaryList, error)
	GetSummariesByEbookID(ctx context.Context, ebookID string, limit, offset int) ([]*entity.EbookSummary, error)
	CountSummaries(ctx context.Context) (int64, error)
	CountSummariesByEbookID(ctx context.Context, ebookID string) (int64, error)
}
