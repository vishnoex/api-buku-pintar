package service

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// EbookService defines the interface for ebook business operations
type EbookService interface {
	CreateEbook(ctx context.Context, ebook *entity.Ebook) error
	GetEbookByID(ctx context.Context, id string) (*entity.Ebook, error)
	GetEbookBySlug(ctx context.Context, slug string) (*entity.Ebook, error)
	UpdateEbook(ctx context.Context, ebook *entity.Ebook) error
	DeleteEbook(ctx context.Context, id string) error
	GetEbookList(ctx context.Context, limit, offset int) ([]*entity.EbookList, error)
	GetEbookListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Ebook, error)
	GetEbookListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Ebook, error)
	GetEbookCount(ctx context.Context) (int64, error)
	GetEbookCountByCategory(ctx context.Context, categoryID string) (int64, error)
	GetEbookCountByAuthor(ctx context.Context, authorID string) (int64, error)
}
