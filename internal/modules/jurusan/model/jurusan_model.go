package model

import (
	"time"

	"gorm.io/gorm"
)

type Jurusan struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	NamaJurusan string         `gorm:"type:varchar(255);uniqueIndex:idx_jurusan_nama_active,where:deleted_at IS NULL" json:"nama_jurusan"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Jurusan) TableName() string {
	return "jurusan"
}
