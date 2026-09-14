package service

import (
	"bytes"
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/peserta/dto"
	"backend/internal/modules/peserta/model"
	"backend/internal/modules/peserta/repository"
	"backend/internal/utils"

	"github.com/go-pdf/fpdf"
	"gorm.io/gorm"
)

func pesertaWithKelasToResponse(p *repository.PesertaWithKelas) *dto.PesertaResponse {
	return &dto.PesertaResponse{
		ID:        p.ID,
		Nama:      p.Nama,
		IDKelas:   p.IDKelas,
		NamaKelas: p.NamaKelas,
		Username:  p.Username,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

type PesertaService interface {
	CreatePeserta(req *dto.CreatePesertaRequest) (*dto.PesertaResponse, error)
	GetPesertaByID(id string) (*dto.PesertaResponse, error)
	GetAllPeserta(page, pageSize int, idKelas string) (*dto.PesertaListResponse, error)
	UpdatePeserta(id string, req *dto.UpdatePesertaRequest) (*dto.PesertaResponse, error)
	DeletePeserta(id string) error
	RestorePeserta(id string) error
	GenerateKartuUjianPDF(idKelas string) ([]byte, error)
}

type pesertaService struct {
	repo repository.PesertaRepository
}

func NewPesertaService(repo repository.PesertaRepository) PesertaService {
	return &pesertaService{repo: repo}
}

func (s *pesertaService) CreatePeserta(req *dto.CreatePesertaRequest) (*dto.PesertaResponse, error) {
	existing, err := s.repo.GetByUsername(req.Username)
	if err == nil && existing != nil {
		return nil, errors.New("username sudah digunakan")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	peserta := &model.Peserta{
		Nama:     req.Nama,
		IDKelas:  req.IDKelas,
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := s.repo.Create(peserta); err != nil {
		return nil, err
	}

	created, err := s.repo.GetByID(peserta.ID)
	if err != nil {
		return nil, err
	}

	return pesertaWithKelasToResponse(created), nil
}

func (s *pesertaService) GetPesertaByID(id string) (*dto.PesertaResponse, error) {
	peserta, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	return pesertaWithKelasToResponse(peserta), nil
}

func (s *pesertaService) GetAllPeserta(page, pageSize int, idKelas string) (*dto.PesertaListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	pesertaList, total, err := s.repo.GetAll(page, pageSize, idKelas)
	if err != nil {
		return nil, err
	}

	var responses []dto.PesertaResponse
	for _, p := range pesertaList {
		responses = append(responses, *pesertaWithKelasToResponse(&p))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.PesertaListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *pesertaService) UpdatePeserta(id string, req *dto.UpdatePesertaRequest) (*dto.PesertaResponse, error) {
	existing, err := s.repo.GetRawByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	if req.Username != existing.Username {
		taken, err := s.repo.GetByUsername(req.Username)
		if err == nil && taken != nil && taken.ID != id {
			return nil, errors.New("username sudah digunakan")
		}
	}

	peserta := &model.Peserta{
		ID:       existing.ID,
		Nama:     req.Nama,
		IDKelas:  req.IDKelas,
		Username: req.Username,
		Password: existing.Password,
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New("gagal memproses password")
		}
		peserta.Password = hashedPassword
	}

	if err := s.repo.Update(peserta); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return pesertaWithKelasToResponse(updated), nil
}

func (s *pesertaService) DeletePeserta(id string) error {
	peserta, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}
	return s.repo.Delete(peserta.ID)
}

func (s *pesertaService) RestorePeserta(id string) error {
	return s.repo.Restore(id)
}

// GenerateKartuUjianPDF membuat PDF kartu peserta ujian untuk satu kelas, ditata sebagai
// grid kartu (2 kolom x 5 baris per halaman A4) dengan garis putus-putus di tiap kartu
// sebagai panduan gunting. Kartu bersifat global (tidak terikat jadwal/ujian tertentu) dan
// tidak menampilkan password — hanya nama, username, dan kelas.
func (s *pesertaService) GenerateKartuUjianPDF(idKelas string) ([]byte, error) {
	pesertaList, total, err := s.repo.GetAll(1, 99999, idKelas)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return nil, errors.New("kelas tidak ditemukan atau belum memiliki peserta")
	}

	return buildKartuUjianPDF(pesertaList)
}

func buildKartuUjianPDF(pesertaList []repository.PesertaWithKelas) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)

	const (
		cols, rows       = 2, 5
		cardW, cardH     = 90.0, 50.0
		marginX, marginY = 10.0, 10.0
		gapX, gapY       = 10.0, 5.0
	)
	perPage := cols * rows

	for i, p := range pesertaList {
		if i%perPage == 0 {
			pdf.AddPage()
		}
		idx := i % perPage
		col := idx % cols
		row := idx / cols
		x := marginX + float64(col)*(cardW+gapX)
		y := marginY + float64(row)*(cardH+gapY)

		pdf.SetDrawColor(120, 120, 120)
		pdf.SetDashPattern([]float64{2, 2}, 0)
		pdf.Rect(x, y, cardW, cardH, "D")
		pdf.SetDashPattern([]float64{}, 0)

		pdf.SetXY(x, y+5)
		pdf.SetFont("Helvetica", "B", 11)
		pdf.CellFormat(cardW, 5, "KARTU PESERTA UJIAN", "", 0, "C", false, 0, "")

		pdf.SetDrawColor(0, 0, 0)
		pdf.Line(x+6, y+12, x+cardW-6, y+12)

		lineY := y + 20
		writeField := func(label, value string) {
			pdf.SetXY(x+6, lineY)
			pdf.SetFont("Helvetica", "", 10)
			pdf.CellFormat(24, 6, label, "", 0, "L", false, 0, "")
			pdf.CellFormat(4, 6, ":", "", 0, "L", false, 0, "")
			pdf.SetFont("Helvetica", "B", 10)
			pdf.CellFormat(cardW-6-24-4-6, 6, value, "", 0, "L", false, 0, "")
			lineY += 7
		}
		writeField("Nama", p.Nama)
		writeField("Username", p.Username)
		writeField("Kelas", p.NamaKelas)
	}

	if !pdf.Ok() {
		return nil, pdf.Error()
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
