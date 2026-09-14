package seeders

import (
	"backend/internal/modules/kelas/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func SeedKelas(db *gorm.DB) error {
	type JurusanRow struct {
		ID          string
		NamaJurusan string
	}

	var jurusanList []JurusanRow
	if err := db.Table("jurusan").
		Where("deleted_at IS NULL").
		Select("id, nama_jurusan").
		Scan(&jurusanList).Error; err != nil {
		return err
	}

	if len(jurusanList) == 0 {
		return nil
	}

	tingkatan := []string{"X", "XI", "XII"}
	var kelasList []model.Kelas

	for _, j := range jurusanList {
		for _, tingkat := range tingkatan {
			kelasList = append(kelasList, model.Kelas{
				IDJurusan: j.ID,
				NamaKelas: fmt.Sprintf("%s - %s", tingkat, j.NamaJurusan),
				Tingkat:   tingkat,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}

	return db.CreateInBatches(kelasList, 100).Error
}
