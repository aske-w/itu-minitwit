package services

import (
	"aske-w/itu-minitwit/models"
	"errors"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) FindByUsername(username string) (*models.User, error) {

	user := &models.User{}
	err := s.DB.First(&user, "username = ?", username).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}
