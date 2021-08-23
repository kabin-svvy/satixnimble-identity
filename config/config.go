package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	err := godotenv.Load(".env-dev")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}
