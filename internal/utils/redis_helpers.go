package utils

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/starks97/alcohol-tracker-api/internal/errors"
)

func GetAndCompareRedisValue(c *fiber.Ctx, redisClient *redis.Client, ctx context.Context, expectedValue string) (string, error) {
	cmd := redisClient.Get(ctx, expectedValue)

	if cmd.Err() != nil {
		return "", errors.NewCustomErrorResponse(c, errors.ErrRedisGet)
	}

	redisValue := cmd.Val()

	if redisValue == "" {
		return "", errors.NewCustomErrorResponse(c, errors.ErrRedisNotFound)
	}

	if redisValue != expectedValue {
		return "", errors.NewCustomErrorResponse(c, errors.ErrTokenMismatch)
	}

	return redisValue, nil
}
