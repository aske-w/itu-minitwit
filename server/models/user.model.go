package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string
	Email     string
	Pw_Hash   string
	Followers []User `gorm:"many2many:follower"`
	// Messages  []Message `gorm:"foreignKey:Author_id"`
}

/*
// gorm.Model definition
type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
*/
