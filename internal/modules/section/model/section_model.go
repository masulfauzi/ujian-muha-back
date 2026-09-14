package model

import (
	"time"
)

type Section struct {
	ID                 string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IDJadwal           string     `gorm:"type:uuid;not null;index" json:"id_jadwal"`
	NamaSection        string     `gorm:"type:varchar(255);not null" json:"nama_section"`
	Urutan             int        `gorm:"type:integer;not null" json:"urutan"`
	NoUrutAwal         int        `gorm:"type:integer;not null" json:"no_urut_awal"`
	NoUrutAkhir        int        `gorm:"type:integer;not null" json:"no_urut_akhir"`
	DurasiMenitMinimal int        `gorm:"type:integer;not null" json:"durasi_menit_minimal"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          *time.Time `gorm:"index" json:"deleted_at"`
}

func (Section) TableName() string {
	return "section"
}
