package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey       string
	Port         string
	DBConnString string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Loading .env file error: %v", err)
	}

	return Config{
		APIKey:       os.Getenv("API_KEY"),
		Port:         os.Getenv("PORT"),
		DBConnString: os.Getenv("DB_CONN_STRING"),
	}
}
