package repository

import (
	"context"

	"backend/internal/modules/peserta/model"

	"gorm.io/gorm"
)

// KelasLookup adalah hasil pencarian kelas berdasarkan nama, dipakai saat
// resolve kolom "kelas" di file excel import peserta ke id_kelas.
type KelasLookup struct {
	ID        string `gorm:"column:id"`
	NamaKelas string `gorm:"column:nama_kelas"`
}

type PesertaWithKelas struct {
	ID        string  `gorm:"column:id"`
	Nama      string  `gorm:"column:nama"`
	IDKelas   string  `gorm:"column:id_kelas"`
	NamaKelas string  `gorm:"column:nama_kelas"`
	Username  string  `gorm:"column:username"`
	Password  string  `gorm:"column:password"`
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
	GetAll(page, pageSize int, idKelas, search string) ([]PesertaWithKelas, int64, error)
	GetRawByID(id string) (*model.Peserta, error)
	GetByUsername(username string) (*model.Peserta, error)
	Update(peserta *model.Peserta) error
	Delete(id string) error
	Restore(id string) error
	DeleteAll(ctx context.Context) (int64, error)
	BulkCreatePeserta(ctx context.Context, pesertaList []model.Peserta) error
	GetKelasByNamaList(ctx context.Context, lowerNamaList []string) ([]KelasLookup, error)
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
		Select("peserta.id, peserta.nama, peserta.id_kelas, peserta.username, peserta.password, peserta.created_at, peserta.updated_at, peserta.deleted_at, kelas.nama_kelas").
		Joins("LEFT JOIN kelas ON peserta.id_kelas = kelas.id").
		Where("peserta.id = ? AND peserta.deleted_at IS NULL", id).
		First(&peserta).Error
	if err != nil {
		return nil, err
	}
	return &peserta, nil
}

func (r *pesertaRepository) GetAll(page, pageSize int, idKelas, search string) ([]PesertaWithKelas, int64, error) {
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
	if search != "" {
		countQuery = countQuery.Where("peserta.nama ILIKE ? OR peserta.username ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.
		Select("peserta.id, peserta.nama, peserta.id_kelas, peserta.username, peserta.password, peserta.created_at, peserta.updated_at, peserta.deleted_at, kelas.nama_kelas").
		Joins("LEFT JOIN kelas ON peserta.id_kelas = kelas.id").
		Where("peserta.deleted_at IS NULL")

	if idKelas != "" {
		query = query.Where("peserta.id_kelas = ?", idKelas)
	}
	if search != "" {
		query = query.Where("peserta.nama ILIKE ? OR peserta.username ILIKE ?", "%"+search+"%", "%"+search+"%")
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

// DeleteAll menghapus PERMANEN seluruh baris di tabel peserta (bukan soft-delete).
// Tidak menyentuh tabel nilai/jawaban yang mungkin masih mereferensikan id peserta
// yang dihapus (tidak ada foreign key constraint di skema saat ini, jadi baris itu
// tidak akan gagal dihapus, tapi record nilai/jawaban terkait jadi yatim/orphan).
func (r *pesertaRepository) DeleteAll(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Unscoped().
		Delete(&model.Peserta{})
	return result.RowsAffected, result.Error
}

func (r *pesertaRepository) BulkCreatePeserta(ctx context.Context, pesertaList []model.Peserta) error {
	if len(pesertaList) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(pesertaList, 100).Error
}

// GetKelasByNamaList mencari kelas berdasarkan nama (case-insensitive, sudah di-lowercase
// oleh pemanggil). Bisa mengembalikan lebih dari satu baris untuk nama yang sama jika ada
// kelas dengan nama kembar di kelas/jurusan berbeda — pemanggil yang menentukan ambigu atau tidak.
func (r *pesertaRepository) GetKelasByNamaList(ctx context.Context, lowerNamaList []string) ([]KelasLookup, error) {
	if len(lowerNamaList) == 0 {
		return nil, nil
	}
	var results []KelasLookup
	err := r.db.WithContext(ctx).
		Table("kelas").
		Select("id, nama_kelas").
		Where("deleted_at IS NULL AND LOWER(nama_kelas) IN ?", lowerNamaList).
		Find(&results).Error
	return results, err
}
