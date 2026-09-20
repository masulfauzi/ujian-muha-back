package configs

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Name        string
	Port        int
	Env         string
	FrontendURL string
	ServerNo    string
	// TrustedProxies adalah daftar IP/CIDR reverse proxy (load balancer) yang
	// boleh dipercaya untuk menentukan IP client asli lewat header
	// X-Forwarded-For (lihat cmd/server/main.go). Kosong berarti tidak ada
	// proxy yang dipercaya — ctx.IP() akan tetap memakai IP koneksi TCP
	// langsung (di belakang load balancer, ini akan jadi IP load balancer-nya).
	TrustedProxies []string
}

func LoadEnv() error {
	return godotenv.Load()
}

func GetAppConfig() *AppConfig {
	return &AppConfig{
		Name:           getEnv("APP_NAME", "Fiber Backend API"),
		Port:           getEnvInt("APP_PORT", 3000),
		Env:            getEnv("APP_ENV", "development"),
		FrontendURL:    getEnv("FRONTEND_URL", "http://localhost:5173"),
		ServerNo:       getEnv("SERVER_NO", "1"),
		TrustedProxies: getEnvList("TRUSTED_PROXIES"),
	}
}

// getEnvList membaca env var berisi daftar dipisah koma (mis. "127.0.0.1,10.0.0.0/8")
// dan mengembalikan slice yang sudah di-trim, tanpa entry kosong.
func getEnvList(key string) []string {
	raw := getEnv(key, "")
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
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
