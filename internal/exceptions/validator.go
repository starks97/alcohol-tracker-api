package exceptions

import (
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidatePassword checks if a password meets certain complexity requirements.
//
// It returns a slice of strings containing error messages for any unmet requirements.
// The requirements include:
// - Password must not be empty.
// - Password length must be between 8 and 100 characters.
// - Password must contain at least one uppercase letter.
// - Password must contain at least one special character.
// - Password must contain at least one digit.
// - Password must contain at least one lowercase letter.
func ValidatePassword(password string) []string {
	passwordErrors := []string{}

	if password == "" {
		passwordErrors = append(passwordErrors, "Password must be provided, please provide a valid password.")
	}

	if len(password) < 8 || len(password) > 100 {
		passwordErrors = append(passwordErrors, "Password must be between 8 and 100 characters")
	}

	specialChars := "!@#$%^&*()_+-=\\"

	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		passwordErrors = append(passwordErrors, "Password must contain at least one uppercase letter")
	}

	if !strings.ContainsAny(password, specialChars) {
		passwordErrors = append(passwordErrors, "Password must contain at least one special character")
	}

	if !strings.ContainsAny(password, "0123456789") {
		passwordErrors = append(passwordErrors, "Password must contain at least one digit")
	}

	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		passwordErrors = append(passwordErrors, "Password must contain at least one lowercase letter")
	}

	return passwordErrors
}

// Init initializes a new validator instance and registers a custom "password" validation rule.
//
// The "password" rule uses the `ValidatePassword` function to check password validity.
// It returns the initialized validator instance.
// If the validator initialization fails, it logs a fatal error and exits.
func Init() *validator.Validate {
	validate := validator.New()
	if validate == nil {
		log.Fatal("Failed to initialize validator")
	}

	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return len(ValidatePassword(password)) == 0
	})

	return validate
}
