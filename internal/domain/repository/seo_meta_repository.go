package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// SeoMetaRepository defines the interface for SEO metadata operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type SeoMetaRepository interface {
	Create(ctx context.Context, seoMeta *entity.SeoMeta) error
	GetByID(ctx context.Context, id string) (*entity.SeoMeta, error)
	GetByEntity(ctx context.Context, entity entity.SeoEntityType, entityID string) (*entity.SeoMeta, error)
	Update(ctx context.Context, seoMeta *entity.SeoMeta) error
	Delete(ctx context.Context, id string) error
	ListByEntity(ctx context.Context, entity entity.SeoEntityType, limit, offset int) ([]*entity.SeoMeta, error)
	CountByEntity(ctx context.Context, entity entity.SeoEntityType) (int64, error)
}
