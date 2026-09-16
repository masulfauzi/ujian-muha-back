package routes

import (
	"backend/internal/middleware"
	jawabanrepo "backend/internal/modules/jawaban/repository"
	"backend/internal/modules/nilai/controller"
	"backend/internal/modules/nilai/repository"
	"backend/internal/modules/nilai/service"
	sectionrepo "backend/internal/modules/section/repository"
	useragentrepo "backend/internal/modules/user_agent/repository"
	useragentservice "backend/internal/modules/user_agent/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupNilaiRoutes(app *fiber.App, db *gorm.DB) {
	repo              := repository.NewNilaiRepository(db)
	jawabanRepository := jawabanrepo.NewJawabanRepository(db)
	sectionRepository := sectionrepo.NewSectionRepository(db)
	svc               := service.NewNilaiService(repo, jawabanRepository, sectionRepository, db)
	ctrl              := controller.NewNilaiController(svc)

	userAgentRepository := useragentrepo.NewUserAgentRepository(db)
	userAgentService    := useragentservice.NewUserAgentService(userAgentRepository)

	api   := app.Group("/api")
	nilai := api.Group("/nilai")

	// Dua-segmen path: aman dari konflik dengan /:id (beda kedalaman)
	nilai.Get("/peserta/:id_peserta", ctrl.GetNilaiByPeserta)
	nilai.Get("/jadwal/:id_jadwal", ctrl.GetNilaiByJadwal)
	nilai.Get("/export/:id_jadwal", middleware.JWTAuth(), ctrl.ExportNilai)
	nilai.Get("/analisis/:id_jadwal", middleware.JWTAuth(), ctrl.AnalisisSoal)
	nilai.Post("/mulai-ujian/:id_jadwal", middleware.JWTAuth(), middleware.CheckAllowedUserAgent(userAgentService), ctrl.MulaiUjian)

	nilai.Post("/", middleware.JWTAuth(), ctrl.CreateNilai)
	nilai.Get("/", ctrl.GetAllNilai)
	nilai.Get("/:id", ctrl.GetNilaiByID)
	nilai.Get("/:id/section-status", middleware.JWTAuth(), ctrl.SectionStatus)
	nilai.Post("/:id/next-section", middleware.JWTAuth(), ctrl.NextSection)
	nilai.Put("/:id", middleware.JWTAuth(), ctrl.UpdateNilai)
	nilai.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteNilai)
	nilai.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreNilai)
}
