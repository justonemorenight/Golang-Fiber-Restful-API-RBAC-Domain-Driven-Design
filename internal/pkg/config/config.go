package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret        string
	JWTRefreshSecret string
	AccessTokenExp   time.Duration
	RefreshTokenExp  time.Duration
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
}

var cfg *Config

func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	cfg = &Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "backend_fiber"),
		JWTSecret:        getEnv("JWT_SECRET", "s3cret_T0ken"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "s3cret_T0ken_Refresh"),
		AccessTokenExp:   parseDuration(getEnv("JWT_ACCESS_EXP", "15m")),
		RefreshTokenExp:  parseDuration(getEnv("JWT_REFRESH_EXP", "7d")),
	}

	// Log loaded configuration
	log.Printf("Loaded configuration: DB_HOST=%s, DB_PORT=%s, DB_USER=%s, DB_NAME=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName)

	return cfg
}

func GetConfig() *Config {
	if cfg == nil {
		return LoadConfig()
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(duration string) time.Duration {
	d, err := time.ParseDuration(duration)
	if err != nil {
		log.Printf("Warning: invalid duration %s, using default", duration)
		return 24 * time.Hour // default duration
	}
	return d
}
