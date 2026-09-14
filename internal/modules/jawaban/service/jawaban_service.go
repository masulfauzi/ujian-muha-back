package service

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"math"
	"math/rand"
	"strings"

	"backend/internal/constants"
	"backend/internal/modules/jawaban/dto"
	"backend/internal/modules/jawaban/model"
	"backend/internal/modules/jawaban/repository"
	nilairepo "backend/internal/modules/nilai/repository"
	sectionrepo "backend/internal/modules/section/repository"

	"gorm.io/gorm"
)

// ErrSectionBelumTerbuka dikembalikan ketika section yang diminta belum dibuka (urutan-nya
// lebih besar dari section aktif/frontier sesi ujian ini) — dipakai controller untuk balas 403.
var ErrSectionBelumTerbuka = errors.New("section ini belum terbuka")

type JawabanService interface {
	CreateJawaban(req *dto.CreateJawabanRequest) (*dto.JawabanResponse, error)
	GetJawabanByID(id string) (*dto.JawabanResponse, error)
	GetAllJawaban(page, pageSize int, idNilai, idPeserta, idSoal string) (*dto.JawabanListResponse, error)
	GetJawabanByNilai(idNilai string) ([]dto.JawabanResponse, error)
	GetJawabanByNilaiSection(idNilai, idSection string) ([]dto.JawabanResponse, error)
	GetJawabanByPeserta(idPeserta string, page, pageSize int) (*dto.JawabanListResponse, error)
	UpdateJawaban(id string, req *dto.UpdateJawabanRequest) (*dto.JawabanResponse, error)
	DeleteJawaban(id string) error
	RestoreJawaban(id string) error
}

type jawabanService struct {
	repo        repository.JawabanRepository
	nilaiRepo   nilairepo.NilaiRepository
	sectionRepo sectionrepo.SectionRepository
}

func NewJawabanService(repo repository.JawabanRepository, nilaiRepo nilairepo.NilaiRepository, sectionRepo sectionrepo.SectionRepository) JawabanService {
	return &jawabanService{repo: repo, nilaiRepo: nilaiRepo, sectionRepo: sectionRepo}
}

func normalizeJawaban(j string) (string, error) {
	j = strings.ToUpper(strings.TrimSpace(j))
	switch j {
	case "A", "B", "C", "D", "E":
		return j, nil
	default:
		return "", errors.New("jawaban harus salah satu dari: A, B, C, D, E")
	}
}

// generateOpsiOrder menghasilkan urutan acak yang sama dengan soal service.
// Wajib identik dengan soalService.generateSeed + rand.Shuffle.
func generateOpsiOrder(pesertaID, soalID string) []string {
	hash := md5.Sum([]byte(pesertaID + "|" + soalID))
	seed := int64(binary.BigEndian.Uint64(hash[:8]))
	rng := rand.New(rand.NewSource(seed))

	order := []string{"A", "B", "C", "D", "E"}
	rng.Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})
	return order
}

// reverseMapJawaban mengonversi posisi acak yang dipilih peserta
// ke posisi asli sebelum diacak, agar bisa dibandingkan dengan kunci di DB.
func reverseMapJawaban(submittedJawaban, pesertaID, soalID string) string {
	order := generateOpsiOrder(pesertaID, soalID)
	idx := int(submittedJawaban[0] - 'A') // "B" → 1
	if idx < 0 || idx >= len(order) {
		return submittedJawaban
	}
	return order[idx] // posisi asli
}

