package controller

import (
	"errors"
	"path/filepath"
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

// DeleteAllPeserta godoc
// @Summary Hapus PERMANEN seluruh peserta di seluruh sistem
// @Description Menghapus semua baris di tabel peserta secara permanen (bukan soft-delete, tidak bisa di-restore). Wajib menyertakan query param confirm=true sebagai pengaman agar tidak terpicu tidak sengaja. Tidak menghapus record nilai/jawaban terkait — record itu jadi yatim (orphan) mereferensikan id peserta yang sudah tidak ada.
// @Tags Peserta
// @Produce json
// @Security BearerAuth
// @Param confirm query string true "Wajib diisi 'true' untuk konfirmasi penghapusan"
// @Success 200 {object} helpers.Response "Delete all peserta successfully"
// @Failure 400 {object} helpers.Response "Konfirmasi tidak ada atau gagal menghapus"
// @Router /peserta [delete]
func (c *PesertaController) DeleteAllPeserta(ctx *fiber.Ctx) error {
	if ctx.Query("confirm") != "true" {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Konfirmasi diperlukan", map[string]string{
			"error": "Tambahkan query param confirm=true untuk menghapus SEMUA peserta secara permanen",
		})
	}

	count, err := c.service.DeleteAllPeserta(ctx.Context())
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Gagal menghapus semua peserta", map[string]string{
			"error": err.Error(),
		})
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete all peserta successfully", map[string]int64{
		"deleted_count": count,
	})
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
// @Description Menghasilkan PDF berisi kartu untuk setiap peserta di kelas tersebut, ditata grid 2x5 kartu per halaman A4 dengan garis putus-putus sebagai panduan gunting. Kartu ini bersifat global (tidak terikat jadwal/ujian tertentu) — menampilkan nama, username, password, dan kelas.
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

// ImportPesertaFromExcel godoc
// @Summary Import peserta dari file Excel
// @Description Upload file .xls/.xlsx (maks 10MB) berisi banyak peserta sekaligus. Kolom: nama, username, password, kelas (diisi nama kelas). Jika ada baris dengan nama kelas yang tidak ditemukan di data master Kelas, seluruh import dibatalkan (tidak ada data yang tersimpan).
// @Tags Peserta
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File Excel (.xls/.xlsx)"
// @Success 200 {object} helpers.Response{data=dto.ImportPesertaResponse} "Import peserta berhasil"
// @Failure 400 {object} helpers.Response "File tidak valid, kolom kelas tidak ditemukan, atau import gagal"
// @Router /peserta/import [post]
func (c *PesertaController) ImportPesertaFromExcel(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "File tidak ditemukan", map[string]string{
			"error": "Silakan upload file excel",
		})
	}

	const maxFileSize = 10 * 1024 * 1024
	if file.Size > maxFileSize {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "File terlalu besar", map[string]string{
			"error": "Max file size adalah 10MB",
		})
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".xls" && ext != ".xlsx" {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Format file tidak valid", map[string]string{
			"error": "File harus berupa .xls atau .xlsx",
		})
	}

	req := &dto.ImportPesertaRequest{File: file}

	resp, err := c.service.ImportPesertaFromExcel(ctx.Context(), req)
	if err != nil {
		var kelasErr *dto.KelasNotFoundError
		if errors.As(err, &kelasErr) {
			return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Import dibatalkan: ada kolom kelas yang tidak ditemukan", kelasErr.Details)
		}
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Import peserta gagal", map[string]string{
			"error": err.Error(),
		})
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Import peserta berhasil", resp)
}

// DownloadTemplate godoc
// @Summary Download template Excel untuk import peserta
// @Description Menghasilkan file .xlsx berisi header + 1 baris contoh (nama, username, password) sesuai urutan kolom yang dibaca endpoint import (POST /peserta/import).
// @Tags Peserta
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success 200 {file} file "File template_import_peserta.xlsx"
// @Failure 500 {object} helpers.Response "Gagal membuat file template"
// @Router /peserta/template [get]
func (c *PesertaController) DownloadTemplate(ctx *fiber.Ctx) error {
	fileBytes, err := c.service.GenerateImportTemplate()
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, "Gagal membuat file template", nil)
	}

	ctx.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Set("Content-Disposition", `attachment; filename="template_import_peserta.xlsx"`)
	return ctx.Send(fileBytes)
}
