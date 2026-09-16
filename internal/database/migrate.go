package database

import (
	banksoalmodel    "backend/internal/modules/bank_soal/model"
	jadwalmodel      "backend/internal/modules/jadwal/model"
	jadwalkelasmodel "backend/internal/modules/jadwal_kelas/model"
	jawabanmodel     "backend/internal/modules/jawaban/model"
	jurusanmodel     "backend/internal/modules/jurusan/model"
	kelasmodel       "backend/internal/modules/kelas/model"
	loginlogmodel    "backend/internal/modules/login_log/model"
	mapelmodel       "backend/internal/modules/mapel/model"
	nilaimodel       "backend/internal/modules/nilai/model"
	pesertamodel     "backend/internal/modules/peserta/model"
	sectionmodel     "backend/internal/modules/section/model"
	soalmodel        "backend/internal/modules/soal/model"
	useragentmodel   "backend/internal/modules/user_agent/model"
	usermodel        "backend/internal/modules/user/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Drop old non-partial unique index on jurusan.nama_jurusan before migrating
	// so AutoMigrate can create the correct partial index (where:deleted_at IS NULL)
	db.Exec("DROP INDEX IF EXISTS idx_jurusans_nama_jurusan")
	db.Exec("DROP INDEX IF EXISTS uni_jurusan_nama_jurusan")
	db.Exec("DROP INDEX IF EXISTS idx_jurusan_nama_jurusan")

	// Drop NOT NULL & DEFAULT dari kolom jawaban supaya bisa di-set NULL saat bulk-insert
	db.Exec("ALTER TABLE jawaban ALTER COLUMN jawaban DROP NOT NULL")
	db.Exec("ALTER TABLE jawaban ALTER COLUMN is_benar DROP NOT NULL")
	db.Exec("ALTER TABLE jawaban ALTER COLUMN is_benar DROP DEFAULT")

	// Alter is_benar column type dari boolean ke smallint (0 atau 1)
	db.Exec("ALTER TABLE jawaban ALTER COLUMN is_benar TYPE smallint USING CASE WHEN is_benar THEN 1 ELSE 0 END")

	// Alter acak_soal dan acak_opsi dari boolean ke smallint (0 atau 1).
	// acak_opsi::int::smallint: works for boolean (true→1, false→0) and int/smallint inputs.
	db.Exec("ALTER TABLE jadwal ALTER COLUMN acak_soal TYPE smallint USING acak_soal::int::smallint")
	db.Exec("ALTER TABLE jadwal ALTER COLUMN acak_opsi TYPE smallint USING acak_opsi::int::smallint")

	// Tambah kolom x_requested_with ke user_agent (awalnya cuma ada kolom user_agent) dengan
	// aman untuk baris lama yang sudah ada (backfill '' dulu sebelum di-set NOT NULL), lalu
	// ganti unique index lama (hanya di kolom user_agent) jadi index gabungan
	// (user_agent, x_requested_with) yang dibuat AutoMigrate di bawah.
	db.Exec("ALTER TABLE user_agent ADD COLUMN IF NOT EXISTS x_requested_with varchar(255)")
	db.Exec("UPDATE user_agent SET x_requested_with = '' WHERE x_requested_with IS NULL")
	db.Exec("ALTER TABLE user_agent ALTER COLUMN x_requested_with SET NOT NULL")
	db.Exec("DROP INDEX IF EXISTS idx_user_agent_value")
	db.Exec("ALTER TABLE user_agent DROP CONSTRAINT IF EXISTS idx_user_agent_user_agent")

	if err := db.AutoMigrate(
		&usermodel.User{},
		&mapelmodel.Mapel{},
		&banksoalmodel.BankSoal{},
		&soalmodel.Soal{},
		&jurusanmodel.Jurusan{},
		&kelasmodel.Kelas{},
		&jadwalmodel.Jadwal{},
		&jadwalkelasmodel.JadwalKelas{},
		&pesertamodel.Peserta{},
		&nilaimodel.Nilai{},
		&jawabanmodel.Jawaban{},
		&sectionmodel.Section{},
		&useragentmodel.UserAgent{},
		&loginlogmodel.LoginLog{},
	); err != nil {
		return err
	}

	// Buat unique constraint untuk mencegah duplikasi assignment kelas ke jadwal
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_jadwal_kelas_unique ON jadwal_kelas(id_jadwal, id_kelas)")
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_nilai_peserta_jadwal_unique ON nilai(id_peserta, id_jadwal) WHERE deleted_at IS NULL")
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_jawaban_nilai_soal_unique ON jawaban(id_nilai, id_soal) WHERE deleted_at IS NULL")
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_soal_bank_no_unique ON soal(id_bank_soal, no_soal) WHERE deleted_at IS NULL")
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_section_jadwal_urutan_unique ON section(id_jadwal, urutan) WHERE deleted_at IS NULL")

	return nil
}
