package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/section/controller"
	"backend/internal/modules/section/repository"
	"backend/internal/modules/section/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupSectionRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewSectionRepository(db)
	svc := service.NewSectionService(repo, db)
	ctrl := controller.NewSectionController(svc)

	api := app.Group("/api")
	section := api.Group("/section")

	// Dua-segmen path: aman dari konflik dengan /:id (beda kedalaman)
	section.Post("/jadwal/:id_jadwal/define", middleware.JWTAuth(), ctrl.DefineSections)
	section.Get("/jadwal/:id_jadwal", ctrl.GetSectionsByJadwal)

	section.Get("/:id", ctrl.GetSectionByID)
	section.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteSection)
	section.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreSection)
}
