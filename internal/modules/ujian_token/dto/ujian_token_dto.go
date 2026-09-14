package dto

type CurrentTokenResponse struct {
	Token        string `json:"token"`
	SisaDetik    int    `json:"sisa_detik"`
	PeriodeDetik int    `json:"periode_detik"`
}
