package model

import (
	"database/sql"
	"time"
)

type BankSoal struct {
	ID           string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	NamaBankSoal string         `gorm:"type:varchar(255);uniqueIndex" json:"nama_bank_soal"`
	IdMapel      string         `gorm:"type:uuid;index" json:"id_mapel"`
	JmlSoal      int            `gorm:"type:integer;default:0" json:"jml_soal"`
	Deskripsi    string         `gorm:"type:text" json:"deskripsi"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time     `gorm:"index" json:"deleted_at"`
	CreatedBy    sql.NullString `gorm:"type:uuid" json:"created_by"`
	UpdatedBy    sql.NullString `gorm:"type:uuid" json:"updated_by"`
}

func (BankSoal) TableName() string {
	return "bank_soal"
}
