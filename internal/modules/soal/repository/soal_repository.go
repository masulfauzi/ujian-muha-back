package repository

import (
	bankSoalModel "backend/internal/modules/bank_soal/model"
	"backend/internal/modules/soal/model"
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SoalRepository interface {
	Create(soal *model.Soal) error
	GetByID(id string) (*model.Soal, error)
	GetAll(page, pageSize int) ([]model.Soal, int64, error)
	GetByBankSoalID(bankSoalID string, page, pageSize int) ([]model.Soal, int64, error)
	CheckDuplicate(bankSoalID string, noSoal int) (bool, error)
	Update(soal *model.Soal) error
	Delete(id string) error
	Restore(id string) error
	HardDelete(id string) error
	BulkCreateSoal(ctx context.Context, soals []model.Soal) error
	GetBankSoalExists(ctx context.Context, bankSoalID string) (bool, error)
}

type soalRepository struct {
	db *gorm.DB
}

func NewSoalRepository(db *gorm.DB) SoalRepository {
	return &soalRepository{db: db}
}

func (r *soalRepository) Create(soal *model.Soal) error {
	return r.db.Create(soal).Error
}

func (r *soalRepository) GetByID(id string) (*model.Soal, error) {
	var soal model.Soal
	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&soal).Error
	if err != nil {
		return nil, err
	}
	return &soal, nil
}

func (r *soalRepository) GetAll(page, pageSize int) ([]model.Soal, int64, error) {
	var soals []model.Soal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Model(&model.Soal{}).
		Where("deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Find(&soals).Error

	return soals, total, err
}

func (r *soalRepository) CheckDuplicate(bankSoalID string, noSoal int) (bool, error) {
	var count int64
	err := r.db.Model(&model.Soal{}).
		Where("id_bank_soal = ? AND no_soal = ? AND deleted_at IS NULL", bankSoalID, noSoal).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *soalRepository) GetByBankSoalID(bankSoalID string, page, pageSize int) ([]model.Soal, int64, error) {
	var soals []model.Soal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Model(&model.Soal{}).
		Where("id_bank_soal = ? AND deleted_at IS NULL", bankSoalID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Where("id_bank_soal = ? AND deleted_at IS NULL", bankSoalID).
		Order("no_soal ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&soals).Error

	return soals, total, err
}

func (r *soalRepository) Update(soal *model.Soal) error {
	return r.db.Save(soal).Error
}

func (r *soalRepository) Delete(id string) error {
	now := time.Now()
	return r.db.Model(&model.Soal{}).Where("id = ?", id).Update("deleted_at", now).Error
}

func (r *soalRepository) Restore(id string) error {
	return r.db.Model(&model.Soal{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NULL")).Error
}

func (r *soalRepository) HardDelete(id string) error {
	return r.db.Unscoped().Delete(&model.Soal{}, "id = ?", id).Error
}

func (r *soalRepository) BulkCreateSoal(ctx context.Context, soals []model.Soal) error {
	if len(soals) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id_bank_soal"},
				{Name: "no_soal"},
			},
			TargetWhere: clause.Where{
				Exprs: []clause.Expression{clause.Expr{SQL: "deleted_at IS NULL"}},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"soal",
				"opsi_a", "opsi_b", "opsi_c", "opsi_d", "opsi_e",
				"gambar_soal", "gambar_a", "gambar_b", "gambar_c", "gambar_d", "gambar_e",
				"kunci",
				"updated_at",
			}),
		}).
		CreateInBatches(soals, 100).Error
}

func (r *soalRepository) GetBankSoalExists(ctx context.Context, bankSoalID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&bankSoalModel.BankSoal{}).
		Where("id = ? AND deleted_at IS NULL", bankSoalID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
