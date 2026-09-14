package repository

import (
	"backend/internal/modules/kelas/model"

	"gorm.io/gorm"
)

// KelasWithJurusan adalah hasil query kelas dengan JOIN ke jurusan
type KelasWithJurusan struct {
	ID          string `gorm:"column:id"`
	IDJurusan   string `gorm:"column:id_jurusan"`
	NamaKelas   string `gorm:"column:nama_kelas"`
	Tingkat     string `gorm:"column:tingkat"`
	NamaJurusan string `gorm:"column:nama_jurusan"`
	CreatedAt   string `gorm:"column:created_at"`
	UpdatedAt   string `gorm:"column:updated_at"`
	DeletedAt   *string `gorm:"column:deleted_at"`
}

// TableName specifies the table name for this struct
func (KelasWithJurusan) TableName() string {
	return "kelas"
}

type KelasRepository interface {
	Create(kelas *model.Kelas) error
	GetByID(id string) (*KelasWithJurusan, error)
	GetAll(page, pageSize int, idJurusan string, tingkat string) ([]KelasWithJurusan, int64, error)
	Update(kelas *model.Kelas) error
	Delete(id string) error
	Restore(id string) error
}

type kelasRepository struct {
	db *gorm.DB
}

func NewKelasRepository(db *gorm.DB) KelasRepository {
	return &kelasRepository{db: db}
}

func (r *kelasRepository) Create(kelas *model.Kelas) error {
	return r.db.Create(kelas).Error
}

func (r *kelasRepository) GetByID(id string) (*KelasWithJurusan, error) {
	var kelas KelasWithJurusan
	err := r.db.
		Select("kelas.id, kelas.id_jurusan, kelas.nama_kelas, kelas.tingkat, kelas.created_at, kelas.updated_at, kelas.deleted_at, jurusan.nama_jurusan").
		Joins("LEFT JOIN jurusan ON kelas.id_jurusan = jurusan.id").
		Where("kelas.id = ? AND kelas.deleted_at IS NULL", id).
		First(&kelas).Error
	if err != nil {
		return nil, err
	}
	return &kelas, nil
}

func (r *kelasRepository) GetAll(page, pageSize int, idJurusan string, tingkat string) ([]KelasWithJurusan, int64, error) {
	var kelasList []KelasWithJurusan
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	countQuery := r.db.Table("kelas").
		Joins("LEFT JOIN jurusan ON kelas.id_jurusan = jurusan.id").
		Where("kelas.deleted_at IS NULL")

	if idJurusan != "" {
		countQuery = countQuery.Where("kelas.id_jurusan = ?", idJurusan)
	}
	if tingkat != "" {
		countQuery = countQuery.Where("kelas.tingkat = ?", tingkat)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.
		Select("kelas.id, kelas.id_jurusan, kelas.nama_kelas, kelas.tingkat, kelas.created_at, kelas.updated_at, kelas.deleted_at, jurusan.nama_jurusan").
		Joins("LEFT JOIN jurusan ON kelas.id_jurusan = jurusan.id").
		Where("kelas.deleted_at IS NULL")

	if idJurusan != "" {
		query = query.Where("kelas.id_jurusan = ?", idJurusan)
	}
	if tingkat != "" {
		query = query.Where("kelas.tingkat = ?", tingkat)
	}

	err := query.Offset(offset).Limit(pageSize).Find(&kelasList).Error

	return kelasList, total, err
}

func (r *kelasRepository) Update(kelas *model.Kelas) error {
	return r.db.Save(kelas).Error
}

func (r *kelasRepository) Delete(id string) error {
	return r.db.Delete(&model.Kelas{}, "id = ?", id).Error
}

func (r *kelasRepository) Restore(id string) error {
	return r.db.Table("kelas").Where("id = ?", id).Update("deleted_at", nil).Error
}
