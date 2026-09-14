package dto

type CreateJurusanRequest struct {
	NamaJurusan string `json:"nama_jurusan" validate:"required"`
}

type UpdateJurusanRequest struct {
	NamaJurusan string `json:"nama_jurusan" validate:"required"`
}

type JurusanResponse struct {
	ID          string `json:"id"`
	NamaJurusan string `json:"nama_jurusan"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type JurusanListResponse struct {
	Data      []JurusanResponse `json:"data"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	TotalPage int               `json:"total_page"`
}
