package service

import (
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/mapel/dto"
	"backend/internal/modules/mapel/model"
	"backend/internal/modules/mapel/repository"

	"gorm.io/gorm"
)

type MapelService interface {
	CreateMapel(req *dto.CreateMapelRequest) (*dto.MapelResponse, error)
	GetMapelByID(id string) (*dto.MapelResponse, error)
	GetAllMapel(page, pageSize int) (*dto.MapelListResponse, error)
	UpdateMapel(id string, req *dto.UpdateMapelRequest) (*dto.MapelResponse, error)
	DeleteMapel(id string) error
	RestoreMapel(id string) error
}

type mapelService struct {
	repo repository.MapelRepository
}

func NewMapelService(repo repository.MapelRepository) MapelService {
	return &mapelService{repo: repo}
}

func (s *mapelService) CreateMapel(req *dto.CreateMapelRequest) (*dto.MapelResponse, error) {
	mapel := &model.Mapel{
		NamaMapel: req.NamaMapel,
		KodeMapel: req.KodeMapel,
		Deskripsi: req.Deskripsi,
	}

	if err := s.repo.Create(mapel); err != nil {
		return nil, err
	}

	return s.modelToResponse(mapel), nil
}

func (s *mapelService) GetMapelByID(id string) (*dto.MapelResponse, error) {
	mapel, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	return s.modelToResponse(mapel), nil
}

func (s *mapelService) GetAllMapel(page, pageSize int) (*dto.MapelListResponse, error) {
	mapels, total, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var responses []dto.MapelResponse
	for _, mapel := range mapels {
		responses = append(responses, *s.modelToResponse(&mapel))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.MapelListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *mapelService) UpdateMapel(id string, req *dto.UpdateMapelRequest) (*dto.MapelResponse, error) {
	mapel, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	mapel.NamaMapel = req.NamaMapel
	mapel.KodeMapel = req.KodeMapel
	mapel.Deskripsi = req.Deskripsi

	if err := s.repo.Update(mapel); err != nil {
		return nil, err
	}

	return s.modelToResponse(mapel), nil
}

func (s *mapelService) DeleteMapel(id string) error {
	mapel, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}

	return s.repo.Delete(mapel.ID)
}

func (s *mapelService) RestoreMapel(id string) error {
	return s.repo.Restore(id)
}

func (s *mapelService) modelToResponse(mapel *model.Mapel) *dto.MapelResponse {
	return &dto.MapelResponse{
		ID:        mapel.ID,
		NamaMapel: mapel.NamaMapel,
		KodeMapel: mapel.KodeMapel,
		Deskripsi: mapel.Deskripsi,
		CreatedAt: mapel.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: mapel.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
