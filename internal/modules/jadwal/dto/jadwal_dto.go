package dto

type CreateJadwalRequest struct {
	IDBankSoal string   `json:"id_bank_soal" validate:"required"`
	NamaUjian  string   `json:"nama_ujian" validate:"required"`
	Tingkat    string   `json:"tingkat" validate:"required"`
	WktMulai   string   `json:"wkt_mulai" validate:"required"`
	WktSelesai string   `json:"wkt_selesai" validate:"required"`
	Durasi     int      `json:"durasi" validate:"required,min=1"`
	AcakSoal   int      `json:"acak_soal"`
	AcakOpsi   int      `json:"acak_opsi"`
	WajibToken int      `json:"wajib_token"`
	IDKelas    []string `json:"id_kelas" validate:"required"`
}

type UpdateJadwalRequest struct {
	IDBankSoal string   `json:"id_bank_soal" validate:"required"`
	NamaUjian  string   `json:"nama_ujian" validate:"required"`
	Tingkat    string   `json:"tingkat" validate:"required"`
	WktMulai   string   `json:"wkt_mulai" validate:"required"`
	WktSelesai string   `json:"wkt_selesai" validate:"required"`
	Durasi     int      `json:"durasi" validate:"required,min=1"`
	AcakSoal   int      `json:"acak_soal"`
	AcakOpsi   int      `json:"acak_opsi"`
	WajibToken int      `json:"wajib_token"`
	IDKelas    []string `json:"id_kelas" validate:"required"`
}

type KelasItem struct {
	ID        string `json:"id"`
	IDKelas   string `json:"id_kelas"`
	NamaKelas string `json:"nama_kelas"`
}

type JurusanItem struct {
	ID          string `json:"id"`
	IDJurusan   string `json:"id_jurusan"`
	NamaJurusan string `json:"nama_jurusan"`
}

type JadwalResponse struct {
	ID           string         `json:"id"`
	IDBankSoal   string         `json:"id_bank_soal"`
	NamaBankSoal string         `json:"nama_bank_soal"`
	NamaUjian    string         `json:"nama_ujian"`
	Tingkat      string         `json:"tingkat"`
	WktMulai     string         `json:"wkt_mulai"`
	WktSelesai   string         `json:"wkt_selesai"`
	Durasi       int            `json:"durasi"`
	AcakSoal     int            `json:"acak_soal"`
	AcakOpsi     int            `json:"acak_opsi"`
	WajibToken   int            `json:"wajib_token"`
	IDKelas      []KelasItem    `json:"id_kelas"`
	IDJurusan    []JurusanItem  `json:"id_jurusan"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
}

type JadwalListResponse struct {
	Data      []JadwalResponse `json:"data"`
	Total     int64            `json:"total"`
	Page      int              `json:"page"`
	PageSize  int              `json:"page_size"`
	TotalPage int              `json:"total_page"`
}

type JadwalAktifResponse struct {
	ID               string  `json:"id"`
	IDBankSoal       string  `json:"id_bank_soal"`
	NamaBankSoal     string  `json:"nama_bank_soal"`
	NamaUjian        string  `json:"nama_ujian"`
	Tingkat          string  `json:"tingkat"`
	WktMulai         string  `json:"wkt_mulai"`
	WktSelesai       string  `json:"wkt_selesai"`
	Durasi           int     `json:"durasi"`
	AcakSoal         int     `json:"acak_soal"`
	AcakOpsi         int     `json:"acak_opsi"`
	WajibToken       int     `json:"wajib_token"`
	IDNilai          *string `json:"id_nilai"`
	StatusPengerjaan string  `json:"status_pengerjaan"`
}
