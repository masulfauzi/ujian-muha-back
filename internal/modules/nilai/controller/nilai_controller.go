package controller

import (
	"fmt"
	"strconv"
	"strings"

	"backend/internal/helpers"
	"backend/internal/modules/nilai/dto"
	"backend/internal/modules/nilai/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type NilaiController struct {
	service service.NilaiService
}

func NewNilaiController(service service.NilaiService) *NilaiController {
	return &NilaiController{service: service}
}

// CreateNilai godoc
// @Summary Buat data nilai secara manual
// @Description Biasanya sesi pengerjaan (nilai) dibuat otomatis lewat endpoint mulai-ujian. Endpoint ini untuk kasus input manual oleh admin.
// @Tags Nilai
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateNilaiRequest true "Data nilai"
// @Success 201 {object} helpers.Response{data=dto.NilaiResponse} "Create nilai successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal atau data sudah ada"
// @Router /nilai [post]
func (c *NilaiController) CreateNilai(ctx *fiber.Ctx) error {
	var req dto.CreateNilaiRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateNilai(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create nilai successfully", resp)
}

// GetAllNilai godoc
// @Summary List semua data nilai/sesi ujian
// @Tags Nilai
// @Produce json
// @Param page query int false "Halaman" default(1)
// @Param page_size query int false "Jumlah per halaman" default(10)
// @Param id_peserta query string false "Filter berdasarkan ID peserta (uuid)"
// @Param id_jadwal query string false "Filter berdasarkan ID jadwal (uuid)"
// @Success 200 {object} helpers.Response{data=dto.NilaiListResponse} "Get all nilai successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Router /nilai [get]
func (c *NilaiController) GetAllNilai(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	idPeserta := ctx.Query("id_peserta", "")
	idJadwal := ctx.Query("id_jadwal", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllNilai(pageNum, pageSizeNum, idPeserta, idJadwal)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all nilai successfully", resp)
}

// GetNilaiByID godoc
// @Summary Detail satu sesi nilai/ujian
// @Description Response menyertakan id_section_aktif & wkt_mulai_section jika jadwal ini memakai fitur section.
// @Tags Nilai
// @Produce json
// @Param id path string true "ID nilai (uuid)"
// @Success 200 {object} helpers.Response{data=dto.NilaiResponse} "Get nilai successfully"
// @Failure 404 {object} helpers.Response "Nilai tidak ditemukan"
// @Router /nilai/{id} [get]
func (c *NilaiController) GetNilaiByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetNilaiByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get nilai successfully", resp)
}

// GetNilaiByPeserta godoc
// @Summary List riwayat nilai satu peserta
// @Tags Nilai
// @Produce json
// @Param id_peserta path string true "ID peserta (uuid)"
// @Param page query int false "Halaman" default(1)
// @Param page_size query int false "Jumlah per halaman" default(10)
// @Success 200 {object} helpers.Response{data=dto.NilaiListResponse} "Get nilai by peserta successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Router /nilai/peserta/{id_peserta} [get]
func (c *NilaiController) GetNilaiByPeserta(ctx *fiber.Ctx) error {
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

	resp, err := c.service.GetNilaiByPeserta(idPeserta, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get nilai by peserta successfully", resp)
}

// GetNilaiByJadwal godoc
// @Summary List nilai semua peserta untuk satu jadwal
// @Tags Nilai
// @Produce json
// @Param id_jadwal path string true "ID jadwal (uuid)"
// @Param page query int false "Halaman" default(1)
// @Param page_size query int false "Jumlah per halaman" default(10)
// @Success 200 {object} helpers.Response{data=dto.NilaiListResponse} "Get nilai by jadwal successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Router /nilai/jadwal/{id_jadwal} [get]
func (c *NilaiController) GetNilaiByJadwal(ctx *fiber.Ctx) error {
	idJadwal := ctx.Params("id_jadwal")
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

	resp, err := c.service.GetNilaiByJadwal(idJadwal, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get nilai by jadwal successfully", resp)
}

// UpdateNilai godoc
// @Summary Update data nilai / selesaikan ujian
// @Description Mengisi wkt_selesai akan otomatis memicu penghitungan ulang nilai akhir berdasarkan jawaban yang sudah masuk.
// @Tags Nilai
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID nilai (uuid)"
// @Param request body dto.UpdateNilaiRequest true "Field yang diubah"
// @Success 200 {object} helpers.Response{data=dto.NilaiResponse} "Update nilai successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal"
// @Router /nilai/{id} [put]
func (c *NilaiController) UpdateNilai(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateNilaiRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateNilai(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update nilai successfully", resp)
}

// DeleteNilai godoc
// @Summary Hapus data nilai (soft delete)
// @Tags Nilai
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID nilai (uuid)"
// @Success 200 {object} helpers.Response "Delete nilai successfully"
// @Failure 400 {object} helpers.Response "Gagal menghapus"
// @Router /nilai/{id} [delete]
func (c *NilaiController) DeleteNilai(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteNilai(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete nilai successfully", nil)
}

// RestoreNilai godoc
// @Summary Restore data nilai yang sudah dihapus
// @Tags Nilai
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID nilai (uuid)"
// @Success 200 {object} helpers.Response "Restore nilai successfully"
// @Failure 400 {object} helpers.Response "Gagal restore"
// @Router /nilai/{id}/restore [patch]
func (c *NilaiController) RestoreNilai(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreNilai(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore nilai successfully", nil)
}

// MulaiUjian godoc
// @Summary Mulai (atau lanjutkan) pengerjaan ujian
// @Description Dipanggil peserta saat membuka jadwal ujian. Jika belum pernah mulai, sesi nilai baru dibuat beserta seluruh baris jawaban kosong (soal diacak jika acak_soal aktif) dan status 201 dikembalikan. Jika sesi sudah ada dan belum wkt_selesai, dianggap resume dan status 200 dikembalikan (token tidak perlu dikirim ulang saat resume). Jika sudah wkt_selesai, request ditolak. Body token bersifat opsional secara umum — hanya wajib diisi jika jadwal ini mengaktifkan wajib_token; ambil kode yang berlaku lewat GET /ujian-token/current. Header User-Agent atau X-Requested-With request ini juga wajib salah satunya terdaftar di whitelist /user-agent, kalau tidak ada yang cocok request ditolak 403.
// @Tags Nilai
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_jadwal path string true "ID jadwal (uuid)"
// @Param request body dto.MulaiUjianRequest false "Token ujian (wajib hanya jika jadwal.wajib_token aktif)"
// @Success 201 {object} helpers.Response{data=dto.NilaiResponse} "Mulai ujian successfully (sesi baru)"
// @Success 200 {object} helpers.Response{data=dto.NilaiResponse} "Lanjutkan ujian successfully (resume sesi yang sudah ada)"
// @Failure 400 {object} helpers.Response "Jadwal tidak ditemukan, ujian sudah pernah selesai dikerjakan, atau token ujian tidak valid/kedaluwarsa"
// @Failure 401 {object} helpers.Response "Token JWT tidak valid"
// @Failure 403 {object} helpers.Response "Kombinasi User-Agent + X-Requested-With tidak terdaftar di whitelist"
// @Router /nilai/mulai-ujian/{id_jadwal} [post]
func (c *NilaiController) MulaiUjian(ctx *fiber.Ctx) error {
	idJadwal := ctx.Params("id_jadwal")
	if idJadwal == "" {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "id_jadwal tidak boleh kosong", nil)
	}

	idPeserta, err := getPesertaIDFromToken(ctx)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, err.Error(), nil)
	}

	var req dto.MulaiUjianRequest
	_ = ctx.BodyParser(&req) // token opsional; body kosong tetap valid

	resp, isNew, err := c.service.MulaiUjian(idPeserta, idJadwal, req.Token)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	if isNew {
		return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Mulai ujian successfully", resp)
	}
	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Lanjutkan ujian successfully", resp)
}

// getPesertaIDFromToken mengambil user_id (= id_peserta) dari JWT yang sudah divalidasi middleware.JWTAuth.
func getPesertaIDFromToken(ctx *fiber.Ctx) (string, error) {
	userToken, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return "", fmt.Errorf("Unauthorized")
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Invalid token claims")
	}
	idPeserta, ok := claims["user_id"].(string)
	if !ok || idPeserta == "" {
		return "", fmt.Errorf("user_id tidak ditemukan di token")
	}
	return idPeserta, nil
}

// NextSection godoc
// @Summary Maju ke section berikutnya
// @Description Memindahkan frontier section peserta ke section berikutnya, hanya jika waktu minimal (durasi_menit_minimal) pada section saat ini sudah terlampaui. Jika belum, response tetap 200 namun boleh_lanjut=false beserta sisa_detik yang harus ditunggu — frontend bisa poll endpoint ini atau /section-status untuk menampilkan countdown. Section yang sudah pernah dibuka tetap bisa diakses bebas (navigasi mundur), endpoint ini hanya menggerakkan batas terjauh (frontier).
// @Tags Nilai
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID nilai (uuid), yaitu ID sesi pengerjaan ujian milik peserta yang login"
// @Success 200 {object} helpers.Response{data=dto.SectionProgressResponse} "Berhasil maju ke section berikutnya, atau ditolak sementara (lihat boleh_lanjut & sisa_detik)"
// @Failure 400 {object} helpers.Response "Ujian tidak memakai section, sudah selesai, atau sudah di section terakhir"
// @Failure 401 {object} helpers.Response "Token tidak valid atau bukan pemilik sesi ini"
// @Router /nilai/{id}/next-section [post]
func (c *NilaiController) NextSection(ctx *fiber.Ctx) error {
	idNilai := ctx.Params("id")

	idPeserta, err := getPesertaIDFromToken(ctx)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, err.Error(), nil)
	}

	resp, err := c.service.AdvanceSection(idNilai, idPeserta)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Next section successfully", resp)
}

