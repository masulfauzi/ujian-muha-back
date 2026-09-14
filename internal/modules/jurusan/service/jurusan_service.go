package service

import (
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/jurusan/dto"
	"backend/internal/modules/jurusan/model"
	"backend/internal/modules/jurusan/repository"

	"gorm.io/gorm"
)

type JurusanService interface {
	CreateJurusan(req *dto.CreateJurusanRequest) (*dto.JurusanResponse, error)
	GetJurusanByID(id string) (*dto.JurusanResponse, error)
	GetAllJurusan(page, pageSize int) (*dto.JurusanListResponse, error)
	UpdateJurusan(id string, req *dto.UpdateJurusanRequest) (*dto.JurusanResponse, error)
	DeleteJurusan(id string) error
	RestoreJurusan(id string) error
}

type jurusanService struct {
	repo repository.JurusanRepository
}

func NewJurusanService(repo repository.JurusanRepository) JurusanService {
	return &jurusanService{repo: repo}
}

func (s *jurusanService) CreateJurusan(req *dto.CreateJurusanRequest) (*dto.JurusanResponse, error) {
	jurusan := &model.Jurusan{
		NamaJurusan: req.NamaJurusan,
	}

	if err := s.repo.Create(jurusan); err != nil {
		return nil, err
	}

	return s.modelToResponse(jurusan), nil
}

func (s *jurusanService) GetJurusanByID(id string) (*dto.JurusanResponse, error) {
	jurusan, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	return s.modelToResponse(jurusan), nil
}

func (s *jurusanService) GetAllJurusan(page, pageSize int) (*dto.JurusanListResponse, error) {
	jurusans, total, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var responses []dto.JurusanResponse
	for _, j := range jurusans {
		responses = append(responses, *s.modelToResponse(&j))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.JurusanListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *jurusanService) UpdateJurusan(id string, req *dto.UpdateJurusanRequest) (*dto.JurusanResponse, error) {
	jurusan, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	jurusan.NamaJurusan = req.NamaJurusan

	if err := s.repo.Update(jurusan); err != nil {
		return nil, err
	}

	return s.modelToResponse(jurusan), nil
}

func (s *jurusanService) DeleteJurusan(id string) error {
	jurusan, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}

	return s.repo.Delete(jurusan.ID)
}

func (s *jurusanService) RestoreJurusan(id string) error {
	return s.repo.Restore(id)
}

func (s *jurusanService) modelToResponse(jurusan *model.Jurusan) *dto.JurusanResponse {
	return &dto.JurusanResponse{
		ID:          jurusan.ID,
		NamaJurusan: jurusan.NamaJurusan,
		CreatedAt:   jurusan.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   jurusan.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
