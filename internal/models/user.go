package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID                   string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email                *string   `gorm:"uniqueIndex;size:255" json:"email,omitempty"` // Nullable
	Password             *string   `gorm:"size:255" json:"-"`
	Name                 *string   `gorm:"size:255" json:"name,omitempty"`
	Provider             *string   `gorm:"size:255" json:"provider,omitempty"`
	ProviderID           *string   `gorm:"size:255;uniqueIndex" json:"provider_id,omitempty"`
	ProfilePicture       *string   `gorm:"size:255" json:"profile_picture,omitempty"`
	ProviderRefreshToken *string   `gorm:"size:255" json:"-"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type RegisterUserSchema struct {
	Email                string    `json:"email" db:"email" validate:"required,email"`
	Password             *string   `json:"password,omitempty" db:"password" validate:"omitempty,min=8,max=100,regexp=^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[!@#$%^&*()_+\\-=]).{8,100}$"`
	Name                 string    `json:"name" db:"name" validate:"required,min=2,max=50"`
	Provider             *string   `json:"provider,omitempty" db:"provider"`
	ProviderID           *string   `json:"provider_id,omitempty" db:"provider_id"`
	ProfilePicture       *string   `json:"profile_picture,omitempty" db:"profile_picture"`
	ProviderRefreshToken *string   `json:"provider_refresh_token,omitempty" db:"provider_refresh_token" `
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

type LoginUserSchema struct {
	Email    string  `json:"email" db:"email" validate:"required,email"`
	Password *string `json:"password,omitempty" db:"password" validate:"omitempty,min=8,max=100,regexp=^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[!@#$%^&*()_+\\-=]).{8,100}$"`
}

type UpdateUserSchema struct {
	Name                 *string `json:"name,omitempty" db:"name" validate:"omitempty,min=2,max=50"`
	ProfilePicture       *string `json:"profile_picture,omitempty" db:"profile_picture"`
	Password             *string `json:"password,omitempty" db:"password" validate:"omitempty,min=8,max=100,regexp=^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[!@#$%^&*()_+\\-=]).{8,100}$"`
	Provider             *string `json:"provider,omitempty" db:"provider"`
	ProviderID           *string `json:"provider_id,omitempty" db:"provider_id"`
	ProviderRefreshToken *string `json:"provider_refresh_token,omitempty" db:"provider_refresh_token"`
}

func (u *RegisterUserSchema) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *LoginUserSchema) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UpdateUserSchema) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
