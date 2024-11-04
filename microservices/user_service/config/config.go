package config

import (
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnvOrDefault("DB_HOST", "user-db"),
		DBPort:     getEnvOrDefault("DB_PORT", "3306"),
		DBName:     getEnvOrDefault("DB_NAME", "users"),
		DBUser:     getEnvOrDefault("DB_USER", "root"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", ""),
	}
}
