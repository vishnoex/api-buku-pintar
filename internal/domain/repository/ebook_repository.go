package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// EbookRepository defines the interface for ebook data operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type EbookRepository interface {
	Create(ctx context.Context, ebook *entity.Ebook) error
	GetByID(ctx context.Context, id string) (*entity.Ebook, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Ebook, error)
	Update(ctx context.Context, ebook *entity.Ebook) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.EbookList, error)
	ListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Ebook, error)
	ListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Ebook, error)
	Count(ctx context.Context) (int64, error)
	CountByCategory(ctx context.Context, categoryID string) (int64, error)
	CountByAuthor(ctx context.Context, authorID string) (int64, error)
}

// EbookRedisRepository defines the interface for ebook Redis operations
type EbookRedisRepository interface {
	GetEbookList(ctx context.Context, limit, offset int) ([]*entity.EbookList, error)
	SetEbookList(ctx context.Context, ebooks []*entity.EbookList, limit, offset int) error
	GetEbookTotal(ctx context.Context) (int64, error)
	SetEbookTotal(ctx context.Context, count int64) error
	GetEbookByID(ctx context.Context, id string) (*entity.Ebook, error)
	SetEbookByID(ctx context.Context, ebook *entity.Ebook) error
	GetEbookBySlug(ctx context.Context, slug string) (*entity.Ebook, error)
	SetEbookBySlug(ctx context.Context, ebook *entity.Ebook) error
	GetEbookListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Ebook, error)
	SetEbookListByCategory(ctx context.Context, ebooks []*entity.Ebook, categoryID string, limit, offset int) error
	GetEbookCountByCategory(ctx context.Context, categoryID string) (int64, error)
	SetEbookCountByCategory(ctx context.Context, categoryID string, count int64) error
	GetEbookListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Ebook, error)
	SetEbookListByAuthor(ctx context.Context, ebooks []*entity.Ebook, authorID string, limit, offset int) error
	GetEbookCountByAuthor(ctx context.Context, authorID string) (int64, error)
	SetEbookCountByAuthor(ctx context.Context, authorID string, count int64) error
	InvalidateEbookCache(ctx context.Context) error
}
