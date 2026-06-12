package config

import (
	"os"
	"strconv"

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
	CacheTTL           int
	Debug              bool
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
		CacheTTL: func() int {
			ttl, err := strconv.Atoi(os.Getenv("CACHE_TTL"))
			if err != nil {
				return 3600
			}
			return ttl
		}(),
		Debug: func() bool {
			debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
			if err != nil {
				return false
			}
			return debug
		}(),
	}, nil
}
