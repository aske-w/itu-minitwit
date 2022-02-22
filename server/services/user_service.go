package services

import (
	"aske-w/itu-minitwit/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) GetById(userId int) (*models.User, error) {

	user := &models.User{}
	err := s.DB.First(&user, userId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
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

func (s *UserService) UserIsFollowing(userId int, followerId int) bool {

	num := s.DB.Find(&models.User{}, userId).Association(
		"Followers",
	).Count()

	return num > 0

}
func (s *UserService) FollowUser(userId int, followerId int) (bool, error) {
	//FIXME DOESNT WORK
	err := s.DB.First(&models.User{
		ID: uint(userId),
	}).Association(
		"Followers",
	).Append(&models.User{
		ID: uint(followerId),
	})
	fmt.Println(err)
	if err != nil {
		return false, err
	}
	return true, err

}
