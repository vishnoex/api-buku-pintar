package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// BannerRepository defines the interface for banner data operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type BannerRepository interface {
	Create(ctx context.Context, banner *entity.Banner) error
	GetByID(ctx context.Context, id string) (*entity.Banner, error)
	Update(ctx context.Context, banner *entity.Banner) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Banner, error)
	ListActive(ctx context.Context, limit, offset int) ([]*entity.Banner, error)
	Count(ctx context.Context) (int64, error)
	CountActive(ctx context.Context) (int64, error)
}

type BannerRedisRepository interface {
	GetBannerTotal(ctx context.Context) (int64, error)
	SetBannerTotal(ctx context.Context, data int64) error
	GetBannerList(ctx context.Context, limit, offset int) ([]*entity.Banner, error)
	SetBannerList(ctx context.Context, data []*entity.Banner, limit, offset int) error
}
