package dto

type SectionItem struct {
	NamaSection        string `json:"nama_section" validate:"required"`
	JmlSoal            int    `json:"jml_soal" validate:"required,min=1"`
	DurasiMenitMinimal int    `json:"durasi_menit_minimal" validate:"min=0"`
}

type DefineSectionRequest struct {
	Sections []SectionItem `json:"sections" validate:"required,min=1,dive"`
}

type SectionResponse struct {
	ID                 string `json:"id"`
	IDJadwal           string `json:"id_jadwal"`
	NamaSection        string `json:"nama_section"`
	Urutan             int    `json:"urutan"`
	NoUrutAwal         int    `json:"no_urut_awal"`
	NoUrutAkhir        int    `json:"no_urut_akhir"`
	JmlSoal            int    `json:"jml_soal"`
	DurasiMenitMinimal int    `json:"durasi_menit_minimal"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

type SectionListResponse struct {
	Data []SectionResponse `json:"data"`
}
