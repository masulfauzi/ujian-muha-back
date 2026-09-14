package repository

import (
	"backend/internal/modules/peserta/model"

	"gorm.io/gorm"
)

type PesertaWithKelas struct {
	ID        string  `gorm:"column:id"`
	Nama      string  `gorm:"column:nama"`
	IDKelas   string  `gorm:"column:id_kelas"`
	NamaKelas string  `gorm:"column:nama_kelas"`
	Username  string  `gorm:"column:username"`
	CreatedAt string  `gorm:"column:created_at"`
	UpdatedAt string  `gorm:"column:updated_at"`
	DeletedAt *string `gorm:"column:deleted_at"`
}

func (PesertaWithKelas) TableName() string {
	return "peserta"
}

type PesertaRepository interface {
	Create(peserta *model.Peserta) error
	GetByID(id string) (*PesertaWithKelas, error)
	GetAll(page, pageSize int, idKelas string) ([]PesertaWithKelas, int64, error)
	GetRawByID(id string) (*model.Peserta, error)
	GetByUsername(username string) (*model.Peserta, error)
	Update(peserta *model.Peserta) error
	Delete(id string) error
	Restore(id string) error
}

type pesertaRepository struct {
	db *gorm.DB
}

func NewPesertaRepository(db *gorm.DB) PesertaRepository {
	return &pesertaRepository{db: db}
}

func (r *pesertaRepository) Create(peserta *model.Peserta) error {
	return r.db.Create(peserta).Error
}

func (r *pesertaRepository) GetByID(id string) (*PesertaWithKelas, error) {
	var peserta PesertaWithKelas
	err := r.db.
		Select("peserta.id, peserta.nama, peserta.id_kelas, peserta.username, peserta.created_at, peserta.updated_at, peserta.deleted_at, kelas.nama_kelas").
		Joins("LEFT JOIN kelas ON peserta.id_kelas = kelas.id").
		Where("peserta.id = ? AND peserta.deleted_at IS NULL", id).
		First(&peserta).Error
	if err != nil {
		return nil, err
	}
	return &peserta, nil
}

func (r *pesertaRepository) GetAll(page, pageSize int, idKelas string) ([]PesertaWithKelas, int64, error) {
	var pesertaList []PesertaWithKelas
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	countQuery := r.db.Table("peserta").
		Joins("LEFT JOIN kelas ON peserta.id_kelas = kelas.id").
		Where("peserta.deleted_at IS NULL")

	if idKelas != "" {
		countQuery = countQuery.Where("peserta.id_kelas = ?", idKelas)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.
		Select("peserta.id, peserta.nama, peserta.id_kelas, peserta.username, peserta.created_at, peserta.updated_at, peserta.deleted_at, kelas.nama_kelas").
		Joins("LEFT JOIN kelas ON peserta.id_kelas = kelas.id").
		Where("peserta.deleted_at IS NULL")

	if idKelas != "" {
		query = query.Where("peserta.id_kelas = ?", idKelas)
	}

	err := query.Offset(offset).Limit(pageSize).Find(&pesertaList).Error
	return pesertaList, total, err
}

func (r *pesertaRepository) GetRawByID(id string) (*model.Peserta, error) {
	var peserta model.Peserta
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&peserta).Error
	if err != nil {
		return nil, err
	}
	return &peserta, nil
}

func (r *pesertaRepository) GetByUsername(username string) (*model.Peserta, error) {
	var peserta model.Peserta
	err := r.db.Where("username = ? AND deleted_at IS NULL", username).First(&peserta).Error
	if err != nil {
		return nil, err
	}
	return &peserta, nil
}

func (r *pesertaRepository) Update(peserta *model.Peserta) error {
	return r.db.Save(peserta).Error
}

func (r *pesertaRepository) Delete(id string) error {
	return r.db.Delete(&model.Peserta{}, "id = ?", id).Error
}

func (r *pesertaRepository) Restore(id string) error {
	return r.db.Table("peserta").Where("id = ?", id).Update("deleted_at", nil).Error
}
