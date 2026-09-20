package service

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"backend/configs"
	"backend/internal/constants"
	"backend/internal/utils"
	jawabanmodel "backend/internal/modules/jawaban/model"
	jawabanrepo "backend/internal/modules/jawaban/repository"
	jadwalmodel "backend/internal/modules/jadwal/model"
	"backend/internal/modules/nilai/dto"
	"backend/internal/modules/nilai/model"
	"backend/internal/modules/nilai/repository"
	sectionmodel "backend/internal/modules/section/model"
	sectionrepo "backend/internal/modules/section/repository"
	soalmodel "backend/internal/modules/soal/model"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

var jakartaLoc, _ = time.LoadLocation("Asia/Jakarta")

func init() {
	if jakartaLoc == nil {
		jakartaLoc = time.FixedZone("WIB", 7*60*60)
	}
}

// wibWallClock mengoreksi time.Time yang dibaca dari kolom "timestamp" (tanpa timezone) di
// Postgres: driver membaca digit jam-nya lalu menandainya sebagai UTC, padahal digit itu
// sebenarnya wall-clock WIB (ditulis via time.Now().In(jakartaLoc), lihat MulaiUjian/AdvanceSection).
// Fungsi ini menafsirkan ulang digit yang sama sebagai WIB agar instant absolutnya benar
// sebelum dibandingkan dengan time.Now().
func wibWallClock(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), jakartaLoc)
}

type ExportResult struct {
	ZipBytes  []byte
	NamaUjian string
}

type NilaiService interface {
	CreateNilai(req *dto.CreateNilaiRequest) (*dto.NilaiResponse, error)
	GetNilaiByID(id string) (*dto.NilaiResponse, error)
	GetAllNilai(page, pageSize int, idPeserta, idJadwal string) (*dto.NilaiListResponse, error)
	GetNilaiByPeserta(idPeserta string, page, pageSize int) (*dto.NilaiListResponse, error)
	GetNilaiByJadwal(idJadwal string, page, pageSize int) (*dto.NilaiListResponse, error)
	UpdateNilai(id string, req *dto.UpdateNilaiRequest) (*dto.NilaiResponse, error)
	DeleteNilai(id string) error
	RestoreNilai(id string) error
	MulaiUjian(idPeserta, idJadwal, token string) (*dto.NilaiResponse, bool, error)
	ExportNilaiByJadwal(idJadwal string) (*ExportResult, error)
	AnalisisSoalByJadwal(idJadwal string) (*ExportResult, error)
	GetMonitoringByJadwal(idJadwal, idKelas string) (*dto.MonitoringResponse, error)
	ForceSelesaikanUjian(id string) (*dto.NilaiResponse, error)
	AdvanceSection(idNilai, idPeserta string) (*dto.SectionProgressResponse, error)
	GetSectionStatus(idNilai, idPeserta string) (*dto.SectionProgressResponse, error)
}

type nilaiService struct {
	repo        repository.NilaiRepository
	jawabanRepo jawabanrepo.JawabanRepository
	sectionRepo sectionrepo.SectionRepository
	db          *gorm.DB
}

func NewNilaiService(repo repository.NilaiRepository, jawabanRepo jawabanrepo.JawabanRepository, sectionRepo sectionrepo.SectionRepository, db *gorm.DB) NilaiService {
	return &nilaiService{
		repo:        repo,
		jawabanRepo: jawabanRepo,
		sectionRepo: sectionRepo,
		db:          db,
	}
}

const timeLayout = "2006-01-02 15:04:05"

func parseTime(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse(timeLayout, *s)
	if err != nil {
		return nil, errors.New("format waktu tidak valid, gunakan: 2006-01-02 15:04:05")
	}
	return &t, nil
}

func (s *nilaiService) CreateNilai(req *dto.CreateNilaiRequest) (*dto.NilaiResponse, error) {
	if req.Nilai < 0 || req.Nilai > 100 {
		return nil, errors.New("nilai harus di antara 0 dan 100")
	}

	exists, err := s.repo.CheckDuplicate(req.IDPeserta, req.IDJadwal)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("nilai untuk peserta dan jadwal ini sudah ada — gunakan endpoint update")
	}

	wktMulai, err := parseTime(req.WktMulai)
	if err != nil {
		return nil, err
	}
	aktivitasTerakhir, err := parseTime(req.AktivitasTerakhir)
	if err != nil {
		return nil, err
	}
	wktSelesai, err := parseTime(req.WktSelesai)
	if err != nil {
		return nil, err
	}

	nilai := &model.Nilai{
		IDPeserta:         req.IDPeserta,
		IDJadwal:          req.IDJadwal,
		Nilai:             req.Nilai,
		WktMulai:          wktMulai,
		AktivitasTerakhir: aktivitasTerakhir,
		WktSelesai:        wktSelesai,
	}

	if err := s.repo.Create(nilai); err != nil {
		return nil, err
	}

	created, err := s.repo.GetByIDWithDetail(nilai.ID)
	if err != nil {
		return nil, err
	}
	return detailToResponse(created), nil
}

