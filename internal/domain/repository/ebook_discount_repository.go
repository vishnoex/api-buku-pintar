package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// EbookDiscountRepository defines the interface for ebook discount data operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type EbookDiscountRepository interface {
	Create(ctx context.Context, discount *entity.EbookDiscount) error
	GetByID(ctx context.Context, id string) (*entity.EbookDiscount, error)
	Update(ctx context.Context, discount *entity.EbookDiscount) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error)
	GetByEbookID(ctx context.Context, ebookID string) ([]*entity.EbookDiscount, error)
	GetActiveDiscounts(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error)
	GetActiveDiscountByEbookID(ctx context.Context, ebookID string) (*entity.EbookDiscount, error)
	Count(ctx context.Context) (int64, error)
	CountByEbookID(ctx context.Context, ebookID string) (int64, error)
	CountActiveDiscounts(ctx context.Context) (int64, error)
}
