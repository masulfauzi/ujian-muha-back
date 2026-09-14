package seeders

import (
	"backend/internal/modules/mapel/model"
	"time"

	"gorm.io/gorm"
)

func SeedMapel(db *gorm.DB) error {
	mapels := []model.Mapel{
		{
			NamaMapel: "Matematika",
			KodeMapel: "MAT",
			Deskripsi: "Pelajaran Matematika",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			NamaMapel: "Bahasa Indonesia",
			KodeMapel: "IND",
			Deskripsi: "Pelajaran Bahasa Indonesia",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			NamaMapel: "Bahasa Inggris",
			KodeMapel: "ENG",
			Deskripsi: "Pelajaran Bahasa Inggris",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			NamaMapel: "IPA",
			KodeMapel: "IPA",
			Deskripsi: "Ilmu Pengetahuan Alam",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			NamaMapel: "IPS",
			KodeMapel: "IPS",
			Deskripsi: "Ilmu Pengetahuan Sosial",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return db.CreateInBatches(mapels, 100).Error
}
