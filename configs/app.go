package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Name        string
	Port        int
	Env         string
	FrontendURL string
	ServerNo    string
}

func LoadEnv() error {
	return godotenv.Load()
}

func GetAppConfig() *AppConfig {
	return &AppConfig{
		Name:        getEnv("APP_NAME", "Fiber Backend API"),
		Port:        getEnvInt("APP_PORT", 3000),
		Env:         getEnv("APP_ENV", "development"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		ServerNo:    getEnv("SERVER_NO", "1"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetEnvOrDefault(key, defaultValue string) string {
	return getEnv(key, defaultValue)
}
