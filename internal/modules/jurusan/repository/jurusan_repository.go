package repository

import (
	"backend/internal/modules/jurusan/model"

	"gorm.io/gorm"
)

type JurusanRepository interface {
	Create(jurusan *model.Jurusan) error
	GetByID(id string) (*model.Jurusan, error)
	GetAll(page, pageSize int) ([]model.Jurusan, int64, error)
	Update(jurusan *model.Jurusan) error
	Delete(id string) error
	Restore(id string) error
}

type jurusanRepository struct {
	db *gorm.DB
}

func NewJurusanRepository(db *gorm.DB) JurusanRepository {
	return &jurusanRepository{db: db}
}

func (r *jurusanRepository) Create(jurusan *model.Jurusan) error {
	return r.db.Create(jurusan).Error
}

func (r *jurusanRepository) GetByID(id string) (*model.Jurusan, error) {
	var jurusan model.Jurusan
	err := r.db.Where("id = ?", id).First(&jurusan).Error
	if err != nil {
		return nil, err
	}
	return &jurusan, nil
}

func (r *jurusanRepository) GetAll(page, pageSize int) ([]model.Jurusan, int64, error) {
	var jurusans []model.Jurusan
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.Jurusan{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(pageSize).Find(&jurusans).Error
	return jurusans, total, err
}

func (r *jurusanRepository) Update(jurusan *model.Jurusan) error {
	return r.db.Save(jurusan).Error
}

func (r *jurusanRepository) Delete(id string) error {
	return r.db.Delete(&model.Jurusan{}, "id = ?", id).Error
}

func (r *jurusanRepository) Restore(id string) error {
	return r.db.Unscoped().
		Model(&model.Jurusan{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}
