package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/jadwal/controller"
	"backend/internal/modules/jadwal/repository"
	"backend/internal/modules/jadwal/service"
	jadwalkelasrepo "backend/internal/modules/jadwal_kelas/repository"
	pesertarepo "backend/internal/modules/peserta/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupJadwalRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewJadwalRepository(db)
	jadwalKelasRepo := jadwalkelasrepo.NewJadwalKelasRepository(db)
	pesertaRepo := pesertarepo.NewPesertaRepository(db)
	svc := service.NewJadwalService(repo, jadwalKelasRepo, pesertaRepo)
	ctrl := controller.NewJadwalController(svc)

	api := app.Group("/api")
	jadwal := api.Group("/jadwal")

	jadwal.Get("/aktif/hari-ini", middleware.JWTAuth(), ctrl.GetJadwalAktifHariIni)

	jadwal.Post("/", middleware.JWTAuth(), ctrl.CreateJadwal)
	jadwal.Get("/", ctrl.GetAllJadwal)
	jadwal.Get("/bank-soal/:bank_soal_id", ctrl.GetJadwalByBankSoal)
	jadwal.Get("/:id", ctrl.GetJadwalByID)
	jadwal.Put("/:id", middleware.JWTAuth(), ctrl.UpdateJadwal)
	jadwal.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteJadwal)
	jadwal.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreJadwal)
}
