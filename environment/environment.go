package environment

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Print("Error loading .env file")
		fmt.Print(err)
	}
}

func Getenv(key string) string {

	return os.Getenv(key)
}
