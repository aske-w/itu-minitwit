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

func (s *MessageService) CreateMessage(userId int, text string) (*models.Message, error) {
	pub_date := int(time.Now().Unix())
	message := models.Message{
		Author_id: userId,
		Text:      text,
		Pub_date:  pub_date,
	}
	err := s.DB.Create(&message).Error
	if err == nil {
		return nil, err
	}
	return &message, nil
}
