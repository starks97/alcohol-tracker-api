package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/starks97/alcohol-tracker-api/config"
)

// NewRedisClient creates a new Redis client and pings the Redis server to ensure a successful connection.
//
// Parameters:
//   - cfg: *config.Config - The application configuration containing Redis connection details.
//   - ctx: context.Context - The context for the Redis ping operation.
//
// Returns:
//   - *redis.Client: A pointer to the initialized Redis client.
//   - error: An error if the Redis connection or ping fails, or nil if successful.
func NewRedisClient(cfg *config.Config, ctx context.Context) (*redis.Client, error) {
	// Create a new Redis client using the provided configuration.
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       0, // Use the default database (0).
	})

	// Ping the Redis server to verify the connection.
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Log a successful Redis connection message.
	fmt.Println("✅ Redis connected successfully")

	// Return the initialized Redis client.
	return client, nil
}
