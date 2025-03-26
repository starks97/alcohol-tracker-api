package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/starks97/alcohol-tracker-api/internal/exceptions"
)

// Validator is an interface that models should implement to provide custom validation logic.
type Validator interface {
	Validate(*validator.Validate) error
}

// ValidationError represents validation errors, mapping field names to lists of error messages.
type ValidationError struct {
	Errors map[string][]string
}

// Error returns a string representation of the validation errors.
// It formats the errors as "field: message, field: message, ...".
func (ve *ValidationError) Error() string {
	var errorMessages []string
	for field, messages := range ve.Errors {
		for _, msg := range messages {
			errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", field, msg))
		}
	}
	return strings.Join(errorMessages, ", ")
}

// errorMessages maps validation tags to human-readable error messages.
// The messages can include placeholders like "{0}" for the field name and "{1}" for parameters.
var errorMessages = map[string]string{
	"required": "Please provide a value for {0}.",
	"name":     "Please enter a valid name for {0}.",
	"email":    "Please enter a valid email address for {0}.",
	"min":      "{0} must be at least {1} characters.",
	"max":      "{0} cannot exceed {1} characters.",
	"password": "{0} error in password.",
}

// ParseValidatorMessage validates a model using the provided validator client and parses the errors.
//
// It takes a Validator interface and a validator.Validate instance.
// If the model's Validate method returns an error, it attempts to parse the error as validator.ValidationErrors.
// It then iterates through the errors, looks up corresponding error messages in `errorMessages`,
// and formats the messages with field names and parameters.
//
// For "password" validation tag, it also executes additional validations from your exceptions package and appends those errors.
//
// It returns a ValidationError if any validation errors occur, or nil if validation succeeds.
func ParseValidatorMessage(model Validator, validatorClient *validator.Validate) error {

	if err := model.Validate(validatorClient); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
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
				if !ok || password == "" {
					errors[field] = append(errors[field], "Password must be provided, please provide a valid password.")
					continue
				}

				passwordErrs := exceptions.ValidatePassword(password)
				errors[field] = append(errors[field], passwordErrs...)
				if len(passwordErrs) > 0 {
					continue
				}
			}

			message = strings.Replace(message, "{0}", field, 1)
			if param != "" {
				message = strings.Replace(message, "{1}", param, 1)
			}

			errors[field] = append(errors[field], message)
		}
		if len(errors) > 0 {
			return &ValidationError{
				Errors: errors,
			}
		}
	}

	return nil
}
