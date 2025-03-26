package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"errors"

	"github.com/starks97/alcohol-tracker-api/internal/exceptions"
	"github.com/starks97/alcohol-tracker-api/internal/models"
	"github.com/starks97/alcohol-tracker-api/internal/repositories"
	"github.com/starks97/alcohol-tracker-api/internal/responses"
	"github.com/starks97/alcohol-tracker-api/internal/state"
	"github.com/starks97/alcohol-tracker-api/internal/utils"
)

func Register(c *fiber.Ctx) error {
	appState := c.Locals("appState").(*state.AppState)
	userQuery := repositories.NewUserRepository(appState.DB)

	var userDataFromReq models.RegisterUserSchema

	if err := c.BodyParser(&userDataFromReq); err != nil {
		return exceptions.HandlerErrorResponse(c, exceptions.ErrRequestBody)
	}

	if err := utils.ParseValidatorMessage(&userDataFromReq, appState.Validator); err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			return exceptions.HandlerValidationErrorResponse(c, exceptions.ErrValidationFailed, validationErr.Errors)
		}
		return exceptions.HandlerErrorResponse(c, err)
	}

	_, err := userQuery.GetUserByEmail(userDataFromReq.Email)
	if err == nil {
		return exceptions.HandlerErrorResponse(c, exceptions.ErrUserAlreadyExists)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return exceptions.HandlerErrorResponse(c, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDataFromReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return exceptions.HandlerErrorResponse(c, err)
	}
	passwordFromBytes := string(hashedPassword)
	userData := &repositories.User{
		Email:    userDataFromReq.Email,
		Name:     userDataFromReq.Name,
		Password: &passwordFromBytes,
	}

	createUser, err := userQuery.CreateUser(userData)
	if err != nil {
		fmt.Println("Error when you create user:", err)
		return exceptions.HandlerErrorResponse(c, exceptions.ErrUserNotCreated)
	}

	message := "User registered successfully"

	return c.JSON(responses.SuccessResponse{
		Status:  "success",
		Message: &message,
		Data:    createUser,
	})

}
