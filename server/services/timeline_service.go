package services

import (
	"aske-w/itu-minitwit/models"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Tweet struct {
	UserId          int
	Username        string
	Email           string
	Message_id      int
	Author_id       int
	Text            string
	Pub_date        int
	Flagged         int
	Gravatar_Url    string
	Format_Datetime string
}

type TimelineService struct {
	DB *gorm.DB
}

func NewTimelineService(db *gorm.DB) *TimelineService {

	return &TimelineService{DB: db}
}

func (s *TimelineService) GetPublicTimeline() (*[]Tweet, error) {

	tweets := []Tweet{}
	err := s.DB.Model(&models.User{}).Select("users.id as UserId", "users.Username", "users.Email", "messages.id", "messages.Author_id", "messages.Text", "messages.Pub_date", "messages.Flagged").Joins("INNER JOIN messages ON messages.author_id = users.id AND messages.flagged = 0").Order("messages.pub_date DESC").Limit(30).Scan(&tweets).Error

	if err != nil {
		return nil, err
	}
	for i, tweet := range tweets {
		tweets[i].Gravatar_Url = gravatar_url(tweet.Email, 48)
		tweets[i].Format_Datetime = format_datetime(tweet.Pub_date)
	}

	return &tweets, nil
}

func (s *TimelineService) GetUserTimeline(userId int) (*[]Tweet, error) {

	// rows, err := c.DB.Raw(`
	// select  user.*, message.* from message, user where
	// user.user_id = message.author_id and user.user_id = ?
	// order by message.pub_date desc limit ?`, userId, 30).Rows()

	tweets := []Tweet{}
	err := s.DB.Model(&models.User{}).Where("users.id = ?", userId).Select("users.id as UserId", "users.Username", "users.Email", "messages.id", "messages.Author_id", "messages.Text", "messages.Pub_date", "messages.Flagged").Joins("INNER JOIN messages ON messages.author_id = users.id AND messages.flagged = 0").Order("messages.pub_date DESC").Limit(30).Scan(&tweets).Error

	if err != nil {
		return nil, err
	}
	for i, tweet := range tweets {
		tweets[i].Gravatar_Url = gravatar_url(tweet.Email, 48)
		tweets[i].Format_Datetime = format_datetime(tweet.Pub_date)
	}

	return &tweets, nil
}

func format_datetime(timestamp int) string {
	unix := time.Unix(int64(timestamp), 0)
	return unix.Format("2006-01-02T15:04:05Z07:00")
}

//     """Return the gravatar image for the given email address."""
func gravatar_url(email string, size int) string {
	stripped := strings.Trim(email, "")
	lowered := strings.ToLower(stripped)
	valid := strings.ToValidUTF8(lowered, "")

	hasher := md5.New()
	data := []byte(valid)
	hasher.Write(data)
	md5Email := hex.EncodeToString(hasher.Sum(nil))

	url := fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", md5Email, size)
	return url
}
