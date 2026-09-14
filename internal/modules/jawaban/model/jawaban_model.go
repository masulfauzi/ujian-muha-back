package model

import (
	"time"
)

type Jawaban struct {
	ID        string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IDNilai   string     `gorm:"type:uuid;not null;index" json:"id_nilai"`
	IDSoal    string     `gorm:"type:uuid;not null;index" json:"id_soal"`
	IDPeserta string     `gorm:"type:uuid;not null;index" json:"id_peserta"`
	NoUrut    int        `gorm:"type:int;not null;default:0" json:"no_urut"`
	Jawaban   *string `gorm:"type:varchar(1)" json:"jawaban"`
	IsBenar   *int    `gorm:"type:smallint" json:"is_benar"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

func (Jawaban) TableName() string {
	return "jawaban"
}
