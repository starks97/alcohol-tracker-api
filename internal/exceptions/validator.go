package exceptions

import (
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

func validatePassword(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if len(password) < 8 || len(password) > 100 {
		return false
	}

	hasUpper := false
	hasSpecial := false
	hasDigit := false
	hasLower := false

	specialChars := "!@#$%^&*()_+-=\\"

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if strings.ContainsRune(specialChars, char) {
			hasSpecial = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		} else if char >= 'a' && char <= 'z' {
			hasLower = true
		}
	}

	return hasUpper && hasSpecial && hasDigit && hasLower
}

func Init() *validator.Validate {
	validate := validator.New()
	if validate == nil {
		log.Fatal("Failed to initialize validator")
	}

	validate.RegisterValidation("password", validatePassword)

	return validate
}
