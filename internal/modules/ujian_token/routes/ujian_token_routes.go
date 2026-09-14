package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/ujian_token/controller"
	"backend/internal/modules/ujian_token/service"

	"github.com/gofiber/fiber/v2"
)

func SetupUjianTokenRoutes(app *fiber.App) {
	svc := service.NewUjianTokenService()
	ctrl := controller.NewUjianTokenController(svc)

	api := app.Group("/api")
	ujianToken := api.Group("/ujian-token")

	ujianToken.Get("/current", middleware.JWTAuth(), ctrl.GetCurrentToken)
}
