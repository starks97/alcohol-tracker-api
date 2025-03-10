package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/starks97/alcohol-tracker-api/internal/errors"
	"github.com/starks97/alcohol-tracker-api/internal/repositories"
	"github.com/starks97/alcohol-tracker-api/internal/responses"
	"github.com/starks97/alcohol-tracker-api/internal/services"
	"github.com/starks97/alcohol-tracker-api/internal/state"
	"github.com/starks97/alcohol-tracker-api/internal/utils"
)

// JWTAuthMiddleware creates a Fiber middleware handler that authenticates requests using JWT tokens.
// It retrieves the token from the "Authorization" header, verifies it, and retrieves the associated user.
//
// The middleware performs the following steps:
// 1. Retrieves the bearer token from the "Authorization" header.
// 2. Verifies the token using the provided public key.
// 3. Retrieves the user ID associated with the token from Redis.
// 4. Retrieves the user from the database using the retrieved user ID.
// 5. Verifies that the user ID from Redis matches the user ID from the database.
// 6. If all steps are successful, it adds the user information to the response and calls the next handler.
// 7. If any step fails, it returns a custom error response.
//
// Returns:
//
//	fiber.Handler: A Fiber middleware handler.
func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the bearer token from the "Authorization" header.
		bearerToken := c.Get("Authorization")

		// Retrieve the AppState and context from the Fiber context.
		appState := c.Locals("appState").(*state.AppState)
		ctx := c.Locals("ctx").(context.Context)

		// Create a new UserRepository instance.
		userRepo := repositories.NewUserRepository(appState.DB)

		// Check if the bearer token is missing.
		if bearerToken == "" {
			// Return a custom error response indicating that the token is missing.
			return errors.NewCustomErrorResponse(c, errors.ErrTokenMissing)
		}

		// Remove the "Bearer " prefix from the token.
		token := strings.TrimPrefix(bearerToken, "Bearer ")

		// Verify the JWT token using the provided public key.
		verifyToken, err := services.VerifyJwtToken(appState.Config.AccessTokenPublicKey, token)
		if err != nil {
			// Return a custom error response indicating that token verification failed.
			return errors.NewCustomErrorResponse(c, errors.ErrTokenVerification)
		}

		// Retrieve the access token UUID from the verified token.
		accessTokenUuid := verifyToken.TokenUUID

		// Retrieve and compare the user ID from Redis.
		redisValue, err := utils.GetAndCompareRedisValue(c, appState.Redis, ctx, accessTokenUuid.String())
		if err != nil {
			// Return the error from GetAndCompareRedisValue, which already handles error responses.
			return err
		}

		// Parse the user ID from the Redis value.
		userID, err := uuid.Parse(redisValue)
		if err != nil {
			// Return a custom error response indicating that parsing the user ID failed.
			return errors.NewCustomErrorResponse(c, errors.ErrUserIDParse)
		}

		// Retrieve the user from the database using the user ID.
		user, err := userRepo.GetUserByID(userID)
		if err != nil {
			// Return a custom error response indicating that the user was not found.
			return errors.NewCustomErrorResponse(c, errors.ErrUserNotFound)
		}

		// Verify that the user ID from Redis matches the user ID from the database.
		if user.ID != userID {
			// Return a custom error response indicating a user ID mismatch.
			return errors.NewCustomErrorResponse(c, errors.ErrUserIDMismatch)
		}

		// Create a JWT middleware response.
		jwtResponse := responses.JwtMiddlewareResponse{
			Token: accessTokenUuid,
			User: responses.UserResponse{
				ID:         user.ID,
				Email:      *user.Email,
				Name:       *user.Name,
				CreatedAt:  user.CreatedAt,
				UpdatedAt:  user.UpdatedAt,
				Provider:   *user.Provider,
				ProviderID: *user.ProviderID,
			},
		}

		// Return a success response with the user information.
		return c.JSON(responses.SuccessResponse{
			Status: "success",
			Data:   jwtResponse,
		})
	}
}
