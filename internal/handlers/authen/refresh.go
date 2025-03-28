package authen

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/starks97/alcohol-tracker-api/internal/exceptions"
	"github.com/starks97/alcohol-tracker-api/internal/repositories"
	"github.com/starks97/alcohol-tracker-api/internal/responses"
	"github.com/starks97/alcohol-tracker-api/internal/services"
	"github.com/starks97/alcohol-tracker-api/internal/state"
	"github.com/starks97/alcohol-tracker-api/internal/utils"
)

func RefreshTokenHandler(c *fiber.Ctx) error {
	appState := c.Locals("appState").(*state.AppState)
	ctx := c.Locals("ctx").(context.Context)
	userRepo := repositories.NewUserRepository(appState.DB)
	tokenService := utils.NewTokenService(appState)

	reCookie := c.Cookies("refresh_token")

	if reCookie == "" {
		return exceptions.HandlerErrorResponse(c, exceptions.ErrTokenMissing)
	}

	verifyToken, err := services.VerifyJwtToken(appState.Config.RefreshTokenPublicKey, reCookie)
	if err != nil {
		return exceptions.HandlerErrorResponse(c, exceptions.ErrTokenVerification)
	}

	accessTokenUuid := verifyToken.TokenUUID

	redisValue, err := tokenService.GetAndCompareRedisValue(c, appState.Redis, ctx, accessTokenUuid.String())
	if err != nil {
		// Return the error from GetAndCompareRedisValue, which already handles error responses.
		return err
	}

	userID, err := uuid.Parse(redisValue)
	if err != nil {
		return exceptions.HandlerErrorResponse(c, exceptions.ErrUserIDParse)
	}

	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		// Return a custom error response indicating that the user was not found.
		return exceptions.HandlerErrorResponse(c, exceptions.ErrUserNotFound)
	}

	// Verify that the user ID from Redis matches the user ID from the database.
	if user.ID != userID {
		// Return a custom error response indicating a user ID mismatch.
		return exceptions.HandlerErrorResponse(c, exceptions.ErrUserIDMismatch)
	}

	tokenResult, err := tokenService.StoreToken(c, ctx, user.ID, "access")
	if err != nil {
		return err
	}

	accessToken := responses.LoginResponse{
		AccessToken: *tokenResult.Token,
	}

	return c.JSON(responses.SuccessResponse{
		Status: "success",
		Data:   accessToken,
	})

}
