package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"backend/internal/assets"
	"backend/internal/constants"
	"backend/internal/modules/peserta/dto"
	"backend/internal/modules/peserta/model"
	"backend/internal/modules/peserta/repository"
	"backend/internal/utils"

	"github.com/go-pdf/fpdf"
	"github.com/xuri/excelize/v2"
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
	DeleteAllPeserta(ctx context.Context) (int64, error)
	GenerateKartuUjianPDF(idKelas string) ([]byte, error)
	ImportPesertaFromExcel(ctx context.Context, req *dto.ImportPesertaRequest) (*dto.ImportPesertaResponse, error)
	GenerateImportTemplate() ([]byte, error)
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

	peserta := &model.Peserta{
		Nama:     req.Nama,
		IDKelas:  req.IDKelas,
		Username: req.Username,
		Password: req.Password,
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
		peserta.Password = req.Password
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

// DeleteAllPeserta menghapus PERMANEN seluruh peserta di seluruh sistem (bukan
// soft-delete, tidak bisa di-restore). Mengembalikan jumlah baris yang terhapus.
func (s *pesertaService) DeleteAllPeserta(ctx context.Context) (int64, error) {
	return s.repo.DeleteAll(ctx)
}

// GenerateKartuUjianPDF membuat PDF kartu peserta ujian untuk satu kelas, ditata sebagai
// grid kartu (2 kolom x 5 baris per halaman A4) dengan garis putus-putus di tiap kartu
// sebagai panduan gunting. Kartu bersifat global (tidak terikat jadwal/ujian tertentu),
// menampilkan logo (lihat internal/assets/logo.png) di pojok kiri atas tiap kartu, serta
// nama, username, password, dan kelas peserta.
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
	// Paginasi kartu dikontrol manual lewat grid (lihat perPage & i%perPage di bawah),
	// bukan lewat fpdf. Auto page break bawaan fpdf.New() (trigger di tinggi_halaman-20mm)
	// harus dimatikan, karena kalau tidak, baris kartu terakhir yang mepet ke trigger itu
	// bisa "terpotong" — border kartu tergambar di satu halaman tapi teksnya otomatis
	// terdorong fpdf ke halaman berikutnya.
	pdf.SetAutoPageBreak(false, 0)

	const (
		cols, rows       = 2, 5
		cardW, cardH     = 90.0, 50.0
		marginX, marginY = 10.0, 10.0
		gapX, gapY       = 10.0, 3.0
		logoSize         = 9.0
	)
	perPage := cols * rows

	logoOpt := fpdf.ImageOptions{ImageType: "PNG"}
	pdf.RegisterImageOptionsReader("kartu-logo", logoOpt, bytes.NewReader(assets.Logo))

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

		pdf.ImageOptions("kartu-logo", x+3, y+2, logoSize, logoSize, false, logoOpt, 0, "")

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
		writeField("Password", p.Password)
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

// ImportPesertaFromExcel membaca file excel (kolom: nama, username, password, kelas) dan
// insert semua peserta valid ke kelas masing-masing sesuai kolom "kelas" di tiap baris.
//
// Kolom kelas divalidasi lebih dulu untuk SEMUA baris sebelum insert apapun dilakukan:
// jika ada baris dengan nama kelas kosong, tidak ditemukan, atau ambigu (nama kelas sama
// dipakai lebih dari satu kelas), seluruh proses import dibatalkan (tidak ada satupun
// baris yang di-insert) dan *dto.KelasNotFoundError dikembalikan berisi detail baris mana
// saja yang bermasalah.
//
// Setelah kolom kelas dipastikan valid, baris dengan username kosong/duplikat (di dalam
// file maupun yang sudah ada di database) atau password kurang dari 6 karakter akan
// ditandai gagal dan dilewati satu-satu, tanpa menggagalkan keseluruhan proses import.
func (s *pesertaService) ImportPesertaFromExcel(ctx context.Context, req *dto.ImportPesertaRequest) (*dto.ImportPesertaResponse, error) {
	file, err := req.File.Open()
	if err != nil {
		return nil, errors.New("gagal membuka file")
	}
	defer file.Close()

	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, errors.New("file bukan format excel yang valid")
	}
	defer xlsx.Close()

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, errors.New("gagal membaca sheet excel")
	}

	var excelRows []*utils.ExcelPesertaRow
	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		row := utils.ParseExcelPesertaRow(rows[rowIndex], rowIndex+1)
		if row.Nama == "" && row.Username == "" && row.Password == "" && row.NamaKelas == "" {
			continue // baris kosong (misal sisa baris kosong di akhir sheet), tidak dihitung
		}
		excelRows = append(excelRows, row)
	}

	// 1. Resolve kolom kelas untuk semua baris lebih dulu; batalkan seluruh import jika
	// ada yang bermasalah.
	kelasByName, ambiguousNames, err := s.resolveKelasNames(ctx, excelRows)
	if err != nil {
		return nil, err
	}

	if kelasErrors := s.validateKelasColumn(excelRows, kelasByName, ambiguousNames); len(kelasErrors) > 0 {
		return nil, &dto.KelasNotFoundError{Details: kelasErrors}
	}

	// 2. Validasi & insert per baris (nama, username, password).
	var pesertaList []model.Peserta
	var errorDetails []dto.ImportPesertaErrorDetail
	successCount := 0

	seenUsernames := make(map[string]int) // username (lowercase) -> row pertama yang memakainya

	for _, excelRow := range excelRows {
		validationErrors := utils.ValidatePesertaRow(excelRow)

		usernameKey := strings.ToLower(excelRow.Username)
		if excelRow.Username != "" {
			if firstRow, dup := seenUsernames[usernameKey]; dup {
				validationErrors = append(validationErrors, fmt.Sprintf("username duplikat dengan row %d di file ini", firstRow))
			} else {
				existing, err := s.repo.GetByUsername(excelRow.Username)
				if err == nil && existing != nil {
					validationErrors = append(validationErrors, "username sudah digunakan")
				}
			}
		}

		if len(validationErrors) > 0 {
			errorDetails = append(errorDetails, dto.ImportPesertaErrorDetail{
				Row:   excelRow.RowIndex,
				Error: strings.Join(validationErrors, "; "),
			})
			continue
		}

		seenUsernames[usernameKey] = excelRow.RowIndex

		pesertaList = append(pesertaList, model.Peserta{
			Nama:     excelRow.Nama,
			IDKelas:  kelasByName[strings.ToLower(strings.TrimSpace(excelRow.NamaKelas))],
			Username: excelRow.Username,
			Password: excelRow.Password,
		})
		successCount++
	}

	if len(pesertaList) > 0 {
		if err := s.repo.BulkCreatePeserta(ctx, pesertaList); err != nil {
			return nil, errors.New("gagal menyimpan data ke database: " + err.Error())
		}
	}

	if len(errorDetails) > 100 {
		errorDetails = errorDetails[:100]
	}

	return &dto.ImportPesertaResponse{
		TotalProcessed: len(excelRows),
		TotalSuccess:   successCount,
		TotalFailed:    len(excelRows) - successCount,
		Timestamp:      time.Now(),
		Summary: map[string]int{
			"inserted": successCount,
			"skipped":  0,
			"errors":   len(excelRows) - successCount,
		},
		Errors: errorDetails,
	}, nil
}

