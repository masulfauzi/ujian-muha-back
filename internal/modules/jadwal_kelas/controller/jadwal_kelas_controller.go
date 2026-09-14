package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/jadwal_kelas/dto"
	"backend/internal/modules/jadwal_kelas/service"

	"github.com/gofiber/fiber/v2"
)

type JadwalKelasController struct {
	service service.JadwalKelasService
}

func NewJadwalKelasController(service service.JadwalKelasService) *JadwalKelasController {
	return &JadwalKelasController{service: service}
}

// CreateJadwalKelas godoc
// @Summary Assign kelas ke jadwal ujian
// @Description Mendaftarkan satu kelas agar bisa mengikuti satu jadwal ujian tertentu.
// @Tags JadwalKelas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateJadwalKelasRequest true "Data jadwal-kelas"
// @Success 201 {object} helpers.Response{data=dto.JadwalKelasResponse} "Create jadwal kelas successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal atau kelas sudah terdaftar di jadwal ini"
// @Router /jadwal-kelas [post]
func (c *JadwalKelasController) CreateJadwalKelas(ctx *fiber.Ctx) error {
	var req dto.CreateJadwalKelasRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateJadwalKelas(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create jadwal kelas successfully", resp)
}

// GetAllJadwalKelas godoc
// @Summary List assignment jadwal-kelas
// @Tags JadwalKelas
// @Produce json
// @Param page query int false "Halaman" default(1)
// @Param page_size query int false "Jumlah per halaman" default(10)
// @Param id_jadwal query string false "Filter berdasarkan ID jadwal (uuid)"
// @Param id_kelas query string false "Filter berdasarkan ID kelas (uuid)"
// @Success 200 {object} helpers.Response{data=dto.JadwalKelasListResponse} "Get all jadwal kelas successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Router /jadwal-kelas [get]
func (c *JadwalKelasController) GetAllJadwalKelas(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	idJadwal := ctx.Query("id_jadwal", "")
	idKelas := ctx.Query("id_kelas", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllJadwalKelas(pageNum, pageSizeNum, idJadwal, idKelas)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all jadwal kelas successfully", resp)
}

// GetJadwalKelasByID godoc
// @Summary Detail assignment jadwal-kelas
// @Tags JadwalKelas
// @Produce json
// @Param id path string true "ID jadwal-kelas (uuid)"
// @Success 200 {object} helpers.Response{data=dto.JadwalKelasResponse} "Get jadwal kelas successfully"
// @Failure 404 {object} helpers.Response "Data tidak ditemukan"
// @Router /jadwal-kelas/{id} [get]
func (c *JadwalKelasController) GetJadwalKelasByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetJadwalKelasByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jadwal kelas successfully", resp)
}

// UpdateJadwalKelas godoc
// @Summary Update assignment jadwal-kelas
// @Tags JadwalKelas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID jadwal-kelas (uuid)"
// @Param request body dto.UpdateJadwalKelasRequest true "Data jadwal-kelas"
// @Success 200 {object} helpers.Response{data=dto.JadwalKelasResponse} "Update jadwal kelas successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal"
// @Router /jadwal-kelas/{id} [put]
func (c *JadwalKelasController) UpdateJadwalKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateJadwalKelasRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateJadwalKelas(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update jadwal kelas successfully", resp)
}

// DeleteJadwalKelas godoc
// @Summary Hapus assignment jadwal-kelas
// @Tags JadwalKelas
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID jadwal-kelas (uuid)"
// @Success 200 {object} helpers.Response "Delete jadwal kelas successfully"
// @Failure 400 {object} helpers.Response "Gagal menghapus"
// @Router /jadwal-kelas/{id} [delete]
func (c *JadwalKelasController) DeleteJadwalKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteJadwalKelas(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete jadwal kelas successfully", nil)
}
