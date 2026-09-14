package seeders

import (
	"backend/internal/modules/peserta/model"
	"backend/internal/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func SeedPeserta(db *gorm.DB) error {
	type KelasRow struct {
		ID        string
		NamaKelas string
	}

	var kelasList []KelasRow
	if err := db.Table("kelas").
		Where("deleted_at IS NULL").
		Select("id, nama_kelas").
		Scan(&kelasList).Error; err != nil {
		return err
	}

	if len(kelasList) == 0 {
		return nil
	}

	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		return err
	}

	var pesertaList []model.Peserta

	for _, k := range kelasList {
		for i := 1; i <= 5; i++ {
			username := fmt.Sprintf("peserta_%s_%d", k.ID[:8], i)
			pesertaList = append(pesertaList, model.Peserta{
				Nama:      fmt.Sprintf("Peserta %d - %s", i, k.NamaKelas),
				IDKelas:   k.ID,
				Username:  username,
				Password:  hashedPassword,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}

	return db.CreateInBatches(pesertaList, 100).Error
}
