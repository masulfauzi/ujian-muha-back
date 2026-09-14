package repository

import (
	"backend/internal/modules/nilai/model"
	"time"

	"gorm.io/gorm"
)

type NilaiWithDetail struct {
	ID                string  `gorm:"column:id"`
	IDPeserta         string  `gorm:"column:id_peserta"`
	NamaPeserta       string  `gorm:"column:nama_peserta"`
	IDJadwal          string  `gorm:"column:id_jadwal"`
	NamaUjian         string  `gorm:"column:nama_ujian"`
	Nilai             float64 `gorm:"column:nilai"`
	WktMulai          *string `gorm:"column:wkt_mulai"`
	AktivitasTerakhir *string `gorm:"column:aktivitas_terakhir"`
	WktSelesai        *string `gorm:"column:wkt_selesai"`
	IDSectionAktif    *string `gorm:"column:id_section_aktif"`
	WktMulaiSection   *string `gorm:"column:wkt_mulai_section"`
	CreatedAt         string  `gorm:"column:created_at"`
	UpdatedAt         string  `gorm:"column:updated_at"`
}

type NilaiExportRow struct {
	NamaPeserta string  `gorm:"column:nama_peserta"`
	Username    string  `gorm:"column:username"`
	Nilai       float64 `gorm:"column:nilai"`
	WktMulai    *string `gorm:"column:wkt_mulai"`
	WktSelesai  *string `gorm:"column:wkt_selesai"`
}

type NilaiRepository interface {
	Create(nilai *model.Nilai) error
	GetByID(id string) (*model.Nilai, error)
	GetByIDWithDetail(id string) (*NilaiWithDetail, error)
	GetAllWithDetail(page, pageSize int, idPeserta, idJadwal string) ([]NilaiWithDetail, int64, error)
	GetByPesertaID(idPeserta string, page, pageSize int) ([]NilaiWithDetail, int64, error)
	GetByJadwalID(idJadwal string, page, pageSize int) ([]NilaiWithDetail, int64, error)
	GetByJadwalAndKelas(idJadwal, idKelas string) ([]NilaiExportRow, error)
	CheckDuplicate(idPeserta, idJadwal string) (bool, error)
	GetByPesertaAndJadwal(idPeserta, idJadwal string) (*model.Nilai, error)
	HitungNilai(idNilai string) (float64, error)
	Update(nilai *model.Nilai) error
	Delete(id string) error
	Restore(id string) error
}

type nilaiRepository struct {
	db *gorm.DB
}

func NewNilaiRepository(db *gorm.DB) NilaiRepository {
	return &nilaiRepository{db: db}
}

func (r *nilaiRepository) Create(nilai *model.Nilai) error {
	return r.db.Create(nilai).Error
}

func (r *nilaiRepository) GetByID(id string) (*model.Nilai, error) {
	var nilai model.Nilai
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&nilai).Error
	if err != nil {
		return nil, err
	}
	return &nilai, nil
}

