package seeders

import (
	"backend/internal/modules/jurusan/model"
	"time"

	"gorm.io/gorm"
)

func SeedJurusan(db *gorm.DB) error {
	jurusans := []model.Jurusan{
		{NamaJurusan: "Teknik Komputer dan Jaringan", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{NamaJurusan: "Rekayasa Perangkat Lunak", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{NamaJurusan: "Multimedia", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{NamaJurusan: "Akuntansi", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{NamaJurusan: "Administrasi Perkantoran", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	return db.CreateInBatches(jurusans, 100).Error
}
