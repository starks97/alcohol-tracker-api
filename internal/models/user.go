package models

import (
	"github.com/go-playground/validator/v10"
)

type RegisterUserSchema struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password" validate:"password"`
	Name     string `json:"name" db:"name" validate:"required,min=2,max=50"`
}

type LoginUserSchema struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password" validate:"password"`
}

type UpdateUserSchema struct {
	Name                 *string `json:"name,omitempty" db:"name" validate:"omitempty,min=2,max=50"`
	ProfilePicture       *string `json:"profile_picture,omitempty" db:"profile_picture"`
	Password             *string `json:"password,omitempty" db:"password" validate:"password"`
	Provider             *string `json:"provider,omitempty" db:"provider"`
	ProviderID           *string `json:"provider_id,omitempty" db:"provider_id"`
	ProviderRefreshToken *string `json:"provider_refresh_token,omitempty" db:"provider_refresh_token"`
}

func (u *RegisterUserSchema) Validate(v *validator.Validate) error {
	return v.Struct(u)
}

func (u *LoginUserSchema) Validate(v *validator.Validate) error {
	return v.Struct(u)
}

func (u *UpdateUserSchema) Validate(v *validator.Validate) error {
	return v.Struct(u)
}
