package authen

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/starks97/alcohol-tracker-api/internal/exceptions"
	"github.com/starks97/alcohol-tracker-api/internal/responses"
	"github.com/starks97/alcohol-tracker-api/internal/services"
	"github.com/starks97/alcohol-tracker-api/internal/state"
	"github.com/starks97/alcohol-tracker-api/internal/utils"
)

func LogOutHandler(c *fiber.Ctx) error {
	userData := c.Locals("mdlData").(*responses.JwtMiddlewareResponse)

	appState := c.Locals("appState").(*state.AppState)
	ctx := c.Locals("ctx").(context.Context)

	tokenService := utils.NewTokenService(appState)

	reCookie := c.Cookies("refresh_token")

	if reCookie == "" {
		return exceptions.HandlerErrorResponse(c, exceptions.ErrTokenMissing)
	}

	tokenDetail, err := services.VerifyJwtToken(appState.Config.RefreshTokenPublicKey, reCookie)
	if err != nil {
		return exceptions.HandlerErrorResponse(c, exceptions.ErrTokenVerification)
	}

	tokenService.RemoveRedisKeys(c, appState.Redis, ctx, userData.AccessToken.String(), tokenDetail.TokenUUID.String())

	c.ClearCookie("refresh_token")
	c.ClearCookie("access_token")

	message := "you are log out, have a nice day"
	return c.JSON(responses.SuccessResponse{
		Status:  "success",
		Message: &message,
	})
}
