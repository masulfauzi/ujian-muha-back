package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend/configs"
	_ "backend/docs/swagger"
	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/storage"
	authroutes "backend/internal/modules/auth/routes"
	banksoalroutes "backend/internal/modules/bank_soal/routes"
	jadwalroutes "backend/internal/modules/jadwal/routes"
	jadwalkelasroutes "backend/internal/modules/jadwal_kelas/routes"
	jawabanroutes "backend/internal/modules/jawaban/routes"
	nilairoutes "backend/internal/modules/nilai/routes"
	jurusanroutes "backend/internal/modules/jurusan/routes"
	kelasroutes "backend/internal/modules/kelas/routes"
	loginlogroutes "backend/internal/modules/login_log/routes"
	mapelroutes "backend/internal/modules/mapel/routes"
	pesertaroutes "backend/internal/modules/peserta/routes"
	sectionroutes "backend/internal/modules/section/routes"
	soalroutes "backend/internal/modules/soal/routes"
	ujiantokenroutes "backend/internal/modules/ujian_token/routes"
	useragentroutes "backend/internal/modules/user_agent/routes"
	userroutes "backend/internal/modules/user/routes"

	"github.com/gofiber/fiber/v2"
	fiberswagger "github.com/gofiber/swagger"
)

// @title Ujian Backend API
// @version 1.0
// @description API untuk sistem ujian online (bank soal, jadwal, pengerjaan ujian per section, penilaian).
// @contact.name Tim Backend Ujian
// @host localhost:3000
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token JWT dengan format: Bearer {token}

func main() {
	if err := configs.LoadEnv(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	appConfig := configs.GetAppConfig()

	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.RunMigrations(database.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := storage.InitMinioClient(); err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: appConfig.Name,
	})

	app.Use(middleware.CORS())
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())

	setupRoutes(app, appConfig)

	go func() {
		addr := fmt.Sprintf(":%d", appConfig.Port)
		log.Printf("Starting server on %s\n", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	_ = app.Shutdown()
	_ = database.Close()
	log.Println("Server shut down successfully")
}

func setupRoutes(app *fiber.App, appConfig *configs.AppConfig) {
	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":     "ok",
			"service":    "Fiber Backend API",
			"server-no":  appConfig.ServerNo,
		})
	})

	app.Static("/uploads", "./uploads")
	app.Get("/swagger/*", fiberswagger.New())

	authroutes.SetupAuthRoutes(app, database.DB)
	userroutes.SetupUserRoutes(app, database.DB)
	mapelroutes.SetupMapelRoutes(app, database.DB)
	banksoalroutes.SetupBankSoalRoutes(app, database.DB)
	soalroutes.SetupSoalRoutes(app, database.DB)
	jurusanroutes.SetupJurusanRoutes(app, database.DB)
	kelasroutes.SetupKelasRoutes(app, database.DB)
	jadwalroutes.SetupJadwalRoutes(app, database.DB)
	jadwalkelasroutes.SetupJadwalKelasRoutes(app, database.DB)
	pesertaroutes.SetupPesertaRoutes(app, database.DB)
	nilairoutes.SetupNilaiRoutes(app, database.DB)
	jawabanroutes.SetupJawabanRoutes(app, database.DB)
	sectionroutes.SetupSectionRoutes(app, database.DB)
	ujiantokenroutes.SetupUjianTokenRoutes(app)
	useragentroutes.SetupUserAgentRoutes(app, database.DB)
	loginlogroutes.SetupLoginLogRoutes(app, database.DB)
}
