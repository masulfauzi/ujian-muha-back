package repository

import (
	"time"

	"backend/internal/modules/section/model"

	"gorm.io/gorm"
)

type SectionRepository interface {
	Create(section *model.Section) error
	CreateManyWithTx(tx *gorm.DB, sections []model.Section) error
	GetByID(id string) (*model.Section, error)
	GetByJadwalID(idJadwal string) ([]model.Section, error)
	CountSoalByJadwalID(idJadwal string) (int64, error)
	SoftDeleteByJadwalIDWithTx(tx *gorm.DB, idJadwal string) error
	Delete(id string) error
	Restore(id string) error
}

type sectionRepository struct {
	db *gorm.DB
}

func NewSectionRepository(db *gorm.DB) SectionRepository {
	return &sectionRepository{db: db}
}

func (r *sectionRepository) Create(section *model.Section) error {
	return r.db.Create(section).Error
}

func (r *sectionRepository) CreateManyWithTx(tx *gorm.DB, sections []model.Section) error {
	if len(sections) == 0 {
		return nil
	}
	return tx.Create(&sections).Error
}

func (r *sectionRepository) GetByID(id string) (*model.Section, error) {
	var section model.Section
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&section).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

func (r *sectionRepository) GetByJadwalID(idJadwal string) ([]model.Section, error) {
	var sections []model.Section
	err := r.db.Where("id_jadwal = ? AND deleted_at IS NULL", idJadwal).
		Order("urutan ASC").
		Find(&sections).Error
	return sections, err
}

// CountSoalByJadwalID menghitung jumlah soal dalam bank_soal milik jadwal ini.
func (r *sectionRepository) CountSoalByJadwalID(idJadwal string) (int64, error) {
	var count int64
	err := r.db.Table("soal").
		Joins("INNER JOIN jadwal ON jadwal.id_bank_soal = soal.id_bank_soal").
		Where("jadwal.id = ? AND jadwal.deleted_at IS NULL AND soal.deleted_at IS NULL", idJadwal).
		Count(&count).Error
	return count, err
}

func (r *sectionRepository) SoftDeleteByJadwalIDWithTx(tx *gorm.DB, idJadwal string) error {
	now := time.Now()
	return tx.Model(&model.Section{}).
		Where("id_jadwal = ? AND deleted_at IS NULL", idJadwal).
		Update("deleted_at", now).Error
}

func (r *sectionRepository) Delete(id string) error {
	now := time.Now()
	return r.db.Model(&model.Section{}).Where("id = ?", id).Update("deleted_at", now).Error
}

func (r *sectionRepository) Restore(id string) error {
	return r.db.Model(&model.Section{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NULL")).Error
}
