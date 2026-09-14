package model

import (
	"time"

	"gorm.io/gorm"
)

type Peserta struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Nama      string         `gorm:"type:varchar(255);not null" json:"nama"`
	IDKelas   string         `gorm:"type:uuid;not null;index" json:"id_kelas"`
	Username  string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_peserta_username,where:deleted_at IS NULL" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Peserta) TableName() string {
	return "peserta"
}