func (s *nilaiService) GetNilaiByID(id string) (*dto.NilaiResponse, error) {
	result, err := s.repo.GetByIDWithDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	return detailToResponse(result), nil
}

func (s *nilaiService) GetAllNilai(page, pageSize int, idPeserta, idJadwal string) (*dto.NilaiListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	results, total, err := s.repo.GetAllWithDetail(page, pageSize, idPeserta, idJadwal)
	if err != nil {
		return nil, err
	}

	responses := []dto.NilaiResponse{}
	for _, r := range results {
		responses = append(responses, *detailToResponse(&r))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.NilaiListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *nilaiService) GetNilaiByPeserta(idPeserta string, page, pageSize int) (*dto.NilaiListResponse, error) {
	return s.GetAllNilai(page, pageSize, idPeserta, "")
}

func (s *nilaiService) GetNilaiByJadwal(idJadwal string, page, pageSize int) (*dto.NilaiListResponse, error) {
	return s.GetAllNilai(page, pageSize, "", idJadwal)
}

func (s *nilaiService) UpdateNilai(id string, req *dto.UpdateNilaiRequest) (*dto.NilaiResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	if req.Nilai != nil {
		if *req.Nilai < 0 || *req.Nilai > 100 {
			return nil, errors.New("nilai harus di antara 0 dan 100")
		}
		existing.Nilai = *req.Nilai
	}

	if req.IDPeserta != nil || req.IDJadwal != nil {
		newPeserta := existing.IDPeserta
		newJadwal  := existing.IDJadwal
		if req.IDPeserta != nil {
			newPeserta = *req.IDPeserta
		}
		if req.IDJadwal != nil {
			newJadwal = *req.IDJadwal
		}
		if newPeserta != existing.IDPeserta || newJadwal != existing.IDJadwal {
			exists, err := s.repo.CheckDuplicate(newPeserta, newJadwal)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("nilai untuk peserta dan jadwal ini sudah ada")
			}
		}
		existing.IDPeserta = newPeserta
		existing.IDJadwal  = newJadwal
	}

	if req.WktMulai != nil {
		t, err := parseTime(req.WktMulai)
		if err != nil {
			return nil, err
		}
		existing.WktMulai = t
	}

	if req.AktivitasTerakhir != nil {
		t, err := parseTime(req.AktivitasTerakhir)
		if err != nil {
			return nil, err
		}
		existing.AktivitasTerakhir = t
	}

	if req.WktSelesai != nil {
		t, err := parseTime(req.WktSelesai)
		if err != nil {
			return nil, err
		}
		existing.WktSelesai = t

		if t != nil {
			nilai, err := s.repo.HitungNilai(id)
			if err != nil {
				return nil, err
			}
			existing.Nilai = nilai
		}
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByIDWithDetail(id)
	if err != nil {
		return nil, err
	}
	return detailToResponse(updated), nil
}

// ForceSelesaikanUjian dipanggil admin untuk memaksa selesaikan sesi ujian peserta
// (mis. peserta lupa submit atau koneksinya terputus). Nilai dihitung ulang dari
// jawaban yang sudah sempat diisi lewat repo.HitungNilai — mekanisme yang sama
// persis dipakai saat peserta submit sendiri via UpdateNilai di atas.
func (s *nilaiService) ForceSelesaikanUjian(id string) (*dto.NilaiResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	if existing.WktSelesai != nil {
		return nil, errors.New("ujian peserta ini sudah selesai")
	}

	now := time.Now().In(jakartaLoc)
	existing.WktSelesai = &now
	existing.AktivitasTerakhir = &now

	nilai, err := s.repo.HitungNilai(id)
	if err != nil {
		return nil, err
	}
	existing.Nilai = nilai

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByIDWithDetail(id)
	if err != nil {
		return nil, err
	}
	return detailToResponse(updated), nil
}

func (s *nilaiService) DeleteNilai(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.jawabanRepo.SoftDeleteByNilaiID(tx, id); err != nil {
			return err
		}
		now := time.Now()
		return tx.Model(&model.Nilai{}).Where("id = ?", id).Update("deleted_at", now).Error
	})
}

func (s *nilaiService) RestoreNilai(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.jawabanRepo.RestoreByNilaiID(tx, id); err != nil {
			return err
		}
		return tx.Model(&model.Nilai{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NULL")).Error
	})
}

