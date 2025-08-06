package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// ContentStatusRepository defines the interface for content status data operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type ContentStatusRepository interface {
	Create(ctx context.Context, contentStatus *entity.ContentStatus) error
	GetByID(ctx context.Context, id string) (*entity.ContentStatus, error)
	GetByName(ctx context.Context, name string) (*entity.ContentStatus, error)
	Update(ctx context.Context, contentStatus *entity.ContentStatus) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.ContentStatus, error)
	Count(ctx context.Context) (int64, error)
}
