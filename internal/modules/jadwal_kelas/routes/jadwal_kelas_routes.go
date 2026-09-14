package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/jadwal_kelas/controller"
	"backend/internal/modules/jadwal_kelas/repository"
	"backend/internal/modules/jadwal_kelas/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupJadwalKelasRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewJadwalKelasRepository(db)
	svc  := service.NewJadwalKelasService(repo)
	ctrl := controller.NewJadwalKelasController(svc)

	api       := app.Group("/api")
	jadwalKelas := api.Group("/jadwal-kelas")

	jadwalKelas.Post("/", middleware.JWTAuth(), ctrl.CreateJadwalKelas)
	jadwalKelas.Get("/", ctrl.GetAllJadwalKelas)
	jadwalKelas.Get("/:id", ctrl.GetJadwalKelasByID)
	jadwalKelas.Put("/:id", middleware.JWTAuth(), ctrl.UpdateJadwalKelas)
	jadwalKelas.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteJadwalKelas)
}
