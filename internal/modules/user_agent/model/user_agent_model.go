package model

import "time"

// UserAgent menyimpan daftar user agent (browser/aplikasi) yang diizinkan
// untuk mengerjakan ujian. Dicocokkan terhadap header User-Agent request
// peserta saat login/mengerjakan ujian.
type UserAgent struct {
	ID         string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserAgent  string     `gorm:"type:varchar(500);not null;uniqueIndex:idx_user_agent_value,where:deleted_at IS NULL" json:"user_agent"`
	Keterangan string     `gorm:"type:text" json:"keterangan"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at"`
}

func (UserAgent) TableName() string {
	return "user_agent"
}
