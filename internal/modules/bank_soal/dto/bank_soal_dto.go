package dto

type CreateBankSoalRequest struct {
	NamaBankSoal string `json:"nama_bank_soal" validate:"required"`
	IdMapel      string `json:"id_mapel" validate:"required"`
	JmlSoal      int    `json:"jml_soal" validate:"required,min=0"`
	Deskripsi    string `json:"deskripsi"`
}

type UpdateBankSoalRequest struct {
	NamaBankSoal string `json:"nama_bank_soal" validate:"required"`
	IdMapel      string `json:"id_mapel" validate:"required"`
	JmlSoal      int    `json:"jml_soal" validate:"required,min=0"`
	Deskripsi    string `json:"deskripsi"`
}

type BankSoalResponse struct {
	ID           string `json:"id"`
	NamaBankSoal string `json:"nama_bank_soal"`
	IdMapel      string `json:"id_mapel"`
	NamaMapel    string `json:"nama_mapel"`
	JmlSoal      int    `json:"jml_soal"`
	Deskripsi    string `json:"deskripsi"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type BankSoalListResponse struct {
	Data      []BankSoalResponse `json:"data"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	TotalPage int                `json:"total_page"`
}
