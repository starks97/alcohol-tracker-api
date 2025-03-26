package exceptions

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrTokenMissing        = fmt.Errorf("No authentication token found. Please log in to get a valid token and try again.")
	ErrTokenVerification   = fmt.Errorf("Your session has expired or the token is invalid. Please log in again to get a new token.")
	ErrRedisGet            = fmt.Errorf("We couldn't retrieve your authentication token. Please try logging in again.")
	ErrRedisNotFound       = fmt.Errorf("We couldn’t find your token. Please log in again to obtain a new one.")
	ErrTokenMismatch       = fmt.Errorf("The token does not match our records. Please log in again.")
	ErrUserNotFound        = fmt.Errorf("No user found with the provided information. Please check your input and try again.")
	ErrUserIDMismatch      = fmt.Errorf("You are not authorized to perform this action. Please check if you're logged in with the correct account.")
	ErrUserIDParse         = fmt.Errorf("The user ID you entered is not valid. Please check your input and try again.")
	ErrUserAlreadyExists   = fmt.Errorf("A user with this email or username already exists. Please choose a different one.")
	ErrUserNotExists       = fmt.Errorf("No user found with the provided details. Please check your input and try again.")
	ErrUserNotCreated      = fmt.Errorf("We couldn't create your account. Please try again later or contact support.")
	ErrUserNotUpdated      = fmt.Errorf("We couldn't update your account. Please try again later or contact support.")
	ErrTokenNotGenerated   = fmt.Errorf("We couldn't generate a token. Please try logging in again.")
	ErrRedisSet            = fmt.Errorf("We couldn't save your session. Please log in again.")
	ErrExchangeToken       = fmt.Errorf("We couldn’t exchange your token. Please try again later or contact support.")
	ErrToReadUserInfo      = fmt.Errorf("We couldn't retrieve your user information. Please try again later or contact support.")
	ErrToUnmarshalUserInfo = fmt.Errorf("We couldn't process your user information. Please try again later or contact support.")
	ErrInvalidCredentials  = fmt.Errorf("Invalid email or password. Please verify your credentials and try again.")
	ErrRequestBody         = fmt.Errorf("The request body is invalid. Please check the request and try again.")
	ErrDatabase            = fmt.Errorf("A database error occurred. Please try again or contact support.")
	ErrPasswordRequired    = fmt.Errorf("Password is required. Please provide a valid password.")
	ErrValidationFailed    = fmt.Errorf("Validation failed. Please check your input and try again.")
)

// ErrorMapping maps error types to HTTP status codes.
var ErrorMapping = map[error]struct {
	StatusCode int
}{
	ErrTokenMissing:        {http.StatusUnauthorized},
	ErrTokenVerification:   {http.StatusUnauthorized},
	ErrRedisGet:            {http.StatusUnauthorized},
	ErrRedisNotFound:       {http.StatusUnauthorized},
	ErrTokenMismatch:       {http.StatusUnauthorized},
	ErrUserIDParse:         {http.StatusUnauthorized},
	ErrUserNotFound:        {http.StatusNotFound},
	ErrUserIDMismatch:      {http.StatusConflict},
	ErrUserNotExists:       {http.StatusNotFound},
	ErrUserNotCreated:      {http.StatusInternalServerError},
	ErrUserNotUpdated:      {http.StatusInternalServerError},
	ErrUserAlreadyExists:   {http.StatusConflict},
	ErrTokenNotGenerated:   {http.StatusInternalServerError},
	ErrRedisSet:            {http.StatusInternalServerError},
	ErrExchangeToken:       {http.StatusInternalServerError},
	ErrToReadUserInfo:      {http.StatusInternalServerError},
	ErrToUnmarshalUserInfo: {http.StatusInternalServerError},
	ErrInvalidCredentials:  {http.StatusUnauthorized},
	ErrRequestBody:         {http.StatusBadRequest},
	ErrDatabase:            {http.StatusInternalServerError},
	ErrPasswordRequired:    {http.StatusBadRequest},
	ErrValidationFailed:    {http.StatusBadRequest},
}

// ErrorResponse represents a JSON error response.
type ErrorResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Errors  *map[string][]string `json:"errors"`
}

// HandlerErrorResponse creates a custom error response for Fiber.
// It maps application-specific errors and standard library errors to appropriate HTTP status codes.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context for sending the error response.
//   - err: error - The error to handle.
//
// Returns:
//   - error: An error indicating that sending the error response failed, or nil if successful.
func HandlerErrorResponse(c *fiber.Ctx, err error) error {
	// Check if the error is in the application-specific error mapping.
	if errInfo, ok := ErrorMapping[err]; ok {
		log.Println("Error:", err)
		return c.Status(errInfo.StatusCode).JSON(ErrorResponse{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	// Check if the error is a Redis "not found" error.
	if errors.Is(err, redis.Nil) {
		return c.Status(http.StatusNotFound).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Resource not found in Redis",
		})
	}

	// Check if the error is a GORM "record not found" error.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Resource not found in database",
		})
	}

	// Handle unexpected errors by logging the error and returning an internal server error.
	log.Printf("Unexpected Error: %+v\n", err) // Include stack trace
	return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
		Status:  "failed",
		Message: "Internal server error",
	})
}

// HandlerValidationErrorResponse handles validation errors and generates appropriate HTTP responses.
//
// It takes a Fiber context, an error, and a map of validation errors.
// It checks if the error is mapped to a custom error response in `ErrorMapping`.
// If a mapping exists, it logs the error and returns a JSON response with the mapped status code and error details.
// If the error is a `validator.ValidationErrors`, it returns a 400 Bad Request response with a generic "Validation failed" message.
// If the error does not match any of the above conditions, it returns nil, allowing the caller to handle the error further.
//
// Example usage:
//
//	if err := someValidationFunction(); err != nil {
//		if validationErr, ok := err.(*validation.ValidationError); ok {
//			return HandlerValidationErrorResponse(c, err, validationErr.Errors)
//		}
//		// handle other errors
//		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
//	}
func HandlerValidationErrorResponse(c *fiber.Ctx, err error, validationErrors map[string][]string) error {
	if errorInfo, ok := ErrorMapping[err]; ok {
		log.Println("Error:", err)
		return c.Status(errorInfo.StatusCode).JSON(ErrorResponse{
			Status:  "failed",
			Message: err.Error(),
			Errors:  &validationErrors,
		})
	}

	if errors.Is(err, validator.ValidationErrors{}) {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Validation failed",
		})
	}
	return nil
}
