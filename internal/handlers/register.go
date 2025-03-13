package handlers

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"errors"

	"github.com/starks97/alcohol-tracker-api/internal/models"
	"github.com/starks97/alcohol-tracker-api/internal/repositories"
	"github.com/starks97/alcohol-tracker-api/internal/responses"
	"github.com/starks97/alcohol-tracker-api/internal/state"
	"github.com/starks97/alcohol-tracker-api/internal/utils"
)

func Register(c *fiber.Ctx) error {
	appState := c.Locals("appState").(*state.AppState)
	userQuery := repositories.NewUserRepository(appState.DB)

	userDataFromReq := new(models.RegisterUserSchema)

	if err := c.BodyParser(userDataFromReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := utils.ParseValidatorMessage(c, userDataFromReq, appState.Validator); err != nil {
		return err
	}

	_, err := userQuery.GetUserByEmail(userDataFromReq.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}

	} else {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*userDataFromReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	passwordFromBytes := string(hashedPassword)

	var userData = &repositories.User{
		Email:    &userDataFromReq.Email,
		Password: &passwordFromBytes,
		Name:     &userDataFromReq.Name,
	}

	createUser, err := userQuery.CreateUser(userData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	message := "User registered successfully"

	return c.JSON(responses.SuccessResponse{
		Status:  "success",
		Message: &message,
		Data:    createUser,
	})

}
