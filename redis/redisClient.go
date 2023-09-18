package redisClient

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCredential struct {
	Username string
	Password string
	Host     string
	Port     string
	DB       int
}
type redisClient struct {
	Redis *redis.Client
}

func (r redisClient) Ping(ctx context.Context) (string, error) {
	return r.Redis.Ping(ctx).Result()
}
func (r redisClient) Get(ctx context.Context, key string) (string, error) {
	stringCmd := r.Redis.Get(ctx, key)
	return stringCmd.Result()
}
func (r redisClient) SetWithExpire(ctx context.Context, key string, value interface{}, time time.Duration) (string, error) {
	statusCmd := r.Redis.Set(ctx, key, value, time)
	return statusCmd.Result()
}
func (r redisClient) Set(ctx context.Context, key string, value interface{}) (string, error) {
	setExpirationTime, _ := strconv.Atoi(os.Getenv("REDIS_EXPIRATION_TIME"))
	statusCmd := r.Redis.Set(ctx, key, value, time.Duration(setExpirationTime)*time.Hour)
	return statusCmd.Result()
}
func (r redisClient) Del(ctx context.Context, key string) (int64, error) {
	intCmd := r.Redis.Del(ctx, key)
	return intCmd.Result()
}
func (r redisClient) SetBit(ctx context.Context, key string, offset int64, value int) (int64, error) {
	intCmd := r.Redis.SetBit(ctx, key, offset, value)
	return intCmd.Result()
}
func (r redisClient) GetAllBits(ctx context.Context, key string) ([]bool, error) {
	re, err := r.Get(ctx, key)
	return bitstringToBool(re), err
}
func (r redisClient) GetAllKeys(ctx context.Context, prefix string) ([]string, uint64, error) {
	var cursor uint64
	intCmd := r.Redis.Scan(ctx, cursor, prefix, 10)
	return intCmd.Result()
}

func (r redisClient) Increment(ctx context.Context, key string) (int64, error) {
	intCmd := r.Redis.Incr(ctx, key)
	return intCmd.Result()
}

func (r redisClient) Decrement(ctx context.Context, key string) (int64, error) {
	intCmd := r.Redis.Decr(ctx, key)
	return intCmd.Result()
}

func (r redisClient) Pipeline(ctx context.Context, key []string, value interface{}) (string, error) {
	pipe := r.Redis.Pipeline()
	var stringCmd string
	var err error
	for _, dataKey := range key {
		stringCmd, err = r.Redis.Get(ctx, dataKey).Result()
		if err != nil {
			return "", err
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return stringCmd, nil
}
func bitstringToBool(str string) []bool {
	s := make([]bool, len(str)*8)
	for i := 0; i < len(str); i++ {
		for bit := 7; bit >= 0; bit-- {
			bitN := uint(i*8 + (7 - bit))
			s[bitN] = (str[i]>>uint(bit))&1 == 1
		}
	}
	return s
}
func NewRedisClient(redis *redis.Client) RedisClient {
	return &redisClient{Redis: redis}
}
