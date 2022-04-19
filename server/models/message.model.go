package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Text      string
	Pub_date  int `gorm:"index:date_index;index:author_date,sort:desc,priority:2"`
	Flagged   int
	Author_id int `gorm:"index:author_date,priority:1"`
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
