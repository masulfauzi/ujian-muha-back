package repository

import (
	"backend/internal/modules/login_log/model"

	"gorm.io/gorm"
)

type LoginLogRepository interface {
	Create(log *model.LoginLog) error
	GetAll(page, pageSize int, username string) ([]model.LoginLog, int64, error)
}

type loginLogRepository struct {
	db *gorm.DB
}

func NewLoginLogRepository(db *gorm.DB) LoginLogRepository {
	return &loginLogRepository{db: db}
}

func (r *loginLogRepository) Create(log *model.LoginLog) error {
	return r.db.Create(log).Error
}

func (r *loginLogRepository) GetAll(page, pageSize int, username string) ([]model.LoginLog, int64, error) {
	var logs []model.LoginLog
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	countQuery := r.db.Model(&model.LoginLog{})
	if username != "" {
		countQuery = countQuery.Where("username ILIKE ?", "%"+username+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Model(&model.LoginLog{})
	if username != "" {
		query = query.Where("username ILIKE ?", "%"+username+"%")
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}