func (s *nilaiService) MulaiUjian(idPeserta, idJadwal, token string) (*dto.NilaiResponse, bool, error) {
	// 1. Cek apakah record sudah ada
	existing, err := s.repo.GetByPesertaAndJadwal(idPeserta, idJadwal)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	// 2. Jika sudah ada → cek wkt_selesai
	if existing != nil {
		if existing.WktSelesai != nil {
			return nil, false, errors.New("Ujian sudah pernah dilakukan")
		}
		// Resume: ambil detail (dengan JOIN) lalu return
		detail, err := s.repo.GetByIDWithDetail(existing.ID)
		if err != nil {
			return nil, false, err
		}
		return detailToResponse(detail), false, nil
	}

	// 3. Belum ada → transaction: insert nilai + bulk insert jawaban
	var newNilaiID string
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 3a. Get jadwal untuk dapatkan id_bank_soal.
		// SELECT eksplisit + cast ke int agar acak_soal/acak_opsi terbaca
		// walau kolom DB masih boolean.
		var jadwal jadwalmodel.Jadwal
		if err := tx.Table("jadwal").
			Select("id, id_bank_soal, nama_ujian, tingkat, wkt_mulai, wkt_selesai, durasi, acak_soal::int AS acak_soal, acak_opsi::int AS acak_opsi, wajib_token::int AS wajib_token, created_at, updated_at, deleted_at").
			Where("id = ? AND deleted_at IS NULL", idJadwal).
			First(&jadwal).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("jadwal tidak ditemukan")
			}
			return err
		}

		// 3a-bis. Jika jadwal mewajibkan token, validasi token TOTP global yang berlaku saat ini
		// (berubah tiap UJIAN_TOKEN_PERIOD, dipakai bersama oleh semua jadwal yang mewajibkannya).
		if jadwal.WajibToken == 1 {
			tokenCfg := configs.GetUjianTokenConfig()
			if token == "" || !utils.ValidateTOTP(tokenCfg.Secret, token, time.Now(), tokenCfg.Period, tokenCfg.Digits, tokenCfg.GraceSteps) {
				return errors.New("token ujian tidak valid atau sudah kedaluwarsa")
			}
		}

		// 3b. Insert nilai baru
		now := time.Now().In(jakartaLoc)
		nilai := &model.Nilai{
			IDPeserta:         idPeserta,
			IDJadwal:          idJadwal,
			Nilai:             0,
			WktMulai:          &now,
			AktivitasTerakhir: &now,
			WktSelesai:        nil,
		}
		if err := tx.Create(nilai).Error; err != nil {
			return err
		}
		newNilaiID = nilai.ID

		// 3c-bis. Jika jadwal ini memakai section, aktifkan section pertama sebagai frontier awal.
		var sections []sectionmodel.Section
		if err := tx.Where("id_jadwal = ? AND deleted_at IS NULL", idJadwal).
			Order("urutan ASC").
			Find(&sections).Error; err != nil {
			return err
		}
		if len(sections) > 0 {
			firstSectionID := sections[0].ID
			nilai.IDSectionAktif = &firstSectionID
			nilai.WktMulaiSection = &now
			if err := tx.Model(&model.Nilai{}).Where("id = ?", nilai.ID).
				Updates(map[string]interface{}{
					"id_section_aktif":  firstSectionID,
					"wkt_mulai_section": now,
				}).Error; err != nil {
				return err
			}
		}

		// 3c. Query soal by bank_soal — acak jika acak_soal=1, urut jika 0
		var soals []soalmodel.Soal
		soalQuery := tx.Where("id_bank_soal = ? AND deleted_at IS NULL", jadwal.IDBankSoal)
		if jadwal.AcakSoal == 1 {
			soalQuery = soalQuery.Order("RANDOM()")
		} else {
			soalQuery = soalQuery.Order("no_soal ASC")
		}
		if err := soalQuery.Find(&soals).Error; err != nil {
			return err
		}

		// 3d. Build & bulk insert jawaban kosong
		if len(soals) > 0 {
			jawabans := make([]jawabanmodel.Jawaban, len(soals))
			if jadwal.AcakSoal == 1 {
				noUrutSequence := rand.Perm(len(soals))
				for i, soal := range soals {
					jawabans[i] = jawabanmodel.Jawaban{
						IDNilai:   nilai.ID,
						IDSoal:    soal.ID,
						IDPeserta: idPeserta,
						NoUrut:    noUrutSequence[i] + 1,
						Jawaban:   nil,
						IsBenar:   nil,
					}
				}
			} else {
				for i, soal := range soals {
					jawabans[i] = jawabanmodel.Jawaban{
						IDNilai:   nilai.ID,
						IDSoal:    soal.ID,
						IDPeserta: idPeserta,
						NoUrut:    i + 1,
						Jawaban:   nil,
						IsBenar:   nil,
					}
				}
			}
			if err := s.jawabanRepo.BulkCreateWithTx(tx, jawabans); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, false, err
	}

	// 4. Ambil detail nilai yang baru di-insert (di luar transaction, read)
	created, err := s.repo.GetByIDWithDetail(newNilaiID)
	if err != nil {
		return nil, false, err
	}
	return detailToResponse(created), true, nil
}

