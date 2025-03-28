package authen

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/starks97/alcohol-tracker-api/internal/dtos"
	"github.com/starks97/alcohol-tracker-api/internal/exceptions"
	"github.com/starks97/alcohol-tracker-api/internal/repositories"
	"github.com/starks97/alcohol-tracker-api/internal/responses"
	"github.com/starks97/alcohol-tracker-api/internal/state"
	"github.com/starks97/alcohol-tracker-api/internal/utils"
)

func LoginHandler(c *fiber.Ctx) error {
	appState := c.Locals("appState").(*state.AppState)

	userRepo := repositories.NewUserRepository(appState.DB)

	tokenService := utils.NewTokenService(appState)

	var userDataFromReq dtos.LoginUserDto

	ctx := c.Locals("ctx").(context.Context)

	if err := c.BodyParser(&userDataFromReq); err != nil {
		return exceptions.HandlerErrorResponse(c, err)
	}

	if err := utils.ParseValidatorMessage(&userDataFromReq, appState.Validator); err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			return exceptions.HandlerValidationErrorResponse(c, exceptions.ErrValidationFailed, validationErr.Errors)
		}
		return exceptions.HandlerErrorResponse(c, err)
	}

	userInDB, err := userRepo.GetUserByEmail(userDataFromReq.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.HandlerErrorResponse(c, exceptions.ErrDatabase)

		} else {
			return exceptions.HandlerErrorResponse(c, exceptions.ErrUserNotFound)
		}

	}

	if err := bcrypt.CompareHashAndPassword([]byte(*userInDB.Password), []byte(userDataFromReq.Password)); err != nil {
		return exceptions.HandlerErrorResponse(c, exceptions.ErrInvalidCredentials)
	}

	tokenResult, err := tokenService.StoreToken(c, ctx, userInDB.ID, "both")
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
