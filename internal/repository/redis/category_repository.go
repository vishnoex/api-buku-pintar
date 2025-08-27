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

type categoryRedisRepository struct {
	client *redis.Client
}

func NewCategoryRedisRepository(client *redis.Client) repository.CategoryRedisRepository {
	return &categoryRedisRepository{
		client: client,
	}
}

func (r *categoryRedisRepository) GetCategoryList(ctx context.Context, limit, offset int) ([]*entity.Category, error) {
	key := fmt.Sprintf("category:list:%d:%d", limit, offset)
	
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var categories []*entity.Category
	err = json.Unmarshal([]byte(data), &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRedisRepository) SetCategoryList(ctx context.Context, categories []*entity.Category, limit, offset int) error {
	key := fmt.Sprintf("category:list:%d:%d", limit, offset)
	
	data, err := json.Marshal(categories)
	if err != nil {
		return err
	}

	// Cache for 10 minutes
	return r.client.Set(ctx, key, data, 10*time.Minute).Err()
}

func (r *categoryRedisRepository) GetCategoryTotal(ctx context.Context) (int64, error) {
	key := "category:count:total"
	
	count, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Cache miss
		}
		return 0, err
	}

	return count, nil
}

func (r *categoryRedisRepository) SetCategoryTotal(ctx context.Context, count int64) error {
	key := "category:count:total"
	
	// Cache for 10 minutes
	return r.client.Set(ctx, key, count, 10*time.Minute).Err()
}

func (r *categoryRedisRepository) GetActiveCategoryList(ctx context.Context, limit, offset int) ([]*entity.Category, error) {
	key := fmt.Sprintf("category:active:list:%d:%d", limit, offset)
	
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var categories []*entity.Category
	err = json.Unmarshal([]byte(data), &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRedisRepository) SetActiveCategoryList(ctx context.Context, categories []*entity.Category, limit, offset int) error {
	key := fmt.Sprintf("category:active:list:%d:%d", limit, offset)
	
	data, err := json.Marshal(categories)
	if err != nil {
		return err
	}

	// Cache for 10 minutes
	return r.client.Set(ctx, key, data, 10*time.Minute).Err()
}

func (r *categoryRedisRepository) GetActiveCategoryTotal(ctx context.Context) (int64, error) {
	key := "category:count:active"
	
	count, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Cache miss
		}
		return 0, err
	}

	return count, nil
}

func (r *categoryRedisRepository) SetActiveCategoryTotal(ctx context.Context, count int64) error {
	key := "category:count:active"
	
	// Cache for 10 minutes
	return r.client.Set(ctx, key, count, 10*time.Minute).Err()
}

func (r *categoryRedisRepository) GetCategoryByID(ctx context.Context, id string) (*entity.Category, error) {
	key := fmt.Sprintf("category:id:%s", id)
	
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var category entity.Category
	err = json.Unmarshal([]byte(data), &category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRedisRepository) SetCategoryByID(ctx context.Context, category *entity.Category) error {
	key := fmt.Sprintf("category:id:%s", category.ID)
	
	data, err := json.Marshal(category)
	if err != nil {
		return err
	}

	// Cache for 10 minutes
	return r.client.Set(ctx, key, data, 10*time.Minute).Err()
}

func (r *categoryRedisRepository) InvalidateCategoryCache(ctx context.Context) error {
	// Get all keys matching category pattern
	pattern := "category:*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	// Delete all category-related cache
	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}
