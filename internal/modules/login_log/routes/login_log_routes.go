package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/login_log/controller"
	"backend/internal/modules/login_log/repository"
	"backend/internal/modules/login_log/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupLoginLogRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewLoginLogRepository(db)
	svc := service.NewLoginLogService(repo)
	ctrl := controller.NewLoginLogController(svc)

	api := app.Group("/api")
	loginLog := api.Group("/login-log")

	loginLog.Get("/", middleware.JWTAuth(), ctrl.GetAllLoginLog)
}
