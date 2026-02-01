package redis

import (
	"context"
	"os"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var dbClient *redis.Client

func InitDatabase() {
	addr := os.Getenv("CACHE_HOST")
	port := os.Getenv("CACHE_PORT")
	pass := os.Getenv("CACHE_PASSWORD")
	if port == "" {
		port = "6379"
	}
	dbClient = redis.NewClient(&redis.Options{
		Addr:     addr + ":" + port,
		Password: pass,
	})
}

func Set(ctx context.Context, key string, value interface{}) error {
	return dbClient.Set(ctx, key, value, 24*60*60*1e9).Err() // 1 day in nanoseconds
}

// SetWithTTL sets a key with custom TTL duration
func SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return dbClient.Set(ctx, key, value, ttl).Err()
}

// Exists checks if a key exists in Redis
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := dbClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func Get(ctx context.Context, key string) (string, error) {
	return dbClient.Get(ctx, key).Result()
}

func Del(ctx context.Context, key string) error {
	return dbClient.Del(ctx, key).Err()
}
