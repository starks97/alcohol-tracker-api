package repositories

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the application.
type User struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email                string    `gorm:"uniqueIndex;size:255;not null" json:"email" validate:"required,email"` // Nullable
	Password             *string   `gorm:"size:255" json:"password,omitempty" validate:"password"`
	Name                 string    `gorm:"size:255;not null" json:"name" validate:"required,min=2,max=50"`
	Provider             *string   `gorm:"size:255" json:"provider,omitempty"`
	ProviderID           *string   `gorm:"size:255;uniqueIndex" json:"provider_id,omitempty"`
	ProfilePicture       *string   `gorm:"size:255" json:"profile_picture,omitempty"`
	ProviderRefreshToken *string   `gorm:"size:255" json:"provider_refresh_token,omitempty"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// UserRepository provides methods to interact with the User model in the database.
type UserRepository struct {
	db *gorm.DB
	//validator *validator.Validate
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

func (u *User) Validate(validatorClient *validator.Validate) error {
	return validatorClient.Struct(u)
}

/*func (ur *User) BeforeCreate(tx *gorm.DB) error {
validatorClient, ok := tx.Statement.Context.Value("validator").(*validator.Validate)
if !ok || validatorClient == nil {
	return errors.New("validator not found in transaction context")
}

if err := ur.Validate(validatorClient); err != nil {
	return err
}

return nil
}*/

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

func (usr *UserRepository) UpdateUser(user *User) (*User, error) {
	result := usr.db.Model(&User{}).
		Where("email = ?", user.Email). // You can change this to another unique identifier
		Updates(map[string]interface{}{
			"provider":               user.Provider,
			"provider_id":            user.ProviderID,
			"profile_picture":        user.ProfilePicture,
			"provider_refresh_token": user.ProviderRefreshToken,
			"name":                   user.Name,
		})

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
