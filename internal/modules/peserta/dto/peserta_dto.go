package dto

type CreatePesertaRequest struct {
	Nama     string `json:"nama" validate:"required"`
	IDKelas  string `json:"id_kelas" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdatePesertaRequest struct {
	Nama     string `json:"nama" validate:"required"`
	IDKelas  string `json:"id_kelas" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password"`
}

type PesertaResponse struct {
	ID        string `json:"id"`
	Nama      string `json:"nama"`
	IDKelas   string `json:"id_kelas"`
	NamaKelas string `json:"nama_kelas"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PesertaListResponse struct {
	Data      []PesertaResponse `json:"data"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	TotalPage int               `json:"total_page"`
}
