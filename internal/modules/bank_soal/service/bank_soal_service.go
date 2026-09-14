package service

import (
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/bank_soal/dto"
	"backend/internal/modules/bank_soal/model"
	"backend/internal/modules/bank_soal/repository"

	"gorm.io/gorm"
)

type BankSoalService interface {
	CreateBankSoal(req *dto.CreateBankSoalRequest) (*dto.BankSoalResponse, error)
	GetBankSoalByID(id string) (*dto.BankSoalResponse, error)
	GetAllBankSoal(page, pageSize int) (*dto.BankSoalListResponse, error)
	GetBankSoalByMapel(mapelID string, page, pageSize int) (*dto.BankSoalListResponse, error)
	UpdateBankSoal(id string, req *dto.UpdateBankSoalRequest) (*dto.BankSoalResponse, error)
	DeleteBankSoal(id string) error
	RestoreBankSoal(id string) error
}

type bankSoalService struct {
	repo repository.BankSoalRepository
}

func NewBankSoalService(repo repository.BankSoalRepository) BankSoalService {
	return &bankSoalService{repo: repo}
}

func (s *bankSoalService) CreateBankSoal(req *dto.CreateBankSoalRequest) (*dto.BankSoalResponse, error) {
	bankSoal := &model.BankSoal{
		NamaBankSoal: req.NamaBankSoal,
		IdMapel:      req.IdMapel,
		JmlSoal:      req.JmlSoal,
		Deskripsi:    req.Deskripsi,
	}

	if err := s.repo.Create(bankSoal); err != nil {
		return nil, err
	}

	return s.modelToResponse(bankSoal), nil
}

func (s *bankSoalService) GetBankSoalByID(id string) (*dto.BankSoalResponse, error) {
	bankSoal, err := s.repo.GetByIDWithMapel(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	return s.joinedToResponse(bankSoal), nil
}

func (s *bankSoalService) GetAllBankSoal(page, pageSize int) (*dto.BankSoalListResponse, error) {
	bankSoals, total, err := s.repo.GetAllWithMapel(page, pageSize)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	responses := []dto.BankSoalResponse{}
	for _, bs := range bankSoals {
		responses = append(responses, *s.joinedToResponse(&bs))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.BankSoalListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *bankSoalService) GetBankSoalByMapel(mapelID string, page, pageSize int) (*dto.BankSoalListResponse, error) {
	bankSoals, total, err := s.repo.GetByMapelIDWithMapel(mapelID, page, pageSize)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	responses := []dto.BankSoalResponse{}
	for _, bs := range bankSoals {
		responses = append(responses, *s.joinedToResponse(&bs))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.BankSoalListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *bankSoalService) UpdateBankSoal(id string, req *dto.UpdateBankSoalRequest) (*dto.BankSoalResponse, error) {
	bankSoal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	bankSoal.NamaBankSoal = req.NamaBankSoal
	bankSoal.IdMapel = req.IdMapel
	bankSoal.JmlSoal = req.JmlSoal
	bankSoal.Deskripsi = req.Deskripsi

	if err := s.repo.Update(bankSoal); err != nil {
		return nil, err
	}

	return s.modelToResponse(bankSoal), nil
}

func (s *bankSoalService) DeleteBankSoal(id string) error {
	bankSoal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}

	return s.repo.Delete(bankSoal.ID)
}

func (s *bankSoalService) RestoreBankSoal(id string) error {
	return s.repo.Restore(id)
}

func (s *bankSoalService) modelToResponse(bankSoal *model.BankSoal) *dto.BankSoalResponse {
	return &dto.BankSoalResponse{
		ID:           bankSoal.ID,
		NamaBankSoal: bankSoal.NamaBankSoal,
		IdMapel:      bankSoal.IdMapel,
		NamaMapel:    "",
		JmlSoal:      bankSoal.JmlSoal,
		Deskripsi:    bankSoal.Deskripsi,
		CreatedAt:    bankSoal.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    bankSoal.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *bankSoalService) joinedToResponse(bankSoal *repository.BankSoalWithMapel) *dto.BankSoalResponse {
	return &dto.BankSoalResponse{
		ID:           bankSoal.ID,
		NamaBankSoal: bankSoal.NamaBankSoal,
		IdMapel:      bankSoal.IdMapel,
		NamaMapel:    bankSoal.NamaMapel,
		JmlSoal:      bankSoal.JmlSoal,
		Deskripsi:    bankSoal.Deskripsi,
		CreatedAt:    bankSoal.CreatedAt,
		UpdatedAt:    bankSoal.UpdatedAt,
	}
}
