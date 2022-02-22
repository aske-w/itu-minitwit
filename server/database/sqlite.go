package database

import (
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
)

// github.com/mattn/go-sqlite3
func ConnectSqlite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	// conn, err := sql.Open("sqlite3", "./db.db")

	if err != nil {
		return nil, err
	}

	return db, nil

}
