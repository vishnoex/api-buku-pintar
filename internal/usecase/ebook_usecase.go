package usecase

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// EbookUsecase defines the interface for ebook use cases
type EbookUsecase interface {
	CreateEbook(ctx context.Context, ebook *entity.Ebook) error
	GetEbookByID(ctx context.Context, id string) (*entity.Ebook, error)
	GetEbookBySlug(ctx context.Context, slug string) (*entity.Ebook, error)
	UpdateEbook(ctx context.Context, ebook *entity.Ebook) error
	DeleteEbook(ctx context.Context, id string) error
	ListEbooks(ctx context.Context, limit, offset int) ([]*entity.Ebook, error)
	ListEbooksByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Ebook, error)
	ListEbooksByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Ebook, error)
	CountEbooks(ctx context.Context) (int64, error)
	CountEbooksByCategory(ctx context.Context, categoryID string) (int64, error)
	CountEbooksByAuthor(ctx context.Context, authorID string) (int64, error)
}
