package middleware

import (
	"backend/internal/helpers"
	useragentservice "backend/internal/modules/user_agent/service"

	"github.com/gofiber/fiber/v2"
)

// CheckAllowedUserAgent menolak request jika header User-Agent maupun
// X-Requested-With yang dikirim, keduanya tidak cocok dengan whitelist tabel
// user_agent (lihat modules/user_agent) — cukup salah satu yang cocok (OR,
// tidak perlu baris yang sama) untuk request diloloskan. Dipasang khusus di
// endpoint yang menandai peserta mulai mengerjakan ujian
// (POST /nilai/mulai-ujian/:id_jadwal), sehingga hanya browser/aplikasi resmi
// yang bisa memulai sesi ujian.
func CheckAllowedUserAgent(svc useragentservice.UserAgentService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userAgent := ctx.Get("User-Agent")
		xRequestedWith := ctx.Get("X-Requested-With")

		allowed, err := svc.IsAllowed(userAgent, xRequestedWith)
		if err != nil {
			return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, "Gagal memvalidasi perangkat/aplikasi", nil)
		}
		if !allowed {
			return helpers.ErrorResponse(ctx, fiber.StatusForbidden, "Perangkat/aplikasi tidak diizinkan untuk mengerjakan ujian", nil)
		}

		return ctx.Next()
	}
}
