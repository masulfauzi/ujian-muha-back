package service

import (
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/kelas/dto"
	"backend/internal/modules/kelas/model"
	"backend/internal/modules/kelas/repository"

	"gorm.io/gorm"
)

// Helper function to convert KelasWithJurusan to KelasResponse
func kelasWithJurusanToResponse(k *repository.KelasWithJurusan) *dto.KelasResponse {
	return &dto.KelasResponse{
		ID:          k.ID,
		IDJurusan:   k.IDJurusan,
		NamaKelas:   k.NamaKelas,
		Tingkat:     k.Tingkat,
		NamaJurusan: k.NamaJurusan,
		CreatedAt:   k.CreatedAt,
		UpdatedAt:   k.UpdatedAt,
	}
}

type KelasService interface {
	CreateKelas(req *dto.CreateKelasRequest) (*dto.KelasResponse, error)
	GetKelasByID(id string) (*dto.KelasResponse, error)
	GetAllKelas(page, pageSize int, idJurusan string, tingkat string) (*dto.KelasListResponse, error)
	UpdateKelas(id string, req *dto.UpdateKelasRequest) (*dto.KelasResponse, error)
	DeleteKelas(id string) error
	RestoreKelas(id string) error
}

type kelasService struct {
	repo repository.KelasRepository
}

func NewKelasService(repo repository.KelasRepository) KelasService {
	return &kelasService{repo: repo}
}

func (s *kelasService) CreateKelas(req *dto.CreateKelasRequest) (*dto.KelasResponse, error) {
	kelas := &model.Kelas{
		IDJurusan: req.IDJurusan,
		NamaKelas: req.NamaKelas,
		Tingkat:   req.Tingkat,
	}

	if err := s.repo.Create(kelas); err != nil {
		return nil, err
	}

	// Fetch the created kelas with jurusan info via JOIN
	createdKelas, err := s.repo.GetByID(kelas.ID)
	if err != nil {
		return nil, err
	}

	return kelasWithJurusanToResponse(createdKelas), nil
}

func (s *kelasService) GetKelasByID(id string) (*dto.KelasResponse, error) {
	kelas, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	return kelasWithJurusanToResponse(kelas), nil
}

func (s *kelasService) GetAllKelas(page, pageSize int, idJurusan string, tingkat string) (*dto.KelasListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	kelasList, total, err := s.repo.GetAll(page, pageSize, idJurusan, tingkat)
	if err != nil {
		return nil, err
	}

	var responses []dto.KelasResponse
	for _, k := range kelasList {
		responses = append(responses, *kelasWithJurusanToResponse(&k))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.KelasListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *kelasService) UpdateKelas(id string, req *dto.UpdateKelasRequest) (*dto.KelasResponse, error) {
	kelasWithJurusan, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	kelas := &model.Kelas{
		ID:        kelasWithJurusan.ID,
		IDJurusan: req.IDJurusan,
		NamaKelas: req.NamaKelas,
		Tingkat:   req.Tingkat,
	}

	if err := s.repo.Update(kelas); err != nil {
		return nil, err
	}

	// Fetch the updated kelas with jurusan info via JOIN
	updatedKelas, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return kelasWithJurusanToResponse(updatedKelas), nil
}

func (s *kelasService) DeleteKelas(id string) error {
	kelas, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}

	return s.repo.Delete(kelas.ID)
}

func (s *kelasService) RestoreKelas(id string) error {
	return s.repo.Restore(id)
}

