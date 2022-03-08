package services

import (
	"aske-w/itu-minitwit/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) CreateUser(username string, email string, password string) error {

	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return err
	}

	creationTime := time.Now()
	s.db.Exec(`INSERT INTO users (created_at,updated_at,deleted_at,username,email,pw_hash) VALUES (?,?,NULL,?,?,?)`, creationTime, creationTime, username, email, string(pwHash))

	return nil
}

func (s *AuthService) CheckPassword(user *models.User, password string) bool {
	pwErr := bcrypt.CompareHashAndPassword([]byte(user.Pw_Hash), []byte(password))
	if pwErr == nil {
		return true
	} else {
		return false
	}
}
