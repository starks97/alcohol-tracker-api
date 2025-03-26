package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/starks97/alcohol-tracker-api/internal/exceptions"
	"github.com/starks97/alcohol-tracker-api/internal/models"
	"github.com/starks97/alcohol-tracker-api/internal/services"
	"github.com/starks97/alcohol-tracker-api/internal/state"
)

// GetAndCompareRedisValue retrieves a value from Redis using the provided key,
// compares it to the expected value, and returns the retrieved value if they match.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context for error responses.
//   - redisClient: *redis.Client - The Redis client.
//   - ctx: context.Context - The context for Redis operations.
//   - expectedValue: string - The key to retrieve from Redis and the value to compare against.
//
// Returns:
//   - string: The retrieved value from Redis if it matches the expected value.
//   - error: An error if the Redis operation fails, the key is not found, or the values don't match.
func GetAndCompareRedisValue(c *fiber.Ctx, redisClient *redis.Client, ctx context.Context, expectedValue string) (string, error) {
	// Retrieve the value from Redis using the provided key.
	cmd := redisClient.Get(ctx, expectedValue)

	// Check for Redis errors.
	if cmd.Err() != nil {
		return "", exceptions.HandlerErrorResponse(c, exceptions.ErrRedisGet)
	}

	// Get the retrieved value.
	redisValue := cmd.Val()

	// Check if the retrieved value is empty (key not found).
	if redisValue == "" {
		return "", exceptions.HandlerErrorResponse(c, exceptions.ErrRedisNotFound)
	}

	// Compare the retrieved value with the expected value.
	if redisValue != expectedValue {
		return "", exceptions.HandlerErrorResponse(c, exceptions.ErrTokenMismatch)
	}

	// Return the retrieved value if it matches the expected value.
	return redisValue, nil
}

// setRedisValue sets a key-value pair in Redis with the specified expiration time.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context for error responses.
//   - redisClient: *redis.Client - The Redis client.
//   - ctx: context.Context - The context for Redis operations.
//   - key: string - The key to set in Redis.
//   - value: string - The value to set in Redis.
//   - expiration: time.Duration - The expiration time for the key-value pair.
//
// Returns:
//   - error: An error if the Redis operation fails.
func setRedisValue(c *fiber.Ctx, redisClient *redis.Client, ctx context.Context, key string, value string, expiration time.Duration) error {
	// Set the key-value pair in Redis with the specified expiration time.
	cmd := redisClient.SetEx(ctx, key, value, expiration)

	// Check for Redis errors.
	if cmd.Err() != nil {
		return exceptions.HandlerErrorResponse(c, exceptions.ErrRedisSet)
	}

	// Return nil if the Redis operation is successful.
	return nil
}

// StoreTokens generates JWT tokens for the given user ID, stores them in Redis,
// and sets a refresh token cookie.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context for error responses.
//   - appState: *state.AppState - The application state containing configuration and Redis client.
//   - ctx: context.Context - The context for token generation and Redis operations.
//   - userID: uuid.UUID - The user ID to generate tokens for.
//
// Returns:
//   - models.TokenDetails: Token details containing the generated access token.
//   - error: An error if any step of the process fails.
func StoreTokens(c *fiber.Ctx, appState *state.AppState, ctx context.Context, userID uuid.UUID) (models.TokenDetails, error) {
	// Calculate access and refresh token expiration times.
	accessTokenMaxAge := time.Duration(appState.Config.AccessTokenMaxAge) * time.Minute
	refreshTokenMaxAge := time.Duration(appState.Config.RefreshTokenMaxAge) * time.Minute

	// Generate access token.
	generateAccessToken, err := services.GenerateJwtToken(userID, appState.Config.AccessTokenMaxAge, appState.Config.AccessTokenPrivateKey)
	if err != nil {
		log.Println("Failed to generate access token:", err)
		return models.TokenDetails{}, fmt.Errorf("StoreTokens: %w", exceptions.HandlerErrorResponse(c, exceptions.ErrTokenNotGenerated))
	}

	// Generate refresh token.
	generateRefreshToken, err := services.GenerateJwtToken(userID, appState.Config.RefreshTokenMaxAge, appState.Config.RefreshTokenPrivateKey)
	if err != nil {
		log.Println("Failed to generate refresh token:", err)
		return models.TokenDetails{}, fmt.Errorf("StoreTokens: %w", exceptions.HandlerErrorResponse(c, exceptions.ErrTokenNotGenerated))
	}

	// Store access token in Redis.
	err = setRedisValue(c, appState.Redis, ctx, generateAccessToken.TokenUUID.String(), generateAccessToken.UserID.String(), accessTokenMaxAge)
	if err != nil {
		return models.TokenDetails{}, fmt.Errorf("StoreTokens: %w", err)
	}

	// Store refresh token in Redis.
	err = setRedisValue(c, appState.Redis, ctx, generateRefreshToken.TokenUUID.String(), generateRefreshToken.UserID.String(), refreshTokenMaxAge)
	if err != nil {
		return models.TokenDetails{}, fmt.Errorf("StoreTokens: %w", err)
	}

	// Create refresh token cookie.
	refreshCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    *generateRefreshToken.Token,
		Expires:  time.Now().Add(time.Duration(appState.Config.RefreshTokenMaxAge) * time.Minute),
		HTTPOnly: true,
		Secure:   false, // Use config
		Path:     "/",
		Domain:   "localhost", // Use config
	}

	// Set refresh token cookie.
	c.Cookie(refreshCookie)

	// Return access token details.
	return models.TokenDetails{
		Token: generateAccessToken.Token,
	}, nil
}
