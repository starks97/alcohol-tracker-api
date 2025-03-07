package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/starks97/alcohol-tracker-api/config"
)

func NewRedisClient(cfg *config.Config, ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	fmt.Println("✅ Redis connected successfully")

	return client, nil
}
