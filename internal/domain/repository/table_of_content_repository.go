package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// TableOfContentRepository defines the interface for table of content operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type TableOfContentRepository interface {
	Create(ctx context.Context, toc *entity.TableOfContent) error
	GetByID(ctx context.Context, id string) (*entity.TableOfContent, error)
	Update(ctx context.Context, toc *entity.TableOfContent) error
	Delete(ctx context.Context, id string) error
	ListByEbook(ctx context.Context, ebookID string) ([]*entity.TableOfContent, error)
	DeleteByEbook(ctx context.Context, ebookID string) error
	CountByEbook(ctx context.Context, ebookID string) (int64, error)
}
