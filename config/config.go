package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerUrl          string
	GoogleClientId     string
	GoogleClientSecret string
	SessionSecret      string
	DatabaseUrl        string
	GiphyApiKey        string
	AuthEndpoint       string
}

func GetConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}

	return &Config{
		ServerUrl:          os.Getenv("SERVER_URL"),
		GoogleClientId:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		SessionSecret:      os.Getenv("SESSION_SECRET"),
		DatabaseUrl:        os.Getenv("DATABASE_URL"),
		GiphyApiKey:        os.Getenv("GIPHY_API_KEY"),
		AuthEndpoint:       os.Getenv("AUTH_ENDPOINT"),
	}, nil
}
