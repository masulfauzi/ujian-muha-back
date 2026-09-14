package seeders

import (
	"backend/internal/modules/bank_soal/model"
	"time"

	"gorm.io/gorm"
)

func SeedBankSoal(db *gorm.DB) error {
	var mapelIDs []string
	if err := db.Model(&struct{}{}).
		Table("mapel").
		Where("deleted_at IS NULL").
		Limit(5).
		Pluck("id", &mapelIDs).Error; err != nil {
		return err
	}

	if len(mapelIDs) == 0 {
		return nil
	}

	bankSoals := []model.BankSoal{
		{
			NamaBankSoal: "Bank Soal Matematika Dasar",
			IdMapel:      mapelIDs[0],
			JmlSoal:      50,
			Deskripsi:    "Kumpulan soal matematika level dasar",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			NamaBankSoal: "Bank Soal Matematika Lanjutan",
			IdMapel:      mapelIDs[0],
			JmlSoal:      75,
			Deskripsi:    "Kumpulan soal matematika level lanjutan",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			NamaBankSoal: "Bank Soal Bahasa Indonesia Umum",
			IdMapel:      mapelIDs[1],
			JmlSoal:      40,
			Deskripsi:    "Soal umum bahasa Indonesia",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			NamaBankSoal: "Bank Soal Grammar Bahasa Inggris",
			IdMapel:      mapelIDs[2],
			JmlSoal:      60,
			Deskripsi:    "Soal grammar bahasa Inggris",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			NamaBankSoal: "Bank Soal IPA Fisika",
			IdMapel:      mapelIDs[3],
			JmlSoal:      45,
			Deskripsi:    "Soal fisika IPA",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	return db.CreateInBatches(bankSoals, 100).Error
}
