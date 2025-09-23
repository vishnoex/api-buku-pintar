package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// SummaryRedisRepository defines the interface for summary Redis operations
type SummaryRedisRepository interface {
	GetSummaryByID(ctx context.Context, id string) (*entity.EbookSummary, error)
	SetSummaryByID(ctx context.Context, id string, summary *entity.EbookSummary) error
	GetSummariesList(ctx context.Context, limit, offset int) ([]*entity.EbookSummaryList, error)
	SetSummariesList(ctx context.Context, limit, offset int, summaries []*entity.EbookSummaryList) error
	GetSummariesByEbookID(ctx context.Context, ebookID string, limit, offset int) ([]*entity.EbookSummary, error)
	SetSummariesByEbookID(ctx context.Context, ebookID string, limit, offset int, summaries []*entity.EbookSummary) error
	GetSummariesCount(ctx context.Context) (int64, error)
	SetSummariesCount(ctx context.Context, count int64) error
	GetSummariesCountByEbookID(ctx context.Context, ebookID string) (int64, error)
	SetSummariesCountByEbookID(ctx context.Context, ebookID string, count int64) error
	ClearCache(ctx context.Context) error
}
