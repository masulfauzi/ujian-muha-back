package repository

import (
	"backend/internal/modules/mapel/model"

	"gorm.io/gorm"
)

type MapelRepository interface {
	Create(mapel *model.Mapel) error
	GetByID(id string) (*model.Mapel, error)
	GetAll(page, pageSize int) ([]model.Mapel, int64, error)
	Update(mapel *model.Mapel) error
	Delete(id string) error
	Restore(id string) error
	HardDelete(id string) error
}

type mapelRepository struct {
	db *gorm.DB
}

func NewMapelRepository(db *gorm.DB) MapelRepository {
	return &mapelRepository{db: db}
}

func (r *mapelRepository) Create(mapel *model.Mapel) error {
	return r.db.Create(mapel).Error
}

func (r *mapelRepository) GetByID(id string) (*model.Mapel, error) {
	var mapel model.Mapel
	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&mapel).Error
	if err != nil {
		return nil, err
	}
	return &mapel, nil
}

func (r *mapelRepository) GetAll(page, pageSize int) ([]model.Mapel, int64, error) {
	var mapels []model.Mapel
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Count total records (excluding soft deleted)
	err := r.db.
		Model(&model.Mapel{}).
		Where("deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get records with pagination (excluding soft deleted)
	err = r.db.
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Find(&mapels).Error

	return mapels, total, err
}

func (r *mapelRepository) Update(mapel *model.Mapel) error {
	return r.db.Save(mapel).Error
}

func (r *mapelRepository) Delete(id string) error {
	// Soft delete - GORM will automatically set deleted_at
	return r.db.Delete(&model.Mapel{}, "id = ?", id).Error
}

func (r *mapelRepository) Restore(id string) error {
	// Restore - clear deleted_at using direct update
	return r.db.Table("mapel").Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *mapelRepository) HardDelete(id string) error {
	return r.db.Unscoped().Delete(&model.Mapel{}, "id = ?", id).Error
}
