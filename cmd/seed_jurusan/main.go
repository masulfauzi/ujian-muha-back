package main

import (
	"log"
	"time"

	"backend/configs"
	"backend/internal/database"
	"backend/internal/modules/jurusan/model"
)

func main() {
	if err := configs.LoadEnv(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	db := database.DB
	jurusans := []model.Jurusan{
		{NamaJurusan: "Teknik Komputer dan Jaringan", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{NamaJurusan: "Rekayasa Perangkat Lunak", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{NamaJurusan: "Multimedia", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{NamaJurusan: "Akuntansi", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{NamaJurusan: "Administrasi Perkantoran", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	// Use CreateInBatches which handles duplicates gracefully
	if err := db.CreateInBatches(jurusans, 100).Error; err != nil {
		log.Printf("Error seeding jurusan (may already exist): %v\n", err)
	}

	// Verify by counting records
	var count int64
	db.Model(&model.Jurusan{}).Where("deleted_at IS NULL").Count(&count)

	log.Printf("✅ Seeding completed! Total jurusan records: %d\n", count)
	_ = database.Close()
}
