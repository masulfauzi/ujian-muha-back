package controller

import (
	"backend/internal/helpers"
	"backend/internal/modules/ujian_token/dto"
	"backend/internal/modules/ujian_token/service"

	"github.com/gofiber/fiber/v2"
)

type UjianTokenController struct {
	service service.UjianTokenService
}

func NewUjianTokenController(service service.UjianTokenService) *UjianTokenController {
	return &UjianTokenController{service: service}
}

// GetCurrentToken godoc
// @Summary Ambil token ujian yang sedang berlaku
// @Description Token bersifat global (bukan per-jadwal) dan berubah otomatis setiap periode (default 15 menit, lihat UJIAN_TOKEN_PERIOD). Kode yang sama berlaku untuk semua jadwal yang mengaktifkan wajib_token, termasuk saat beberapa jadwal berjalan bersamaan — tinggal umumkan/tampilkan kode ini ke peserta di awal sesi ujian, lalu peserta mengirimkannya kembali saat memanggil /nilai/mulai-ujian/{id_jadwal}.
// @Tags UjianToken
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helpers.Response{data=dto.CurrentTokenResponse} "Get current token successfully"
// @Router /ujian-token/current [get]
func (c *UjianTokenController) GetCurrentToken(ctx *fiber.Ctx) error {
	token, sisaDetik, periodeDetik := c.service.GetCurrentToken()

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get current token successfully", dto.CurrentTokenResponse{
		Token:        token,
		SisaDetik:    sisaDetik,
		PeriodeDetik: periodeDetik,
	})
}
