package services

import (
	"errors"

	"github.com/alesio/gestion-actividades-deportivas/models"
	"github.com/alesio/gestion-actividades-deportivas/security"
	"gorm.io/gorm"
)

// UserService manages CRUD logic for users.
type UserService struct {
	db *gorm.DB
}

var ErrEmailAlreadyExists = errors.New("email already exists")

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) CreateUser(name, email, password, role string) (*models.User, error) {
	var count int64
	if err := s.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         role,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
