package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/constants"
	"backend/internal/modules/soal/dto"
	"backend/internal/modules/soal/model"
	"backend/internal/modules/soal/repository"
	"backend/internal/storage"
	"backend/internal/utils"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type SoalService interface {
	CreateSoal(req *dto.CreateSoalRequest) (*dto.SoalResponse, error)
	GetSoalByID(id string) (*dto.SoalResponse, error)
	GetAllSoal(page, pageSize int) (*dto.SoalListResponse, error)
	GetSoalByBankSoal(bankSoalID string, page, pageSize int) (*dto.SoalListResponse, error)
	UpdateSoal(id string, req *dto.UpdateSoalRequest) (*dto.SoalResponse, error)
	DeleteSoal(id string) error
	RestoreSoal(id string) error
	ImportSoalFromExcel(ctx context.Context, req *dto.ImportSoalRequest) (*dto.ImportSoalResponse, error)
	RandomizeOpsiForSoal(soal *dto.SoalResponse, pesertaID string) *dto.SoalResponse
	GenerateImportTemplate() ([]byte, error)
}

type soalService struct {
	repo repository.SoalRepository
}

func NewSoalService(repo repository.SoalRepository) SoalService {
	return &soalService{repo: repo}
}

func (s *soalService) CreateSoal(req *dto.CreateSoalRequest) (*dto.SoalResponse, error) {
	if err := s.validateKunci(req.Kunci, req.OpsiA, req.OpsiB, req.OpsiC, req.OpsiD, req.OpsiE); err != nil {
		return nil, err
	}

	exists, err := s.repo.CheckDuplicate(req.IdBankSoal, req.NoSoal)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("soal dengan no_soal " + fmt.Sprintf("%d", req.NoSoal) + " sudah ada di bank soal ini")
	}

	var gambarSoalName string
	if req.GambarSoal != nil {
		filename, err := utils.SaveImage(req.GambarSoal, "soal")
		if err != nil {
			return nil, errors.New("gagal upload gambar_soal: " + err.Error())
		}
		gambarSoalName = filename
	}

	var gambarAName string
	if req.GambarA != nil {
		filename, err := utils.SaveImage(req.GambarA, "opsi")
		if err != nil {
			utils.DeleteImage(gambarSoalName, "soal")
			return nil, errors.New("gagal upload gambar_a: " + err.Error())
		}
		gambarAName = filename
	}

	var gambarBName string
	if req.GambarB != nil {
		filename, err := utils.SaveImage(req.GambarB, "opsi")
		if err != nil {
			utils.DeleteImage(gambarSoalName, "soal")
			utils.DeleteImage(gambarAName, "opsi")
			return nil, errors.New("gagal upload gambar_b: " + err.Error())
		}
		gambarBName = filename
	}

	var gambarCName string
	if req.GambarC != nil {
		filename, err := utils.SaveImage(req.GambarC, "opsi")
		if err != nil {
			utils.DeleteImage(gambarSoalName, "soal")
			utils.DeleteImage(gambarAName, "opsi")
			utils.DeleteImage(gambarBName, "opsi")
			return nil, errors.New("gagal upload gambar_c: " + err.Error())
		}
		gambarCName = filename
	}

	var gambarDName string
	if req.GambarD != nil {
		filename, err := utils.SaveImage(req.GambarD, "opsi")
		if err != nil {
			utils.DeleteImage(gambarSoalName, "soal")
			utils.DeleteImage(gambarAName, "opsi")
			utils.DeleteImage(gambarBName, "opsi")
			utils.DeleteImage(gambarCName, "opsi")
			return nil, errors.New("gagal upload gambar_d: " + err.Error())
		}
		gambarDName = filename
	}

	var gambarEName string
	if req.GambarE != nil {
		filename, err := utils.SaveImage(req.GambarE, "opsi")
		if err != nil {
			utils.DeleteImage(gambarSoalName, "soal")
			utils.DeleteImage(gambarAName, "opsi")
			utils.DeleteImage(gambarBName, "opsi")
			utils.DeleteImage(gambarCName, "opsi")
			utils.DeleteImage(gambarDName, "opsi")
			return nil, errors.New("gagal upload gambar_e: " + err.Error())
		}
		gambarEName = filename
	}

	soal := &model.Soal{
		IdBankSoal: req.IdBankSoal,
		NoSoal:     req.NoSoal,
		Soal:       req.Soal,
		GambarSoal: gambarSoalName,
		OpsiA:      req.OpsiA,
		OpsiB:      req.OpsiB,
		OpsiC:      req.OpsiC,
		OpsiD:      req.OpsiD,
		OpsiE:      req.OpsiE,
		GambarA:    gambarAName,
		GambarB:    gambarBName,
		GambarC:    gambarCName,
		GambarD:    gambarDName,
		GambarE:    gambarEName,
		Kunci:      req.Kunci,
	}

	if err := s.repo.Create(soal); err != nil {
		utils.DeleteImage(gambarSoalName, "soal")
		utils.DeleteImage(gambarAName, "opsi")
		utils.DeleteImage(gambarBName, "opsi")
		utils.DeleteImage(gambarCName, "opsi")
		utils.DeleteImage(gambarDName, "opsi")
		utils.DeleteImage(gambarEName, "opsi")
		return nil, err
	}

	return s.modelToResponse(soal), nil
}

