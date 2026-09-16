package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/user_agent/dto"
	"backend/internal/modules/user_agent/service"

	"github.com/gofiber/fiber/v2"
)

type UserAgentController struct {
	service service.UserAgentService
}

func NewUserAgentController(service service.UserAgentService) *UserAgentController {
	return &UserAgentController{service: service}
}

// CreateUserAgent godoc
// @Summary Daftarkan user agent yang boleh mengerjakan ujian
// @Description Saat peserta memanggil POST /nilai/mulai-ujian/{id_jadwal}, request diloloskan jika header User-Agent COCOK dengan user_agent di baris manapun, ATAU header X-Requested-With cocok dengan x_requested_with di baris manapun (exact match, tidak harus baris yang sama). Lihat middleware.CheckAllowedUserAgent.
// @Tags User Agent
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUserAgentRequest true "Data user agent"
// @Success 201 {object} helpers.Response{data=dto.UserAgentResponse} "Create user agent successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal"
// @Router /user-agent [post]
func (c *UserAgentController) CreateUserAgent(ctx *fiber.Ctx) error {
	var req dto.CreateUserAgentRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateUserAgent(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create user agent successfully", resp)
}

// GetAllUserAgent godoc
// @Summary List user agent yang boleh mengerjakan ujian
// @Tags User Agent
// @Produce json
// @Param page query int false "Halaman" default(1)
// @Param page_size query int false "Jumlah per halaman" default(10)
// @Success 200 {object} helpers.Response{data=dto.UserAgentListResponse} "Get all user agent successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Router /user-agent [get]
func (c *UserAgentController) GetAllUserAgent(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllUserAgent(pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all user agent successfully", resp)
}

// GetUserAgentByID godoc
// @Summary Detail user agent
// @Tags User Agent
// @Produce json
// @Param id path string true "ID user agent (uuid)"
// @Success 200 {object} helpers.Response{data=dto.UserAgentResponse} "Get user agent successfully"
// @Failure 404 {object} helpers.Response "User agent tidak ditemukan"
// @Router /user-agent/{id} [get]
func (c *UserAgentController) GetUserAgentByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetUserAgentByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get user agent successfully", resp)
}

// UpdateUserAgent godoc
// @Summary Update user agent
// @Tags User Agent
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID user agent (uuid)"
// @Param request body dto.UpdateUserAgentRequest true "Data user agent"
// @Success 200 {object} helpers.Response{data=dto.UserAgentResponse} "Update user agent successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal"
// @Router /user-agent/{id} [put]
func (c *UserAgentController) UpdateUserAgent(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateUserAgentRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateUserAgent(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update user agent successfully", resp)
}

// DeleteUserAgent godoc
// @Summary Hapus user agent (soft delete)
// @Tags User Agent
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID user agent (uuid)"
// @Success 200 {object} helpers.Response "Delete user agent successfully"
// @Failure 400 {object} helpers.Response "Gagal menghapus"
// @Router /user-agent/{id} [delete]
func (c *UserAgentController) DeleteUserAgent(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteUserAgent(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete user agent successfully", nil)
}

// RestoreUserAgent godoc
// @Summary Restore user agent yang sudah dihapus
// @Tags User Agent
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID user agent (uuid)"
// @Success 200 {object} helpers.Response "Restore user agent successfully"
// @Failure 400 {object} helpers.Response "Gagal restore"
// @Router /user-agent/{id}/restore [patch]
func (c *UserAgentController) RestoreUserAgent(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreUserAgent(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore user agent successfully", nil)
}
