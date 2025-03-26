package config

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfg struct {
	Addr     string
	Password string
}

var Client *redis.Client
var TotalPosts int64

func ConnectCache() error {
	cfg := loadRedisConfig()

	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       1,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := Client.Ping(ctx).Result()
	return err
}

func SetTotalPosts(count int64) error {
	ctx := context.Background()
	return Client.Set(ctx, "total_posts", count, 0).Err()
}

func GetTotalPosts() (int64, error) {
	ctx := context.Background()
	return Client.Get(ctx, "total_posts").Int64()
}

func IncrementTotalPosts() error {
	ctx := context.Background()
	return Client.Incr(ctx, "total_posts").Err()
}

func DecrementTotalPosts() error {
	ctx := context.Background()
	return Client.Decr(ctx, "total_posts").Err()
}

func loadRedisConfig() RedisConfg {
	return RedisConfg{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
}
