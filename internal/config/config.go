package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment      string
	DBHost           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBPort           string
	JWTSecret        string
	JWTRefreshSecret string
	AccessTokenExp   time.Duration
	RefreshTokenExp  time.Duration
}

func LoadConfig() *Config {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// Load file env
	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: %s not found, falling back to .env\n", envFile)
		// Fallback to default .env
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	accessExp, _ := time.ParseDuration(getEnvWithDefault("JWT_ACCESS_EXP", "15m"))
	refreshExp, _ := time.ParseDuration(getEnvWithDefault("JWT_REFRESH_EXP", "7d"))

	return &Config{
		Environment:      env,
		DBHost:           os.Getenv("DB_HOST"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		DBPort:           os.Getenv("DB_PORT"),
		JWTSecret:        getEnvWithDefault("JWT_SECRET", "your-secret-key"),
		JWTRefreshSecret: getEnvWithDefault("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
		AccessTokenExp:   accessExp,
		RefreshTokenExp:  refreshExp,
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
