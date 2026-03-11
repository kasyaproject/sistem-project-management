package repositories

import (
	"github.com/kasyaproject/sistem-project-management/config"
	"github.com/kasyaproject/sistem-project-management/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByPublicID(publicID string) (*models.User, error)
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

// Func untuk Create user
func (r *userRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

// Func untuk Find user by Email
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error

	return &user, err
}

// Func untuk Find user by Internal ID
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User

	err := config.DB.First(&user, id).Error

	return &user, err
}

// Func untuk Find user by Public ID
func (r *userRepository) FindByPublicID(publicID string) (*models.User, error) {
	var user models.User

	err := config.DB.Where("public_id = ?", publicID).First(&user).Error

	return &user, err
}
