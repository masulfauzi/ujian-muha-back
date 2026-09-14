package service

import (
	"errors"
	"math"
	"time"

	"backend/internal/constants"
	"backend/internal/modules/jadwal/dto"
	"backend/internal/modules/jadwal/model"
	"backend/internal/modules/jadwal/repository"
	jadwalkelasmodel "backend/internal/modules/jadwal_kelas/model"
	jadwalkelasrepo "backend/internal/modules/jadwal_kelas/repository"
	pesertarepo "backend/internal/modules/peserta/repository"

	"gorm.io/gorm"
)

const timeLayout = "2006-01-02 15:04:05"

var jakartaLoc, _ = time.LoadLocation("Asia/Jakarta")

type JadwalService interface {
	CreateJadwal(req *dto.CreateJadwalRequest) (*dto.JadwalResponse, error)
	GetJadwalByID(id string) (*dto.JadwalResponse, error)
	GetAllJadwal(page, pageSize int) (*dto.JadwalListResponse, error)
	GetJadwalByBankSoal(bankSoalID string, page, pageSize int) (*dto.JadwalListResponse, error)
	UpdateJadwal(id string, req *dto.UpdateJadwalRequest) (*dto.JadwalResponse, error)
	DeleteJadwal(id string) error
	RestoreJadwal(id string) error
	GetJadwalAktifHariIniByUser(userID string) ([]dto.JadwalAktifResponse, error)
}

type jadwalService struct {
	repo            repository.JadwalRepository
	jadwalKelasRepo jadwalkelasrepo.JadwalKelasRepository
	pesertaRepo     pesertarepo.PesertaRepository
}

func NewJadwalService(repo repository.JadwalRepository, jadwalKelasRepo jadwalkelasrepo.JadwalKelasRepository, pesertaRepo pesertarepo.PesertaRepository) JadwalService {
	return &jadwalService{repo: repo, jadwalKelasRepo: jadwalKelasRepo, pesertaRepo: pesertaRepo}
}

func (s *jadwalService) CreateJadwal(req *dto.CreateJadwalRequest) (*dto.JadwalResponse, error) {
	wktMulai, err := time.ParseInLocation(timeLayout, req.WktMulai, jakartaLoc)
	if err != nil {
		return nil, errors.New("format wkt_mulai tidak valid, gunakan: 2006-01-02 15:04:05")
	}

	wktSelesai, err := time.ParseInLocation(timeLayout, req.WktSelesai, jakartaLoc)
	if err != nil {
		return nil, errors.New("format wkt_selesai tidak valid, gunakan: 2006-01-02 15:04:05")
	}

	if !wktSelesai.After(wktMulai) {
		return nil, errors.New("wkt_selesai harus setelah wkt_mulai")
	}

	if len(req.IDKelas) == 0 {
		return nil, errors.New("minimal harus ada satu kelas yang didaftarkan")
	}

	jadwal := &model.Jadwal{
		IDBankSoal: req.IDBankSoal,
		NamaUjian:  req.NamaUjian,
		Tingkat:    req.Tingkat,
		WktMulai:   wktMulai,
		WktSelesai: wktSelesai,
		Durasi:     req.Durasi,
		AcakSoal:   req.AcakSoal,
		AcakOpsi:   req.AcakOpsi,
		WajibToken: req.WajibToken,
	}

	if err := s.repo.Create(jadwal); err != nil {
		return nil, err
	}

	jadwalKelasList := make([]*jadwalkelasmodel.JadwalKelas, 0, len(req.IDKelas))
	for _, idKelas := range req.IDKelas {
		exists, err := s.jadwalKelasRepo.CheckDuplicate(jadwal.ID, idKelas)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("kelas dengan ID " + idKelas + " sudah terdaftar di jadwal tersebut")
		}

		jadwalKelasList = append(jadwalKelasList, &jadwalkelasmodel.JadwalKelas{
			IDJadwal: jadwal.ID,
			IDKelas:  idKelas,
		})
	}

	if err := s.jadwalKelasRepo.CreateBulk(jadwalKelasList); err != nil {
		return nil, err
	}

	created, err := s.repo.GetByIDWithBankSoal(jadwal.ID)
	if err != nil {
		return nil, err
	}

	return joinedToResponse(created), nil
}

