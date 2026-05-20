package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

func (r *CacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *CacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (r *CacheRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *CacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

func (r *CacheRepository) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return r.client.SetNX(ctx, key, data, expiration).Result()
}

func (r *CacheRepository) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *CacheRepository) SetExpiration(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

type SessionCache struct {
	cache *CacheRepository
}

func NewSessionCache(client *redis.Client) *SessionCache {
	return &SessionCache{
		cache: NewCacheRepository(client),
	}
}

func (s *SessionCache) StoreToken(ctx context.Context, userID string, token string, expiration time.Duration) error {
	key := "session:" + token
	return s.cache.Set(ctx, key, userID, expiration)
}

func (s *SessionCache) GetUserIDFromToken(ctx context.Context, token string) (string, error) {
	var userID string
	key := "session:" + token
	err := s.cache.Get(ctx, key, &userID)
	return userID, err
}

func (s *SessionCache) InvalidateToken(ctx context.Context, token string) error {
	key := "session:" + token
	return s.cache.Delete(ctx, key)
}

type RateLimitCache struct {
	cache *CacheRepository
}

func NewRateLimitCache(client *redis.Client) *RateLimitCache {
	return &RateLimitCache{
		cache: NewCacheRepository(client),
	}
}

func (r *RateLimitCache) IncrementRequest(ctx context.Context, key string, window time.Duration) (int64, error) {
	count, err := r.cache.Increment(ctx, key)
	if err != nil {
		return 0, err
	}

	if count == 1 {
		r.cache.SetExpiration(ctx, key, window)
	}

	return count, nil
}

func (r *RateLimitCache) GetRequestCount(ctx context.Context, key string) (int64, error) {
	val, err := r.cache.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}
