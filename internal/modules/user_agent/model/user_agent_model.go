package model

import "time"

// UserAgent menyimpan daftar kombinasi header User-Agent + X-Requested-With yang
// diizinkan untuk mengerjakan ujian. Dicocokkan terhadap request peserta saat
// memulai pengerjaan ujian (lihat middleware.CheckAllowedUserAgent).
type UserAgent struct {
	ID             string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserAgent      string     `gorm:"type:varchar(500);not null;uniqueIndex:idx_user_agent_combo,where:deleted_at IS NULL" json:"user_agent"`
	XRequestedWith string     `gorm:"column:x_requested_with;type:varchar(255);not null;uniqueIndex:idx_user_agent_combo,where:deleted_at IS NULL" json:"x_requested_with"`
	Keterangan     string     `gorm:"type:text" json:"keterangan"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at"`
}

func (UserAgent) TableName() string {
	return "user_agent"
}