// SectionStatus godoc
// @Summary Cek status section aktif (read-only, untuk polling countdown)
// @Description Sama seperti next-section tapi tidak mengubah apapun — cocok dipanggil berkala oleh frontend untuk menampilkan sisa waktu sebelum tombol "lanjut ke section berikutnya" aktif.
// @Tags Nilai
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID nilai (uuid)"
// @Success 200 {object} helpers.Response{data=dto.SectionProgressResponse} "Get section status successfully"
// @Failure 400 {object} helpers.Response "Ujian tidak memakai section"
// @Failure 401 {object} helpers.Response "Token tidak valid atau bukan pemilik sesi ini"
// @Router /nilai/{id}/section-status [get]
func (c *NilaiController) SectionStatus(ctx *fiber.Ctx) error {
	idNilai := ctx.Params("id")

	idPeserta, err := getPesertaIDFromToken(ctx)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, err.Error(), nil)
	}

	resp, err := c.service.GetSectionStatus(idNilai, idPeserta)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get section status successfully", resp)
}

// ExportNilai godoc
// @Summary Export nilai satu jadwal ke Excel (per kelas, dalam ZIP)
// @Description Menghasilkan file ZIP berisi satu file .xlsx per kelas yang terdaftar pada jadwal tersebut.
// @Tags Nilai
// @Produce application/zip
// @Security BearerAuth
// @Param id_jadwal path string true "ID jadwal (uuid)"
// @Success 200 {file} file "File ZIP berisi export nilai per kelas"
// @Failure 400 {object} helpers.Response "Jadwal tidak ditemukan atau tidak ada kelas terdaftar"
// @Router /nilai/export/{id_jadwal} [get]
func (c *NilaiController) ExportNilai(ctx *fiber.Ctx) error {
	idJadwal := ctx.Params("id_jadwal")
	if idJadwal == "" {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "id_jadwal tidak boleh kosong", nil)
	}

	result, err := c.service.ExportNilaiByJadwal(idJadwal)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	safeNama := strings.ReplaceAll(result.NamaUjian, " ", "_")
	filename := fmt.Sprintf("export_nilai_%s.zip", safeNama)

	ctx.Set("Content-Type", "application/zip")
	ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return ctx.Send(result.ZipBytes)
}

