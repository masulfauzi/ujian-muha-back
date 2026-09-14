package model

import (
	"time"

	"gorm.io/gorm"
)

type Kelas struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IDJurusan string         `gorm:"type:uuid;not null;index" json:"id_jurusan"`
	NamaKelas string         `gorm:"type:varchar(255);not null" json:"nama_kelas"`
	Tingkat   string         `gorm:"type:varchar(10);not null;index" json:"tingkat"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Kelas) TableName() string {
	return "kelas"
}
