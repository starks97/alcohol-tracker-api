package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(*validator.Validate) error // Pass the validator instance
}

var errorMessages = map[string]string{
	"required": "Please provide a value for {0}.",
	"email":    "Please enter a valid email address for {0}.",
	"min":      "{0} must be at least {1} characters.",
	"max":      "{0} cannot exceed {1} characters.",
	"regexp":   "{0} does not match the required format.",
	"password": "{0} must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character.",
}

func ParseValidatorMessage(c *fiber.Ctx, model Validator, validatorClient *validator.Validate) error {

	if valid := model.Validate(validatorClient); valid != nil {
		errs, ok := valid.(validator.ValidationErrors)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Validation failed",
			})
		}

		errors := make(map[string]string)
		for _, e := range errs {
			tag := e.Tag()
			field := e.Field()
			param := e.Param()
			message, found := errorMessages[tag]
			if !found {
				message = "Validation failed for " + field // Default message
			}
			message = strings.Replace(message, "{0}", field, 1)
			message = strings.Replace(message, "{1}", param, 1)

			errors[field] = message
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	return nil
}
