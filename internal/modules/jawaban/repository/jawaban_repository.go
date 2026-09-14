package repository

import (
	"backend/internal/modules/jawaban/model"
	"time"

	"gorm.io/gorm"
)

type JawabanWithDetail struct {
	ID          string  `gorm:"column:id"`
	IDNilai     string  `gorm:"column:id_nilai"`
	IDSoal      string  `gorm:"column:id_soal"`
	NoUrut      int     `gorm:"column:no_urut"`
	NoSoal      int     `gorm:"column:no_soal"`
	SoalText    string  `gorm:"column:soal_text"`
	Kunci       string  `gorm:"column:kunci"`
	OpsiA       string  `gorm:"column:opsi_a"`
	OpsiB       string  `gorm:"column:opsi_b"`
	OpsiC       string  `gorm:"column:opsi_c"`
	OpsiD       string  `gorm:"column:opsi_d"`
	OpsiE       string  `gorm:"column:opsi_e"`
	GambarA     string  `gorm:"column:gambar_a"`
	GambarB     string  `gorm:"column:gambar_b"`
	GambarC     string  `gorm:"column:gambar_c"`
	GambarD     string  `gorm:"column:gambar_d"`
	GambarE     string  `gorm:"column:gambar_e"`
	IDPeserta   string  `gorm:"column:id_peserta"`
	NamaPeserta string  `gorm:"column:nama_peserta"`
	Jawaban     *string `gorm:"column:jawaban"`
	IsBenar     *int    `gorm:"column:is_benar"`
	AcakOpsi    int     `gorm:"column:acak_opsi"`
	CreatedAt   string  `gorm:"column:created_at"`
	UpdatedAt   string  `gorm:"column:updated_at"`
}

type JawabanRepository interface {
	Create(jawaban *model.Jawaban) error
	GetByID(id string) (*model.Jawaban, error)
	GetByIDWithDetail(id string) (*JawabanWithDetail, error)
	GetAllWithDetail(page, pageSize int, idNilai, idPeserta, idSoal string) ([]JawabanWithDetail, int64, error)
	GetByNilaiID(idNilai string, page, pageSize int) ([]JawabanWithDetail, int64, error)
	GetByNilaiIDInRange(idNilai string, noUrutAwal, noUrutAkhir int) ([]JawabanWithDetail, error)
	GetByPesertaID(idPeserta string, page, pageSize int) ([]JawabanWithDetail, int64, error)
	GetSoalKunci(idSoal string) (string, error)
	CheckDuplicate(idNilai, idSoal string) (bool, error)
	Update(jawaban *model.Jawaban) error
	Delete(id string) error
	Restore(id string) error
	BulkCreateWithTx(tx *gorm.DB, jawabans []model.Jawaban) error
	SoftDeleteByNilaiID(tx *gorm.DB, nilaiID string) error
	RestoreByNilaiID(tx *gorm.DB, nilaiID string) error
	GetAcakOpsiByNilaiID(nilaiID string) (int, error)
}

type jawabanRepository struct {
	db *gorm.DB
}

func NewJawabanRepository(db *gorm.DB) JawabanRepository {
	return &jawabanRepository{db: db}
}

func (r *jawabanRepository) Create(jawaban *model.Jawaban) error {
	return r.db.Create(jawaban).Error
}

func (r *jawabanRepository) GetByID(id string) (*model.Jawaban, error) {
	var jawaban model.Jawaban
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&jawaban).Error
	if err != nil {
		return nil, err
	}
	return &jawaban, nil
}

func (r *jawabanRepository) GetAcakOpsiByNilaiID(nilaiID string) (int, error) {
	var acakOpsi int
	err := r.db.Table("nilai").
		Select("jadwal.acak_opsi::int AS acak_opsi").
		Joins("INNER JOIN jadwal ON nilai.id_jadwal = jadwal.id").
		Where("nilai.id = ?", nilaiID).
		Scan(&acakOpsi).Error
	return acakOpsi, err
}

func (r *jawabanRepository) GetByIDWithDetail(id string) (*JawabanWithDetail, error) {
	var result JawabanWithDetail
	err := r.db.
		Table("jawaban").
		Select(`
			jawaban.id,
			jawaban.id_nilai,
			jawaban.id_soal,
			jawaban.no_urut,
			soal.no_soal,
			soal.soal AS soal_text,
			soal.kunci,
			soal.opsi_a,
			soal.opsi_b,
			soal.opsi_c,
			soal.opsi_d,
			soal.opsi_e,
			soal.gambar_a,
			soal.gambar_b,
			soal.gambar_c,
			soal.gambar_d,
			soal.gambar_e,
			jawaban.id_peserta,
			peserta.nama AS nama_peserta,
			jawaban.jawaban,
			jawaban.is_benar,
			jadwal.acak_opsi::int AS acak_opsi,
			jawaban.created_at,
			jawaban.updated_at
		`).
		Joins("INNER JOIN soal ON jawaban.id_soal = soal.id").
		Joins("INNER JOIN peserta ON jawaban.id_peserta = peserta.id").
		Joins("INNER JOIN nilai ON jawaban.id_nilai = nilai.id").
		Joins("INNER JOIN jadwal ON nilai.id_jadwal = jadwal.id").
		Where("jawaban.id = ? AND jawaban.deleted_at IS NULL", id).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *jawabanRepository) GetAllWithDetail(page, pageSize int, idNilai, idPeserta, idSoal string) ([]JawabanWithDetail, int64, error) {
	var results []JawabanWithDetail
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	buildBase := func() *gorm.DB {
		q := r.db.Table("jawaban").
			Joins("INNER JOIN soal ON jawaban.id_soal = soal.id").
			Joins("INNER JOIN peserta ON jawaban.id_peserta = peserta.id").
			Joins("INNER JOIN nilai ON jawaban.id_nilai = nilai.id").
			Joins("INNER JOIN jadwal ON nilai.id_jadwal = jadwal.id").
			Where("jawaban.deleted_at IS NULL")
		if idNilai != "" {
			q = q.Where("jawaban.id_nilai = ?", idNilai)
		}
		if idPeserta != "" {
			q = q.Where("jawaban.id_peserta = ?", idPeserta)
		}
		if idSoal != "" {
			q = q.Where("jawaban.id_soal = ?", idSoal)
		}
		return q
	}

	if err := buildBase().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := buildBase().
		Select(`
			jawaban.id,
			jawaban.id_nilai,
			jawaban.id_soal,
			jawaban.no_urut,
			soal.no_soal,
			soal.soal AS soal_text,
			soal.kunci,
			soal.opsi_a,
			soal.opsi_b,
			soal.opsi_c,
			soal.opsi_d,
			soal.opsi_e,
			soal.gambar_a,
			soal.gambar_b,
			soal.gambar_c,
			soal.gambar_d,
			soal.gambar_e,
			jawaban.id_peserta,
			peserta.nama AS nama_peserta,
			jawaban.jawaban,
			jawaban.is_benar,
			jadwal.acak_opsi::int AS acak_opsi,
			jawaban.created_at,
			jawaban.updated_at
		`).
		Order("jawaban.no_urut ASC").
		Offset(offset).Limit(pageSize).
		Scan(&results).Error

	return results, total, err
}

