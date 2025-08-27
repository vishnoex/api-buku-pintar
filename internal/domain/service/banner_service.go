package service

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// BannerService defines the interface for banner business operations
type BannerService interface {
	GetBannerList(ctx context.Context, limit, offset int) ([]*entity.Banner, error)
	GetBannerByID(ctx context.Context, id string) (*entity.Banner, error)
	CreateBanner(ctx context.Context, banner *entity.Banner) error
	UpdateBanner(ctx context.Context, banner *entity.Banner) error
	DeleteBanner(ctx context.Context, id string) error
	GetBannerCount(ctx context.Context) (int64, error)
	GetActiveBannerList(ctx context.Context, limit, offset int) ([]*entity.Banner, error)
	GetActiveBannerCount(ctx context.Context) (int64, error)
}
