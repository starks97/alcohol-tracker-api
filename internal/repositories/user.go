package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the application.
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

// UserRepository provides methods to interact with the User model in the database.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance.
//
// Parameters:
//   - db: *gorm.DB - The GORM database connection.
//
// Returns:
//   - *UserRepository: A new UserRepository instance.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database.
//
// Parameters:
//   - user: *User - The user to create.
//
// Returns:
//   - error: An error if the creation fails.
func (ur *UserRepository) CreateUser(user *User) (*User, error) {
	if err := ur.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user from the database by ID.
//
// Parameters:
//   - id: uuid.UUID - The user ID.
//
// Returns:
//   - *User: The user if found, or nil if not found.
//   - error: An error if the retrieval fails.
func (ur *UserRepository) GetUserByID(id uuid.UUID) (*User, error) {
	var user User
	if err := ur.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByProvider retrieves a user from the database by provider and provider ID.
//
// Parameters:
//   - provider: string - The authentication provider (e.g., "google").
//   - providerID: string - The provider's user ID.
//
// Returns:
//   - *User: The user if found, or nil if not found.
//   - error: An error if the retrieval fails.
func (usr *UserRepository) GetUserByProvider(provider string, providerID string) (*User, error) {
	var user User
	if err := usr.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user from the database by email.
//
// Parameters:
//   - email: string - The user's email address.
//
// Returns:
//   - *User: The user if found, or nil if not found.
//   - error: An error if the retrieval fails.
func (usr *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := usr.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
