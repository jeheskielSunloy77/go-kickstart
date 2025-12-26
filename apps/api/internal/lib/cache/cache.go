package cache

import (
	"context"
	"errors"
	"time"

	"github.com/jeheskielSunloy77/go-kickstart/internal/config"
	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl ...time.Duration) error
	Delete(ctx context.Context, keys ...string) error
}

type RedisCache struct {
	client *redis.Client
	cfg    *config.CacheConfig
}

func NewRedisCache(client *redis.Client, cfg *config.CacheConfig) *RedisCache {
	return &RedisCache{client: client, cfg: cfg}
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	return data, err
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl ...time.Duration) error {
	var expiration time.Duration
	if len(ttl) > 0 {
		expiration = ttl[0]
	} else {
		expiration = c.cfg.TTL
	}
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}