func (r *jawabanRepository) GetByNilaiID(idNilai string, page, pageSize int) ([]JawabanWithDetail, int64, error) {
	return r.GetAllWithDetail(page, pageSize, idNilai, "", "")
}

// GetByNilaiIDInRange mengambil jawaban satu sesi yang no_urut-nya berada dalam range section tertentu.
func (r *jawabanRepository) GetByNilaiIDInRange(idNilai string, noUrutAwal, noUrutAkhir int) ([]JawabanWithDetail, error) {
	var results []JawabanWithDetail
	err := r.db.Table("jawaban").
		Select(`
			jawaban.id,
			jawaban.id_nilai,
			jawaban.id_soal,
			jawaban.no_urut,
			soal.no_soal,
			soal.soal AS soal_text,
			soal.kunci,
			soal.opsi_a,
			soal.opsi_b,
			soal.opsi_c,
			soal.opsi_d,
			soal.opsi_e,
			soal.gambar_a,
			soal.gambar_b,
			soal.gambar_c,
			soal.gambar_d,
			soal.gambar_e,
			jawaban.id_peserta,
			peserta.nama AS nama_peserta,
			jawaban.jawaban,
			jawaban.is_benar,
			jadwal.acak_opsi::int AS acak_opsi,
			jawaban.created_at,
			jawaban.updated_at
		`).
		Joins("INNER JOIN soal ON jawaban.id_soal = soal.id").
		Joins("INNER JOIN peserta ON jawaban.id_peserta = peserta.id").
		Joins("INNER JOIN nilai ON jawaban.id_nilai = nilai.id").
		Joins("INNER JOIN jadwal ON nilai.id_jadwal = jadwal.id").
		Where("jawaban.id_nilai = ? AND jawaban.deleted_at IS NULL AND jawaban.no_urut BETWEEN ? AND ?", idNilai, noUrutAwal, noUrutAkhir).
		Order("jawaban.no_urut ASC").
		Scan(&results).Error
	return results, err
}

func (r *jawabanRepository) GetByPesertaID(idPeserta string, page, pageSize int) ([]JawabanWithDetail, int64, error) {
	return r.GetAllWithDetail(page, pageSize, "", idPeserta, "")
}

func (r *jawabanRepository) GetSoalKunci(idSoal string) (string, error) {
	var kunci string
	err := r.db.Table("soal").
		Select("kunci").
		Where("id = ? AND deleted_at IS NULL", idSoal).
		Scan(&kunci).Error
	return kunci, err
}

func (r *jawabanRepository) CheckDuplicate(idNilai, idSoal string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Jawaban{}).
		Where("id_nilai = ? AND id_soal = ? AND deleted_at IS NULL", idNilai, idSoal).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *jawabanRepository) Update(jawaban *model.Jawaban) error {
	return r.db.Save(jawaban).Error
}

func (r *jawabanRepository) Delete(id string) error {
	now := time.Now()
	return r.db.Model(&model.Jawaban{}).Where("id = ?", id).Update("deleted_at", now).Error
}

func (r *jawabanRepository) Restore(id string) error {
	return r.db.Model(&model.Jawaban{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NULL")).Error
}

func (r *jawabanRepository) BulkCreateWithTx(tx *gorm.DB, jawabans []model.Jawaban) error {
	if len(jawabans) == 0 {
		return nil
	}
	return tx.Create(&jawabans).Error
}

func (r *jawabanRepository) SoftDeleteByNilaiID(tx *gorm.DB, nilaiID string) error {
	now := time.Now()
	return tx.Model(&model.Jawaban{}).
		Where("id_nilai = ? AND deleted_at IS NULL", nilaiID).
		Update("deleted_at", now).Error
}

func (r *jawabanRepository) RestoreByNilaiID(tx *gorm.DB, nilaiID string) error {
	return tx.Model(&model.Jawaban{}).
		Where("id_nilai = ? AND deleted_at IS NOT NULL", nilaiID).
		Update("deleted_at", gorm.Expr("NULL")).Error
}
