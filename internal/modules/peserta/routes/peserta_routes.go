package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/peserta/controller"
	"backend/internal/modules/peserta/repository"
	"backend/internal/modules/peserta/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupPesertaRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewPesertaRepository(db)
	svc := service.NewPesertaService(repo)
	ctrl := controller.NewPesertaController(svc)

	api := app.Group("/api")
	peserta := api.Group("/peserta")

	peserta.Post("/", middleware.JWTAuth(), ctrl.CreatePeserta)
	peserta.Get("/", ctrl.GetAllPeserta)
	peserta.Get("/kartu-ujian/:id_kelas", middleware.JWTAuth(), ctrl.DownloadKartuUjian)
	peserta.Get("/:id", ctrl.GetPesertaByID)
	peserta.Put("/:id", middleware.JWTAuth(), ctrl.UpdatePeserta)
	peserta.Delete("/:id", middleware.JWTAuth(), ctrl.DeletePeserta)
	peserta.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestorePeserta)
}
