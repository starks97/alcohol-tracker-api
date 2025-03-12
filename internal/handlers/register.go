package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/starks97/alcohol-tracker-api/internal/models"
	"github.com/starks97/alcohol-tracker-api/internal/state"
	"github.com/starks97/alcohol-tracker-api/internal/utils"
)

func Register(c *fiber.Ctx) error {
	appState := c.Locals("appState").(*state.AppState)

	user := new(models.RegisterUserSchema)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := utils.ParseValidatorMessage(c, user, appState.Validator); err != nil {
		return err
	}

	return nil

}