func randomizeOpsi(detail *repository.JawabanWithDetail) (opsiA, opsiB, opsiC, opsiD, opsiE, gambarA, gambarB, gambarC, gambarD, gambarE, kunci string) {
	opsiOrder := generateOpsiOrder(detail.IDPeserta, detail.IDSoal)

	opsiMap := map[string]string{
		"A": detail.OpsiA, "B": detail.OpsiB, "C": detail.OpsiC,
		"D": detail.OpsiD, "E": detail.OpsiE,
	}
	gambarMap := map[string]string{
		"A": detail.GambarA, "B": detail.GambarB, "C": detail.GambarC,
		"D": detail.GambarD, "E": detail.GambarE,
	}

	newKunciPos := 0
	for i, key := range opsiOrder {
		if key == detail.Kunci {
			newKunciPos = i
			break
		}
	}
	newKunci := string(rune('A' + newKunciPos))

	return opsiMap[opsiOrder[0]], opsiMap[opsiOrder[1]], opsiMap[opsiOrder[2]], opsiMap[opsiOrder[3]], opsiMap[opsiOrder[4]],
		gambarMap[opsiOrder[0]], gambarMap[opsiOrder[1]], gambarMap[opsiOrder[2]], gambarMap[opsiOrder[3]], gambarMap[opsiOrder[4]],
		newKunci
}

func (s *jawabanService) CreateJawaban(req *dto.CreateJawabanRequest) (*dto.JawabanResponse, error) {
	jawaban, err := normalizeJawaban(req.Jawaban)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.CheckDuplicate(req.IDNilai, req.IDSoal)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("jawaban untuk soal ini di attempt tersebut sudah ada — gunakan endpoint update")
	}

	kunci, err := s.repo.GetSoalKunci(req.IDSoal)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("soal tidak ditemukan")
		}
		return nil, err
	}
	if kunci == "" {
		return nil, errors.New("soal tidak ditemukan")
	}

	acakOpsi, err := s.repo.GetAcakOpsiByNilaiID(req.IDNilai)
	if err != nil {
		return nil, err
	}

	evaluatedJawaban := jawaban
	if acakOpsi == 1 {
		evaluatedJawaban = reverseMapJawaban(jawaban, req.IDPeserta, req.IDSoal)
	}
	isBenarInt := 0
	if strings.EqualFold(evaluatedJawaban, strings.TrimSpace(kunci)) {
		isBenarInt = 1
	}

	row := &model.Jawaban{
		IDNilai:   req.IDNilai,
		IDSoal:    req.IDSoal,
		IDPeserta: req.IDPeserta,
		NoUrut:    req.NoUrut,
		Jawaban:   &jawaban,
		IsBenar:   &isBenarInt,
	}

	if err := s.repo.Create(row); err != nil {
		return nil, err
	}

	created, err := s.repo.GetByIDWithDetail(row.ID)
	if err != nil {
		return nil, err
	}
	return detailToResponse(created), nil
}

func (s *jawabanService) GetJawabanByID(id string) (*dto.JawabanResponse, error) {
	result, err := s.repo.GetByIDWithDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	return detailToResponse(result), nil
}

func (s *jawabanService) GetAllJawaban(page, pageSize int, idNilai, idPeserta, idSoal string) (*dto.JawabanListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	results, total, err := s.repo.GetAllWithDetail(page, pageSize, idNilai, idPeserta, idSoal)
	if err != nil {
		return nil, err
	}

	responses := []dto.JawabanResponse{}
	for _, r := range results {
		responses = append(responses, *detailToResponse(&r))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.JawabanListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *jawabanService) GetJawabanByNilai(idNilai string) ([]dto.JawabanResponse, error) {
	resp, err := s.GetAllJawaban(1, 99999, idNilai, "", "")
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetJawabanByNilaiSection mengambil soal satu sesi ujian yang berada dalam satu section.
// Section dengan urutan di atas frontier (nilai.IDSectionAktif) ditolak dengan ErrSectionBelumTerbuka,
// supaya siswa tidak bisa melihat soal section berikutnya sebelum waktu minimal terlampaui.
func (s *jawabanService) GetJawabanByNilaiSection(idNilai, idSection string) ([]dto.JawabanResponse, error) {
	section, err := s.sectionRepo.GetByID(idSection)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("section tidak ditemukan")
		}
		return nil, err
	}

	nilai, err := s.nilaiRepo.GetByID(idNilai)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	if section.IDJadwal != nilai.IDJadwal {
		return nil, errors.New("section tidak sesuai dengan jadwal sesi ujian ini")
	}
	if nilai.IDSectionAktif == nil {
		return nil, errors.New("sesi ujian ini tidak menggunakan section")
	}

	activeSection, err := s.sectionRepo.GetByID(*nilai.IDSectionAktif)
	if err != nil {
		return nil, err
	}
	if section.Urutan > activeSection.Urutan {
		return nil, ErrSectionBelumTerbuka
	}

	results, err := s.repo.GetByNilaiIDInRange(idNilai, section.NoUrutAwal, section.NoUrutAkhir)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.JawabanResponse, len(results))
	for i, r := range results {
		responses[i] = *detailToResponse(&r)
	}
	return responses, nil
}

