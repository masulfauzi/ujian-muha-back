package model

import "time"

// LoginLog mencatat SETIAP percobaan login (berhasil maupun gagal) beserta
// konteks perangkat/aplikasi yang dipakai, untuk kebutuhan audit/keamanan.
// Dicatat otomatis oleh AuthController.Login — tidak ada endpoint publik
// untuk menulis baris ini secara langsung.
type LoginLog struct {
	ID             string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username       string    `gorm:"type:varchar(255);not null;index" json:"username"`
	IDUser         string    `gorm:"type:varchar(255)" json:"id_user"`
	Role           string    `gorm:"type:varchar(50)" json:"role"`
	Success        int       `gorm:"not null;default:0" json:"success"`
	Keterangan     string    `gorm:"type:text" json:"keterangan"`
	UserAgent      string    `gorm:"type:varchar(500)" json:"user_agent"`
	XRequestedWith string    `gorm:"column:x_requested_with;type:varchar(255)" json:"x_requested_with"`
	IPAddress      string    `gorm:"type:varchar(100)" json:"ip_address"`
	CreatedAt      time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

func (LoginLog) TableName() string {
	return "login_log"
}
