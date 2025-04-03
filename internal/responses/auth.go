package responses

import (
	"time"

	"github.com/google/uuid"
	"github.com/starks97/alcohol-tracker-api/internal/entities"
)

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Provider   string    `json:"provider"`
	ProviderID string    `json:"provider_id"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type JwtMiddlewareResponse struct {
	AccessToken uuid.UUID     `json:"access_token"`
	User        entities.User `json:"user"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type SuccessResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message *string     `json:"message,omitempty"`
}
