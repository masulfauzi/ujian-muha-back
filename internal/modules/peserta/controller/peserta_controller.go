package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/peserta/dto"
	"backend/internal/modules/peserta/service"

	"github.com/gofiber/fiber/v2"
)

type PesertaController struct {
	service service.PesertaService
}

func NewPesertaController(service service.PesertaService) *PesertaController {
	return &PesertaController{service: service}
}

// CreatePeserta godoc
// @Summary Buat peserta (siswa) baru
// @Tags Peserta
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePesertaRequest true "Data peserta"
// @Success 201 {object} helpers.Response{data=dto.PesertaResponse} "Create peserta successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal"
// @Router /peserta [post]
func (c *PesertaController) CreatePeserta(ctx *fiber.Ctx) error {
	var req dto.CreatePesertaRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreatePeserta(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create peserta successfully", resp)
}

// GetAllPeserta godoc
// @Summary List peserta
// @Tags Peserta
// @Produce json
// @Param page query int false "Halaman" default(1)
// @Param page_size query int false "Jumlah per halaman" default(10)
// @Param id_kelas query string false "Filter berdasarkan ID kelas (uuid)"
// @Success 200 {object} helpers.Response{data=dto.PesertaListResponse} "Get all peserta successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Router /peserta [get]
func (c *PesertaController) GetAllPeserta(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	idKelas := ctx.Query("id_kelas", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllPeserta(pageNum, pageSizeNum, idKelas)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all peserta successfully", resp)
}

// GetPesertaByID godoc
// @Summary Detail peserta
// @Tags Peserta
// @Produce json
// @Param id path string true "ID peserta (uuid)"
// @Success 200 {object} helpers.Response{data=dto.PesertaResponse} "Get peserta successfully"
// @Failure 404 {object} helpers.Response "Peserta tidak ditemukan"
// @Router /peserta/{id} [get]
func (c *PesertaController) GetPesertaByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetPesertaByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get peserta successfully", resp)
}

// UpdatePeserta godoc
// @Summary Update peserta
// @Tags Peserta
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID peserta (uuid)"
// @Param request body dto.UpdatePesertaRequest true "Data peserta"
// @Success 200 {object} helpers.Response{data=dto.PesertaResponse} "Update peserta successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal"
// @Router /peserta/{id} [put]
func (c *PesertaController) UpdatePeserta(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdatePesertaRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdatePeserta(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update peserta successfully", resp)
}

// DeletePeserta godoc
// @Summary Hapus peserta (soft delete)
// @Tags Peserta
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID peserta (uuid)"
// @Success 200 {object} helpers.Response "Delete peserta successfully"
// @Failure 400 {object} helpers.Response "Gagal menghapus"
// @Router /peserta/{id} [delete]
func (c *PesertaController) DeletePeserta(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeletePeserta(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete peserta successfully", nil)
}

// RestorePeserta godoc
// @Summary Restore peserta yang sudah dihapus
// @Tags Peserta
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID peserta (uuid)"
// @Success 200 {object} helpers.Response "Restore peserta successfully"
// @Failure 400 {object} helpers.Response "Gagal restore"
// @Router /peserta/{id}/restore [patch]
func (c *PesertaController) RestorePeserta(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestorePeserta(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore peserta successfully", nil)
}

// DownloadKartuUjian godoc
// @Summary Download kartu peserta ujian satu kelas (PDF, siap cetak & gunting)
// @Description Menghasilkan PDF berisi kartu untuk setiap peserta di kelas tersebut, ditata grid 2x5 kartu per halaman A4 dengan garis putus-putus sebagai panduan gunting. Kartu ini bersifat global (tidak terikat jadwal/ujian tertentu) — hanya menampilkan nama, username, dan kelas, tanpa password.
// @Tags Peserta
// @Produce application/pdf
// @Security BearerAuth
// @Param id_kelas path string true "ID kelas (uuid)"
// @Success 200 {file} file "File kartu_ujian_<kelas>.pdf"
// @Failure 400 {object} helpers.Response "Kelas tidak ditemukan atau belum memiliki peserta"
// @Router /peserta/kartu-ujian/{id_kelas} [get]
func (c *PesertaController) DownloadKartuUjian(ctx *fiber.Ctx) error {
	idKelas := ctx.Params("id_kelas")

	fileBytes, err := c.service.GenerateKartuUjianPDF(idKelas)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	ctx.Set("Content-Type", "application/pdf")
	ctx.Set("Content-Disposition", `attachment; filename="kartu_ujian.pdf"`)
	return ctx.Send(fileBytes)
}
