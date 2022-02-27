package services

import (
	"aske-w/itu-minitwit/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) CreateUser(username string, email string, password string) (uint, error) {

	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return 0, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Pw_Hash:  string(pwHash),
	}
	s.db.Create(user)
	return user.ID, nil
}

func (s *AuthService) CheckPassword(user *models.User, password string) bool {
	pwErr := bcrypt.CompareHashAndPassword([]byte(user.Pw_Hash), []byte(password))
	if pwErr == nil {
		return true
	} else {
		return false
	}
}
