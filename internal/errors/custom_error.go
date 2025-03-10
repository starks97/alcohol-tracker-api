package errors

import (
	"errors"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var (
	ErrTokenMissing      = errors.New("Token not provided or not found")
	ErrTokenVerification = errors.New("Failed to verify token")
	ErrRedisGet          = errors.New("Failed to get token from redis, please provide the correct token")
	ErrRedisNotFound     = errors.New("Token not found in redis, please provide the correct token")
	ErrTokenMismatch     = errors.New("The token provided does not match the expected token")
	ErrUserNotFound      = errors.New("User could not be found with the provided ID, please provide a valid ID")
	ErrUserIDMismatch    = errors.New("User ID mismatch")
	ErrUserIDParse       = errors.New("Failed to parse user ID")
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var ErrorMapping = map[error]struct {
	StatusCode int
}{
	ErrTokenMissing:      {http.StatusUnauthorized},
	ErrTokenVerification: {http.StatusUnauthorized},
	ErrRedisGet:          {http.StatusUnauthorized},
	ErrRedisNotFound:     {http.StatusUnauthorized},
	ErrTokenMismatch:     {http.StatusUnauthorized},
	ErrUserIDParse:       {http.StatusUnauthorized},
	ErrUserNotFound:      {http.StatusNotFound},
	ErrUserIDMismatch:    {http.StatusConflict},
}

func NewCustomErrorResponse(c *fiber.Ctx, err error) error {
	if errInfo, ok := ErrorMapping[err]; ok {
		log.Println("Error:", err)
		return c.Status(errInfo.StatusCode).JSON(ErrorResponse{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	if errors.Is(err, redis.Nil) {
		return c.Status(http.StatusNotFound).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Resource not found in Redis",
		})
	}

	log.Println("Unexpected Error:", err)
	return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
		Status:  "failed",
		Message: "Internal server error",
	})
}
