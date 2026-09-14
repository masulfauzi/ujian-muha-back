package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/kelas/controller"
	"backend/internal/modules/kelas/repository"
	"backend/internal/modules/kelas/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupKelasRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewKelasRepository(db)
	svc := service.NewKelasService(repo)
	ctrl := controller.NewKelasController(svc)

	api := app.Group("/api")
	kelas := api.Group("/kelas")

	kelas.Post("/", middleware.JWTAuth(), ctrl.CreateKelas)
	kelas.Get("/", ctrl.GetAllKelas)
	kelas.Get("/:id", ctrl.GetKelasByID)
	kelas.Put("/:id", middleware.JWTAuth(), ctrl.UpdateKelas)
	kelas.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteKelas)
	kelas.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreKelas)
}
