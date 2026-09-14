package dto

import (
	"mime/multipart"
	"time"
)

type CreateSoalRequest struct {
	IdBankSoal  string                 `json:"id_bank_soal" form:"id_bank_soal" validate:"required"`
	NoSoal      int                    `json:"no_soal" form:"no_soal"`
	Soal        string                 `json:"soal" form:"soal"`
	GambarSoal  *multipart.FileHeader  `form:"gambar_soal"`
	OpsiA       string                 `json:"opsi_a" form:"opsi_a"`
	OpsiB       string                 `json:"opsi_b" form:"opsi_b"`
	OpsiC       string                 `json:"opsi_c" form:"opsi_c"`
	OpsiD       string                 `json:"opsi_d" form:"opsi_d"`
	OpsiE       string                 `json:"opsi_e" form:"opsi_e"`
	GambarA     *multipart.FileHeader  `form:"gambar_a"`
	GambarB     *multipart.FileHeader  `form:"gambar_b"`
	GambarC     *multipart.FileHeader  `form:"gambar_c"`
	GambarD     *multipart.FileHeader  `form:"gambar_d"`
	GambarE     *multipart.FileHeader  `form:"gambar_e"`
	Kunci       string                 `json:"kunci" form:"kunci" validate:"required,len=1"`
}

type UpdateSoalRequest struct {
	NoSoal     int                    `json:"no_soal" form:"no_soal"`
	Soal       string                 `json:"soal" form:"soal"`
	GambarSoal *multipart.FileHeader  `form:"gambar_soal"`
	OpsiA      string                 `json:"opsi_a" form:"opsi_a"`
	OpsiB      string                 `json:"opsi_b" form:"opsi_b"`
	OpsiC      string                 `json:"opsi_c" form:"opsi_c"`
	OpsiD      string                 `json:"opsi_d" form:"opsi_d"`
	OpsiE      string                 `json:"opsi_e" form:"opsi_e"`
	GambarA    *multipart.FileHeader  `form:"gambar_a"`
	GambarB    *multipart.FileHeader  `form:"gambar_b"`
	GambarC    *multipart.FileHeader  `form:"gambar_c"`
	GambarD    *multipart.FileHeader  `form:"gambar_d"`
	GambarE    *multipart.FileHeader  `form:"gambar_e"`
	Kunci      string                 `json:"kunci" form:"kunci" validate:"required,len=1"`
}

type SoalResponse struct {
	ID         string `json:"id"`
	IdBankSoal string `json:"id_bank_soal"`
	NoSoal     int    `json:"no_soal"`
	Soal       string `json:"soal"`
	GambarSoal string `json:"gambar_soal"`
	OpsiA      string `json:"opsi_a"`
	OpsiB      string `json:"opsi_b"`
	OpsiC      string `json:"opsi_c"`
	OpsiD      string `json:"opsi_d"`
	OpsiE      string `json:"opsi_e"`
	GambarA    string `json:"gambar_a"`
	GambarB    string `json:"gambar_b"`
	GambarC    string `json:"gambar_c"`
	GambarD    string `json:"gambar_d"`
	GambarE    string `json:"gambar_e"`
	Kunci      string `json:"kunci"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type SoalListResponse struct {
	Data      []SoalResponse `json:"data"`
	Total     int64          `json:"total"`
	Page      int            `json:"page"`
	PageSize  int            `json:"page_size"`
	TotalPage int            `json:"total_page"`
}

// ImportSoalRequest adalah request untuk import soal dari excel
type ImportSoalRequest struct {
	IdBankSoal string                `form:"id_bank_soal" validate:"required"`
	File       *multipart.FileHeader `form:"file" validate:"required"`
}

// ImportSoalErrorDetail adalah detail error per row
type ImportSoalErrorDetail struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

// ImportSoalResponse adalah response dari import
type ImportSoalResponse struct {
	TotalProcessed int                     `json:"total_processed"`
	TotalSuccess   int                     `json:"total_success"`
	TotalFailed    int                     `json:"total_failed"`
	ImportID       string                  `json:"import_id"`
	Timestamp      time.Time               `json:"timestamp"`
	Summary        map[string]int          `json:"summary"`
	Errors         []ImportSoalErrorDetail `json:"errors"`
}
