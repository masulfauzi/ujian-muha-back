package dto

type CreateMapelRequest struct {
	NamaMapel string `json:"nama_mapel" validate:"required"`
	KodeMapel string `json:"kode_mapel" validate:"required,max=20"`
	Deskripsi string `json:"deskripsi"`
}

type UpdateMapelRequest struct {
	NamaMapel string `json:"nama_mapel" validate:"required"`
	KodeMapel string `json:"kode_mapel" validate:"required,max=20"`
	Deskripsi string `json:"deskripsi"`
}

type MapelResponse struct {
	ID        string `json:"id"`
	NamaMapel string `json:"nama_mapel"`
	KodeMapel string `json:"kode_mapel"`
	Deskripsi string `json:"deskripsi"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type MapelListResponse struct {
	Data      []MapelResponse `json:"data"`
	Total     int64           `json:"total"`
	Page      int             `json:"page"`
	PageSize  int             `json:"page_size"`
	TotalPage int             `json:"total_page"`
}
