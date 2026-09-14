package main

import (
	"log"

	"backend/configs"
	"backend/internal/database"
)

func main() {
	if err := configs.LoadEnv(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.RunMigrations(database.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := database.RunSeeders(database.DB); err != nil {
		log.Fatalf("Failed to run seeders: %v", err)
	}

	log.Println("Seeders ran successfully!")
	_ = database.Close()
}
