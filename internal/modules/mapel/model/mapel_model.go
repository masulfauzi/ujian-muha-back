package model

import (
	"database/sql"
	"time"
)

type Mapel struct {
	ID        string       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	NamaMapel string       `gorm:"type:varchar(255);uniqueIndex" json:"nama_mapel"`
	KodeMapel string       `gorm:"type:varchar(20);uniqueIndex" json:"kode_mapel"`
	Deskripsi string       `gorm:"type:text" json:"deskripsi"`
	CreatedAt time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time   `gorm:"index" json:"deleted_at"`
	CreatedBy sql.NullString `gorm:"type:uuid" json:"created_by"`
	UpdatedBy sql.NullString `gorm:"type:uuid" json:"updated_by"`
}

func (Mapel) TableName() string {
	return "mapel"
}
