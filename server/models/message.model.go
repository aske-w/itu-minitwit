package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Text      string
	Pub_date  int
	Flagged   int
	Author_id int
	Author    User
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