func (s *nilaiService) loadNilaiForOwner(idNilai, idPeserta string) (*model.Nilai, error) {
	nilai, err := s.repo.GetByID(idNilai)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	if nilai.IDPeserta != idPeserta {
		return nil, errors.New("Anda tidak memiliki akses ke sesi ujian ini")
	}
	return nilai, nil
}

// buildSectionProgress menghitung status section aktif (frontier) sesi ujian ini,
// tanpa mengubah state apa pun.
func (s *nilaiService) buildSectionProgress(nilai *model.Nilai) (*dto.SectionProgressResponse, *sectionmodel.Section, error) {
	if nilai.IDSectionAktif == nil {
		return nil, nil, errors.New("ujian ini tidak menggunakan section")
	}

	current, err := s.sectionRepo.GetByID(*nilai.IDSectionAktif)
	if err != nil {
		return nil, nil, err
	}

	elapsedDetik := 0
	if nilai.WktMulaiSection != nil {
		elapsedDetik = int(time.Now().Sub(wibWallClock(*nilai.WktMulaiSection)).Seconds())
	}
	minDetik := current.DurasiMenitMinimal * 60
	sisaDetik := minDetik - elapsedDetik
	if sisaDetik < 0 {
		sisaDetik = 0
	}

	sections, err := s.sectionRepo.GetByJadwalID(nilai.IDJadwal)
	if err != nil {
		return nil, nil, err
	}
	isTerakhir := true
	for _, sec := range sections {
		if sec.Urutan == current.Urutan+1 {
			isTerakhir = false
			break
		}
	}

	return &dto.SectionProgressResponse{
		IDNilai:            nilai.ID,
		IDSection:          current.ID,
		NamaSection:        current.NamaSection,
		Urutan:             current.Urutan,
		DurasiMenitMinimal: current.DurasiMenitMinimal,
		BolehLanjut:        elapsedDetik >= minDetik,
		SisaDetik:          sisaDetik,
		IsSectionTerakhir:  isTerakhir,
	}, current, nil
}

func (s *nilaiService) GetSectionStatus(idNilai, idPeserta string) (*dto.SectionProgressResponse, error) {
	nilai, err := s.loadNilaiForOwner(idNilai, idPeserta)
	if err != nil {
		return nil, err
	}
	resp, _, err := s.buildSectionProgress(nilai)
	return resp, err
}

func (s *nilaiService) AdvanceSection(idNilai, idPeserta string) (*dto.SectionProgressResponse, error) {
	nilai, err := s.loadNilaiForOwner(idNilai, idPeserta)
	if err != nil {
		return nil, err
	}
	if nilai.WktSelesai != nil {
		return nil, errors.New("ujian sudah selesai")
	}

	progress, current, err := s.buildSectionProgress(nilai)
	if err != nil {
		return nil, err
	}
	if progress.IsSectionTerakhir {
		return nil, errors.New("sudah berada di section terakhir")
	}
	if !progress.BolehLanjut {
		return progress, nil
	}

	sections, err := s.sectionRepo.GetByJadwalID(nilai.IDJadwal)
	if err != nil {
		return nil, err
	}
	var next *sectionmodel.Section
	for i := range sections {
		if sections[i].Urutan == current.Urutan+1 {
			next = &sections[i]
			break
		}
	}
	if next == nil {
		return nil, errors.New("sudah berada di section terakhir")
	}

	now := time.Now().In(jakartaLoc)
	nilai.IDSectionAktif = &next.ID
	nilai.WktMulaiSection = &now
	if err := s.repo.Update(nilai); err != nil {
		return nil, err
	}

	return &dto.SectionProgressResponse{
		IDNilai:            nilai.ID,
		IDSection:          next.ID,
		NamaSection:        next.NamaSection,
		Urutan:             next.Urutan,
		DurasiMenitMinimal: next.DurasiMenitMinimal,
		SudahLanjut:        true,
		BolehLanjut:        true,
		SisaDetik:          0,
		IsSectionTerakhir:  false,
	}, nil
}

