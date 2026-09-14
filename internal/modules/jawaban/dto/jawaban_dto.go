package dto

type CreateJawabanRequest struct {
	IDNilai   string `json:"id_nilai" validate:"required"`
	IDSoal    string `json:"id_soal" validate:"required"`
	IDPeserta string `json:"id_peserta" validate:"required"`
	NoUrut    int    `json:"no_urut" validate:"required"`
	Jawaban   string `json:"jawaban" validate:"required,oneof=A B C D E"`
}

type UpdateJawabanRequest struct {
	Jawaban string `json:"jawaban" validate:"required,oneof=A B C D E"`
}

type JawabanResponse struct {
	ID          string `json:"id"`
	IDNilai     string `json:"id_nilai"`
	IDSoal      string `json:"id_soal"`
	NoUrut      int    `json:"no_urut"`
	NoSoal      int    `json:"no_soal"`
	SoalText    string `json:"soal"`
	Kunci       string `json:"kunci"`
	OpsiA       string `json:"opsi_a"`
	OpsiB       string `json:"opsi_b"`
	OpsiC       string `json:"opsi_c"`
	OpsiD       string `json:"opsi_d"`
	OpsiE       string `json:"opsi_e"`
	GambarA     string `json:"gambar_a"`
	GambarB     string `json:"gambar_b"`
	GambarC     string `json:"gambar_c"`
	GambarD     string `json:"gambar_d"`
	GambarE     string `json:"gambar_e"`
	IDPeserta   string `json:"id_peserta"`
	NamaPeserta string `json:"nama_peserta"`
	Jawaban     *string `json:"jawaban"`
	IsBenar     *int   `json:"is_benar"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type JawabanListResponse struct {
	Data      []JawabanResponse `json:"data"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	TotalPage int               `json:"total_page"`
}
