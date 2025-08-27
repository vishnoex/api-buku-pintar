package usecase

import (
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/entity"
	"context"
)

type BannerUsecase interface {
	ListBanner(ctx context.Context, limit, offset int) ([]*response.BannerResponse, error)
	GetBannerByID(ctx context.Context, id string) (*response.BannerResponse, error)
	CreateBanner(ctx context.Context, banner *entity.Banner) error
	UpdateBanner(ctx context.Context, banner *entity.Banner) error
	DeleteBanner(ctx context.Context, id string) error
	CountBanner(ctx context.Context) (int64, error)
	ListActiveBanner(ctx context.Context, limit, offset int) ([]*response.BannerResponse, error)
	CountActiveBanner(ctx context.Context) (int64, error)
}
