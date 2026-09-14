package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/kelas/dto"
	"backend/internal/modules/kelas/service"

	"github.com/gofiber/fiber/v2"
)

type KelasController struct {
	service service.KelasService
}

func NewKelasController(service service.KelasService) *KelasController {
	return &KelasController{service: service}
}

// CreateKelas godoc
// @Summary Buat kelas baru
// @Tags Kelas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateKelasRequest true "Data kelas"
// @Success 201 {object} helpers.Response{data=dto.KelasResponse} "Create kelas successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal"
// @Router /kelas [post]
func (c *KelasController) CreateKelas(ctx *fiber.Ctx) error {
	var req dto.CreateKelasRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateKelas(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create kelas successfully", resp)
}

// GetAllKelas godoc
// @Summary List kelas
// @Tags Kelas
// @Produce json
// @Param page query int false "Halaman" default(1)
// @Param page_size query int false "Jumlah per halaman" default(10)
// @Param id_jurusan query string false "Filter berdasarkan ID jurusan (uuid)"
// @Param tingkat query string false "Filter berdasarkan tingkat"
// @Success 200 {object} helpers.Response{data=dto.KelasListResponse} "Get all kelas successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Router /kelas [get]
func (c *KelasController) GetAllKelas(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	idJurusan := ctx.Query("id_jurusan", "")
	tingkat := ctx.Query("tingkat", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllKelas(pageNum, pageSizeNum, idJurusan, tingkat)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all kelas successfully", resp)
}

// GetKelasByID godoc
// @Summary Detail kelas
// @Tags Kelas
// @Produce json
// @Param id path string true "ID kelas (uuid)"
// @Success 200 {object} helpers.Response{data=dto.KelasResponse} "Get kelas successfully"
// @Failure 404 {object} helpers.Response "Kelas tidak ditemukan"
// @Router /kelas/{id} [get]
func (c *KelasController) GetKelasByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetKelasByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get kelas successfully", resp)
}

// UpdateKelas godoc
// @Summary Update kelas
// @Tags Kelas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID kelas (uuid)"
// @Param request body dto.UpdateKelasRequest true "Data kelas"
// @Success 200 {object} helpers.Response{data=dto.KelasResponse} "Update kelas successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal"
// @Router /kelas/{id} [put]
func (c *KelasController) UpdateKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateKelasRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateKelas(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update kelas successfully", resp)
}

// DeleteKelas godoc
// @Summary Hapus kelas (soft delete)
// @Tags Kelas
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID kelas (uuid)"
// @Success 200 {object} helpers.Response "Delete kelas successfully"
// @Failure 400 {object} helpers.Response "Gagal menghapus"
// @Router /kelas/{id} [delete]
func (c *KelasController) DeleteKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteKelas(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete kelas successfully", nil)
}

// RestoreKelas godoc
// @Summary Restore kelas yang sudah dihapus
// @Tags Kelas
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID kelas (uuid)"
// @Success 200 {object} helpers.Response "Restore kelas successfully"
// @Failure 400 {object} helpers.Response "Gagal restore"
// @Router /kelas/{id}/restore [patch]
func (c *KelasController) RestoreKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreKelas(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore kelas successfully", nil)
}
