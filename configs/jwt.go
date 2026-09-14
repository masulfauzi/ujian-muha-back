package configs

import (
	"time"
)

type JWTConfig struct {
	Secret   string
	Expired  time.Duration
}

func GetJWTConfig() *JWTConfig {
	expiredStr := getEnv("JWT_EXPIRED", "72h")
	duration, err := time.ParseDuration(expiredStr)
	if err != nil {
		duration = 72 * time.Hour
	}

	return &JWTConfig{
		Secret:   getEnv("JWT_SECRET", "supersecretkey"),
		Expired:  duration,
	}
}