// resolveKelasNames mengambil semua nama kelas unik yang dipakai di excelRows lalu
// mencocokkannya (case-insensitive) ke tabel kelas. Mengembalikan:
//   - kelasByName: map nama kelas (lowercase, trimmed) -> id_kelas, untuk nama yang
//     cocok dengan TEPAT SATU kelas.
//   - ambiguousNames: set nama kelas (lowercase, trimmed) yang cocok dengan LEBIH DARI
//     SATU kelas, sehingga tidak bisa ditentukan otomatis.
func (s *pesertaService) resolveKelasNames(ctx context.Context, excelRows []*utils.ExcelPesertaRow) (map[string]string, map[string]bool, error) {
	uniqueNames := make(map[string]bool)
	for _, row := range excelRows {
		name := strings.ToLower(strings.TrimSpace(row.NamaKelas))
		if name != "" {
			uniqueNames[name] = true
		}
	}

	if len(uniqueNames) == 0 {
		return map[string]string{}, map[string]bool{}, nil
	}

	lowerNamaList := make([]string, 0, len(uniqueNames))
	for name := range uniqueNames {
		lowerNamaList = append(lowerNamaList, name)
	}

	matches, err := s.repo.GetKelasByNamaList(ctx, lowerNamaList)
	if err != nil {
		return nil, nil, errors.New("gagal memvalidasi kolom kelas: " + err.Error())
	}

	countByName := make(map[string]int)
	kelasByName := make(map[string]string)
	for _, m := range matches {
		key := strings.ToLower(strings.TrimSpace(m.NamaKelas))
		countByName[key]++
		kelasByName[key] = m.ID
	}

	ambiguousNames := make(map[string]bool)
	for key, count := range countByName {
		if count > 1 {
			delete(kelasByName, key)
			ambiguousNames[key] = true
		}
	}

	return kelasByName, ambiguousNames, nil
}