func detailToResponse(r *repository.NilaiWithDetail) *dto.NilaiResponse {
	return &dto.NilaiResponse{
		ID:                r.ID,
		IDPeserta:         r.IDPeserta,
		NamaPeserta:       r.NamaPeserta,
		IDJadwal:          r.IDJadwal,
		NamaUjian:         r.NamaUjian,
		Nilai:             r.Nilai,
		WktMulai:          r.WktMulai,
		AktivitasTerakhir: r.AktivitasTerakhir,
		WktSelesai:        r.WktSelesai,
		IDSectionAktif:    r.IDSectionAktif,
		WktMulaiSection:   r.WktMulaiSection,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}

func (s *nilaiService) ExportNilaiByJadwal(idJadwal string) (*ExportResult, error) {
	// 1. Ambil nama ujian dari tabel jadwal
	var namaUjian string
	if err := s.db.Table("jadwal").
		Select("nama_ujian").
		Where("id = ? AND deleted_at IS NULL", idJadwal).
		Scan(&namaUjian).Error; err != nil || namaUjian == "" {
		return nil, errors.New("jadwal tidak ditemukan")
	}

	// 2. Ambil nama mapel dari jadwal
	type JadwalDetail struct {
		IDBankSoal string `gorm:"column:id_bank_soal"`
	}
	var jadwalDetail JadwalDetail
	if err := s.db.Table("jadwal").
		Select("id_bank_soal").
		Where("id = ? AND deleted_at IS NULL", idJadwal).
		Scan(&jadwalDetail).Error; err != nil {
		return nil, err
	}

	var namaMapel string
	if err := s.db.Table("bank_soal").
		Select("mapel.nama_mapel").
		Joins("INNER JOIN mapel ON bank_soal.id_mapel = mapel.id").
		Where("bank_soal.id = ? AND bank_soal.deleted_at IS NULL", jadwalDetail.IDBankSoal).
		Scan(&namaMapel).Error; err != nil || namaMapel == "" {
		return nil, errors.New("mapel tidak ditemukan untuk jadwal ini")
	}

	// 3. Ambil semua kelas yang terdaftar pada jadwal ini
	type KelasRow struct {
		IDKelas   string `gorm:"column:id_kelas"`
		NamaKelas string `gorm:"column:nama_kelas"`
	}
	var kelasList []KelasRow
	if err := s.db.Table("jadwal_kelas").
		Select("jadwal_kelas.id_kelas, kelas.nama_kelas").
		Joins("INNER JOIN kelas ON jadwal_kelas.id_kelas = kelas.id").
		Where("jadwal_kelas.id_jadwal = ?", idJadwal).
		Scan(&kelasList).Error; err != nil {
		return nil, err
	}
	if len(kelasList) == 0 {
		return nil, errors.New("tidak ada kelas yang terdaftar pada jadwal ini")
	}

	// 4. Buat ZIP di memory
	var zipBuf bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuf)

	for _, kelas := range kelasList {
		// 4a. Ambil data nilai peserta untuk kelas ini
		rows, err := s.repo.GetByJadwalAndKelas(idJadwal, kelas.IDKelas)
		if err != nil {
			return nil, fmt.Errorf("gagal ambil data kelas %s: %w", kelas.NamaKelas, err)
		}

		// 4b. Buat file Excel
		xlsx := excelize.NewFile()
		sheet := "Nilai Ujian"
		xlsx.SetSheetName("Sheet1", sheet)

		// Header
		headers := []string{"No", "Nama Peserta", "Username", "Nilai", "Waktu Mulai", "Waktu Selesai"}
		for col, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			xlsx.SetCellValue(sheet, cell, h)
		}

		// Data rows
		for i, row := range rows {
			rowNum := i + 2
			xlsx.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), i+1)
			xlsx.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), row.NamaPeserta)
			xlsx.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), row.Username)
			xlsx.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), row.Nilai)

			wktMulai := "-"
			if row.WktMulai != nil {
				wktMulai = *row.WktMulai
			}
			xlsx.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), wktMulai)

			wktSelesai := "-"
			if row.WktSelesai != nil {
				wktSelesai = *row.WktSelesai
			}
			xlsx.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), wktSelesai)
		}

		// 4c. Tulis Excel ke buffer lalu masukkan ke ZIP
		var xlsBuf bytes.Buffer
		if err := xlsx.Write(&xlsBuf); err != nil {
			return nil, fmt.Errorf("gagal tulis excel kelas %s: %w", kelas.NamaKelas, err)
		}

		// Nama file: {MAPEL}_{KELAS}.xlsx (spasi diganti underscore)
		safeMapel := strings.ReplaceAll(namaMapel, " ", "_")
		safeKelas := strings.ReplaceAll(kelas.NamaKelas, " ", "_")
		filename := fmt.Sprintf("%s_%s.xlsx", safeMapel, safeKelas)

		zipEntry, err := zipWriter.Create(filename)
		if err != nil {
			return nil, err
		}
		if _, err := zipEntry.Write(xlsBuf.Bytes()); err != nil {
			return nil, err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return &ExportResult{
		ZipBytes:  zipBuf.Bytes(),
		NamaUjian: namaUjian,
	}, nil
}

