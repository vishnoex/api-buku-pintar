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
