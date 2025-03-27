package repositories

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/starks97/alcohol-tracker-api/internal/entities"
)

type UserRepository interface {
	CreateUser(user *entities.User) (*entities.User, error)
	GetUserByID(id uuid.UUID) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByProvider(provider string, providerID string) (*entities.User, error)
	UpdateUser(user *entities.User) (*entities.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (usr *userRepository) Validate(validatorClient *validator.Validate) error {
	return validatorClient.Struct(usr)
}

func (usr *userRepository) CreateUser(user *entities.User) (*entities.User, error) {

	if err := usr.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (usr *userRepository) GetUserByID(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	if err := usr.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (usr *userRepository) GetUserByProvider(provider string, providerID string) (*entities.User, error) {
	var user entities.User
	if err := usr.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (usr *userRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := usr.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (usr *userRepository) UpdateUser(user *entities.User) (*entities.User, error) {
	result := usr.db.Model(&entities.User{}).
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
