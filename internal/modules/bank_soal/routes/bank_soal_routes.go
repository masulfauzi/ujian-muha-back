package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/bank_soal/controller"
	"backend/internal/modules/bank_soal/repository"
	"backend/internal/modules/bank_soal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupBankSoalRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewBankSoalRepository(db)
	svc := service.NewBankSoalService(repo)
	ctrl := controller.NewBankSoalController(svc)

	api := app.Group("/api")
	bankSoal := api.Group("/bank-soal")

	bankSoal.Post("/", middleware.JWTAuth(), ctrl.CreateBankSoal)
	bankSoal.Get("/", ctrl.GetAllBankSoal)
	bankSoal.Get("/mapel/:mapel_id", ctrl.GetBankSoalByMapel)
	bankSoal.Get("/:id", ctrl.GetBankSoalByID)
	bankSoal.Put("/:id", middleware.JWTAuth(), ctrl.UpdateBankSoal)
	bankSoal.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteBankSoal)
	bankSoal.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreBankSoal)
}
