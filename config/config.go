package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	DBUrl              string
	JwtSecret          string
	ServerPort         string
	HTTPTimeout        time.Duration
	SpoonacularApiKey  string
	SpoonacularBaseUrl string
}

func LoadConfig() (*Config, error) {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL harus diset")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET harus diset")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "2001"
	}

	timeoutStr := os.Getenv("HTTP_TIMEOUT")
	timeout := 10 * time.Second
	if timeoutStr != "" {
		if t, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = t
		}
	}

	spoonacularApiKey := os.Getenv("SPOONACULAR_API_KEY")
	if spoonacularApiKey == "" {
		return nil, fmt.Errorf("SPOONACULAR_API_KEY harus diset")
	}

	spoonacularBaseUrl := os.Getenv("SPOONACULAR_BASE_URL")
	if spoonacularBaseUrl == "" {
		spoonacularBaseUrl = "https://api.spoonacular.com"
	}

	return &Config{
		DBUrl:              dbUrl,
		JwtSecret:          jwtSecret,
		ServerPort:         port,
		HTTPTimeout:        timeout,
		SpoonacularApiKey:  spoonacularApiKey,
		SpoonacularBaseUrl: spoonacularBaseUrl,
	}, nil
}
