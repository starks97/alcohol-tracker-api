package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
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

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) CreateUser(user *User) error {
	return ur.db.Create(user).Error
}

func (ur *UserRepository) GetUserByID(id uuid.UUID) (*User, error) {
	var user User
	if err := ur.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
