package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// InspirationRepository defines the interface for inspiration data operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type InspirationRepository interface {
	Create(ctx context.Context, inspiration *entity.Inspiration) error
	GetByID(ctx context.Context, id string) (*entity.Inspiration, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Inspiration, error)
	Update(ctx context.Context, inspiration *entity.Inspiration) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Inspiration, error)
	ListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Inspiration, error)
	ListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Inspiration, error)
	ListPublished(ctx context.Context, limit, offset int) ([]*entity.Inspiration, error)
	Count(ctx context.Context) (int64, error)
	CountByAuthor(ctx context.Context, authorID string) (int64, error)
	CountByCategory(ctx context.Context, categoryID string) (int64, error)
	CountPublished(ctx context.Context) (int64, error)
}
