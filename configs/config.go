package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// type of Config *AppConfig
type AppConfig struct {
	AppName                  string
	AppEnv                   string
	Port                     string
	DBHost                   string
	DBPort                   string
	DBUser                   string
	DBPassword               string
	DBName                   string
	RedisHost                string
	RedisPort                string
	RedisPassword            string
	JWTSecretKey             string
	AccessTokenExpireMinutes int
	RefreshTokenExpireHours  int
}

// var
var Config *AppConfig

// func for load
func LoadEnv() {
	// Load .env ffile
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found. Using system envs.")
	}

	Config = &AppConfig{
		AppName:                  MustGetEnvOrDefault("APP_NAME", "TaskManagmentApis"),
		AppEnv:                   MustGetEnvOrDefault("APP_ENV", "development"),
		Port:                     MustGetEnvOrDefault("PORT", "8080"),
		DBHost:                   MustGetEnvOrDefault("DB_HOST", "localhost"),
		DBPort:                   MustGetEnvOrDefault("DB_PORT", "5432"),
		DBUser:                   MustGetEnvOrDefault("DB_USER", "Bishalkoirala"),
		DBPassword:               MustGetEnvOrDefault("DB_PASSWORD", "bishal1212"),
		DBName:                   MustGetEnvOrDefault("DB_NAME", "task_db"),
		RedisHost:                MustGetEnvOrDefault("REDIS_HOST", "localhost"),
		RedisPort:                MustGetEnvOrDefault("REDIS_PORT", "6379"),
		RedisPassword:            MustGetEnvOrDefault("REDIS_PASSWORD", ""),
		JWTSecretKey:             MustGetEnvOrDefault("JWT_SECRET", "mysecretkey"),
		AccessTokenExpireMinutes: mustGetEnvASInt("ACCESS_TOKEN_EXPIRE_MINUTES", 15),
		RefreshTokenExpireHours:  mustGetEnvASInt("REFRESH_TOKEN_EXPIRE_HOURS", 24),
	}
}

// mustGetEnvorDefault
func MustGetEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// // mustGetEnvASInt
func mustGetEnvASInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	// Try to conver the value to an integer
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}