// GetMonitoringByJadwal mengembalikan status pengerjaan semua peserta yang terdaftar
// pada jadwal ini (termasuk yang belum pernah mulai ujian sama sekali — beda dengan
// GetNilaiByJadwal yang hanya baca dari tabel nilai). idKelas kosong berarti semua
// kelas yang terdaftar di jadwal ini digabung jadi satu list.
func (s *nilaiService) GetMonitoringByJadwal(idJadwal, idKelas string) (*dto.MonitoringResponse, error) {
	var namaUjian string
	if err := s.db.Table("jadwal").
		Select("nama_ujian").
		Where("id = ? AND deleted_at IS NULL", idJadwal).
		Scan(&namaUjian).Error; err != nil || namaUjian == "" {
		return nil, errors.New("jadwal tidak ditemukan")
	}

	rows, err := s.repo.GetMonitoringByJadwal(idJadwal, idKelas)
	if err != nil {
		return nil, err
	}

	summary := dto.MonitoringSummary{}
	data := make([]dto.MonitoringPesertaResponse, 0, len(rows))

	for _, r := range rows {
		status := "belum_mulai"
		switch {
		case r.WktSelesai != nil:
			status = "selesai"
			summary.Selesai++
		case r.WktMulai != nil:
			status = "sedang_mengerjakan"
			summary.SedangMengerjakan++
		default:
			summary.BelumMulai++
		}
		summary.Total++

		data = append(data, dto.MonitoringPesertaResponse{
			IDPeserta:         r.IDPeserta,
			NamaPeserta:       r.NamaPeserta,
			Username:          r.Username,
			IDKelas:           r.IDKelas,
			NamaKelas:         r.NamaKelas,
			IDNilai:           r.IDNilai,
			Status:            status,
			Nilai:             r.Nilai,
			WktMulai:          r.WktMulai,
			AktivitasTerakhir: r.AktivitasTerakhir,
			WktSelesai:        r.WktSelesai,
		})
	}

	return &dto.MonitoringResponse{
		NamaUjian: namaUjian,
		Summary:   summary,
		Data:      data,
	}, nil
}

// analisisJawabanRow adalah baris flat hasil LEFT JOIN peserta->nilai->jawaban->soal.
// NoSoal/Jawaban/IsBenar bernilai nil jika peserta belum pernah mulai ujian ini sama sekali.
type analisisJawabanRow struct {
	IDPeserta   string   `gorm:"column:id_peserta"`
	NamaPeserta string   `gorm:"column:nama_peserta"`
	NoSoal      *int     `gorm:"column:no_soal"`
	Jawaban     *string  `gorm:"column:jawaban"`
	IsBenar     *int     `gorm:"column:is_benar"`
	NilaiAkhir  *float64 `gorm:"column:nilai_akhir"`
}

