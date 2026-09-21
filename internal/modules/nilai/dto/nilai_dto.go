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

// MonitoringPesertaResponse adalah status pengerjaan satu peserta untuk satu jadwal,
// dikembalikan endpoint GET /nilai/monitoring/{id_jadwal}. IDNilai nil jika peserta
// belum pernah memulai ujian ini sama sekali (Status = "belum_mulai").
type MonitoringPesertaResponse struct {
	IDPeserta         string  `json:"id_peserta"`
	NamaPeserta       string  `json:"nama_peserta"`
	Username          string  `json:"username"`
	IDKelas           string  `json:"id_kelas"`
	NamaKelas         string  `json:"nama_kelas"`
	IDNilai           *string `json:"id_nilai"`
	Status            string  `json:"status"` // belum_mulai | sedang_mengerjakan | selesai
	Nilai             float64 `json:"nilai"`
	WktMulai          *string `json:"wkt_mulai"`
	AktivitasTerakhir *string `json:"aktivitas_terakhir"`
	WktSelesai        *string `json:"wkt_selesai"`
}

// MonitoringSummary adalah ringkasan jumlah peserta per status.
type MonitoringSummary struct {
	Total             int `json:"total"`
	BelumMulai        int `json:"belum_mulai"`
	SedangMengerjakan int `json:"sedang_mengerjakan"`
	Selesai           int `json:"selesai"`
}

// MonitoringResponse adalah response endpoint GET /nilai/monitoring/{id_jadwal}.
type MonitoringResponse struct {
	NamaUjian string                      `json:"nama_ujian"`
	Summary   MonitoringSummary           `json:"summary"`
	Data      []MonitoringPesertaResponse `json:"data"`
}

// BulkSelesaikanErrorDetail adalah detail kegagalan satu sesi nilai saat proses
// "selesaikan semua" (dilaporkan tapi tidak menggagalkan sesi lain yang lolos).
type BulkSelesaikanErrorDetail struct {
	IDNilai string `json:"id_nilai"`
	Error   string `json:"error"`
}

// BulkSelesaikanResponse adalah response endpoint POST /nilai/monitoring/{id_jadwal}/selesaikan-semua.
type BulkSelesaikanResponse struct {
	TotalDiproses int                         `json:"total_diproses"`
	TotalBerhasil int                         `json:"total_berhasil"`
	TotalGagal    int                         `json:"total_gagal"`
	Errors        []BulkSelesaikanErrorDetail `json:"errors"`
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
