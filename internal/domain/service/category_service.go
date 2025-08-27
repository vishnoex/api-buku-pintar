package service

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// CategoryService defines the interface for category business operations
type CategoryService interface {
	GetCategoryList(ctx context.Context, limit, offset int) ([]*entity.Category, error)
	GetActiveCategoryList(ctx context.Context, limit, offset int) ([]*entity.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*entity.Category, error)
	GetCategoryByName(ctx context.Context, name string) (*entity.Category, error)
	CreateCategory(ctx context.Context, category *entity.Category) error
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategoryCount(ctx context.Context) (int64, error)
	GetActiveCategoryCount(ctx context.Context) (int64, error)
	GetCategoriesByParent(ctx context.Context, parentID string, limit, offset int) ([]*entity.Category, error)
	GetCategoryCountByParent(ctx context.Context, parentID string) (int64, error)
}
