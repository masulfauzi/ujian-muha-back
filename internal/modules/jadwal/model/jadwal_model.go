package model

import (
	"time"
)

type Jadwal struct {
	ID          string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IDBankSoal  string     `gorm:"type:uuid;not null;index" json:"id_bank_soal"`
	NamaUjian   string     `gorm:"type:varchar(255);not null" json:"nama_ujian"`
	Tingkat     string     `gorm:"type:varchar(10);not null" json:"tingkat"`
	WktMulai    time.Time  `gorm:"type:timestamp;not null" json:"wkt_mulai"`
	WktSelesai  time.Time  `gorm:"type:timestamp;not null" json:"wkt_selesai"`
	Durasi      int        `gorm:"type:integer;not null" json:"durasi"`
	AcakSoal    int        `gorm:"type:smallint;not null;default:0" json:"acak_soal"`
	AcakOpsi    int        `gorm:"type:smallint;not null;default:0" json:"acak_opsi"`
	WajibToken  int        `gorm:"type:smallint;not null;default:0" json:"wajib_token"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
}

func (Jadwal) TableName() string {
	return "jadwal"
}