func (r *nilaiRepository) GetByIDWithDetail(id string) (*NilaiWithDetail, error) {
	var result NilaiWithDetail
	err := r.db.
		Table("nilai").
		Select(`
			nilai.id,
			nilai.id_peserta,
			peserta.nama AS nama_peserta,
			nilai.id_jadwal,
			jadwal.nama_ujian,
			nilai.nilai,
			TO_CHAR(nilai.wkt_mulai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_mulai,
			TO_CHAR(nilai.aktivitas_terakhir, 'YYYY-MM-DD HH24:MI:SS') AS aktivitas_terakhir,
			TO_CHAR(nilai.wkt_selesai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_selesai,
			nilai.id_section_aktif,
			TO_CHAR(nilai.wkt_mulai_section, 'YYYY-MM-DD HH24:MI:SS') AS wkt_mulai_section,
			TO_CHAR(nilai.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at,
			TO_CHAR(nilai.updated_at, 'YYYY-MM-DD HH24:MI:SS') AS updated_at
		`).
		Joins("INNER JOIN peserta ON nilai.id_peserta = peserta.id").
		Joins("INNER JOIN jadwal ON nilai.id_jadwal = jadwal.id").
		Where("nilai.id = ? AND nilai.deleted_at IS NULL", id).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *nilaiRepository) GetAllWithDetail(page, pageSize int, idPeserta, idJadwal string) ([]NilaiWithDetail, int64, error) {
	var results []NilaiWithDetail
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	countQuery := r.db.Table("nilai").
		Joins("INNER JOIN peserta ON nilai.id_peserta = peserta.id").
		Joins("INNER JOIN jadwal ON nilai.id_jadwal = jadwal.id").
		Where("nilai.deleted_at IS NULL")

	if idPeserta != "" {
		countQuery = countQuery.Where("nilai.id_peserta = ?", idPeserta)
	}
	if idJadwal != "" {
		countQuery = countQuery.Where("nilai.id_jadwal = ?", idJadwal)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.
		Table("nilai").
		Select(`
			nilai.id,
			nilai.id_peserta,
			peserta.nama AS nama_peserta,
			nilai.id_jadwal,
			jadwal.nama_ujian,
			nilai.nilai,
			TO_CHAR(nilai.wkt_mulai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_mulai,
			TO_CHAR(nilai.aktivitas_terakhir, 'YYYY-MM-DD HH24:MI:SS') AS aktivitas_terakhir,
			TO_CHAR(nilai.wkt_selesai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_selesai,
			nilai.id_section_aktif,
			TO_CHAR(nilai.wkt_mulai_section, 'YYYY-MM-DD HH24:MI:SS') AS wkt_mulai_section,
			TO_CHAR(nilai.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at,
			TO_CHAR(nilai.updated_at, 'YYYY-MM-DD HH24:MI:SS') AS updated_at
		`).
		Joins("INNER JOIN peserta ON nilai.id_peserta = peserta.id").
		Joins("INNER JOIN jadwal ON nilai.id_jadwal = jadwal.id").
		Where("nilai.deleted_at IS NULL")

	if idPeserta != "" {
		query = query.Where("nilai.id_peserta = ?", idPeserta)
	}
	if idJadwal != "" {
		query = query.Where("nilai.id_jadwal = ?", idJadwal)
	}

	err := query.Offset(offset).Limit(pageSize).Scan(&results).Error
	return results, total, err
}

func (r *nilaiRepository) GetByPesertaID(idPeserta string, page, pageSize int) ([]NilaiWithDetail, int64, error) {
	return r.GetAllWithDetail(page, pageSize, idPeserta, "")
}

func (r *nilaiRepository) GetByJadwalID(idJadwal string, page, pageSize int) ([]NilaiWithDetail, int64, error) {
	return r.GetAllWithDetail(page, pageSize, "", idJadwal)
}

func (r *nilaiRepository) GetByJadwalAndKelas(idJadwal, idKelas string) ([]NilaiExportRow, error) {
	var results []NilaiExportRow
	err := r.db.
		Table("peserta").
		Select(`
			peserta.nama AS nama_peserta,
			peserta.username,
			COALESCE(nilai.nilai, 0) AS nilai,
			TO_CHAR(nilai.wkt_mulai AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD HH24:MI:SS') AS wkt_mulai,
			TO_CHAR(nilai.wkt_selesai AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD HH24:MI:SS') AS wkt_selesai
		`).
		Joins("LEFT JOIN nilai ON peserta.id = nilai.id_peserta AND nilai.id_jadwal = ? AND nilai.deleted_at IS NULL", idJadwal).
		Where("peserta.id_kelas = ? AND peserta.deleted_at IS NULL", idKelas).
		Order("peserta.nama ASC").
		Scan(&results).Error
	return results, err
}

func (r *nilaiRepository) CheckDuplicate(idPeserta, idJadwal string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Nilai{}).
		Where("id_peserta = ? AND id_jadwal = ? AND deleted_at IS NULL", idPeserta, idJadwal).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *nilaiRepository) GetByPesertaAndJadwal(idPeserta, idJadwal string) (*model.Nilai, error) {
	var nilai model.Nilai
	err := r.db.Where("id_peserta = ? AND id_jadwal = ? AND deleted_at IS NULL", idPeserta, idJadwal).First(&nilai).Error
	if err != nil {
		return nil, err
	}
	return &nilai, nil
}

func (r *nilaiRepository) HitungNilai(idNilai string) (float64, error) {
	var hasil float64
	err := r.db.Raw(`
		SELECT COALESCE(
			COUNT(*) FILTER (WHERE is_benar = 1) * 100.0 / NULLIF(COUNT(*), 0),
			0
		)
		FROM jawaban
		WHERE id_nilai = ? AND deleted_at IS NULL
	`, idNilai).Scan(&hasil).Error
	return hasil, err
}

func (r *nilaiRepository) Update(nilai *model.Nilai) error {
	return r.db.Save(nilai).Error
}

func (r *nilaiRepository) Delete(id string) error {
	now := time.Now()
	return r.db.Model(&model.Nilai{}).Where("id = ?", id).Update("deleted_at", now).Error
}

func (r *nilaiRepository) Restore(id string) error {
	return r.db.Model(&model.Nilai{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NULL")).Error
}
