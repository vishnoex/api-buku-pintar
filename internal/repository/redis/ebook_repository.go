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

type ebookRedisRepository struct {
	client *redis.Client
}

func NewEbookRedisRepository(client *redis.Client) repository.EbookRedisRepository {
	return &ebookRedisRepository{
		client: client,
	}
}

func (r *ebookRedisRepository) GetEbookList(ctx context.Context, limit, offset int) ([]*entity.EbookList, error) {
	key := fmt.Sprintf("ebook:list:%d:%d", limit, offset)
	
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var ebooks []*entity.EbookList
	err = json.Unmarshal([]byte(data), &ebooks)
	if err != nil {
		return nil, err
	}

	return ebooks, nil
}

func (r *ebookRedisRepository) SetEbookList(ctx context.Context, ebooks []*entity.EbookList, limit, offset int) error {
	key := fmt.Sprintf("ebook:list:%d:%d", limit, offset)
	
	data, err := json.Marshal(ebooks)
	if err != nil {
		return err
	}

	// Cache for 15 minutes
	return r.client.Set(ctx, key, data, 15*time.Minute).Err()
}

func (r *ebookRedisRepository) GetEbookTotal(ctx context.Context) (int64, error) {
	key := "ebook:count:total"
	
	count, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Cache miss
		}
		return 0, err
	}

	return count, nil
}

func (r *ebookRedisRepository) SetEbookTotal(ctx context.Context, count int64) error {
	key := "ebook:count:total"
	
	// Cache for 15 minutes
	return r.client.Set(ctx, key, count, 15*time.Minute).Err()
}

func (r *ebookRedisRepository) GetEbookByID(ctx context.Context, id string) (*entity.Ebook, error) {
	key := fmt.Sprintf("ebook:id:%s", id)
	
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var ebook entity.Ebook
	err = json.Unmarshal([]byte(data), &ebook)
	if err != nil {
		return nil, err
	}

	return &ebook, nil
}

func (r *ebookRedisRepository) SetEbookByID(ctx context.Context, ebook *entity.Ebook) error {
	key := fmt.Sprintf("ebook:id:%s", ebook.ID)
	
	data, err := json.Marshal(ebook)
	if err != nil {
		return err
	}

	// Cache for 15 minutes
	return r.client.Set(ctx, key, data, 15*time.Minute).Err()
}

func (r *ebookRedisRepository) GetEbookBySlug(ctx context.Context, slug string) (*entity.Ebook, error) {
	key := fmt.Sprintf("ebook:slug:%s", slug)
	
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var ebook entity.Ebook
	err = json.Unmarshal([]byte(data), &ebook)
	if err != nil {
		return nil, err
	}

	return &ebook, nil
}

func (r *ebookRedisRepository) SetEbookBySlug(ctx context.Context, ebook *entity.Ebook) error {
	key := fmt.Sprintf("ebook:slug:%s", ebook.Slug)
	
	data, err := json.Marshal(ebook)
	if err != nil {
		return err
	}

	// Cache for 15 minutes
	return r.client.Set(ctx, key, data, 15*time.Minute).Err()
}

func (r *ebookRedisRepository) GetEbookListByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Ebook, error) {
	key := fmt.Sprintf("ebook:category:%s:list:%d:%d", categoryID, limit, offset)
	
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var ebooks []*entity.Ebook
	err = json.Unmarshal([]byte(data), &ebooks)
	if err != nil {
		return nil, err
	}

	return ebooks, nil
}

func (r *ebookRedisRepository) SetEbookListByCategory(ctx context.Context, ebooks []*entity.Ebook, categoryID string, limit, offset int) error {
	key := fmt.Sprintf("ebook:category:%s:list:%d:%d", categoryID, limit, offset)
	
	data, err := json.Marshal(ebooks)
	if err != nil {
		return err
	}

	// Cache for 15 minutes
	return r.client.Set(ctx, key, data, 15*time.Minute).Err()
}

func (r *ebookRedisRepository) GetEbookCountByCategory(ctx context.Context, categoryID string) (int64, error) {
	key := fmt.Sprintf("ebook:category:%s:count", categoryID)
	
	count, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Cache miss
		}
		return 0, err
	}

	return count, nil
}

func (r *ebookRedisRepository) SetEbookCountByCategory(ctx context.Context, categoryID string, count int64) error {
	key := fmt.Sprintf("ebook:category:%s:count", categoryID)
	
	// Cache for 15 minutes
	return r.client.Set(ctx, key, count, 15*time.Minute).Err()
}

func (r *ebookRedisRepository) GetEbookListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Ebook, error) {
	key := fmt.Sprintf("ebook:author:%s:list:%d:%d", authorID, limit, offset)
	
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var ebooks []*entity.Ebook
	err = json.Unmarshal([]byte(data), &ebooks)
	if err != nil {
		return nil, err
	}

	return ebooks, nil
}

func (r *ebookRedisRepository) SetEbookListByAuthor(ctx context.Context, ebooks []*entity.Ebook, authorID string, limit, offset int) error {
	key := fmt.Sprintf("ebook:author:%s:list:%d:%d", authorID, limit, offset)
	
	data, err := json.Marshal(ebooks)
	if err != nil {
		return err
	}

	// Cache for 15 minutes
	return r.client.Set(ctx, key, data, 15*time.Minute).Err()
}

func (r *ebookRedisRepository) GetEbookCountByAuthor(ctx context.Context, authorID string) (int64, error) {
	key := fmt.Sprintf("ebook:author:%s:count", authorID)
	
	count, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Cache miss
		}
		return 0, err
	}

	return count, nil
}

func (r *ebookRedisRepository) SetEbookCountByAuthor(ctx context.Context, authorID string, count int64) error {
	key := fmt.Sprintf("ebook:author:%s:count", authorID)
	
	// Cache for 15 minutes
	return r.client.Set(ctx, key, count, 15*time.Minute).Err()
}

func (r *ebookRedisRepository) InvalidateEbookCache(ctx context.Context) error {
	// Get all keys matching ebook pattern
	pattern := "ebook:*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	// Delete all ebook-related cache
	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}
