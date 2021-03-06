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

func (s *UserService) GetById(userId int) (*models.User, error) {

	user := &models.User{
		Model: gorm.Model{
			ID: uint(userId),
		},
	}
	err := s.DB.First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserService) UsernameToId(username string) (int, error) {

	id := -1
	err := s.DB.Table("users").Where("username = ?", username).Select("id").Scan(&id).Error
	if err != nil {
		return id, err
	}

	return id, nil
}

func (s *UserService) CheckUsernameExists(username string) (bool, error) {
	var exists bool
	err := s.DB.Model(&models.User{}).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&exists).
		Error
	if err != nil {
		return exists, err
	}
	return exists, err
}

func (s *UserService) FindByUsername(username string) (*models.User, error) {

	user := &models.User{}
	err := s.DB.First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UserIsFollowing(userId, userToFollowId uint) bool {
	num := s.DB.Model(&models.User{
		Model: gorm.Model{
			ID: userId,
		},
	}).Where("follower_id = ?", userToFollowId).Association(
		"Followers",
	).Count()

	return num > 0

}
func (s *UserService) FollowUser(userId, userToFollowId uint) (bool, error) {
	err := s.DB.Model(&models.User{
		Model: gorm.Model{
			ID: userId,
		},
	}).Association(
		"Followers",
	).Append(&models.User{
		Model: gorm.Model{
			ID: userToFollowId,
		},
	})

	if err != nil {
		return false, err
	}
	return true, err

}
func (s *UserService) UnfollowUser(userId, userToUnfollowId uint) (bool, error) {
	err := s.DB.Model(&models.User{
		Model: gorm.Model{
			ID: userId,
		},
	}).Association(
		"Followers",
	).Delete(&models.User{
		Model: gorm.Model{
			ID: userToUnfollowId,
		},
	})
	if err != nil {
		return false, err
	}
	return true, err

}

func (s *UserService) CountUsers() int64 {
	var count int64
	s.DB.Table("users").Count(&count)
	return count
}

func (s *UserService) GetFollowersByUsername(username string, limit int) []string {

	var names []string

	s.DB.Raw(`SELECT users.username FROM users INNER JOIN followers ON followers.follower_id=users.id WHERE followers.user_id = (SELECT id from users where username = ?) LIMIT ?`, username, limit).Scan(&names)

	return names
}
