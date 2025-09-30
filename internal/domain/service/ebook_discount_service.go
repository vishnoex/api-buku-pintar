package service

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// EbookDiscountService defines the interface for ebook discount business operations
type EbookDiscountService interface {
	CreateDiscount(ctx context.Context, discount *entity.EbookDiscount) error
	GetDiscountByID(ctx context.Context, id string) (*entity.EbookDiscount, error)
	UpdateDiscount(ctx context.Context, discount *entity.EbookDiscount) error
	DeleteDiscount(ctx context.Context, id string) error
	GetDiscountList(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error)
	GetDiscountsByEbookID(ctx context.Context, ebookID string) ([]*entity.EbookDiscount, error)
	GetActiveDiscounts(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error)
	GetActiveDiscountByEbookID(ctx context.Context, ebookID string) (*entity.EbookDiscount, error)
	GetDiscountCount(ctx context.Context) (int64, error)
	GetDiscountCountByEbookID(ctx context.Context, ebookID string) (int64, error)
	GetActiveDiscountCount(ctx context.Context) (int64, error)
}
