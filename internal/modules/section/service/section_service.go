package service

import (
	"errors"
	"time"

	"backend/internal/constants"
	"backend/internal/modules/section/dto"
	"backend/internal/modules/section/model"
	"backend/internal/modules/section/repository"

	"gorm.io/gorm"
)

// jakartaLoc dipakai saat menulis wkt_mulai_section (kolom "timestamp" tanpa
// zona di Postgres, yang di seluruh codebase ini isinya sengaja disimpan
// sebagai wall-clock WIB apa adanya — lihat nilai_service.go untuk konvensi
// yang sama).
var jakartaLoc, _ = time.LoadLocation("Asia/Jakarta")

func init() {
	if jakartaLoc == nil {
		jakartaLoc = time.FixedZone("WIB", 7*60*60)
	}
}

type SectionService interface {
	DefineSections(idJadwal string, req *dto.DefineSectionRequest) ([]dto.SectionResponse, error)
	GetSectionsByJadwal(idJadwal string) ([]dto.SectionResponse, error)
	GetSectionByID(id string) (*dto.SectionResponse, error)
	DeleteSection(id string) error
	RestoreSection(id string) error
}

type sectionService struct {
	repo repository.SectionRepository
	db   *gorm.DB
}

func NewSectionService(repo repository.SectionRepository, db *gorm.DB) SectionService {
	return &sectionService{repo: repo, db: db}
}

func (s *sectionService) DefineSections(idJadwal string, req *dto.DefineSectionRequest) ([]dto.SectionResponse, error) {
	if len(req.Sections) == 0 {
		return nil, errors.New("minimal harus ada 1 section")
	}

	totalSoal, err := s.repo.CountSoalByJadwalID(idJadwal)
	if err != nil {
		return nil, err
	}
	if totalSoal == 0 {
		return nil, errors.New("jadwal tidak ditemukan atau belum memiliki soal")
	}

	var sumJmlSoal int
	for _, item := range req.Sections {
		if item.JmlSoal <= 0 {
			return nil, errors.New("jml_soal setiap section harus lebih dari 0")
		}
		if item.NamaSection == "" {
			return nil, errors.New("nama_section tidak boleh kosong")
		}
		sumJmlSoal += item.JmlSoal
	}
	if int64(sumJmlSoal) != totalSoal {
		return nil, errors.New("total jml_soal pada semua section harus sama dengan jumlah soal pada jadwal ini")
	}

	sections := make([]model.Section, len(req.Sections))
	noUrutAwal := 1
	for i, item := range req.Sections {
		noUrutAkhir := noUrutAwal + item.JmlSoal - 1
		sections[i] = model.Section{
			IDJadwal:           idJadwal,
			NamaSection:        item.NamaSection,
			Urutan:             i + 1,
			NoUrutAwal:         noUrutAwal,
			NoUrutAkhir:        noUrutAkhir,
			DurasiMenitMinimal: item.DurasiMenitMinimal,
		}
		noUrutAwal = noUrutAkhir + 1
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.SoftDeleteByJadwalIDWithTx(tx, idJadwal); err != nil {
			return err
		}
		if err := s.repo.CreateManyWithTx(tx, sections); err != nil {
			return err
		}

		// Section lama baru saja di-soft-delete di atas, jadi peserta yang sedang
		// mengerjakan (belum wkt_selesai) dan sebelumnya sudah punya section aktif
		// sekarang menggantung ke section yang sudah terhapus. Migrasikan mereka ke
		// section pertama (urutan 1) yang baru supaya tidak macet — timer durasi
		// minimal-nya juga direset dari sekarang untuk section baru itu. Peserta yang
		// belum pernah punya section aktif (belum mulai, atau jadwal ini sebelumnya
		// tidak pakai section) sengaja tidak disentuh.
		firstSectionID := sections[0].ID // urutan=1, lihat konstruksi slice sections di atas
		now := time.Now().In(jakartaLoc)
		return tx.Table("nilai").
			Where("id_jadwal = ? AND wkt_selesai IS NULL AND id_section_aktif IS NOT NULL", idJadwal).
			Updates(map[string]interface{}{
				"id_section_aktif":  firstSectionID,
				"wkt_mulai_section": now,
			}).Error
	})
	if err != nil {
		return nil, err
	}

	return s.GetSectionsByJadwal(idJadwal)
}

func (s *sectionService) GetSectionsByJadwal(idJadwal string) ([]dto.SectionResponse, error) {
	sections, err := s.repo.GetByJadwalID(idJadwal)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.SectionResponse, len(sections))
	for i, sec := range sections {
		responses[i] = toResponse(&sec)
	}
	return responses, nil
}

func (s *sectionService) GetSectionByID(id string) (*dto.SectionResponse, error) {
	section, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	resp := toResponse(section)
	return &resp, nil
}

func (s *sectionService) DeleteSection(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}
	return s.repo.Delete(id)
}

func (s *sectionService) RestoreSection(id string) error {
	return s.repo.Restore(id)
}

func toResponse(s *model.Section) dto.SectionResponse {
	return dto.SectionResponse{
		ID:                 s.ID,
		IDJadwal:           s.IDJadwal,
		NamaSection:        s.NamaSection,
		Urutan:             s.Urutan,
		NoUrutAwal:         s.NoUrutAwal,
		NoUrutAkhir:        s.NoUrutAkhir,
		JmlSoal:            s.NoUrutAkhir - s.NoUrutAwal + 1,
		DurasiMenitMinimal: s.DurasiMenitMinimal,
		CreatedAt:          s.CreatedAt.Format(timeLayout),
		UpdatedAt:          s.UpdatedAt.Format(timeLayout),
	}
}

const timeLayout = "2006-01-02 15:04:05"
