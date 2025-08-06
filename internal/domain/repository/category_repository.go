package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// CategoryRepository defines the interface for category data operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id string) (*entity.Category, error)
	GetByName(ctx context.Context, name string) (*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Category, error)
	ListActive(ctx context.Context, limit, offset int) ([]*entity.Category, error)
	ListByParent(ctx context.Context, parentID string, limit, offset int) ([]*entity.Category, error)
	Count(ctx context.Context) (int64, error)
	CountActive(ctx context.Context) (int64, error)
	CountByParent(ctx context.Context, parentID string) (int64, error)
}
