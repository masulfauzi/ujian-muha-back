package repository

import (
	"backend/internal/modules/bank_soal/model"
	"time"

	"gorm.io/gorm"
)

type BankSoalWithMapel struct {
	ID           string `gorm:"column:id"`
	NamaBankSoal string `gorm:"column:nama_bank_soal"`
	IdMapel      string `gorm:"column:id_mapel"`
	NamaMapel    string `gorm:"column:nama_mapel"`
	JmlSoal      int    `gorm:"column:jml_soal"`
	Deskripsi    string `gorm:"column:deskripsi"`
	CreatedAt    string `gorm:"column:created_at"`
	UpdatedAt    string `gorm:"column:updated_at"`
}

type BankSoalRepository interface {
	Create(bankSoal *model.BankSoal) error
	GetByID(id string) (*model.BankSoal, error)
	GetByIDWithMapel(id string) (*BankSoalWithMapel, error)
	GetAll(page, pageSize int) ([]model.BankSoal, int64, error)
	GetByMapelID(mapelID string, page, pageSize int) ([]model.BankSoal, int64, error)
	GetAllWithMapel(page, pageSize int) ([]BankSoalWithMapel, int64, error)
	GetByMapelIDWithMapel(mapelID string, page, pageSize int) ([]BankSoalWithMapel, int64, error)
	Update(bankSoal *model.BankSoal) error
	Delete(id string) error
	Restore(id string) error
	HardDelete(id string) error
}

type bankSoalRepository struct {
	db *gorm.DB
}

func NewBankSoalRepository(db *gorm.DB) BankSoalRepository {
	return &bankSoalRepository{db: db}
}

func (r *bankSoalRepository) Create(bankSoal *model.BankSoal) error {
	return r.db.Create(bankSoal).Error
}

func (r *bankSoalRepository) GetByID(id string) (*model.BankSoal, error) {
	var bankSoal model.BankSoal
	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&bankSoal).Error
	if err != nil {
		return nil, err
	}
	return &bankSoal, nil
}

func (r *bankSoalRepository) GetByIDWithMapel(id string) (*BankSoalWithMapel, error) {
	var bankSoal BankSoalWithMapel
	err := r.db.
		Table("bank_soal").
		Select("bank_soal.id, bank_soal.nama_bank_soal, bank_soal.id_mapel, mapel.nama_mapel, bank_soal.jml_soal, bank_soal.deskripsi, bank_soal.created_at, bank_soal.updated_at").
		Joins("INNER JOIN mapel ON bank_soal.id_mapel = mapel.id").
		Where("bank_soal.id = ? AND bank_soal.deleted_at IS NULL", id).
		First(&bankSoal).Error
	if err != nil {
		return nil, err
	}
	return &bankSoal, nil
}

func (r *bankSoalRepository) GetAll(page, pageSize int) ([]model.BankSoal, int64, error) {
	var bankSoals []model.BankSoal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Model(&model.BankSoal{}).
		Where("deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Find(&bankSoals).Error

	return bankSoals, total, err
}

func (r *bankSoalRepository) GetByMapelID(mapelID string, page, pageSize int) ([]model.BankSoal, int64, error) {
	var bankSoals []model.BankSoal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Model(&model.BankSoal{}).
		Where("id_mapel = ? AND deleted_at IS NULL", mapelID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Where("id_mapel = ? AND deleted_at IS NULL", mapelID).
		Offset(offset).
		Limit(pageSize).
		Find(&bankSoals).Error

	return bankSoals, total, err
}

func (r *bankSoalRepository) Update(bankSoal *model.BankSoal) error {
	return r.db.Save(bankSoal).Error
}

func (r *bankSoalRepository) Delete(id string) error {
	now := time.Now()
	return r.db.Model(&model.BankSoal{}).Where("id = ?", id).Update("deleted_at", now).Error
}

func (r *bankSoalRepository) Restore(id string) error {
	return r.db.Model(&model.BankSoal{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NULL")).Error
}

func (r *bankSoalRepository) HardDelete(id string) error {
	return r.db.Unscoped().Delete(&model.BankSoal{}, "id = ?", id).Error
}

func (r *bankSoalRepository) GetAllWithMapel(page, pageSize int) ([]BankSoalWithMapel, int64, error) {
	var bankSoals []BankSoalWithMapel
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Table("bank_soal").
		Select("bank_soal.id, bank_soal.nama_bank_soal, bank_soal.id_mapel, mapel.nama_mapel, bank_soal.jml_soal, bank_soal.deskripsi, bank_soal.created_at, bank_soal.updated_at").
		Joins("INNER JOIN mapel ON bank_soal.id_mapel = mapel.id").
		Where("bank_soal.deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Table("bank_soal").
		Select("bank_soal.id, bank_soal.nama_bank_soal, bank_soal.id_mapel, mapel.nama_mapel, bank_soal.jml_soal, bank_soal.deskripsi, bank_soal.created_at, bank_soal.updated_at").
		Joins("INNER JOIN mapel ON bank_soal.id_mapel = mapel.id").
		Where("bank_soal.deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Scan(&bankSoals).Error

	return bankSoals, total, err
}

func (r *bankSoalRepository) GetByMapelIDWithMapel(mapelID string, page, pageSize int) ([]BankSoalWithMapel, int64, error) {
	var bankSoals []BankSoalWithMapel
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Table("bank_soal").
		Select("bank_soal.id, bank_soal.nama_bank_soal, bank_soal.id_mapel, mapel.nama_mapel, bank_soal.jml_soal, bank_soal.deskripsi, bank_soal.created_at, bank_soal.updated_at").
		Joins("INNER JOIN mapel ON bank_soal.id_mapel = mapel.id").
		Where("bank_soal.id_mapel = ? AND bank_soal.deleted_at IS NULL", mapelID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Table("bank_soal").
		Select("bank_soal.id, bank_soal.nama_bank_soal, bank_soal.id_mapel, mapel.nama_mapel, bank_soal.jml_soal, bank_soal.deskripsi, bank_soal.created_at, bank_soal.updated_at").
		Joins("INNER JOIN mapel ON bank_soal.id_mapel = mapel.id").
		Where("bank_soal.id_mapel = ? AND bank_soal.deleted_at IS NULL", mapelID).
		Offset(offset).
		Limit(pageSize).
		Scan(&bankSoals).Error

	return bankSoals, total, err
}
