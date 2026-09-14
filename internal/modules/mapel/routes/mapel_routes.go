package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/mapel/controller"
	"backend/internal/modules/mapel/repository"
	"backend/internal/modules/mapel/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupMapelRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewMapelRepository(db)
	svc := service.NewMapelService(repo)
	ctrl := controller.NewMapelController(svc)

	api := app.Group("/api")
	mapel := api.Group("/mapel")

	mapel.Post("/", middleware.JWTAuth(), ctrl.CreateMapel)
	mapel.Get("/", ctrl.GetAllMapel)
	mapel.Get("/:id", ctrl.GetMapelByID)
	mapel.Put("/:id", middleware.JWTAuth(), ctrl.UpdateMapel)
	mapel.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteMapel)
	mapel.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreMapel)
}
