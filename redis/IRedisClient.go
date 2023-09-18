package redisClient

import (
	"context"
	"time"
)

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	SetWithExpire(ctx context.Context, key string, value interface{}, second time.Duration) (string, error)
	Set(ctx context.Context, key string, value interface{}) (string, error)
	Del(ctx context.Context, key string) (int64, error)
	Ping(ctx context.Context) (string, error)
	Pipeline(ctx context.Context, key []string, value interface{}) (string, error)
	GetAllKeys(ctx context.Context, prefix string) ([]string, uint64, error)
	SetBit(ctx context.Context, key string, offset int64, value int) (int64, error)
	GetAllBits(ctx context.Context, key string) ([]bool, error)
	Increment(ctx context.Context, key string) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)
}
