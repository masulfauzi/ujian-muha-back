package service

import (
	"time"

	"backend/configs"
	"backend/internal/utils"
)

// UjianTokenService menghitung token ujian global (TOTP) yang berlaku saat ini.
// Token ini murni fungsi dari waktu + secret, jadi tidak butuh scheduler/cron untuk
// "merotasi" dan otomatis konsisten dipakai bersama oleh semua jadwal yang mewajibkannya,
// termasuk saat beberapa jadwal berjalan bersamaan.
type UjianTokenService interface {
	GetCurrentToken() (token string, sisaDetik int, periodeDetik int)
}

type ujianTokenService struct{}

func NewUjianTokenService() UjianTokenService {
	return &ujianTokenService{}
}

func (s *ujianTokenService) GetCurrentToken() (string, int, int) {
	cfg := configs.GetUjianTokenConfig()
	now := time.Now()

	token := utils.GenerateTOTP(cfg.Secret, now, cfg.Period, cfg.Digits)
	sisaDetik := utils.RemainingSeconds(now, cfg.Period)

	return token, sisaDetik, int(cfg.Period.Seconds())
}
