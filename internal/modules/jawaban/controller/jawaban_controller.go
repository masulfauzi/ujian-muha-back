package controller

import (
	"errors"
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/jawaban/dto"
	"backend/internal/modules/jawaban/service"

	"github.com/gofiber/fiber/v2"
)

type JawabanController struct {
	service service.JawabanService
}

func NewJawabanController(service service.JawabanService) *JawabanController {
	return &JawabanController{service: service}
}

// CreateJawaban godoc
// @Summary Buat baris jawaban secara manual
// @Description Pada flow normal, baris jawaban kosong sudah dibuat otomatis saat peserta mulai ujian (lihat /nilai/mulai-ujian). Endpoint ini untuk kasus insert manual.
// @Tags Jawaban
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateJawabanRequest true "Data jawaban"
// @Success 201 {object} helpers.Response{data=dto.JawabanResponse} "Create jawaban successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal atau data sudah ada"
// @Router /jawaban [post]
func (c *JawabanController) CreateJawaban(ctx *fiber.Ctx) error {
	var req dto.CreateJawabanRequest
	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateJawaban(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create jawaban successfully", resp)
}

// GetAllJawaban godoc
// @Summary List semua jawaban
// @Tags Jawaban
// @Produce json
// @Param page query int false "Halaman" default(1)
// @Param page_size query int false "Jumlah per halaman" default(10)
// @Param id_nilai query string false "Filter berdasarkan ID nilai/sesi ujian (uuid)"
// @Param id_peserta query string false "Filter berdasarkan ID peserta (uuid)"
// @Param id_soal query string false "Filter berdasarkan ID soal (uuid)"
// @Success 200 {object} helpers.Response{data=dto.JawabanListResponse} "Get all jawaban successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Router /jawaban [get]
func (c *JawabanController) GetAllJawaban(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	idNilai := ctx.Query("id_nilai", "")
	idPeserta := ctx.Query("id_peserta", "")
	idSoal := ctx.Query("id_soal", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}
	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllJawaban(pageNum, pageSizeNum, idNilai, idPeserta, idSoal)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all jawaban successfully", resp)
}

// GetJawabanByID godoc
// @Summary Detail satu jawaban
// @Tags Jawaban
// @Produce json
// @Param id path string true "ID jawaban (uuid)"
// @Success 200 {object} helpers.Response{data=dto.JawabanResponse} "Get jawaban successfully"
// @Failure 404 {object} helpers.Response "Jawaban tidak ditemukan"
// @Router /jawaban/{id} [get]
func (c *JawabanController) GetJawabanByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetJawabanByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jawaban successfully", resp)
}

// GetJawabanByNilai godoc
// @Summary List semua soal & jawaban dalam satu sesi ujian
// @Description Dipanggil setelah mulai-ujian untuk mengambil seluruh soal (via baris jawaban) dalam satu sesi. Opsi jawaban otomatis diacak per-peserta jika acak_opsi aktif pada jadwal. Hasil terurut berdasarkan no_soal di database — untuk urutan tampil yang benar (terutama jika acak_soal aktif), frontend harus re-sort berdasarkan field no_urut. Untuk ujian yang memakai section, gunakan endpoint /jawaban/nilai/{id_nilai}/section/{id_section} agar hanya soal dalam section yang sedang aktif yang dikirim.
// @Tags Jawaban
// @Produce json
// @Param id_nilai path string true "ID nilai/sesi ujian (uuid)"
// @Success 200 {object} helpers.Response{data=[]dto.JawabanResponse} "Get jawaban by nilai successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Router /jawaban/nilai/{id_nilai} [get]
func (c *JawabanController) GetJawabanByNilai(ctx *fiber.Ctx) error {
	idNilai := ctx.Params("id_nilai")

	resp, err := c.service.GetJawabanByNilai(idNilai)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jawaban by nilai successfully", resp)
}

