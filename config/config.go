package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	value := os.Getenv(key)
	if value == "" {
		// Load the .env file if the environment variable is not set
		err := godotenv.Load(".env")
		if err != nil {
			log.Printf("Error loading .env file: %v", err)
		}
		// Get the value from the loaded .env file
		value = os.Getenv(key)
	}
	return value
}