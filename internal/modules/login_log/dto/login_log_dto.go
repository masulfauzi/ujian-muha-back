package dto

type LoginLogResponse struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	IDUser         string `json:"id_user"`
	Role           string `json:"role"`
	Success        bool   `json:"success"`
	Keterangan     string `json:"keterangan"`
	UserAgent      string `json:"user_agent"`
	XRequestedWith string `json:"x_requested_with"`
	IPAddress      string `json:"ip_address"`
	CreatedAt      string `json:"created_at"`
}

type LoginLogListResponse struct {
	Data      []LoginLogResponse `json:"data"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	TotalPage int                `json:"total_page"`
}
