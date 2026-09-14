package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/jawaban/controller"
	"backend/internal/modules/jawaban/repository"
	"backend/internal/modules/jawaban/service"
	nilairepo "backend/internal/modules/nilai/repository"
	sectionrepo "backend/internal/modules/section/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupJawabanRoutes(app *fiber.App, db *gorm.DB) {
	repo              := repository.NewJawabanRepository(db)
	nilaiRepository   := nilairepo.NewNilaiRepository(db)
	sectionRepository := sectionrepo.NewSectionRepository(db)
	svc               := service.NewJawabanService(repo, nilaiRepository, sectionRepository)
	ctrl              := controller.NewJawabanController(svc)

	api     := app.Group("/api")
	jawaban := api.Group("/jawaban")

	// Dua-segmen path: aman dari konflik dengan /:id (beda kedalaman)
	jawaban.Get("/nilai/:id_nilai", ctrl.GetJawabanByNilai)
	jawaban.Get("/nilai/:id_nilai/section/:id_section", ctrl.GetJawabanByNilaiSection)
	jawaban.Get("/peserta/:id_peserta", ctrl.GetJawabanByPeserta)

	jawaban.Post("/", middleware.JWTAuth(), ctrl.CreateJawaban)
	jawaban.Get("/", ctrl.GetAllJawaban)
	jawaban.Get("/:id", ctrl.GetJawabanByID)
	jawaban.Put("/:id", middleware.JWTAuth(), ctrl.UpdateJawaban)
	jawaban.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteJawaban)
	jawaban.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreJawaban)
}
