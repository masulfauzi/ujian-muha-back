package dto

type CreateJadwalKelasRequest struct {
	IDJadwal string `json:"id_jadwal" validate:"required"`
	IDKelas  string `json:"id_kelas" validate:"required"`
}

type UpdateJadwalKelasRequest struct {
	IDJadwal string `json:"id_jadwal" validate:"required"`
	IDKelas  string `json:"id_kelas" validate:"required"`
}

type JadwalKelasResponse struct {
	ID           string `json:"id"`
	IDJadwal     string `json:"id_jadwal"`
	IDKelas      string `json:"id_kelas"`
	NamaKelas    string `json:"nama_kelas"`
	NamaBankSoal string `json:"nama_bank_soal"`
	WktMulai     string `json:"wkt_mulai"`
	WktSelesai   string `json:"wkt_selesai"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type JadwalKelasListResponse struct {
	Data      []JadwalKelasResponse `json:"data"`
	Total     int64                 `json:"total"`
	Page      int                   `json:"page"`
	PageSize  int                   `json:"page_size"`
	TotalPage int                   `json:"total_page"`
}
