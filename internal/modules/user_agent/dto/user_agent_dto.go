package dto

type CreateUserAgentRequest struct {
	UserAgent  string `json:"user_agent" validate:"required,max=500"`
	Keterangan string `json:"keterangan"`
}

type UpdateUserAgentRequest struct {
	UserAgent  string `json:"user_agent" validate:"required,max=500"`
	Keterangan string `json:"keterangan"`
}

type UserAgentResponse struct {
	ID         string `json:"id"`
	UserAgent  string `json:"user_agent"`
	Keterangan string `json:"keterangan"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type UserAgentListResponse struct {
	Data      []UserAgentResponse `json:"data"`
	Total     int64               `json:"total"`
	Page      int                 `json:"page"`
	PageSize  int                 `json:"page_size"`
	TotalPage int                 `json:"total_page"`
}
