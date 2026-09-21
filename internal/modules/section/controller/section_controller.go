package controller

import (
	"backend/internal/helpers"
	"backend/internal/modules/section/dto"
	"backend/internal/modules/section/service"

	"github.com/gofiber/fiber/v2"
)

type SectionController struct {
	service service.SectionService
}

func NewSectionController(service service.SectionService) *SectionController {
	return &SectionController{service: service}
}

// DefineSections godoc
// @Summary Definisikan/ubah pembagian section untuk satu jadwal
// @Description Membagi soal-soal satu jadwal ujian ke dalam beberapa section berurutan, tiap section punya jumlah soal & durasi minimal sendiri sebelum peserta bisa lanjut ke section berikutnya. Server otomatis menghitung range no_urut tiap section secara berurutan (sequential) — total jml_soal dari semua section harus sama dengan jumlah soal pada bank soal jadwal ini. Memanggil endpoint ini akan MENGGANTI seluruh definisi section lama untuk jadwal tersebut (soft-delete lalu buat baru). Peserta yang sedang mengerjakan (belum selesai) dan sebelumnya sudah punya section aktif akan OTOMATIS dimigrasikan ke section pertama (urutan 1) yang baru, dengan timer durasi minimal direset dari sekarang — peserta yang belum pernah mulai section sama sekali tidak disentuh.
// @Tags Section
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_jadwal path string true "ID jadwal (uuid)"
// @Param request body dto.DefineSectionRequest true "Daftar section berurutan"
// @Success 200 {object} helpers.Response{data=[]dto.SectionResponse} "Define section successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal, jadwal tidak ditemukan, atau total jml_soal tidak sama dengan jumlah soal jadwal"
// @Router /section/jadwal/{id_jadwal}/define [post]
func (c *SectionController) DefineSections(ctx *fiber.Ctx) error {
	idJadwal := ctx.Params("id_jadwal")
	var req dto.DefineSectionRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.DefineSections(idJadwal, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Define section successfully", resp)
}

// GetSectionsByJadwal godoc
// @Summary List section satu jadwal
// @Description Mengembalikan daftar section (urut berdasarkan urutan) beserta range no_urut dan durasi minimal masing-masing. Jadwal yang belum memakai section akan mengembalikan data kosong.
// @Tags Section
// @Produce json
// @Param id_jadwal path string true "ID jadwal (uuid)"
// @Success 200 {object} helpers.Response{data=[]dto.SectionResponse} "Get section by jadwal successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Router /section/jadwal/{id_jadwal} [get]
func (c *SectionController) GetSectionsByJadwal(ctx *fiber.Ctx) error {
	idJadwal := ctx.Params("id_jadwal")

	resp, err := c.service.GetSectionsByJadwal(idJadwal)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get section by jadwal successfully", resp)
}

// GetSectionByID godoc
// @Summary Detail satu section
// @Tags Section
// @Produce json
// @Param id path string true "ID section (uuid)"
// @Success 200 {object} helpers.Response{data=dto.SectionResponse} "Get section successfully"
// @Failure 404 {object} helpers.Response "Section tidak ditemukan"
// @Router /section/{id} [get]
func (c *SectionController) GetSectionByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetSectionByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get section successfully", resp)
}

// DeleteSection godoc
// @Summary Hapus satu section (soft delete)
// @Description Utilitas admin biasa. Untuk mengubah pembagian section (jumlah/urutan/range), gunakan endpoint define ulang daripada menghapus satu-satu, agar range no_urut tetap konsisten.
// @Tags Section
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID section (uuid)"
// @Success 200 {object} helpers.Response "Delete section successfully"
// @Failure 400 {object} helpers.Response "Gagal menghapus"
// @Router /section/{id} [delete]
func (c *SectionController) DeleteSection(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.DeleteSection(id); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete section successfully", nil)
}

// RestoreSection godoc
// @Summary Restore section yang sudah dihapus
// @Tags Section
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID section (uuid)"
// @Success 200 {object} helpers.Response "Restore section successfully"
// @Failure 400 {object} helpers.Response "Gagal restore"
// @Router /section/{id}/restore [patch]
func (c *SectionController) RestoreSection(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.RestoreSection(id); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore section successfully", nil)
}
