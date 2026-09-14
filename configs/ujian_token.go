package configs

import (
	"time"
)

type UjianTokenConfig struct {
	Secret     string
	Period     time.Duration
	Digits     int
	GraceSteps int
}

func GetUjianTokenConfig() *UjianTokenConfig {
	periodStr := getEnv("UJIAN_TOKEN_PERIOD", "15m")
	period, err := time.ParseDuration(periodStr)
	if err != nil {
		period = 15 * time.Minute
	}

	return &UjianTokenConfig{
		Secret:     getEnv("UJIAN_TOKEN_SECRET", "change-this-ujian-token-secret"),
		Period:     period,
		Digits:     6,
		GraceSteps: 1,
	}
}
