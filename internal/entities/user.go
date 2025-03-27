package entities

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the application.
type User struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email                string    `gorm:"uniqueIndex;size:255;not null" validate:"required,email"` // Nullable
	Password             *string   `gorm:"size:255" validate:"password"`
	Name                 string    `gorm:"size:255;not null" validate:"required,min=2,max=50"`
	Provider             *string   `gorm:"size:255"`
	ProviderID           *string   `gorm:"size:255;uniqueIndex"`
	ProfilePicture       *string   `gorm:"size:255"`
	ProviderRefreshToken *string   `gorm:"size:255"`
	CreatedAt            time.Time `gorm:"autoCreateTime"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`
}
