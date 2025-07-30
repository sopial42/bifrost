package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/bifrost/pkg/logger"
)

type Config struct {
	Cors   Cors
	DB     DBConfig
	Logger logger.Config
	Port   string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Unsecure bool
}

type Cors struct {
	AllowOrigin string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Logger: logger.Config{
			Level:         mustGet("LOG_LEVEL"),
			IsDevelopment: mustGetBool("LOG_IS_DEVELOPMENT"),
		},
		DB: DBConfig{
			Host:     mustGet("DB_HOST"),
			Port:     mustGet("DB_PORT"),
			User:     mustGet("DB_USER"),
			Password: mustGet("DB_PASSWORD"),
			DBName:   mustGet("DB_NAME"),
			Unsecure: mustGetBool("DB_UNSECURE_MODE"),
		},
		Port: mustGet("PORT"),
		Cors: Cors{
			AllowOrigin: mustGet("CORS_ALLOW_ORIGIN"),
		},
	}
}

func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}

// mustGetBool is a helper function to get a boolean value from an environment variable
// It panics if the environment variable is not set or cannot be converted to a boolean
func mustGetBool(key string) bool {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		log.Fatalf("unable to cast env var %s to bool: %v", key, err)
	}
	return boolVal
}
