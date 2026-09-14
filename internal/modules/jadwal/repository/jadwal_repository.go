package repository

import (
	"backend/internal/modules/jadwal/model"
	"time"

	"gorm.io/gorm"
)

type JadwalWithBankSoal struct {
	ID           string  `gorm:"column:id"`
	IDBankSoal   string  `gorm:"column:id_bank_soal"`
	NamaBankSoal string  `gorm:"column:nama_bank_soal"`
	NamaUjian    string  `gorm:"column:nama_ujian"`
	Tingkat      string  `gorm:"column:tingkat"`
	WktMulai     string  `gorm:"column:wkt_mulai"`
	WktSelesai   string  `gorm:"column:wkt_selesai"`
	Durasi       int     `gorm:"column:durasi"`
	AcakSoal     int     `gorm:"column:acak_soal"`
	AcakOpsi     int     `gorm:"column:acak_opsi"`
	WajibToken   int     `gorm:"column:wajib_token"`
	CreatedAt    string  `gorm:"column:created_at"`
	UpdatedAt    string  `gorm:"column:updated_at"`
}

type KelasDetail struct {
	ID        string `gorm:"column:id"`
	IDKelas   string `gorm:"column:id_kelas"`
	NamaKelas string `gorm:"column:nama_kelas"`
	IDJurusan string `gorm:"column:id_jurusan"`
}

type JurusanDetail struct {
	ID          string `gorm:"column:id"`
	IDJurusan   string `gorm:"column:id_jurusan"`
	NamaJurusan string `gorm:"column:nama_jurusan"`
}

