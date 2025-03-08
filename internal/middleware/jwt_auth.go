package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/starks97/alcohol-tracker-api/internal/services"
	"github.com/starks97/alcohol-tracker-api/internal/state"
)

func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		bearerToken := c.Get("Authorization")
		appState := c.Locals("cfg").(*state.AppState)
		ctx := c.Locals("ctx").(context.Context)

		if bearerToken == "" {
			log.Println("Failed to exchange token:")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  "failed",
				"message": "Token not provided or not found, please provide a valid token",
			})
		}

		token := strings.TrimPrefix(bearerToken, "Bearer ")

		verifyToken, err := services.VerifyJwtToken(appState.Config.AccessTokenPublicKey, token)
		if err != nil {
			log.Println("Failed to verify token:", err)
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  "failed",
				"message": "Invalid token",
			})
		}
		accessTokenUuid := verifyToken.TokenUUID

		redisResult := appState.Redis.Get(ctx, accessTokenUuid)

		if redisResult.Err() != nil {
			log.Println("Failed to get token from redis:", redisResult.Err())
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  "failed",
				"message": "Invalid token",
			})
		}
		if redisResult.Val() == "" {
			log.Println("Token not found in redis")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  "failed",
				"message": "Invalid token",
			})
		}
		if redisResult.Val() != accessTokenUuid {
			log.Println("Token mismatch")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  "failed",
				"message": "Invalid token",
			})
		}

		return c.Next()
	}
}
