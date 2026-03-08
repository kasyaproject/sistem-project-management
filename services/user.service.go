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
