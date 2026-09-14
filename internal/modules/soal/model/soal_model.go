package model

import (
	"database/sql"
	"time"
)

type Soal struct {
	ID           string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IdBankSoal   string         `gorm:"type:uuid;index" json:"id_bank_soal"`
	NoSoal       int            `gorm:"type:integer" json:"no_soal"`
	Soal         string         `gorm:"type:text" json:"soal"`
	GambarSoal   string         `gorm:"type:varchar(500)" json:"gambar_soal"`
	OpsiA        string         `gorm:"type:text" json:"opsi_a"`
	OpsiB        string         `gorm:"type:text" json:"opsi_b"`
	OpsiC        string         `gorm:"type:text" json:"opsi_c"`
	OpsiD        string         `gorm:"type:text" json:"opsi_d"`
	OpsiE        string         `gorm:"type:text" json:"opsi_e"`
	GambarA      string         `gorm:"type:varchar(500)" json:"gambar_a"`
	GambarB      string         `gorm:"type:varchar(500)" json:"gambar_b"`
	GambarC      string         `gorm:"type:varchar(500)" json:"gambar_c"`
	GambarD      string         `gorm:"type:varchar(500)" json:"gambar_d"`
	GambarE      string         `gorm:"type:varchar(500)" json:"gambar_e"`
	Kunci        string         `gorm:"type:varchar(1)" json:"kunci"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time     `gorm:"index" json:"deleted_at"`
	CreatedBy    sql.NullString `gorm:"type:uuid" json:"created_by"`
	UpdatedBy    sql.NullString `gorm:"type:uuid" json:"updated_by"`
}

func (Soal) TableName() string {
	return "soal"
}
