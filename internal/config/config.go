package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	DBHost      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBPort      string
}

func LoadConfig() *Config {
	// Đọc môi trường từ biến GO_ENV, mặc định là "development"
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// Load file env tương ứng
	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: %s not found, falling back to .env\n", envFile)
		// Fallback to default .env
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return &Config{
		Environment: env,
		DBHost:      os.Getenv("DB_HOST"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		DBPort:      os.Getenv("DB_PORT"),
	}
}
