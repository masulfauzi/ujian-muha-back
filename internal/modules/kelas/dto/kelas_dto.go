package dto

type CreateKelasRequest struct {
	IDJurusan string `json:"id_jurusan" validate:"required"`
	NamaKelas string `json:"nama_kelas" validate:"required"`
	Tingkat   string `json:"tingkat" validate:"required"`
}

type UpdateKelasRequest struct {
	IDJurusan string `json:"id_jurusan" validate:"required"`
	NamaKelas string `json:"nama_kelas" validate:"required"`
	Tingkat   string `json:"tingkat" validate:"required"`
}

type KelasResponse struct {
	ID          string `json:"id"`
	IDJurusan   string `json:"id_jurusan"`
	NamaKelas   string `json:"nama_kelas"`
	Tingkat     string `json:"tingkat"`
	NamaJurusan string `json:"nama_jurusan"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type KelasListResponse struct {
	Data      []KelasResponse `json:"data"`
	Total     int64           `json:"total"`
	Page      int             `json:"page"`
	PageSize  int             `json:"page_size"`
	TotalPage int             `json:"total_page"`
}
