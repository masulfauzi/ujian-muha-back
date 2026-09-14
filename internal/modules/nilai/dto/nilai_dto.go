package dto

// MulaiUjianRequest adalah body opsional untuk endpoint mulai-ujian.
// Token hanya wajib diisi jika jadwal terkait mengaktifkan wajib_token.
type MulaiUjianRequest struct {
	Token string `json:"token"`
}

type CreateNilaiRequest struct {
	IDPeserta         string  `json:"id_peserta" validate:"required"`
	IDJadwal          string  `json:"id_jadwal" validate:"required"`
	Nilai             float64 `json:"nilai" validate:"required,min=0,max=100"`
	WktMulai          *string `json:"wkt_mulai"`
	AktivitasTerakhir *string `json:"aktivitas_terakhir"`
	WktSelesai        *string `json:"wkt_selesai"`
}

type UpdateNilaiRequest struct {
	IDPeserta         *string  `json:"id_peserta"`
	IDJadwal          *string  `json:"id_jadwal"`
	Nilai             *float64 `json:"nilai"`
	WktMulai          *string  `json:"wkt_mulai"`
	AktivitasTerakhir *string  `json:"aktivitas_terakhir"`
	WktSelesai        *string  `json:"wkt_selesai"`
}

type NilaiResponse struct {
	ID                string  `json:"id"`
	IDPeserta         string  `json:"id_peserta"`
	NamaPeserta       string  `json:"nama_peserta"`
	IDJadwal          string  `json:"id_jadwal"`
	NamaUjian         string  `json:"nama_ujian"`
	Nilai             float64 `json:"nilai"`
	WktMulai          *string `json:"wkt_mulai"`
	AktivitasTerakhir *string `json:"aktivitas_terakhir"`
	WktSelesai        *string `json:"wkt_selesai"`
	IDSectionAktif    *string `json:"id_section_aktif"`
	WktMulaiSection   *string `json:"wkt_mulai_section"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

type NilaiListResponse struct {
	Data      []NilaiResponse `json:"data"`
	Total     int64           `json:"total"`
	Page      int             `json:"page"`
	PageSize  int             `json:"page_size"`
	TotalPage int             `json:"total_page"`
}

// SectionProgressResponse dikembalikan oleh endpoint next-section & section-status.
type SectionProgressResponse struct {
	IDNilai            string `json:"id_nilai"`
	IDSection          string `json:"id_section"`
	NamaSection        string `json:"nama_section"`
	Urutan             int    `json:"urutan"`
	DurasiMenitMinimal int    `json:"durasi_menit_minimal"`
	SudahLanjut        bool   `json:"sudah_lanjut"`
	BolehLanjut        bool   `json:"boleh_lanjut"`
	SisaDetik          int    `json:"sisa_detik"`
	IsSectionTerakhir  bool   `json:"is_section_terakhir"`
}