// validateKelasColumn menandai baris dengan kolom kelas kosong, tidak ditemukan, atau
// ambigu (dua kelas berbeda memakai nama yang sama).
func (s *pesertaService) validateKelasColumn(excelRows []*utils.ExcelPesertaRow, kelasByName map[string]string, ambiguousNames map[string]bool) []dto.ImportPesertaErrorDetail {
	var errorsFound []dto.ImportPesertaErrorDetail

	for _, row := range excelRows {
		trimmed := strings.TrimSpace(row.NamaKelas)
		name := strings.ToLower(trimmed)

		if trimmed == "" {
			errorsFound = append(errorsFound, dto.ImportPesertaErrorDetail{
				Row:   row.RowIndex,
				Error: "kolom kelas tidak boleh kosong",
			})
			continue
		}

		if ambiguousNames[name] {
			errorsFound = append(errorsFound, dto.ImportPesertaErrorDetail{
				Row:   row.RowIndex,
				Error: fmt.Sprintf("kelas '%s' ambigu (ditemukan lebih dari satu kelas dengan nama sama)", trimmed),
			})
			continue
		}

		if _, ok := kelasByName[name]; !ok {
			errorsFound = append(errorsFound, dto.ImportPesertaErrorDetail{
				Row:   row.RowIndex,
				Error: fmt.Sprintf("kelas '%s' tidak ditemukan", trimmed),
			})
		}
	}

	return errorsFound
}

// GenerateImportTemplate membuat file .xlsx berisi header kolom yang sama persis dengan
// urutan yang dibaca utils.ParseExcelPesertaRow saat import: nama, username, password, kelas.
// Nilai kolom "kelas" harus persis sama (tidak case-sensitive) dengan nama_kelas yang
// sudah ada di data master kelas — kalau tidak cocok, seluruh import akan dibatalkan.
func (s *pesertaService) GenerateImportTemplate() ([]byte, error) {
	xlsx := excelize.NewFile()
	defer xlsx.Close()

	sheet := "Template Peserta"
	xlsx.SetSheetName("Sheet1", sheet)

	headers := []string{"nama", "username", "password", "kelas"}
	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		xlsx.SetCellValue(sheet, cell, h)
	}

	xlsx.SetCellValue(sheet, "A2", "Budi Santoso")
	xlsx.SetCellValue(sheet, "B2", "budi01")
	xlsx.SetCellValue(sheet, "C2", "rahasia123")
	xlsx.SetCellValue(sheet, "D2", "X TKJ 1")

	if err := xlsx.AddComment(sheet, excelize.Comment{
		Cell:   "D1",
		Author: "System",
		Text:   "Isi dengan nama kelas persis seperti di data master Kelas. Jika tidak ditemukan, seluruh import akan dibatalkan.",
	}); err != nil {
		return nil, err
	}

	widths := []float64{30, 20, 20, 20}
	for col, w := range widths {
		colName, _ := excelize.ColumnNumberToName(col + 1)
		xlsx.SetColWidth(sheet, colName, colName, w)
	}

	var buf bytes.Buffer
	if err := xlsx.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
