package database

import (
	"backend/internal/database/seeders"
	"fmt"

	"gorm.io/gorm"
)

func RunSeeders(db *gorm.DB) error {
	if err := seeders.SeedMapel(db); err != nil {
		return fmt.Errorf("failed to seed mapel: %w", err)
	}

	if err := seeders.SeedBankSoal(db); err != nil {
		return fmt.Errorf("failed to seed bank_soal: %w", err)
	}

	if err := seeders.SeedJurusan(db); err != nil {
		return fmt.Errorf("failed to seed jurusan: %w", err)
	}

	if err := seeders.SeedKelas(db); err != nil {
		return fmt.Errorf("failed to seed kelas: %w", err)
	}

	if err := seeders.SeedPeserta(db); err != nil {
		return fmt.Errorf("failed to seed peserta: %w", err)
	}

	return nil
}