// GetJawabanByNilaiSection godoc
// @Summary List soal & jawaban dalam satu section
// @Description Mengambil hanya soal yang termasuk dalam satu section tertentu pada sesi ujian ini. Ditolak dengan 403 jika section yang diminta urutannya lebih besar dari section aktif (frontier) sesi ini — mencegah peserta mengambil soal section berikutnya sebelum waktu minimal section saat ini terlampaui, walau memanggil API secara langsung. Section yang sudah pernah dibuka (urutan <= frontier) selalu bisa diakses ulang.
// @Tags Jawaban
// @Produce json
// @Param id_nilai path string true "ID nilai/sesi ujian (uuid)"
// @Param id_section path string true "ID section (uuid)"
// @Success 200 {object} helpers.Response{data=[]dto.JawabanResponse} "Get jawaban by nilai dan section successfully"
// @Failure 400 {object} helpers.Response "Section/nilai tidak ditemukan, atau tidak cocok dengan jadwalnya"
// @Failure 403 {object} helpers.Response "Section ini belum terbuka (melebihi frontier sesi ini)"
// @Router /jawaban/nilai/{id_nilai}/section/{id_section} [get]
func (c *JawabanController) GetJawabanByNilaiSection(ctx *fiber.Ctx) error {
	idNilai := ctx.Params("id_nilai")
	idSection := ctx.Params("id_section")

	resp, err := c.service.GetJawabanByNilaiSection(idNilai, idSection)
	if err != nil {
		if errors.Is(err, service.ErrSectionBelumTerbuka) {
			return helpers.ErrorResponse(ctx, fiber.StatusForbidden, err.Error(), nil)
		}
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jawaban by nilai dan section successfully", resp)
}

// GetJawabanByPeserta godoc
// @Summary List seluruh jawaban satu peserta (lintas sesi)
// @Tags Jawaban
// @Produce json
// @Param id_peserta path string true "ID peserta (uuid)"
// @Param page query int false "Halaman" default(1)
// @Param page_size query int false "Jumlah per halaman" default(10)
// @Success 200 {object} helpers.Response{data=dto.JawabanListResponse} "Get jawaban by peserta successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Router /jawaban/peserta/{id_peserta} [get]
func (c *JawabanController) GetJawabanByPeserta(ctx *fiber.Ctx) error {
	idPeserta := ctx.Params("id_peserta")
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

	resp, err := c.service.GetJawabanByPeserta(idPeserta, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jawaban by peserta successfully", resp)
}

// UpdateJawaban godoc
// @Summary Submit/update jawaban satu soal
// @Description Dipakai peserta untuk mengisi/mengganti pilihan jawaban satu soal dalam sesi ujiannya. Server langsung menghitung is_benar dengan membandingkan ke kunci (memperhitungkan unshuffle opsi jika acak_opsi aktif).
// @Tags Jawaban
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID jawaban (uuid)"
// @Param request body dto.UpdateJawabanRequest true "Pilihan jawaban (A-E)"
// @Success 200 {object} helpers.Response{data=dto.JawabanResponse} "Update jawaban successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal"
// @Router /jawaban/{id} [put]
func (c *JawabanController) UpdateJawaban(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateJawabanRequest
	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateJawaban(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update jawaban successfully", resp)
}

// DeleteJawaban godoc
// @Summary Hapus jawaban (soft delete)
// @Tags Jawaban
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID jawaban (uuid)"
// @Success 200 {object} helpers.Response "Delete jawaban successfully"
// @Failure 400 {object} helpers.Response "Gagal menghapus"
// @Router /jawaban/{id} [delete]
func (c *JawabanController) DeleteJawaban(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteJawaban(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete jawaban successfully", nil)
}

// RestoreJawaban godoc
// @Summary Restore jawaban yang sudah dihapus
// @Tags Jawaban
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID jawaban (uuid)"
// @Success 200 {object} helpers.Response "Restore jawaban successfully"
// @Failure 400 {object} helpers.Response "Gagal restore"
// @Router /jawaban/{id}/restore [patch]
func (c *JawabanController) RestoreJawaban(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreJawaban(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore jawaban successfully", nil)
}
