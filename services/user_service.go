package services

import (
    "github.com/alesio/gestion-actividades-deportivas/models"
    "gorm.io/gorm"
)

// UserService manages CRUD logic for users.
type UserService struct {
    db *gorm.DB
}

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
