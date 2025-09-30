package redis

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type ebookDiscountRedisRepository struct {
	client *redis.Client
	ttl    time.Duration
}

// NewEbookDiscountRedisRepository creates a new instance of EbookDiscountRedisRepository
func NewEbookDiscountRedisRepository(client *redis.Client) repository.EbookDiscountRedisRepository {
	return &ebookDiscountRedisRepository{
		client: client,
		ttl:    10 * time.Minute, // 10 minutes TTL
	}
}

func (r *ebookDiscountRedisRepository) GetDiscountByID(ctx context.Context, id string) (*entity.EbookDiscount, error) {
	key := fmt.Sprintf("discount:id:%s", id)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}
	
	var discount entity.EbookDiscount
	err = json.Unmarshal([]byte(data), &discount)
	if err != nil {
		return nil, err
	}
	
	return &discount, nil
}

func (r *ebookDiscountRedisRepository) SetDiscountByID(ctx context.Context, discount *entity.EbookDiscount) error {
	key := fmt.Sprintf("discount:id:%s", discount.ID)
	data, err := json.Marshal(discount)
	if err != nil {
		return err
	}
	
	return r.client.Set(ctx, key, string(data), r.ttl).Err()
}

func (r *ebookDiscountRedisRepository) GetDiscountList(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error) {
	key := fmt.Sprintf("discount:list:%d:%d", limit, offset)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}
	
	var discounts []*entity.EbookDiscount
	err = json.Unmarshal([]byte(data), &discounts)
	if err != nil {
		return nil, err
	}
	
	return discounts, nil
}

func (r *ebookDiscountRedisRepository) SetDiscountList(ctx context.Context, discounts []*entity.EbookDiscount, limit, offset int) error {
	key := fmt.Sprintf("discount:list:%d:%d", limit, offset)
	data, err := json.Marshal(discounts)
	if err != nil {
		return err
	}
	
	return r.client.Set(ctx, key, string(data), r.ttl).Err()
}

func (r *ebookDiscountRedisRepository) GetDiscountsByEbookID(ctx context.Context, ebookID string) ([]*entity.EbookDiscount, error) {
	key := fmt.Sprintf("discount:ebook:%s", ebookID)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}
	
	var discounts []*entity.EbookDiscount
	err = json.Unmarshal([]byte(data), &discounts)
	if err != nil {
		return nil, err
	}
	
	return discounts, nil
}

func (r *ebookDiscountRedisRepository) SetDiscountsByEbookID(ctx context.Context, ebookID string, discounts []*entity.EbookDiscount) error {
	key := fmt.Sprintf("discount:ebook:%s", ebookID)
	data, err := json.Marshal(discounts)
	if err != nil {
		return err
	}
	
	return r.client.Set(ctx, key, string(data), r.ttl).Err()
}

func (r *ebookDiscountRedisRepository) GetActiveDiscountByEbookID(ctx context.Context, ebookID string) (*entity.EbookDiscount, error) {
	key := fmt.Sprintf("discount:active:ebook:%s", ebookID)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}
	
	var discount entity.EbookDiscount
	err = json.Unmarshal([]byte(data), &discount)
	if err != nil {
		return nil, err
	}
	
	return &discount, nil
}

func (r *ebookDiscountRedisRepository) SetActiveDiscountByEbookID(ctx context.Context, ebookID string, discount *entity.EbookDiscount) error {
	key := fmt.Sprintf("discount:active:ebook:%s", ebookID)
	data, err := json.Marshal(discount)
	if err != nil {
		return err
	}
	
	return r.client.Set(ctx, key, string(data), r.ttl).Err()
}

func (r *ebookDiscountRedisRepository) GetActiveDiscounts(ctx context.Context, limit, offset int) ([]*entity.EbookDiscount, error) {
	key := fmt.Sprintf("discount:active:list:%d:%d", limit, offset)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}
	
	var discounts []*entity.EbookDiscount
	err = json.Unmarshal([]byte(data), &discounts)
	if err != nil {
		return nil, err
	}
	
	return discounts, nil
}

func (r *ebookDiscountRedisRepository) SetActiveDiscounts(ctx context.Context, discounts []*entity.EbookDiscount, limit, offset int) error {
	key := fmt.Sprintf("discount:active:list:%d:%d", limit, offset)
	data, err := json.Marshal(discounts)
	if err != nil {
		return err
	}
	
	return r.client.Set(ctx, key, string(data), r.ttl).Err()
}

func (r *ebookDiscountRedisRepository) GetDiscountCount(ctx context.Context) (int64, error) {
	key := "discount:count:total"
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Cache miss
		}
		return 0, err
	}
	
	var count int64
	err = json.Unmarshal([]byte(data), &count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

func (r *ebookDiscountRedisRepository) SetDiscountCount(ctx context.Context, count int64) error {
	key := "discount:count:total"
	data, err := json.Marshal(count)
	if err != nil {
		return err
	}
	
	return r.client.Set(ctx, key, string(data), r.ttl).Err()
}

func (r *ebookDiscountRedisRepository) GetDiscountCountByEbookID(ctx context.Context, ebookID string) (int64, error) {
	key := fmt.Sprintf("discount:count:ebook:%s", ebookID)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Cache miss
		}
		return 0, err
	}
	
	var count int64
	err = json.Unmarshal([]byte(data), &count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

func (r *ebookDiscountRedisRepository) SetDiscountCountByEbookID(ctx context.Context, ebookID string, count int64) error {
	key := fmt.Sprintf("discount:count:ebook:%s", ebookID)
	data, err := json.Marshal(count)
	if err != nil {
		return err
	}
	
	return r.client.Set(ctx, key, string(data), r.ttl).Err()
}

func (r *ebookDiscountRedisRepository) GetActiveDiscountCount(ctx context.Context) (int64, error) {
	key := "discount:count:active"
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Cache miss
		}
		return 0, err
	}
	
	var count int64
	err = json.Unmarshal([]byte(data), &count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

func (r *ebookDiscountRedisRepository) SetActiveDiscountCount(ctx context.Context, count int64) error {
	key := "discount:count:active"
	data, err := json.Marshal(count)
	if err != nil {
		return err
	}
	
	return r.client.Set(ctx, key, string(data), r.ttl).Err()
}

func (r *ebookDiscountRedisRepository) InvalidateDiscountCache(ctx context.Context) error {
	pattern := "discount:*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	
	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}
	
	return nil
}