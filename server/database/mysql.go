package database

import (
	"aske-w/itu-minitwit/models"
	"log"
	"os"
	"time"

	// Sqlite driver based on GGO

	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// github.com/mattn/go-sqlite3
func ConnectMySql() (*gorm.DB, error) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: false,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,         // Disable color
		},
	)

	dsn := "user:password@tcp(127.0.0.1:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	// conn, err := sql.Open("sqlite3", "./db.db")

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Message{})

	return db, nil

}