func (s *jawabanService) GetJawabanByPeserta(idPeserta string, page, pageSize int) (*dto.JawabanListResponse, error) {
	return s.GetAllJawaban(page, pageSize, "", idPeserta, "")
}

func (s *jawabanService) UpdateJawaban(id string, req *dto.UpdateJawabanRequest) (*dto.JawabanResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	jawaban, err := normalizeJawaban(req.Jawaban)
	if err != nil {
		return nil, err
	}

	kunci, err := s.repo.GetSoalKunci(existing.IDSoal)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("soal tidak ditemukan")
		}
		return nil, err
	}
	if kunci == "" {
		return nil, errors.New("soal tidak ditemukan")
	}

	acakOpsi, err := s.repo.GetAcakOpsiByNilaiID(existing.IDNilai)
	if err != nil {
		return nil, err
	}

	evaluatedJawaban := jawaban
	if acakOpsi == 1 {
		evaluatedJawaban = reverseMapJawaban(jawaban, existing.IDPeserta, existing.IDSoal)
	}
	isBenarInt := 0
	if strings.EqualFold(evaluatedJawaban, strings.TrimSpace(kunci)) {
		isBenarInt = 1
	}

	existing.Jawaban = &jawaban
	existing.IsBenar = &isBenarInt

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByIDWithDetail(id)
	if err != nil {
		return nil, err
	}
	return detailToResponse(updated), nil
}

func (s *jawabanService) DeleteJawaban(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}
	return s.repo.Delete(id)
}

func (s *jawabanService) RestoreJawaban(id string) error {
	return s.repo.Restore(id)
}

func detailToResponse(r *repository.JawabanWithDetail) *dto.JawabanResponse {
	opsiA, opsiB, opsiC, opsiD, opsiE := r.OpsiA, r.OpsiB, r.OpsiC, r.OpsiD, r.OpsiE
	gambarA, gambarB, gambarC, gambarD, gambarE := r.GambarA, r.GambarB, r.GambarC, r.GambarD, r.GambarE
	kunci := r.Kunci

	if r.AcakOpsi == 1 {
		opsiA, opsiB, opsiC, opsiD, opsiE,
			gambarA, gambarB, gambarC, gambarD, gambarE,
			kunci = randomizeOpsi(r)
	}

	return &dto.JawabanResponse{
		ID:          r.ID,
		IDNilai:     r.IDNilai,
		IDSoal:      r.IDSoal,
		NoUrut:      r.NoUrut,
		NoSoal:      r.NoSoal,
		SoalText:    r.SoalText,
		Kunci:       kunci,
		OpsiA:       opsiA,
		OpsiB:       opsiB,
		OpsiC:       opsiC,
		OpsiD:       opsiD,
		OpsiE:       opsiE,
		GambarA:     gambarA,
		GambarB:     gambarB,
		GambarC:     gambarC,
		GambarD:     gambarD,
		GambarE:     gambarE,
		IDPeserta:   r.IDPeserta,
		NamaPeserta: r.NamaPeserta,
		Jawaban:     r.Jawaban,
		IsBenar:     r.IsBenar,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
