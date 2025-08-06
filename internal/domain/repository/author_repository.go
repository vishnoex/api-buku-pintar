package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// AuthorRepository defines the interface for author data operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type AuthorRepository interface {
	Create(ctx context.Context, author *entity.Author) error
	GetByID(ctx context.Context, id string) (*entity.Author, error)
	GetByName(ctx context.Context, name string) (*entity.Author, error)
	Update(ctx context.Context, author *entity.Author) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Author, error)
	Count(ctx context.Context) (int64, error)
}
