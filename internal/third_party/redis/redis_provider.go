package redis

import (
	"app/global"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisProvider struct {
	client *redis.Client
}

func NewRedisProvider() *RedisProvider {
	return &RedisProvider{
		client: global.Redis,
	}
}

func (r *RedisProvider) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisProvider) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisProvider) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisProvider) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
