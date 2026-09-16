package service

import (
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/user_agent/dto"
	"backend/internal/modules/user_agent/model"
	"backend/internal/modules/user_agent/repository"

	"gorm.io/gorm"
)

type UserAgentService interface {
	CreateUserAgent(req *dto.CreateUserAgentRequest) (*dto.UserAgentResponse, error)
	GetUserAgentByID(id string) (*dto.UserAgentResponse, error)
	GetAllUserAgent(page, pageSize int) (*dto.UserAgentListResponse, error)
	UpdateUserAgent(id string, req *dto.UpdateUserAgentRequest) (*dto.UserAgentResponse, error)
	DeleteUserAgent(id string) error
	RestoreUserAgent(id string) error
	// IsAllowed mengecek apakah header User-Agent ATAU X-Requested-With cocok
	// dengan salah satu baris whitelist manapun (tidak harus baris yang sama).
	// Dipakai middleware.CheckAllowedUserAgent.
	IsAllowed(userAgent, xRequestedWith string) (bool, error)
}

type userAgentService struct {
	repo repository.UserAgentRepository
}

func NewUserAgentService(repo repository.UserAgentRepository) UserAgentService {
	return &userAgentService{repo: repo}
}

func (s *userAgentService) CreateUserAgent(req *dto.CreateUserAgentRequest) (*dto.UserAgentResponse, error) {
	exists, err := s.repo.Exists(req.UserAgent, req.XRequestedWith)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("kombinasi user agent & x-requested-with sudah terdaftar")
	}

	userAgent := &model.UserAgent{
		UserAgent:      req.UserAgent,
		XRequestedWith: req.XRequestedWith,
		Keterangan:     req.Keterangan,
	}

	if err := s.repo.Create(userAgent); err != nil {
		return nil, err
	}

	return s.modelToResponse(userAgent), nil
}

func (s *userAgentService) GetUserAgentByID(id string) (*dto.UserAgentResponse, error) {
	userAgent, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	return s.modelToResponse(userAgent), nil
}

func (s *userAgentService) GetAllUserAgent(page, pageSize int) (*dto.UserAgentListResponse, error) {
	userAgents, total, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	responses := []dto.UserAgentResponse{}
	for _, ua := range userAgents {
		responses = append(responses, *s.modelToResponse(&ua))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.UserAgentListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *userAgentService) UpdateUserAgent(id string, req *dto.UpdateUserAgentRequest) (*dto.UserAgentResponse, error) {
	userAgent, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	if req.UserAgent != userAgent.UserAgent || req.XRequestedWith != userAgent.XRequestedWith {
		exists, err := s.repo.Exists(req.UserAgent, req.XRequestedWith)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("kombinasi user agent & x-requested-with sudah terdaftar")
		}
	}

	userAgent.UserAgent = req.UserAgent
	userAgent.XRequestedWith = req.XRequestedWith
	userAgent.Keterangan = req.Keterangan

	if err := s.repo.Update(userAgent); err != nil {
		return nil, err
	}

	return s.modelToResponse(userAgent), nil
}

func (s *userAgentService) DeleteUserAgent(id string) error {
	userAgent, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}

	return s.repo.Delete(userAgent.ID)
}

func (s *userAgentService) RestoreUserAgent(id string) error {
	return s.repo.Restore(id)
}

func (s *userAgentService) IsAllowed(userAgent, xRequestedWith string) (bool, error) {
	return s.repo.IsAllowed(userAgent, xRequestedWith)
}

func (s *userAgentService) modelToResponse(userAgent *model.UserAgent) *dto.UserAgentResponse {
	return &dto.UserAgentResponse{
		ID:             userAgent.ID,
		UserAgent:      userAgent.UserAgent,
		XRequestedWith: userAgent.XRequestedWith,
		Keterangan:     userAgent.Keterangan,
		CreatedAt:      userAgent.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      userAgent.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