type JadwalWithKelas struct {
	ID           string          `json:"id"`
	IDBankSoal   string          `json:"id_bank_soal"`
	NamaBankSoal string          `json:"nama_bank_soal"`
	NamaUjian    string          `json:"nama_ujian"`
	Tingkat      string          `json:"tingkat"`
	WktMulai     string          `json:"wkt_mulai"`
	WktSelesai   string          `json:"wkt_selesai"`
	Durasi       int             `json:"durasi"`
	AcakSoal     int             `json:"acak_soal"`
	AcakOpsi     int             `json:"acak_opsi"`
	WajibToken   int             `json:"wajib_token"`
	IDKelas      []KelasDetail   `json:"id_kelas"`
	IDJurusan    []JurusanDetail `json:"id_jurusan"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

type JadwalAktifWithStatus struct {
	ID              string  `gorm:"column:id"`
	IDBankSoal      string  `gorm:"column:id_bank_soal"`
	NamaBankSoal    string  `gorm:"column:nama_bank_soal"`
	NamaUjian       string  `gorm:"column:nama_ujian"`
	Tingkat         string  `gorm:"column:tingkat"`
	WktMulai        string  `gorm:"column:wkt_mulai"`
	WktSelesai      string  `gorm:"column:wkt_selesai"`
	Durasi          int     `gorm:"column:durasi"`
	AcakSoal        int     `gorm:"column:acak_soal"`
	AcakOpsi        int     `gorm:"column:acak_opsi"`
	WajibToken      int     `gorm:"column:wajib_token"`
	IDNilai         *string `gorm:"column:id_nilai"`
	NilaiWktSelesai *string `gorm:"column:nilai_wkt_selesai"`
}

type JadwalRepository interface {
	Create(jadwal *model.Jadwal) error
	GetByID(id string) (*model.Jadwal, error)
	GetByIDWithBankSoal(id string) (*JadwalWithBankSoal, error)
	GetByIDWithKelas(id string) (*JadwalWithKelas, error)
	GetAllWithBankSoal(page, pageSize int) ([]JadwalWithBankSoal, int64, error)
	GetByBankSoalID(bankSoalID string, page, pageSize int) ([]JadwalWithBankSoal, int64, error)
	GetAktifHariIniByKelas(idKelas, idPeserta string) ([]JadwalAktifWithStatus, error)
	GetAcakOpsiForPesertaSoal(pesertaID, soalID string) (int, error)
	Update(jadwal *model.Jadwal) error
	Delete(id string) error
	Restore(id string) error
}

type jadwalRepository struct {
	db *gorm.DB
}

func NewJadwalRepository(db *gorm.DB) JadwalRepository {
	return &jadwalRepository{db: db}
}

func (r *jadwalRepository) Create(jadwal *model.Jadwal) error {
	return r.db.Create(jadwal).Error
}

func (r *jadwalRepository) GetByID(id string) (*model.Jadwal, error) {
	var jadwal model.Jadwal
	err := r.db.
		Table("jadwal").
		Select("id, id_bank_soal, nama_ujian, tingkat, wkt_mulai, wkt_selesai, durasi, acak_soal::int AS acak_soal, acak_opsi::int AS acak_opsi, wajib_token::int AS wajib_token, created_at, updated_at, deleted_at").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&jadwal).Error
	if err != nil {
		return nil, err
	}
	return &jadwal, nil
}

func (r *jadwalRepository) GetByIDWithBankSoal(id string) (*JadwalWithBankSoal, error) {
	var jadwal JadwalWithBankSoal
	err := r.db.
		Table("jadwal").
		Select("jadwal.id, jadwal.id_bank_soal, bank_soal.nama_bank_soal, jadwal.nama_ujian, jadwal.tingkat, TO_CHAR(jadwal.wkt_mulai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_mulai, TO_CHAR(jadwal.wkt_selesai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_selesai, jadwal.durasi, jadwal.acak_soal::int AS acak_soal, jadwal.acak_opsi::int AS acak_opsi, jadwal.wajib_token::int AS wajib_token, TO_CHAR(jadwal.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at, TO_CHAR(jadwal.updated_at, 'YYYY-MM-DD HH24:MI:SS') AS updated_at").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.id = ? AND jadwal.deleted_at IS NULL", id).
		First(&jadwal).Error
	if err != nil {
		return nil, err
	}
	return &jadwal, nil
}

func (r *jadwalRepository) GetByIDWithKelas(id string) (*JadwalWithKelas, error) {
	var jadwal JadwalWithBankSoal
	err := r.db.
		Table("jadwal").
		Select("jadwal.id, jadwal.id_bank_soal, bank_soal.nama_bank_soal, jadwal.nama_ujian, jadwal.tingkat, TO_CHAR(jadwal.wkt_mulai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_mulai, TO_CHAR(jadwal.wkt_selesai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_selesai, jadwal.durasi, jadwal.acak_soal::int AS acak_soal, jadwal.acak_opsi::int AS acak_opsi, jadwal.wajib_token::int AS wajib_token, TO_CHAR(jadwal.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at, TO_CHAR(jadwal.updated_at, 'YYYY-MM-DD HH24:MI:SS') AS updated_at").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.id = ? AND jadwal.deleted_at IS NULL", id).
		First(&jadwal).Error
	if err != nil {
		return nil, err
	}

	var kelasList []KelasDetail
	err = r.db.
		Table("jadwal_kelas").
		Select("jadwal_kelas.id, jadwal_kelas.id_kelas, kelas.nama_kelas, kelas.id_jurusan").
		Joins("INNER JOIN kelas ON jadwal_kelas.id_kelas = kelas.id").
		Where("jadwal_kelas.id_jadwal = ?", id).
		Scan(&kelasList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if kelasList == nil {
		kelasList = []KelasDetail{}
	}

	var jurusanList []JurusanDetail
	err = r.db.
		Table("jurusan").
		Select("DISTINCT jurusan.id, kelas.id_jurusan, jurusan.nama_jurusan").
		Joins("INNER JOIN kelas ON jurusan.id = kelas.id_jurusan").
		Joins("INNER JOIN jadwal_kelas ON kelas.id = jadwal_kelas.id_kelas").
		Where("jadwal_kelas.id_jadwal = ? AND jurusan.deleted_at IS NULL", id).
		Scan(&jurusanList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if jurusanList == nil {
		jurusanList = []JurusanDetail{}
	}

	return &JadwalWithKelas{
		ID:           jadwal.ID,
		IDBankSoal:   jadwal.IDBankSoal,
		NamaBankSoal: jadwal.NamaBankSoal,
		NamaUjian:    jadwal.NamaUjian,
		Tingkat:      jadwal.Tingkat,
		WktMulai:     jadwal.WktMulai,
		WktSelesai:   jadwal.WktSelesai,
		Durasi:       jadwal.Durasi,
		AcakSoal:     jadwal.AcakSoal,
		AcakOpsi:     jadwal.AcakOpsi,
		WajibToken:   jadwal.WajibToken,
		IDKelas:      kelasList,
		IDJurusan:    jurusanList,
		CreatedAt:    jadwal.CreatedAt,
		UpdatedAt:    jadwal.UpdatedAt,
	}, nil
}

func (r *jadwalRepository) GetAllWithBankSoal(page, pageSize int) ([]JadwalWithBankSoal, int64, error) {
	var jadwalList []JadwalWithBankSoal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Table("jadwal").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Table("jadwal").
		Select("jadwal.id, jadwal.id_bank_soal, bank_soal.nama_bank_soal, jadwal.nama_ujian, jadwal.tingkat, TO_CHAR(jadwal.wkt_mulai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_mulai, TO_CHAR(jadwal.wkt_selesai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_selesai, jadwal.durasi, jadwal.acak_soal::int AS acak_soal, jadwal.acak_opsi::int AS acak_opsi, jadwal.wajib_token::int AS wajib_token, TO_CHAR(jadwal.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at, TO_CHAR(jadwal.updated_at, 'YYYY-MM-DD HH24:MI:SS') AS updated_at").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Scan(&jadwalList).Error

	return jadwalList, total, err
}

func (r *jadwalRepository) GetByBankSoalID(bankSoalID string, page, pageSize int) ([]JadwalWithBankSoal, int64, error) {
	var jadwalList []JadwalWithBankSoal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Table("jadwal").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.id_bank_soal = ? AND jadwal.deleted_at IS NULL", bankSoalID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Table("jadwal").
		Select("jadwal.id, jadwal.id_bank_soal, bank_soal.nama_bank_soal, jadwal.nama_ujian, jadwal.tingkat, TO_CHAR(jadwal.wkt_mulai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_mulai, TO_CHAR(jadwal.wkt_selesai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_selesai, jadwal.durasi, jadwal.acak_soal::int AS acak_soal, jadwal.acak_opsi::int AS acak_opsi, jadwal.wajib_token::int AS wajib_token, TO_CHAR(jadwal.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at, TO_CHAR(jadwal.updated_at, 'YYYY-MM-DD HH24:MI:SS') AS updated_at").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.id_bank_soal = ? AND jadwal.deleted_at IS NULL", bankSoalID).
		Offset(offset).
		Limit(pageSize).
		Scan(&jadwalList).Error

	return jadwalList, total, err
}

func (r *jadwalRepository) GetAktifHariIniByKelas(idKelas, idPeserta string) ([]JadwalAktifWithStatus, error) {
	var jadwalList []JadwalAktifWithStatus

	err := r.db.
		Table("jadwal").
		Select(`
			jadwal.id,
			jadwal.id_bank_soal,
			bank_soal.nama_bank_soal,
			jadwal.nama_ujian,
			jadwal.tingkat,
			TO_CHAR(jadwal.wkt_mulai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_mulai,
			TO_CHAR(jadwal.wkt_selesai, 'YYYY-MM-DD HH24:MI:SS') AS wkt_selesai,
			jadwal.durasi,
			jadwal.acak_soal::int AS acak_soal,
			jadwal.acak_opsi::int AS acak_opsi,
			jadwal.wajib_token::int AS wajib_token,
			nilai.id AS id_nilai,
			TO_CHAR(nilai.wkt_selesai, 'YYYY-MM-DD HH24:MI:SS') AS nilai_wkt_selesai
		`).
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Joins("INNER JOIN jadwal_kelas ON jadwal_kelas.id_jadwal = jadwal.id").
		Joins("LEFT JOIN nilai ON nilai.id_jadwal = jadwal.id AND nilai.id_peserta = ? AND nilai.deleted_at IS NULL", idPeserta).
		Where("jadwal_kelas.id_kelas = ?", idKelas).
		Where("jadwal.deleted_at IS NULL").
		Where("DATE(jadwal.wkt_mulai) <= CURRENT_DATE AND DATE(jadwal.wkt_selesai) >= CURRENT_DATE").
		Order("jadwal.wkt_mulai ASC").
		Scan(&jadwalList).Error

	if err != nil {
		return nil, err
	}
	return jadwalList, nil
}

// GetAcakOpsiForPesertaSoal mencari nilai aktif peserta yang jadwalnya
// mengandung soal ini, lalu ambil acak_opsi dari jadwal tersebut.
// Return 0 jika peserta tidak punya ujian aktif untuk soal ini.
func (r *jadwalRepository) GetAcakOpsiForPesertaSoal(pesertaID, soalID string) (int, error) {
	var acakOpsi int
	err := r.db.Table("nilai").
		Select("jadwal.acak_opsi::int AS acak_opsi").
		Joins("INNER JOIN jadwal ON nilai.id_jadwal = jadwal.id AND jadwal.deleted_at IS NULL").
		Joins("INNER JOIN soal ON soal.id_bank_soal = jadwal.id_bank_soal AND soal.id = ? AND soal.deleted_at IS NULL", soalID).
		Where("nilai.id_peserta = ? AND nilai.deleted_at IS NULL AND nilai.wkt_selesai IS NULL", pesertaID).
		Limit(1).
		Scan(&acakOpsi).Error
	return acakOpsi, err
}

func (r *jadwalRepository) Update(jadwal *model.Jadwal) error {
	return r.db.Save(jadwal).Error
}

func (r *jadwalRepository) Delete(id string) error {
	now := time.Now()
	return r.db.Model(&model.Jadwal{}).Where("id = ?", id).Update("deleted_at", now).Error
}

func (r *jadwalRepository) Restore(id string) error {
	return r.db.Model(&model.Jadwal{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NULL")).Error
}