// AnalisisSoalByJadwal membuat ZIP berisi satu .xlsx analisis per kelas: baris = peserta,
// kolom = soal (diurutkan berdasarkan no_soal — bukan no_urut — supaya tetap sebanding
// antar peserta walau acak_soal aktif dan urutan tampil tiap peserta berbeda-beda).
// Cell diwarnai hijau (jawaban benar), merah (jawaban salah), atau kuning (belum dijawab).
func (s *nilaiService) AnalisisSoalByJadwal(idJadwal string) (*ExportResult, error) {
	// 1. Ambil nama_ujian & id_bank_soal dari jadwal
	type jadwalMeta struct {
		NamaUjian  string `gorm:"column:nama_ujian"`
		IDBankSoal string `gorm:"column:id_bank_soal"`
	}
	var meta jadwalMeta
	if err := s.db.Table("jadwal").
		Select("nama_ujian, id_bank_soal").
		Where("id = ? AND deleted_at IS NULL", idJadwal).
		Scan(&meta).Error; err != nil || meta.NamaUjian == "" {
		return nil, errors.New("jadwal tidak ditemukan")
	}

	// 2. Nama mapel (untuk nama file)
	var namaMapel string
	if err := s.db.Table("bank_soal").
		Select("mapel.nama_mapel").
		Joins("INNER JOIN mapel ON bank_soal.id_mapel = mapel.id").
		Where("bank_soal.id = ? AND bank_soal.deleted_at IS NULL", meta.IDBankSoal).
		Scan(&namaMapel).Error; err != nil || namaMapel == "" {
		return nil, errors.New("mapel tidak ditemukan untuk jadwal ini")
	}

	// 3. Daftar no_soal — kolom tetap & urutannya sama untuk semua peserta
	var noSoalList []int
	if err := s.db.Table("soal").
		Select("no_soal").
		Where("id_bank_soal = ? AND deleted_at IS NULL", meta.IDBankSoal).
		Order("no_soal ASC").
		Scan(&noSoalList).Error; err != nil {
		return nil, err
	}
	if len(noSoalList) == 0 {
		return nil, errors.New("jadwal ini belum memiliki soal")
	}

	// 4. Daftar kelas yang terdaftar pada jadwal ini
	type kelasRow struct {
		IDKelas   string `gorm:"column:id_kelas"`
		NamaKelas string `gorm:"column:nama_kelas"`
	}
	var kelasList []kelasRow
	if err := s.db.Table("jadwal_kelas").
		Select("jadwal_kelas.id_kelas, kelas.nama_kelas").
		Joins("INNER JOIN kelas ON jadwal_kelas.id_kelas = kelas.id").
		Where("jadwal_kelas.id_jadwal = ?", idJadwal).
		Scan(&kelasList).Error; err != nil {
		return nil, err
	}
	if len(kelasList) == 0 {
		return nil, errors.New("tidak ada kelas yang terdaftar pada jadwal ini")
	}

	var zipBuf bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuf)

	for _, kelas := range kelasList {
		var rows []analisisJawabanRow
		if err := s.db.Table("peserta").
			Select(`
				peserta.id AS id_peserta,
				peserta.nama AS nama_peserta,
				soal.no_soal,
				jawaban.jawaban,
				jawaban.is_benar,
				nilai.nilai AS nilai_akhir
			`).
			Joins("LEFT JOIN nilai ON nilai.id_peserta = peserta.id AND nilai.id_jadwal = ? AND nilai.deleted_at IS NULL", idJadwal).
			Joins("LEFT JOIN jawaban ON jawaban.id_nilai = nilai.id AND jawaban.deleted_at IS NULL").
			Joins("LEFT JOIN soal ON soal.id = jawaban.id_soal AND soal.deleted_at IS NULL").
			Where("peserta.id_kelas = ? AND peserta.deleted_at IS NULL", kelas.IDKelas).
			Order("peserta.nama ASC, soal.no_soal ASC").
			Scan(&rows).Error; err != nil {
			return nil, fmt.Errorf("gagal ambil data jawaban kelas %s: %w", kelas.NamaKelas, err)
		}

		xlsxFile, err := buildAnalisisSoalWorkbook(rows, noSoalList)
		if err != nil {
			return nil, fmt.Errorf("gagal membuat sheet analisis kelas %s: %w", kelas.NamaKelas, err)
		}

		var xlsBuf bytes.Buffer
		if err := xlsxFile.Write(&xlsBuf); err != nil {
			return nil, err
		}

		safeMapel := strings.ReplaceAll(namaMapel, " ", "_")
		safeKelas := strings.ReplaceAll(kelas.NamaKelas, " ", "_")
		filename := fmt.Sprintf("Analisis_%s_%s.xlsx", safeMapel, safeKelas)

		zipEntry, err := zipWriter.Create(filename)
		if err != nil {
			return nil, err
		}
		if _, err := zipEntry.Write(xlsBuf.Bytes()); err != nil {
			return nil, err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return &ExportResult{
		ZipBytes:  zipBuf.Bytes(),
		NamaUjian: meta.NamaUjian,
	}, nil
}

type analisisAnswer struct {
	Jawaban *string
	IsBenar *int
}

type analisisPesertaAgg struct {
	Nama       string
	Started    bool
	NilaiAkhir *float64
	Answers    map[int]analisisAnswer
}

// buildAnalisisSoalWorkbook memivot baris flat (peserta x soal) menjadi satu sheet:
// baris = peserta (urut sesuai kemunculan pertama di `rows`, yaitu nama ASC),
// kolom = soal urut noSoalList, plus ringkasan Benar/Salah/Belum Dijawab dan Nilai Akhir
// (kolom paling kanan, dari tabel nilai — kosong jika peserta belum pernah mulai ujian)
// per peserta, dan baris "Jumlah Benar per Soal" di paling bawah untuk analisis tingkat
// kesulitan soal.
func buildAnalisisSoalWorkbook(rows []analisisJawabanRow, noSoalList []int) (*excelize.File, error) {
	order := make([]string, 0)
	agg := make(map[string]*analisisPesertaAgg)

	for _, r := range rows {
		a, ok := agg[r.IDPeserta]
		if !ok {
			a = &analisisPesertaAgg{Nama: r.NamaPeserta, Answers: make(map[int]analisisAnswer)}
			agg[r.IDPeserta] = a
			order = append(order, r.IDPeserta)
		}
		if r.NoSoal != nil {
			a.Started = true
			a.Answers[*r.NoSoal] = analisisAnswer{Jawaban: r.Jawaban, IsBenar: r.IsBenar}
		}
		if r.NilaiAkhir != nil {
			a.NilaiAkhir = r.NilaiAkhir
		}
	}

	f := excelize.NewFile()
	sheet := "Analisis Soal"
	f.SetSheetName("Sheet1", sheet)

	greenStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"C6EFCE"}, Pattern: 1},
		Font: &excelize.Font{Color: "006100"},
	})
	if err != nil {
		return nil, err
	}
	redStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FFC7CE"}, Pattern: 1},
		Font: &excelize.Font{Color: "9C0006"},
	})
	if err != nil {
		return nil, err
	}
	yellowStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FFEB9C"}, Pattern: 1},
		Font: &excelize.Font{Color: "9C6500"},
	})
	if err != nil {
		return nil, err
	}
	headerStyle, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	if err != nil {
		return nil, err
	}

	// Header
	headers := []string{"No", "Nama Peserta"}
	for _, noSoal := range noSoalList {
		headers = append(headers, fmt.Sprintf("Soal %d", noSoal))
	}
	headers = append(headers, "Benar", "Salah", "Belum Dijawab", "Nilai Akhir")
	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	summaryColStart := len(noSoalList) + 3 // kolom setelah "No", "Nama Peserta", dan semua soal
	benarPerSoal := make([]int, len(noSoalList))

	for i, idPeserta := range order {
		a := agg[idPeserta]
		rowNum := i + 2

		noCell, _ := excelize.CoordinatesToCellName(1, rowNum)
		namaCell, _ := excelize.CoordinatesToCellName(2, rowNum)
		f.SetCellValue(sheet, noCell, i+1)
		f.SetCellValue(sheet, namaCell, a.Nama)

		var benar, salah, belum int
		for j, noSoal := range noSoalList {
			cellName, _ := excelize.CoordinatesToCellName(j+3, rowNum)

			if !a.Started {
				continue // peserta belum pernah mulai ujian ini — cell dibiarkan kosong tanpa warna
			}

			ans, answered := a.Answers[noSoal]
			if !answered || ans.Jawaban == nil {
				f.SetCellStyle(sheet, cellName, cellName, yellowStyle)
				belum++
				continue
			}

			f.SetCellValue(sheet, cellName, strings.ToUpper(*ans.Jawaban))
			if ans.IsBenar != nil && *ans.IsBenar == 1 {
				f.SetCellStyle(sheet, cellName, cellName, greenStyle)
				benar++
				benarPerSoal[j]++
			} else {
				f.SetCellStyle(sheet, cellName, cellName, redStyle)
				salah++
			}
		}

		benarCell, _ := excelize.CoordinatesToCellName(summaryColStart, rowNum)
		salahCell, _ := excelize.CoordinatesToCellName(summaryColStart+1, rowNum)
		belumCell, _ := excelize.CoordinatesToCellName(summaryColStart+2, rowNum)
		nilaiCell, _ := excelize.CoordinatesToCellName(summaryColStart+3, rowNum)
		f.SetCellValue(sheet, benarCell, benar)
		f.SetCellValue(sheet, salahCell, salah)
		f.SetCellValue(sheet, belumCell, belum)
		if a.NilaiAkhir != nil {
			f.SetCellValue(sheet, nilaiCell, *a.NilaiAkhir)
		}
	}

	// Baris total: jumlah peserta yang menjawab benar per soal — untuk lihat soal mana yang paling banyak dijawab salah
	totalRow := len(order) + 2
	totalLabelCell, _ := excelize.CoordinatesToCellName(2, totalRow)
	f.SetCellValue(sheet, totalLabelCell, "Jumlah Benar per Soal")
	f.SetCellStyle(sheet, totalLabelCell, totalLabelCell, headerStyle)
	for j := range noSoalList {
		cellName, _ := excelize.CoordinatesToCellName(j+3, totalRow)
		f.SetCellValue(sheet, cellName, benarPerSoal[j])
		f.SetCellStyle(sheet, cellName, cellName, headerStyle)
	}

	f.SetColWidth(sheet, "A", "A", 5)
	f.SetColWidth(sheet, "B", "B", 28)

	return f, nil
}
