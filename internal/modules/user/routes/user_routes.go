package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/user/controller"
	"backend/internal/modules/user/repository"
	"backend/internal/modules/user/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	ctrl := controller.NewUserController(svc)

	api := app.Group("/api")
	users := api.Group("/users")

	users.Get("/", ctrl.GetAll)
	users.Get("/:id", ctrl.GetByID)
	users.Post("/", middleware.JWTAuth(), ctrl.Create)
	users.Put("/:id", middleware.JWTAuth(), ctrl.Update)
	users.Delete("/:id", middleware.JWTAuth(), ctrl.Delete)
}
