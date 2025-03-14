package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/go-playground/validator/v10"
	"github.com/starks97/alcohol-tracker-api/internal/exceptions"
)

type Validator interface {
	Validate(*validator.Validate) error // Pass the validator instance
}

var errorMessages = map[string]string{
	"required": "Please provide a value for {0}.",
	"email":    "Please enter a valid email address for {0}.",
	"min":      "{0} must be at least {1} characters.",
	"max":      "{0} cannot exceed {1} characters.",
	"password": "{0} error in password.",
}

func ParseValidatorMessage(c *fiber.Ctx, model Validator, validatorClient *validator.Validate) error {

	if valid := model.Validate(validatorClient); valid != nil {
		errs, ok := valid.(validator.ValidationErrors)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Validation failed",
			})
		}

		errors := make(map[string][]string)
		for _, e := range errs {
			tag := e.Tag()
			field := e.Field()
			param := e.Param()
			message, found := errorMessages[tag]
			if !found {
				message = "Validation failed for " + field // Default message
			}

			if tag == "password" {
				password, ok := e.Value().(string)
				if !ok {
					errors[field] = append(errors[field], "Invalid password type")
					continue
				}

				passwordErrs := exceptions.ValidatePassword(password)
				errors[field] = append(errors[field], passwordErrs...)
				if len(passwordErrs) > 0 {
					continue
				}
			}

			message = strings.Replace(message, "{0}", field, 1)
			message = strings.Replace(message, "{1}", param, 1)

			errors[field] = append(errors[field], message)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	return nil
}
