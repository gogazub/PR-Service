package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPHost string
	HTTPPort string

	PGHost     string
	PGPort     string
	PGUser     string
	PGPassword string
	PGDatabase string
	SSLMode    string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		HTTPHost: getEnv("HTTP_HOST", "0.0.0.0"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),

		PGHost:     getEnv("POSTGRES_HOST", "localhost"),
		PGPort:     getEnv("POSTGRES_PORT", "5432"),
		PGUser:     getEnv("POSTGRES_USER", "postgres"),
		PGPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PGDatabase: getEnv("POSTGRES_DB", "postgres"),
		SSLMode:    getEnv("SSL_MODE", "disable"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
