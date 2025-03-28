package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/starks97/alcohol-tracker-api/internal/dtos"
	"github.com/starks97/alcohol-tracker-api/internal/exceptions"
	"github.com/starks97/alcohol-tracker-api/internal/services"
	"github.com/starks97/alcohol-tracker-api/internal/state"
)

type RedisCmdMethos interface {
	StoreToken(c *fiber.Ctx, ctx context.Context, userID uuid.UUID, tokenMethodKey string) (dtos.TokenDetailsDto, error)
	GetAndCompareRedisValue(c *fiber.Ctx, redisClient *redis.Client, ctx context.Context, expectedValue string) (string, error)
	SetRedisValue(c *fiber.Ctx, redisClient *redis.Client, ctx context.Context, key string, value string, expiration time.Duration) error
}

// TokenService struct to manage token operations
type TokenService struct {
	AppState *state.AppState
}

// NewTokenService initializes a new token service
func NewTokenService(appState *state.AppState) RedisCmdMethos {
	return &TokenService{AppState: appState}
}

func (ts *TokenService) StoreToken(c *fiber.Ctx, ctx context.Context, userID uuid.UUID, tokenMethodKey string) (dtos.TokenDetailsDto, error) {
	var accessToken, refreshToken string
	var accessUUID, refreshUUID uuid.UUID
	var accessMaxAge, refreshMaxAge time.Duration
	var accessPrivateKey, refreshPrivateKey string
	var accessMaxAgeInt64, refreshMaxAgeInt64 int64

	// Fetch configurations
	accessMaxAge = time.Duration(ts.AppState.Config.AccessTokenMaxAge) * time.Minute
	accessPrivateKey = ts.AppState.Config.AccessTokenPrivateKey
	accessMaxAgeInt64 = ts.AppState.Config.AccessTokenMaxAge

	refreshMaxAge = time.Duration(ts.AppState.Config.RefreshTokenMaxAge) * time.Minute
	refreshPrivateKey = ts.AppState.Config.RefreshTokenPrivateKey
	refreshMaxAgeInt64 = ts.AppState.Config.RefreshTokenMaxAge

	if tokenMethodKey == "access" || tokenMethodKey == "both" {
		// Generate Access Token
		generatedAccessToken, err := services.GenerateJwtToken(userID, accessMaxAgeInt64, accessPrivateKey)
		if err != nil {
			log.Println("Failed to generate access token:", err)
			return dtos.TokenDetailsDto{}, fmt.Errorf("StoreTokens: %w", exceptions.HandlerErrorResponse(c, exceptions.ErrTokenNotGenerated))
		}

		accessToken = *generatedAccessToken.Token
		accessUUID = generatedAccessToken.TokenUUID

		// Store Access Token in Redis
		err = ts.SetRedisValue(c, ts.AppState.Redis, ctx, accessUUID.String(), userID.String(), accessMaxAge)
		if err != nil {
			return dtos.TokenDetailsDto{}, fmt.Errorf("StoreTokens: %w", err)
		}
	}

	if tokenMethodKey == "refresh" || tokenMethodKey == "both" {
		// Generate Refresh Token
		generatedRefreshToken, err := services.GenerateJwtToken(userID, refreshMaxAgeInt64, refreshPrivateKey)
		if err != nil {
			log.Println("Failed to generate refresh token:", err)
			return dtos.TokenDetailsDto{}, fmt.Errorf("StoreTokens: %w", exceptions.HandlerErrorResponse(c, exceptions.ErrTokenNotGenerated))
		}

		refreshToken = *generatedRefreshToken.Token
		refreshUUID = generatedRefreshToken.TokenUUID

		// Store Refresh Token in Redis
		err = ts.SetRedisValue(c, ts.AppState.Redis, ctx, refreshUUID.String(), userID.String(), refreshMaxAge)
		if err != nil {
			return dtos.TokenDetailsDto{}, fmt.Errorf("StoreTokens: %w", err)
		}

		// Set refresh token as a cookie
		refreshCookie := &fiber.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  time.Now().Add(refreshMaxAge),
			HTTPOnly: true,
			Secure:   false, // Fetch from config
			Path:     "/",
			Domain:   "localhost", // Fetch from config
		}
		c.Cookie(refreshCookie)
	}

	return dtos.TokenDetailsDto{
		Token: &accessToken,
	}, nil
}

func (ts *TokenService) GetAndCompareRedisValue(c *fiber.Ctx, redisClient *redis.Client, ctx context.Context, expectedValue string) (string, error) {
	// Retrieve the value from Redis using the provided key.
	cmd := ts.AppState.Redis.Get(ctx, expectedValue)

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

func (ts *TokenService) SetRedisValue(c *fiber.Ctx, redisClient *redis.Client, ctx context.Context, key string, value string, expiration time.Duration) error {
	// Set the value in Redis with the provided key and expiration.
	cmd := ts.AppState.Redis.Set(ctx, key, value, expiration)

	// Check for Redis errors.
	if cmd.Err() != nil {
		return exceptions.HandlerErrorResponse(c, exceptions.ErrRedisSet)
	}

	return nil
}
