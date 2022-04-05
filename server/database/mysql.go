package database

import (
	"aske-w/itu-minitwit/models"
	"fmt"
	"log"
	"os"
	"time"

	// Sqlite driver based on GGO

	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
)

// github.com/mattn/go-sqlite3
func ConnectMySql(mode string) (*gorm.DB, error) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: false,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,         // Disable color
		},
	)
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	address := os.Getenv("MYSQL_ADDRESS")
	port := os.Getenv("MYSQL_PORT")
	db_name := os.Getenv("MYSQL_DATABASE")

	var db *gorm.DB
	var err error
	if mode == "production" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, address, port, db_name)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		db.Use(prometheus.New(prometheusConfiguration(db_name)))
	} else {
		db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
			Logger: newLogger,
		})
	}

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.Latest{})

	return db, nil

}

func prometheusConfiguration(dbName string) prometheus.Config {
	return prometheus.Config{
		DBName:          dbName,
		RefreshInterval: 60,
		HTTPServerPort:  8080, // Use the port as the Iris server
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"Threads_running", "Slow_queries", "Uptime"},
			},
		}, // user defined metrics
	}
}