func (s *jadwalService) GetJadwalByID(id string) (*dto.JadwalResponse, error) {
	jadwal, err := s.repo.GetByIDWithKelas(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	return jelasToResponse(jadwal), nil
}

func (s *jadwalService) GetAllJadwal(page, pageSize int) (*dto.JadwalListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	jadwalList, total, err := s.repo.GetAllWithBankSoal(page, pageSize)
	if err != nil {
		return nil, err
	}

	responses := []dto.JadwalResponse{}
	for _, j := range jadwalList {
		responses = append(responses, *joinedToResponse(&j))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.JadwalListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *jadwalService) GetJadwalByBankSoal(bankSoalID string, page, pageSize int) (*dto.JadwalListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	jadwalList, total, err := s.repo.GetByBankSoalID(bankSoalID, page, pageSize)
	if err != nil {
		return nil, err
	}

	responses := []dto.JadwalResponse{}
	for _, j := range jadwalList {
		responses = append(responses, *joinedToResponse(&j))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.JadwalListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *jadwalService) UpdateJadwal(id string, req *dto.UpdateJadwalRequest) (*dto.JadwalResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	wktMulai, err := time.ParseInLocation(timeLayout, req.WktMulai, jakartaLoc)
	if err != nil {
		return nil, errors.New("format wkt_mulai tidak valid, gunakan: 2006-01-02 15:04:05")
	}

	wktSelesai, err := time.ParseInLocation(timeLayout, req.WktSelesai, jakartaLoc)
	if err != nil {
		return nil, errors.New("format wkt_selesai tidak valid, gunakan: 2006-01-02 15:04:05")
	}

	if !wktSelesai.After(wktMulai) {
		return nil, errors.New("wkt_selesai harus setelah wkt_mulai")
	}

	if len(req.IDKelas) == 0 {
		return nil, errors.New("minimal harus ada satu kelas yang didaftarkan")
	}

	existing.IDBankSoal = req.IDBankSoal
	existing.NamaUjian  = req.NamaUjian
	existing.Tingkat    = req.Tingkat
	existing.WktMulai   = wktMulai
	existing.WktSelesai = wktSelesai
	existing.Durasi     = req.Durasi
	existing.AcakSoal   = req.AcakSoal
	existing.AcakOpsi   = req.AcakOpsi
	existing.WajibToken = req.WajibToken

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	if err := s.jadwalKelasRepo.DeleteByJadwalID(id); err != nil {
		return nil, err
	}

	jadwalKelasList := make([]*jadwalkelasmodel.JadwalKelas, 0, len(req.IDKelas))
	for _, idKelas := range req.IDKelas {
		jadwalKelasList = append(jadwalKelasList, &jadwalkelasmodel.JadwalKelas{
			IDJadwal: id,
			IDKelas:  idKelas,
		})
	}

	if err := s.jadwalKelasRepo.CreateBulk(jadwalKelasList); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByIDWithKelas(id)
	if err != nil {
		return nil, err
	}

	return jelasToResponse(updated), nil
}

func (s *jadwalService) DeleteJadwal(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}
	return s.repo.Delete(id)
}

func (s *jadwalService) RestoreJadwal(id string) error {
	return s.repo.Restore(id)
}

func (s *jadwalService) GetJadwalAktifHariIniByUser(userID string) ([]dto.JadwalAktifResponse, error) {
	peserta, err := s.pesertaRepo.GetRawByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("peserta tidak ditemukan")
		}
		return nil, err
	}

	jadwalList, err := s.repo.GetAktifHariIniByKelas(peserta.IDKelas, peserta.ID)
	if err != nil {
		return nil, err
	}

	responses := []dto.JadwalAktifResponse{}
	for _, j := range jadwalList {
		status := "belum_mengerjakan"
		if j.IDNilai != nil {
			if j.NilaiWktSelesai != nil {
				status = "sudah_mengerjakan"
			} else {
				status = "sedang_mengerjakan"
			}
		}

		responses = append(responses, dto.JadwalAktifResponse{
			ID:               j.ID,
			IDBankSoal:       j.IDBankSoal,
			NamaBankSoal:     j.NamaBankSoal,
			NamaUjian:        j.NamaUjian,
			Tingkat:          j.Tingkat,
			WktMulai:         j.WktMulai,
			WktSelesai:       j.WktSelesai,
			Durasi:           j.Durasi,
			AcakSoal:         j.AcakSoal,
			AcakOpsi:         j.AcakOpsi,
			WajibToken:       j.WajibToken,
			IDNilai:          j.IDNilai,
			StatusPengerjaan: status,
		})
	}
	return responses, nil
}

func joinedToResponse(j *repository.JadwalWithBankSoal) *dto.JadwalResponse {
	return &dto.JadwalResponse{
		ID:           j.ID,
		IDBankSoal:   j.IDBankSoal,
		NamaBankSoal: j.NamaBankSoal,
		NamaUjian:    j.NamaUjian,
		Tingkat:      j.Tingkat,
		WktMulai:     j.WktMulai,
		WktSelesai:   j.WktSelesai,
		Durasi:       j.Durasi,
		AcakSoal:     j.AcakSoal,
		AcakOpsi:     j.AcakOpsi,
		WajibToken:   j.WajibToken,
		CreatedAt:    j.CreatedAt,
		UpdatedAt:    j.UpdatedAt,
	}
}

func jelasToResponse(j *repository.JadwalWithKelas) *dto.JadwalResponse {
	kelasList := []dto.KelasItem{}
	for _, k := range j.IDKelas {
		kelasList = append(kelasList, dto.KelasItem{
			ID:        k.ID,
			IDKelas:   k.IDKelas,
			NamaKelas: k.NamaKelas,
		})
	}

	jurusanList := []dto.JurusanItem{}
	for _, jur := range j.IDJurusan {
		jurusanList = append(jurusanList, dto.JurusanItem{
			ID:          jur.ID,
			IDJurusan:   jur.IDJurusan,
			NamaJurusan: jur.NamaJurusan,
		})
	}

	return &dto.JadwalResponse{
		ID:           j.ID,
		IDBankSoal:   j.IDBankSoal,
		NamaBankSoal: j.NamaBankSoal,
		NamaUjian:    j.NamaUjian,
		Tingkat:      j.Tingkat,
		WktMulai:     j.WktMulai,
		WktSelesai:   j.WktSelesai,
		Durasi:       j.Durasi,
		AcakSoal:     j.AcakSoal,
		AcakOpsi:     j.AcakOpsi,
		WajibToken:   j.WajibToken,
		IDKelas:      kelasList,
		IDJurusan:    jurusanList,
		CreatedAt:    j.CreatedAt,
		UpdatedAt:    j.UpdatedAt,
	}
}
