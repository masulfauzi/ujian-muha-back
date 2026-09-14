package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/jurusan/controller"
	"backend/internal/modules/jurusan/repository"
	"backend/internal/modules/jurusan/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupJurusanRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewJurusanRepository(db)
	svc := service.NewJurusanService(repo)
	ctrl := controller.NewJurusanController(svc)

	api := app.Group("/api")
	jurusan := api.Group("/jurusan")

	jurusan.Post("/", middleware.JWTAuth(), ctrl.CreateJurusan)
	jurusan.Get("/", ctrl.GetAllJurusan)
	jurusan.Get("/:id", ctrl.GetJurusanByID)
	jurusan.Put("/:id", middleware.JWTAuth(), ctrl.UpdateJurusan)
	jurusan.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteJurusan)
	jurusan.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreJurusan)
}
