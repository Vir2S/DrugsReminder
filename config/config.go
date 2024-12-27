package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnString string
	Port         string
	TwilioConfig TwilioConfig
}

type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

func LoadConfig() Config {
	godotenv.Load()

	return Config{
		DBConnString: os.Getenv("DATABASE_URL"),
		Port:         os.Getenv("PORT"),
		TwilioConfig: TwilioConfig{
			AccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
			AuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
			FromNumber: os.Getenv("TWILIO_FROM_NUMBER"),
		},
	}
}
