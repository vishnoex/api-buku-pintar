package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// EbookDiscountRedisRepository defines the interface for ebook discount Redis operations
type EbookDiscountRedisRepository interface {
	GetDiscountByID(ctx context.Context, id string) (*entity.EbookDiscount, error)
	SetDiscountByID(ctx context.Context, discount *entity.EbookDiscount) error
	GetDiscountList(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error)
	SetDiscountList(ctx context.Context, discounts []*entity.EbookDiscount, limit, offset int) error
	GetDiscountsByEbookID(ctx context.Context, ebookID string) ([]*entity.EbookDiscount, error)
	SetDiscountsByEbookID(ctx context.Context, ebookID string, discounts []*entity.EbookDiscount) error
	GetActiveDiscountByEbookID(ctx context.Context, ebookID string) (*entity.EbookDiscount, error)
	SetActiveDiscountByEbookID(ctx context.Context, ebookID string, discount *entity.EbookDiscount) error
	GetActiveDiscounts(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error)
	SetActiveDiscounts(ctx context.Context, discounts []*entity.EbookDiscount, limit, offset int) error
	GetDiscountCount(ctx context.Context) (int64, error)
	SetDiscountCount(ctx context.Context, count int64) error
	GetDiscountCountByEbookID(ctx context.Context, ebookID string) (int64, error)
	SetDiscountCountByEbookID(ctx context.Context, ebookID string, count int64) error
	GetActiveDiscountCount(ctx context.Context) (int64, error)
	SetActiveDiscountCount(ctx context.Context, count int64) error
	InvalidateDiscountCache(ctx context.Context) error
}
