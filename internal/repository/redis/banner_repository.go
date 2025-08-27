package redis

import (
	"buku-pintar/internal/constant"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type bannerRedisRepository struct {
    client *redis.Client
}

func NewBannerRedisRepository(client *redis.Client) repository.BannerRedisRepository {
    return &bannerRedisRepository{client: client}
}

func (r *bannerRedisRepository) GetBannerTotal(ctx context.Context) (int64, error) {
    return 0, errors.New(constant.ERR_NOT_IMPLEMENTED)
}

func (r *bannerRedisRepository) SetBannerTotal(ctx context.Context, data int64) error {
    return errors.New(constant.ERR_NOT_IMPLEMENTED)
}

func (r *bannerRedisRepository) GetBannerList(ctx context.Context, limit, offset int) ([]*entity.Banner, error) {
    return nil, errors.New(constant.ERR_NOT_IMPLEMENTED)
}

func (r *bannerRedisRepository) SetBannerList(ctx context.Context, data []*entity.Banner, limit, offset int) error {
    return errors.New(constant.ERR_NOT_IMPLEMENTED)
}
