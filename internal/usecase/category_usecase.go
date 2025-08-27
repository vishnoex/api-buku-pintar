package usecase

import (
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/entity"
	"context"
)

type CategoryUsecase interface {
	ListCategory(ctx context.Context, limit, offset int) ([]*response.CategoryResponse, error)
	GetCategoryByID(ctx context.Context, id string) (*response.CategoryResponse, error)
	CreateCategory(ctx context.Context, category *entity.Category) error
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id string) error
	CountCategory(ctx context.Context) (int64, error)
	ListAllCategories(ctx context.Context, limit, offset int) ([]*response.CategoryResponse, error)
	CountAllCategories(ctx context.Context) (int64, error)
	ListCategoriesByParent(ctx context.Context, parentID string, limit, offset int) ([]*response.CategoryResponse, error)
	CountCategoriesByParent(ctx context.Context, parentID string) (int64, error)
}
