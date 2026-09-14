package middleware

import (
	"backend/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORS() fiber.Handler {
	appConfig := configs.GetAppConfig()

	return cors.New(cors.Config{
		AllowOrigins:     appConfig.FrontendURL,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Content-Type,Authorization",
		ExposeHeaders:    "Content-Type,Authorization",
		AllowCredentials: true,
		MaxAge:           3600,
	})
}
