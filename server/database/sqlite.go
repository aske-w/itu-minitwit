package database

import (
	"aske-w/itu-minitwit/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm/logger"

	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
)

// github.com/mattn/go-sqlite3
func ConnectSqlite() (*gorm.DB, error) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
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
