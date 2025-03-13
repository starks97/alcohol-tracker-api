package exceptions

import (
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidatePassword(password string) []string {
	passwordErrors := []string{}

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
