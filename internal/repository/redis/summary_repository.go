package redis

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type SummaryRedisRepositoryImpl struct {
	client *redis.Client
	ttl    time.Duration
}

func NewSummaryRedisRepositoryImpl(client *redis.Client) repository.SummaryRedisRepository {
	return &SummaryRedisRepositoryImpl{
		client: client,
		ttl:    15 * time.Minute, // 15 minutes TTL
	}
}

func (r *SummaryRedisRepositoryImpl) GetSummaryByID(ctx context.Context, id string) (*entity.EbookSummary, error) {
	key := fmt.Sprintf("summary:id:%s", id)
	
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var summary entity.EbookSummary
	err = json.Unmarshal([]byte(val), &summary)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func (r *SummaryRedisRepositoryImpl) SetSummaryByID(ctx context.Context, id string, summary *entity.EbookSummary) error {
	key := fmt.Sprintf("summary:id:%s", id)
	
	data, err := json.Marshal(summary)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *SummaryRedisRepositoryImpl) GetSummariesList(ctx context.Context, limit, offset int) ([]*entity.EbookSummaryList, error) {
	key := fmt.Sprintf("summary:list:%d:%d", limit, offset)
	
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var summaries []*entity.EbookSummaryList
	err = json.Unmarshal([]byte(val), &summaries)
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

func (r *SummaryRedisRepositoryImpl) SetSummariesList(ctx context.Context, limit, offset int, summaries []*entity.EbookSummaryList) error {
	key := fmt.Sprintf("summary:list:%d:%d", limit, offset)
	
	data, err := json.Marshal(summaries)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *SummaryRedisRepositoryImpl) GetSummariesByEbookID(ctx context.Context, ebookID string, limit, offset int) ([]*entity.EbookSummary, error) {
	key := fmt.Sprintf("summary:ebook:%s:%d:%d", ebookID, limit, offset)
	
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var summaries []*entity.EbookSummary
	err = json.Unmarshal([]byte(val), &summaries)
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

func (r *SummaryRedisRepositoryImpl) SetSummariesByEbookID(ctx context.Context, ebookID string, limit, offset int, summaries []*entity.EbookSummary) error {
	key := fmt.Sprintf("summary:ebook:%s:%d:%d", ebookID, limit, offset)
	
	data, err := json.Marshal(summaries)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *SummaryRedisRepositoryImpl) GetSummariesCount(ctx context.Context) (int64, error) {
	key := "summary:count:total"
	
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Cache miss
		}
		return 0, err
	}

	count, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *SummaryRedisRepositoryImpl) SetSummariesCount(ctx context.Context, count int64) error {
	key := "summary:count:total"
	
	return r.client.Set(ctx, key, count, r.ttl).Err()
}

func (r *SummaryRedisRepositoryImpl) GetSummariesCountByEbookID(ctx context.Context, ebookID string) (int64, error) {
	key := fmt.Sprintf("summary:count:ebook:%s", ebookID)
	
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Cache miss
		}
		return 0, err
	}

	count, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *SummaryRedisRepositoryImpl) SetSummariesCountByEbookID(ctx context.Context, ebookID string, count int64) error {
	key := fmt.Sprintf("summary:count:ebook:%s", ebookID)
	
	return r.client.Set(ctx, key, count, r.ttl).Err()
}

func (r *SummaryRedisRepositoryImpl) ClearCache(ctx context.Context) error {
	// Clear all summary-related cache keys
	pattern := "summary:*"
	
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}
