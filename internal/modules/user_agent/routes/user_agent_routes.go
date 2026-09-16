package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/user_agent/controller"
	"backend/internal/modules/user_agent/repository"
	"backend/internal/modules/user_agent/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserAgentRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewUserAgentRepository(db)
	svc := service.NewUserAgentService(repo)
	ctrl := controller.NewUserAgentController(svc)

	api := app.Group("/api")
	userAgent := api.Group("/user-agent")

	userAgent.Post("/", middleware.JWTAuth(), ctrl.CreateUserAgent)
	userAgent.Get("/", middleware.JWTAuth(), ctrl.GetAllUserAgent)
	userAgent.Get("/:id", middleware.JWTAuth(), ctrl.GetUserAgentByID)
	userAgent.Put("/:id", middleware.JWTAuth(), ctrl.UpdateUserAgent)
	userAgent.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteUserAgent)
	userAgent.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreUserAgent)
}
