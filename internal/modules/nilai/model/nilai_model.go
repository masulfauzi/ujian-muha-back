package model

import (
	"time"
)

type Nilai struct {
	ID        string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IDPeserta          string     `gorm:"type:uuid;not null;index" json:"id_peserta"`
	IDJadwal           string     `gorm:"type:uuid;not null;index" json:"id_jadwal"`
	Nilai              float64    `gorm:"type:float;not null" json:"nilai"`
	WktMulai           *time.Time `gorm:"type:timestamp" json:"wkt_mulai"`
	AktivitasTerakhir  *time.Time `gorm:"type:timestamp" json:"aktivitas_terakhir"`
	WktSelesai         *time.Time `gorm:"type:timestamp" json:"wkt_selesai"`
	IDSectionAktif     *string    `gorm:"type:uuid" json:"id_section_aktif"`
	WktMulaiSection    *time.Time `gorm:"type:timestamp" json:"wkt_mulai_section"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

func (Nilai) TableName() string {
	return "nilai"
}
