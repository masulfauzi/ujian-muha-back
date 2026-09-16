package repository

import (
	"backend/internal/modules/user_agent/model"

	"gorm.io/gorm"
)

type UserAgentRepository interface {
	Create(userAgent *model.UserAgent) error
	GetByID(id string) (*model.UserAgent, error)
	GetAll(page, pageSize int) ([]model.UserAgent, int64, error)
	Update(userAgent *model.UserAgent) error
	Delete(id string) error
	Restore(id string) error
	// Exists mengecek apakah kombinasi user_agent + x_requested_with PERSIS sama
	// sudah terdaftar (aktif, belum dihapus). Dipakai untuk cegah duplikat baris
	// saat create/update.
	Exists(userAgent, xRequestedWith string) (bool, error)
	// IsAllowed mengecek apakah user_agent ATAU x_requested_with cocok dengan
	// salah satu baris whitelist manapun (tidak harus baris yang sama). Dipakai
	// untuk pengecekan akses saat peserta mulai mengerjakan ujian.
	IsAllowed(userAgent, xRequestedWith string) (bool, error)
}

type userAgentRepository struct {
	db *gorm.DB
}

func NewUserAgentRepository(db *gorm.DB) UserAgentRepository {
	return &userAgentRepository{db: db}
}

func (r *userAgentRepository) Create(userAgent *model.UserAgent) error {
	return r.db.Create(userAgent).Error
}

func (r *userAgentRepository) GetByID(id string) (*model.UserAgent, error) {
	var userAgent model.UserAgent
	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&userAgent).Error
	if err != nil {
		return nil, err
	}
	return &userAgent, nil
}

func (r *userAgentRepository) GetAll(page, pageSize int) ([]model.UserAgent, int64, error) {
	var userAgents []model.UserAgent
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Model(&model.UserAgent{}).
		Where("deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&userAgents).Error

	return userAgents, total, err
}

func (r *userAgentRepository) Update(userAgent *model.UserAgent) error {
	return r.db.Save(userAgent).Error
}

func (r *userAgentRepository) Delete(id string) error {
	return r.db.Delete(&model.UserAgent{}, "id = ?", id).Error
}

func (r *userAgentRepository) Restore(id string) error {
	return r.db.Table("user_agent").Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *userAgentRepository) Exists(userAgent, xRequestedWith string) (bool, error) {
	var count int64
	err := r.db.
		Model(&model.UserAgent{}).
		Where("user_agent = ? AND x_requested_with = ? AND deleted_at IS NULL", userAgent, xRequestedWith).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userAgentRepository) IsAllowed(userAgent, xRequestedWith string) (bool, error) {
	var count int64
	err := r.db.
		Model(&model.UserAgent{}).
		Where("deleted_at IS NULL AND (user_agent = ? OR x_requested_with = ?)", userAgent, xRequestedWith).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
