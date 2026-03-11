package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/kasyaproject/sistem-project-management/models"
	"github.com/kasyaproject/sistem-project-management/repositories"
	"github.com/kasyaproject/sistem-project-management/utils"
)

type UserService interface {
	Register(user *models.User) error
	Login(email, password string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByPublicID(id string) (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(user *models.User) error {
	// Cek email terdaftar
	existingUser, _ := s.repo.FindByEmail(user.Email)
	if existingUser.InternalID != 0 {
		return errors.New("Email already registered!")
	}

	// Hash password
	hased, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hased      // Set Hash password
	user.Role = "user"         // Set role
	user.PublicID = uuid.New() // Set UUID

	// Simpan user
	return s.repo.Create(user)
}

func (s *userService) Login(email, password string) (*models.User, error) {
	// Cek email terdaftar
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("Invalide Credential")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("Invalide Credential")
	}

	return user, nil
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) GetByPublicID(id string) (*models.User, error) {
	return s.repo.FindByPublicID(id)
}
