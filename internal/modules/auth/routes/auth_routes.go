package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/auth/controller"
	authrepo "backend/internal/modules/auth/repository"
	authsvc "backend/internal/modules/auth/service"
	pesertarepo "backend/internal/modules/peserta/repository"
	pesertasvc "backend/internal/modules/peserta/service"
	userrepo "backend/internal/modules/user/repository"
	usersvc "backend/internal/modules/user/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupAuthRoutes(app *fiber.App, db *gorm.DB) {
	authRepo := authrepo.NewAuthRepository(db)
	authSvc := authsvc.NewAuthService(authRepo)

	userRepo := userrepo.NewUserRepository(db)
	userSvc := usersvc.NewUserService(userRepo)

	pesertaRepo := pesertarepo.NewPesertaRepository(db)
	pesertaSvc := pesertasvc.NewPesertaService(pesertaRepo)

	ctrl := controller.NewAuthController(authSvc, userSvc, pesertaSvc)

	api := app.Group("/api")
	auth := api.Group("/auth")

	auth.Post("/register", ctrl.Register)
	auth.Post("/login", ctrl.Login)
	auth.Get("/me", middleware.JWTAuth(), ctrl.GetCurrentUser)
}
