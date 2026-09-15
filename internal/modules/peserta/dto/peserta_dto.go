package dto

import (
	"mime/multipart"
	"time"
)

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

// ImportPesertaRequest adalah request untuk import peserta dari excel.
// Kelas ditentukan per baris lewat kolom "kelas" (nama kelas) di dalam file,
// bukan lewat field terpisah, karena satu file bisa berisi peserta dari
// beberapa kelas sekaligus.
type ImportPesertaRequest struct {
	File *multipart.FileHeader `form:"file" validate:"required"`
}

// ImportPesertaErrorDetail adalah detail error per row
type ImportPesertaErrorDetail struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

// ImportPesertaResponse adalah response dari import
type ImportPesertaResponse struct {
	TotalProcessed int                        `json:"total_processed"`
	TotalSuccess   int                        `json:"total_success"`
	TotalFailed    int                        `json:"total_failed"`
	Timestamp      time.Time                  `json:"timestamp"`
	Summary        map[string]int             `json:"summary"`
	Errors         []ImportPesertaErrorDetail `json:"errors"`
}

// KelasNotFoundError dikembalikan saat ada satu atau lebih baris di file excel
// yang kolom kelas-nya kosong/tidak cocok dengan kelas manapun di database.
// Seluruh proses import dibatalkan (tidak ada satupun baris yang di-insert)
// begitu error ini terjadi.
type KelasNotFoundError struct {
	Details []ImportPesertaErrorDetail
}

func (e *KelasNotFoundError) Error() string {
	return "import dibatalkan: ada baris dengan kelas yang tidak ditemukan"
}
