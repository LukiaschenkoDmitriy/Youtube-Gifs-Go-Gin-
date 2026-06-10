package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GoogleClientId     string
	GoogleClientSecret string
	SessionSecret      string
	DatabaseUrl        string
	GiphyApiKey        string
}

func GetConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		GoogleClientId:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		SessionSecret:      os.Getenv("SESSION_SECRET"),
		DatabaseUrl:        os.Getenv("DATABASE_URL"),
		GiphyApiKey:        os.Getenv("GIPHY_API_KEY"),
	}, nil
}
