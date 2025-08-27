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

// CategoryRedisRepository defines the interface for category Redis operations
type CategoryRedisRepository interface {
	GetCategoryList(ctx context.Context, limit, offset int) ([]*entity.Category, error)
	SetCategoryList(ctx context.Context, categories []*entity.Category, limit, offset int) error
	GetCategoryTotal(ctx context.Context) (int64, error)
	SetCategoryTotal(ctx context.Context, count int64) error
	GetActiveCategoryList(ctx context.Context, limit, offset int) ([]*entity.Category, error)
	SetActiveCategoryList(ctx context.Context, categories []*entity.Category, limit, offset int) error
	GetActiveCategoryTotal(ctx context.Context) (int64, error)
	SetActiveCategoryTotal(ctx context.Context, count int64) error
	GetCategoryByID(ctx context.Context, id string) (*entity.Category, error)
	SetCategoryByID(ctx context.Context, category *entity.Category) error
	InvalidateCategoryCache(ctx context.Context) error
}