// AnalisisSoal godoc
// @Summary Export analisis jawaban per soal ke Excel (per kelas, dalam ZIP)
// @Description Menghasilkan ZIP berisi satu .xlsx per kelas: baris = peserta, kolom = soal (diurutkan berdasarkan no_soal, konsisten walau acak_soal aktif). Cell berisi huruf jawaban peserta, diwarnai hijau (benar), merah (salah), atau kuning (belum dijawab) — cell kosong tanpa warna berarti peserta belum pernah mulai ujian ini. Baris paling bawah menunjukkan jumlah peserta yang benar per soal, untuk melihat soal mana yang paling banyak dijawab salah.
// @Tags Nilai
// @Produce application/zip
// @Security BearerAuth
// @Param id_jadwal path string true "ID jadwal (uuid)"
// @Success 200 {file} file "File ZIP berisi analisis jawaban per kelas"
// @Failure 400 {object} helpers.Response "Jadwal tidak ditemukan, belum memiliki soal, atau tidak ada kelas terdaftar"
// @Router /nilai/analisis/{id_jadwal} [get]
func (c *NilaiController) AnalisisSoal(ctx *fiber.Ctx) error {
	idJadwal := ctx.Params("id_jadwal")
	if idJadwal == "" {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "id_jadwal tidak boleh kosong", nil)
	}

	result, err := c.service.AnalisisSoalByJadwal(idJadwal)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	safeNama := strings.ReplaceAll(result.NamaUjian, " ", "_")
	filename := fmt.Sprintf("analisis_soal_%s.zip", safeNama)

	ctx.Set("Content-Type", "application/zip")
	ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return ctx.Send(result.ZipBytes)
}
