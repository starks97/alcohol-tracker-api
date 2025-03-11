package exceptions

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrTokenMissing        = fmt.Errorf("Token not provided or not found")
	ErrTokenVerification   = fmt.Errorf("Failed to verify token")
	ErrRedisGet            = fmt.Errorf("Failed to get token from redis, please provide the correct token")
	ErrRedisNotFound       = fmt.Errorf("Token not found in redis, please provide the correct token")
	ErrTokenMismatch       = fmt.Errorf("The token provided does not match the expected token")
	ErrUserNotFound        = fmt.Errorf("User could not be found with the provided ID, please provide a valid ID")
	ErrUserIDMismatch      = fmt.Errorf("User ID mismatch")
	ErrUserIDParse         = fmt.Errorf("Failed to parse user ID")
	ErrUserAlreadyExists   = fmt.Errorf("User already exists")
	ErrUserNotExists       = fmt.Errorf("User does not exist")
	ErrUserNotCreated      = fmt.Errorf("User could not be created")
	ErrTokenNotGenerated   = fmt.Errorf("Failed to generate token")
	ErrRedisSet            = fmt.Errorf("Failed to set token in redis")
	ErrExchangeToken       = fmt.Errorf("Failed to exchange token")
	ErrToReadUserInfo      = fmt.Errorf("Failed to read user info")
	ErrToUnmarshalUserInfo = fmt.Errorf("Failed to unmarshal user info")
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
	ErrUserAlreadyExists:   {http.StatusConflict},
	ErrTokenNotGenerated:   {http.StatusInternalServerError},
	ErrRedisSet:            {http.StatusInternalServerError},
	ErrExchangeToken:       {http.StatusInternalServerError},
	ErrToReadUserInfo:      {http.StatusInternalServerError},
	ErrToUnmarshalUserInfo: {http.StatusInternalServerError},
}

// ErrorResponse represents a JSON error response.
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// NewCustomErrorResponse creates a custom error response for Fiber.
// It maps application-specific errors and standard library errors to appropriate HTTP status codes.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context for sending the error response.
//   - err: error - The error to handle.
//
// Returns:
//   - error: An error indicating that sending the error response failed, or nil if successful.
func NewCustomErrorResponse(c *fiber.Ctx, err error) error {
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
