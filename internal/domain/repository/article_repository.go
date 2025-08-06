package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// ArticleRepository defines the interface for article data operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type ArticleRepository interface {
	Create(ctx context.Context, article *entity.Article) error
	GetByID(ctx context.Context, id string) (*entity.Article, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Article, error)
	Update(ctx context.Context, article *entity.Article) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Article, error)
	ListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Article, error)
	ListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Article, error)
	ListPublished(ctx context.Context, limit, offset int) ([]*entity.Article, error)
	Count(ctx context.Context) (int64, error)
	CountByAuthor(ctx context.Context, authorID string) (int64, error)
	CountByCategory(ctx context.Context, categoryID string) (int64, error)
	CountPublished(ctx context.Context) (int64, error)
}
