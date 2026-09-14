package service

import (
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/jadwal_kelas/dto"
	"backend/internal/modules/jadwal_kelas/model"
	"backend/internal/modules/jadwal_kelas/repository"

	"gorm.io/gorm"
)

type JadwalKelasService interface {
	CreateJadwalKelas(req *dto.CreateJadwalKelasRequest) (*dto.JadwalKelasResponse, error)
	GetJadwalKelasByID(id string) (*dto.JadwalKelasResponse, error)
	GetAllJadwalKelas(page, pageSize int, idJadwal string, idKelas string) (*dto.JadwalKelasListResponse, error)
	UpdateJadwalKelas(id string, req *dto.UpdateJadwalKelasRequest) (*dto.JadwalKelasResponse, error)
	DeleteJadwalKelas(id string) error
}

type jadwalKelasService struct {
	repo repository.JadwalKelasRepository
}

func NewJadwalKelasService(repo repository.JadwalKelasRepository) JadwalKelasService {
	return &jadwalKelasService{repo: repo}
}

func (s *jadwalKelasService) CreateJadwalKelas(req *dto.CreateJadwalKelasRequest) (*dto.JadwalKelasResponse, error) {
	exists, err := s.repo.CheckDuplicate(req.IDJadwal, req.IDKelas)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("kelas ini sudah terdaftar di jadwal tersebut")
	}

	jadwalKelas := &model.JadwalKelas{
		IDJadwal: req.IDJadwal,
		IDKelas:  req.IDKelas,
	}

	if err := s.repo.Create(jadwalKelas); err != nil {
		return nil, err
	}

	created, err := s.repo.GetByIDWithDetail(jadwalKelas.ID)
	if err != nil {
		return nil, err
	}

	return detailToResponse(created), nil
}

func (s *jadwalKelasService) GetJadwalKelasByID(id string) (*dto.JadwalKelasResponse, error) {
	result, err := s.repo.GetByIDWithDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	return detailToResponse(result), nil
}

func (s *jadwalKelasService) GetAllJadwalKelas(page, pageSize int, idJadwal string, idKelas string) (*dto.JadwalKelasListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	results, total, err := s.repo.GetAllWithDetail(page, pageSize, idJadwal, idKelas)
	if err != nil {
		return nil, err
	}

	responses := []dto.JadwalKelasResponse{}
	for _, r := range results {
		responses = append(responses, *detailToResponse(&r))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.JadwalKelasListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *jadwalKelasService) UpdateJadwalKelas(id string, req *dto.UpdateJadwalKelasRequest) (*dto.JadwalKelasResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	if req.IDJadwal != existing.IDJadwal || req.IDKelas != existing.IDKelas {
		exists, err := s.repo.CheckDuplicate(req.IDJadwal, req.IDKelas)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("kelas ini sudah terdaftar di jadwal tersebut")
		}
	}

	existing.IDJadwal = req.IDJadwal
	existing.IDKelas  = req.IDKelas

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByIDWithDetail(id)
	if err != nil {
		return nil, err
	}

	return detailToResponse(updated), nil
}

func (s *jadwalKelasService) DeleteJadwalKelas(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}
	return s.repo.Delete(id)
}

func detailToResponse(r *repository.JadwalKelasWithDetail) *dto.JadwalKelasResponse {
	return &dto.JadwalKelasResponse{
		ID:           r.ID,
		IDJadwal:     r.IDJadwal,
		IDKelas:      r.IDKelas,
		NamaKelas:    r.NamaKelas,
		NamaBankSoal: r.NamaBankSoal,
		WktMulai:     r.WktMulai,
		WktSelesai:   r.WktSelesai,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
