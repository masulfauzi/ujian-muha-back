package repository

import (
	"backend/internal/modules/jadwal_kelas/model"

	"gorm.io/gorm"
)

type JadwalKelasWithDetail struct {
	ID           string `gorm:"column:id"`
	IDJadwal     string `gorm:"column:id_jadwal"`
	IDKelas      string `gorm:"column:id_kelas"`
	NamaKelas    string `gorm:"column:nama_kelas"`
	NamaBankSoal string `gorm:"column:nama_bank_soal"`
	WktMulai     string `gorm:"column:wkt_mulai"`
	WktSelesai   string `gorm:"column:wkt_selesai"`
	CreatedAt    string `gorm:"column:created_at"`
	UpdatedAt    string `gorm:"column:updated_at"`
}

type JadwalKelasRepository interface {
	Create(jadwalKelas *model.JadwalKelas) error
	CreateBulk(jadwalKelasList []*model.JadwalKelas) error
	GetByID(id string) (*model.JadwalKelas, error)
	GetByIDWithDetail(id string) (*JadwalKelasWithDetail, error)
	GetAllWithDetail(page, pageSize int, idJadwal string, idKelas string) ([]JadwalKelasWithDetail, int64, error)
	CheckDuplicate(idJadwal, idKelas string) (bool, error)
	Update(jadwalKelas *model.JadwalKelas) error
	Delete(id string) error
	DeleteByJadwalID(idJadwal string) error
}

type jadwalKelasRepository struct {
	db *gorm.DB
}

func NewJadwalKelasRepository(db *gorm.DB) JadwalKelasRepository {
	return &jadwalKelasRepository{db: db}
}

func (r *jadwalKelasRepository) Create(jadwalKelas *model.JadwalKelas) error {
	return r.db.Create(jadwalKelas).Error
}

func (r *jadwalKelasRepository) CreateBulk(jadwalKelasList []*model.JadwalKelas) error {
	return r.db.CreateInBatches(jadwalKelasList, 100).Error
}

func (r *jadwalKelasRepository) GetByID(id string) (*model.JadwalKelas, error) {
	var jadwalKelas model.JadwalKelas
	err := r.db.Where("id = ?", id).First(&jadwalKelas).Error
	if err != nil {
		return nil, err
	}
	return &jadwalKelas, nil
}

func (r *jadwalKelasRepository) GetByIDWithDetail(id string) (*JadwalKelasWithDetail, error) {
	var result JadwalKelasWithDetail
	err := r.db.
		Table("jadwal_kelas").
		Select(`
			jadwal_kelas.id,
			jadwal_kelas.id_jadwal,
			jadwal_kelas.id_kelas,
			kelas.nama_kelas,
			bank_soal.nama_bank_soal,
			jadwal.wkt_mulai,
			jadwal.wkt_selesai,
			jadwal_kelas.created_at,
			jadwal_kelas.updated_at
		`).
		Joins("INNER JOIN jadwal ON jadwal_kelas.id_jadwal = jadwal.id").
		Joins("INNER JOIN kelas ON jadwal_kelas.id_kelas = kelas.id").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal_kelas.id = ?", id).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *jadwalKelasRepository) GetAllWithDetail(page, pageSize int, idJadwal string, idKelas string) ([]JadwalKelasWithDetail, int64, error) {
	var results []JadwalKelasWithDetail
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	countQuery := r.db.Table("jadwal_kelas").
		Joins("INNER JOIN jadwal ON jadwal_kelas.id_jadwal = jadwal.id").
		Joins("INNER JOIN kelas ON jadwal_kelas.id_kelas = kelas.id").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id")

	if idJadwal != "" {
		countQuery = countQuery.Where("jadwal_kelas.id_jadwal = ?", idJadwal)
	}
	if idKelas != "" {
		countQuery = countQuery.Where("jadwal_kelas.id_kelas = ?", idKelas)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.
		Table("jadwal_kelas").
		Select(`
			jadwal_kelas.id,
			jadwal_kelas.id_jadwal,
			jadwal_kelas.id_kelas,
			kelas.nama_kelas,
			bank_soal.nama_bank_soal,
			jadwal.wkt_mulai,
			jadwal.wkt_selesai,
			jadwal_kelas.created_at,
			jadwal_kelas.updated_at
		`).
		Joins("INNER JOIN jadwal ON jadwal_kelas.id_jadwal = jadwal.id").
		Joins("INNER JOIN kelas ON jadwal_kelas.id_kelas = kelas.id").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id")

	if idJadwal != "" {
		query = query.Where("jadwal_kelas.id_jadwal = ?", idJadwal)
	}
	if idKelas != "" {
		query = query.Where("jadwal_kelas.id_kelas = ?", idKelas)
	}

	err := query.Offset(offset).Limit(pageSize).Scan(&results).Error
	return results, total, err
}

func (r *jadwalKelasRepository) CheckDuplicate(idJadwal, idKelas string) (bool, error) {
	var count int64
	err := r.db.Model(&model.JadwalKelas{}).
		Where("id_jadwal = ? AND id_kelas = ?", idJadwal, idKelas).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *jadwalKelasRepository) Update(jadwalKelas *model.JadwalKelas) error {
	return r.db.Save(jadwalKelas).Error
}

func (r *jadwalKelasRepository) Delete(id string) error {
	return r.db.Delete(&model.JadwalKelas{}, "id = ?", id).Error
}

func (r *jadwalKelasRepository) DeleteByJadwalID(idJadwal string) error {
	return r.db.Delete(&model.JadwalKelas{}, "id_jadwal = ?", idJadwal).Error
}
