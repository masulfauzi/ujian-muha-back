package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/login_log/service"

	"github.com/gofiber/fiber/v2"
)

type LoginLogController struct {
	service service.LoginLogService
}

func NewLoginLogController(service service.LoginLogService) *LoginLogController {
	return &LoginLogController{service: service}
}

// GetAllLoginLog godoc
// @Summary List riwayat percobaan login (berhasil maupun gagal)
// @Description Berisi username yang dipakai, status berhasil/gagal, user agent, X-Requested-With, dan IP address setiap percobaan login. Baris dicatat otomatis oleh POST /auth/login, tidak ada endpoint untuk menulis manual.
// @Tags Login Log
// @Produce json
// @Security BearerAuth
// @Param page query int false "Halaman" default(1)
// @Param page_size query int false "Jumlah per halaman" default(10)
// @Param username query string false "Filter berdasarkan username (partial match)"
// @Success 200 {object} helpers.Response{data=dto.LoginLogListResponse} "Get all login log successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Router /login-log [get]
func (c *LoginLogController) GetAllLoginLog(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	username := ctx.Query("username", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllLoginLog(pageNum, pageSizeNum, username)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all login log successfully", resp)
}
