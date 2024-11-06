package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	JwtSecret  string
	Production bool
}

// Using type constraints to limit T to supported types
func getEnvOrDefault[T string | int | float64 | bool](key string, defaultValue T) T {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// Use type switch on the type parameter T
	switch any(defaultValue).(type) {
	case string:
		return any(value).(T)
	case int:
		if v, err := strconv.Atoi(value); err == nil {
			return any(v).(T)
		}
	case float64:
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			return any(v).(T)
		}
	case bool:
		if v, err := strconv.ParseBool(value); err == nil {
			return any(v).(T)
		}
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
		JwtSecret:  getEnvOrDefault("JWT_SECRET", ""),
		Production: getEnvOrDefault("PRODUCTION", false),
	}
}
