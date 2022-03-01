package models

import "gorm.io/gorm"

type Latest struct {
	gorm.Model
	ID     uint `gorm:"primaryKey"`
	Latest uint
}
