package environment

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func Getenv(key string) string {

	return os.Getenv(key)
}
