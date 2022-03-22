package services

import (
	"aske-w/itu-minitwit/models"
	"time"

	"gorm.io/gorm"
)

type MessageService struct {
	DB *gorm.DB
}

func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{DB: db}
}

func (s *MessageService) CreateMessage(userId int, text string) error {
	message := models.Message{
		Author_id: userId,
		Text:      text,
		Pub_date:  int(time.Now().Unix()),
	}

	result := s.DB.Create(&message)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
