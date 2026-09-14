package model

import (
	"time"
)

type JadwalKelas struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IDJadwal  string    `gorm:"type:uuid;not null;index" json:"id_jadwal"`
	IDKelas   string    `gorm:"type:uuid;not null;index" json:"id_kelas"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (JadwalKelas) TableName() string {
	return "jadwal_kelas"
}