func (s *soalService) GetSoalByID(id string) (*dto.SoalResponse, error) {
	soal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	return s.modelToResponse(soal), nil
}

func (s *soalService) GetAllSoal(page, pageSize int) (*dto.SoalListResponse, error) {
	soals, total, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	responses := []dto.SoalResponse{}
	for _, soal := range soals {
		responses = append(responses, *s.modelToResponse(&soal))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.SoalListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *soalService) GetSoalByBankSoal(bankSoalID string, page, pageSize int) (*dto.SoalListResponse, error) {
	soals, total, err := s.repo.GetByBankSoalID(bankSoalID, page, pageSize)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	responses := []dto.SoalResponse{}
	for _, soal := range soals {
		responses = append(responses, *s.modelToResponse(&soal))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.SoalListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *soalService) UpdateSoal(id string, req *dto.UpdateSoalRequest) (*dto.SoalResponse, error) {
	soal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	if err := s.validateKunci(req.Kunci, req.OpsiA, req.OpsiB, req.OpsiC, req.OpsiD, req.OpsiE); err != nil {
		return nil, err
	}

	if req.GambarSoal != nil {
		newFilename, err := utils.SaveImage(req.GambarSoal, "soal")
		if err != nil {
			return nil, errors.New("gagal upload gambar_soal: " + err.Error())
		}
		if soal.GambarSoal != "" {
			utils.DeleteImage(soal.GambarSoal, "soal")
		}
		soal.GambarSoal = newFilename
	}

	if req.GambarA != nil {
		newFilename, err := utils.SaveImage(req.GambarA, "opsi")
		if err != nil {
			return nil, errors.New("gagal upload gambar_a: " + err.Error())
		}
		if soal.GambarA != "" {
			utils.DeleteImage(soal.GambarA, "opsi")
		}
		soal.GambarA = newFilename
	}

	if req.GambarB != nil {
		newFilename, err := utils.SaveImage(req.GambarB, "opsi")
		if err != nil {
			return nil, errors.New("gagal upload gambar_b: " + err.Error())
		}
		if soal.GambarB != "" {
			utils.DeleteImage(soal.GambarB, "opsi")
		}
		soal.GambarB = newFilename
	}

	if req.GambarC != nil {
		newFilename, err := utils.SaveImage(req.GambarC, "opsi")
		if err != nil {
			return nil, errors.New("gagal upload gambar_c: " + err.Error())
		}
		if soal.GambarC != "" {
			utils.DeleteImage(soal.GambarC, "opsi")
		}
		soal.GambarC = newFilename
	}

	if req.GambarD != nil {
		newFilename, err := utils.SaveImage(req.GambarD, "opsi")
		if err != nil {
			return nil, errors.New("gagal upload gambar_d: " + err.Error())
		}
		if soal.GambarD != "" {
			utils.DeleteImage(soal.GambarD, "opsi")
		}
		soal.GambarD = newFilename
	}

	if req.GambarE != nil {
		newFilename, err := utils.SaveImage(req.GambarE, "opsi")
		if err != nil {
			return nil, errors.New("gagal upload gambar_e: " + err.Error())
		}
		if soal.GambarE != "" {
			utils.DeleteImage(soal.GambarE, "opsi")
		}
		soal.GambarE = newFilename
	}

	soal.NoSoal = req.NoSoal
	soal.Soal = req.Soal
	soal.OpsiA = req.OpsiA
	soal.OpsiB = req.OpsiB
	soal.OpsiC = req.OpsiC
	soal.OpsiD = req.OpsiD
	soal.OpsiE = req.OpsiE
	soal.Kunci = req.Kunci

	if err := s.repo.Update(soal); err != nil {
		return nil, err
	}

	return s.modelToResponse(soal), nil
}

func (s *soalService) DeleteSoal(id string) error {
	soal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}

	utils.DeleteImage(soal.GambarSoal, "soal")
	utils.DeleteImage(soal.GambarA, "opsi")
	utils.DeleteImage(soal.GambarB, "opsi")
	utils.DeleteImage(soal.GambarC, "opsi")
	utils.DeleteImage(soal.GambarD, "opsi")
	utils.DeleteImage(soal.GambarE, "opsi")

	return s.repo.Delete(soal.ID)
}

func (s *soalService) RestoreSoal(id string) error {
	return s.repo.Restore(id)
}

func (s *soalService) RandomizeOpsiForSoal(soal *dto.SoalResponse, pesertaID string) *dto.SoalResponse {
	seed := s.generateSeed(pesertaID, soal.ID)
	rng := rand.New(rand.NewSource(seed))

	opsiMap := map[string]string{
		"A": soal.OpsiA,
		"B": soal.OpsiB,
		"C": soal.OpsiC,
		"D": soal.OpsiD,
		"E": soal.OpsiE,
	}
	gambarMap := map[string]string{
		"A": soal.GambarA,
		"B": soal.GambarB,
		"C": soal.GambarC,
		"D": soal.GambarD,
		"E": soal.GambarE,
	}

	opsiOrder := []string{"A", "B", "C", "D", "E"}
	rng.Shuffle(len(opsiOrder), func(i, j int) {
		opsiOrder[i], opsiOrder[j] = opsiOrder[j], opsiOrder[i]
	})

	newKunciPos := 0
	for i, key := range opsiOrder {
		if key == soal.Kunci {
			newKunciPos = i
			break
		}
	}

	newKunci := string(rune('A' + newKunciPos))

	soal.OpsiA = opsiMap[opsiOrder[0]]
	soal.OpsiB = opsiMap[opsiOrder[1]]
	soal.OpsiC = opsiMap[opsiOrder[2]]
	soal.OpsiD = opsiMap[opsiOrder[3]]
	soal.OpsiE = opsiMap[opsiOrder[4]]
	soal.GambarA = gambarMap[opsiOrder[0]]
	soal.GambarB = gambarMap[opsiOrder[1]]
	soal.GambarC = gambarMap[opsiOrder[2]]
	soal.GambarD = gambarMap[opsiOrder[3]]
	soal.GambarE = gambarMap[opsiOrder[4]]
	soal.Kunci = newKunci

	return soal
}

func (s *soalService) generateSeed(pesertaID, soalID string) int64 {
	hash := md5.Sum([]byte(pesertaID + "|" + soalID))
	seed := int64(binary.BigEndian.Uint64(hash[:8]))
	return seed
}

func (s *soalService) validateKunci(kunci string, opsiA, opsiB, opsiC, opsiD, opsiE string) error {
	validKeys := map[string]bool{
		"A": true, "B": true, "C": true, "D": true, "E": true,
	}
	if !validKeys[kunci] {
		return errors.New("kunci harus A, B, C, D, atau E")
	}
	return nil
}

func (s *soalService) modelToResponse(soal *model.Soal) *dto.SoalResponse {
	return &dto.SoalResponse{
		ID:         soal.ID,
		IdBankSoal: soal.IdBankSoal,
		NoSoal:     soal.NoSoal,
		Soal:       soal.Soal,
		GambarSoal: s.buildImageURL(soal.GambarSoal, "soal"),
		OpsiA:      soal.OpsiA,
		OpsiB:      soal.OpsiB,
		OpsiC:      soal.OpsiC,
		OpsiD:      soal.OpsiD,
		OpsiE:      soal.OpsiE,
		GambarA:    s.buildImageURL(soal.GambarA, "opsi"),
		GambarB:    s.buildImageURL(soal.GambarB, "opsi"),
		GambarC:    s.buildImageURL(soal.GambarC, "opsi"),
		GambarD:    s.buildImageURL(soal.GambarD, "opsi"),
		GambarE:    s.buildImageURL(soal.GambarE, "opsi"),
		Kunci:      soal.Kunci,
		CreatedAt:  soal.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  soal.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *soalService) downloadGambar(ctx context.Context, sourceBase, filename, folder string) string {
	if filename == "" {
		return ""
	}
	stored, err := storage.DownloadAndUploadFromURL(ctx, sourceBase+filename, folder)
	if err != nil {
		return ""
	}
	return stored
}

func (s *soalService) buildImageURL(filename, folder string) string {
	if filename == "" {
		return ""
	}
	cfg := config.GetUploadConfig()
	return fmt.Sprintf("%s/%s/%s", cfg.ImageBaseURL, folder, filename)
}

func (s *soalService) ImportSoalFromExcel(ctx context.Context, req *dto.ImportSoalRequest) (*dto.ImportSoalResponse, error) {
	// 1. Validasi bank_soal exists
	exists, err := s.repo.GetBankSoalExists(ctx, req.IdBankSoal)
	if err != nil || !exists {
		return nil, errors.New("bank_soal tidak ditemukan")
	}

	// 2. Buka file dari request
	file, err := req.File.Open()
	if err != nil {
		return nil, errors.New("gagal membuka file")
	}
	defer file.Close()

	// 3. Parse excel file
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, errors.New("file bukan format excel yang valid")
	}
	defer xlsx.Close()

	// 4. Get sheet pertama
	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, errors.New("gagal membaca sheet excel")
	}

	// 5. Parse dan validasi data
	var soals []model.Soal
	var errorDetails []dto.ImportSoalErrorDetail
	var processedCount, successCount, failedCount int

	// Skip header row (start dari index 1)
	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		row := rows[rowIndex]
		processedCount++

		// Parse row
		excelRow, err := utils.ParseExcelRow(row, rowIndex+1)
		if err != nil {
			failedCount++
			errorDetails = append(errorDetails, dto.ImportSoalErrorDetail{
				Row:   rowIndex + 1,
				Error: err.Error(),
			})
			continue
		}

		// Validate row
		validationErrors := utils.ValidateSoalRow(excelRow)
		if len(validationErrors) > 0 {
			failedCount++
			errorDetails = append(errorDetails, dto.ImportSoalErrorDetail{
				Row:   rowIndex + 1,
				Error: strings.Join(validationErrors, "; "),
			})
			continue
		}

		// Upload gambar dari source URL ke MinIO jika ada
		sourceBase := config.GetSoalImageSourceURL()
		gambarSoal := s.downloadGambar(ctx, sourceBase, excelRow.GambarSoal, "soal")
		gambarA    := s.downloadGambar(ctx, sourceBase, excelRow.GambarA,    "opsi")
		gambarB    := s.downloadGambar(ctx, sourceBase, excelRow.GambarB,    "opsi")
		gambarC    := s.downloadGambar(ctx, sourceBase, excelRow.GambarC,    "opsi")
		gambarD    := s.downloadGambar(ctx, sourceBase, excelRow.GambarD,    "opsi")
		gambarE    := s.downloadGambar(ctx, sourceBase, excelRow.GambarE,    "opsi")

		// Create soal model
		soal := model.Soal{
			IdBankSoal: req.IdBankSoal,
			NoSoal:     excelRow.NoSoal,
			Soal:       excelRow.Soal,
			OpsiA:      excelRow.OpsiA,
			OpsiB:      excelRow.OpsiB,
			OpsiC:      excelRow.OpsiC,
			OpsiD:      excelRow.OpsiD,
			OpsiE:      excelRow.OpsiE,
			Kunci:      excelRow.Kunci,
			GambarSoal: gambarSoal,
			GambarA:    gambarA,
			GambarB:    gambarB,
			GambarC:    gambarC,
			GambarD:    gambarD,
			GambarE:    gambarE,
		}

		soals = append(soals, soal)
		successCount++
	}

	// 6. Bulk insert ke database
	if len(soals) > 0 {
		err = s.repo.BulkCreateSoal(ctx, soals)
		if err != nil {
			return nil, errors.New("gagal menyimpan data ke database: " + err.Error())
		}
	}

	// Limit error details ke max 100
	if len(errorDetails) > 100 {
		errorDetails = errorDetails[:100]
	}

	// 7. Return response
	return &dto.ImportSoalResponse{
		TotalProcessed: processedCount,
		TotalSuccess:   successCount,
		TotalFailed:    failedCount,
		ImportID:       req.IdBankSoal,
		Timestamp:      time.Now(),
		Summary: map[string]int{
			"inserted": successCount,
			"skipped":  0,
			"errors":   failedCount,
		},
		Errors: errorDetails,
	}, nil
}

// GenerateImportTemplate membuat file .xlsx berisi header kolom yang sama persis dengan
// urutan yang dibaca utils.ParseExcelRow saat import (kolom I dan K sengaja dikosongkan/
// tidak dipakai, tapi tetap disediakan agar kolom sesudahnya tidak bergeser index).
func (s *soalService) GenerateImportTemplate() ([]byte, error) {
	xlsx := excelize.NewFile()
	defer xlsx.Close()

	sheet := "Template Soal"
	xlsx.SetSheetName("Sheet1", sheet)

	headers := []string{
		"no_soal", "soal", "opsi_a", "opsi_b", "opsi_c", "opsi_d", "opsi_e", "kunci",
		"(tidak digunakan)", "gambar_soal", "(tidak digunakan)",
		"gambar_a", "gambar_b", "gambar_c", "gambar_d", "gambar_e",
	}
	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		xlsx.SetCellValue(sheet, cell, h)
	}

	widths := []float64{8, 40, 25, 25, 25, 25, 25, 8, 18, 20, 18, 18, 18, 18, 18, 18}
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
