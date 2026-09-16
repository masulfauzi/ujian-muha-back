package service

import (
	"math"

	"backend/internal/modules/login_log/dto"
	"backend/internal/modules/login_log/model"
	"backend/internal/modules/login_log/repository"
)

// RecordLoginInput adalah data yang direkam untuk satu percobaan login,
// dikumpulkan oleh AuthController dari request HTTP (header, IP) dan hasil
// pemanggilan AuthService.Login (berhasil/gagal).
type RecordLoginInput struct {
	Username       string
	IDUser         string
	Role           string
	Success        bool
	Keterangan     string
	UserAgent      string
	XRequestedWith string
	IPAddress      string
}

type LoginLogService interface {
	RecordLogin(input *RecordLoginInput) error
	GetAllLoginLog(page, pageSize int, username string) (*dto.LoginLogListResponse, error)
}

type loginLogService struct {
	repo repository.LoginLogRepository
}

func NewLoginLogService(repo repository.LoginLogRepository) LoginLogService {
	return &loginLogService{repo: repo}
}

func (s *loginLogService) RecordLogin(input *RecordLoginInput) error {
	successInt := 0
	if input.Success {
		successInt = 1
	}

	logEntry := &model.LoginLog{
		Username:       input.Username,
		IDUser:         input.IDUser,
		Role:           input.Role,
		Success:        successInt,
		Keterangan:     input.Keterangan,
		UserAgent:      input.UserAgent,
		XRequestedWith: input.XRequestedWith,
		IPAddress:      input.IPAddress,
	}

	return s.repo.Create(logEntry)
}

func (s *loginLogService) GetAllLoginLog(page, pageSize int, username string) (*dto.LoginLogListResponse, error) {
	logs, total, err := s.repo.GetAll(page, pageSize, username)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	responses := []dto.LoginLogResponse{}
	for _, l := range logs {
		responses = append(responses, dto.LoginLogResponse{
			ID:             l.ID,
			Username:       l.Username,
			IDUser:         l.IDUser,
			Role:           l.Role,
			Success:        l.Success == 1,
			Keterangan:     l.Keterangan,
			UserAgent:      l.UserAgent,
			XRequestedWith: l.XRequestedWith,
			IPAddress:      l.IPAddress,
			CreatedAt:      l.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.LoginLogListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}
