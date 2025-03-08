package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/starks97/alcohol-tracker-api/internal/repositories"
	"github.com/starks97/alcohol-tracker-api/internal/responses"
	"github.com/starks97/alcohol-tracker-api/internal/services"
	"github.com/starks97/alcohol-tracker-api/internal/state"
)

var (
	ErrTokenMissing      = errors.New("token not provided or not found")
	ErrTokenVerification = errors.New("failed to verify token")
	ErrRedisGet          = errors.New("failed to get token from redis")
	ErrRedisNotFound     = errors.New("token not found in redis")
	ErrTokenMismatch     = errors.New("token mismatch")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserIDMismatch    = errors.New("user ID mismatch")
)

//todo improve error handling

func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		bearerToken := c.Get("Authorization")
		appState := c.Locals("cfg").(*state.AppState)
		ctx := c.Locals("ctx").(context.Context)
		userRepo := repositories.NewUserRepository(appState.DB)

		if bearerToken == "" {
			log.Println("Failed to exchange token:")
			return c.Status(http.StatusUnauthorized).JSON(responses.ErrorResponse{
				Message: "Token not provided or not found, please provide a valid token",
				Status:  "failed",
			})
		}

		token := strings.TrimPrefix(bearerToken, "Bearer ")

		verifyToken, err := services.VerifyJwtToken(appState.Config.AccessTokenPublicKey, token)
		if err != nil {
			log.Println("Failed to verify token:", err)
			return c.Status(http.StatusUnauthorized).JSON(responses.ErrorResponse{
				Message: "Failed to verify token:",
				Status:  "failed",
			})
		}
		accessTokenUuid := verifyToken.TokenUUID

		findUserIDInRedis := appState.Redis.Get(ctx, accessTokenUuid.String())

		if findUserIDInRedis.Err() != nil {
			log.Println("Failed to get token from redis:", findUserIDInRedis.Err())
			return c.Status(http.StatusUnauthorized).JSON(responses.ErrorResponse{
				Message: "Failed to get token from redis",
				Status:  "failed",
			})
		}
		if findUserIDInRedis.Val() == "" {
			log.Println("Token not found in redis")
			return c.Status(http.StatusUnauthorized).JSON(responses.ErrorResponse{
				Message: "Token not found in redis",
				Status:  "failed",
			})
		}
		if findUserIDInRedis.Val() != accessTokenUuid.String() {
			log.Println("Token mismatch")
			return c.Status(http.StatusUnauthorized).JSON(responses.ErrorResponse{
				Message: "Token mismatch",
				Status:  "failed",
			})
		}

		userID, err := uuid.Parse(findUserIDInRedis.Val())
		if err != nil {
			log.Println("Failed to parse user ID from redis:", err)
			return c.Status(http.StatusUnauthorized).JSON(responses.ErrorResponse{
				Message: "Failed to parse user ID from redis",
				Status:  "failed",
			})
		}

		userIDInDB, err := userRepo.GetUserByID(userID)
		if err != nil {
			log.Println("Failed to get user from database:", err)
			return c.Status(http.StatusNotFound).JSON(responses.ErrorResponse{
				Status:  "failed",
				Message: "The user could not be found with the provided ID, please provide a valid ID",
			})
		}
		if userIDInDB.ID != userID {
			log.Println("User ID mismatch")
			return c.Status(http.StatusConflict).JSON(responses.ErrorResponse{
				Status:  "failed",
				Message: "User ID mismatch",
			})
		}

		jwtResponse := responses.JwtMiddlewareResponse{
			Token: accessTokenUuid,
			User: responses.UserResponse{
				ID:         userIDInDB.ID,
				Email:      *userIDInDB.Email,
				Name:       *userIDInDB.Name,
				CreatedAt:  userIDInDB.CreatedAt,
				UpdatedAt:  userIDInDB.UpdatedAt,
				Provider:   *userIDInDB.Provider,
				ProviderID: *userIDInDB.ProviderID,
			},
		}

		c.JSON(responses.SuccessResponse{
			Status: "success",
			Data:   jwtResponse,
		})
		return c.Next()
	}
}
