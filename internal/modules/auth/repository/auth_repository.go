package repository

import (
	"backend/internal/modules/peserta/model"
	usermodel "backend/internal/modules/user/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	GetUserByEmail(email string) (*usermodel.User, error)
	GetUserByUsername(username string) (*usermodel.User, error)
	GetPesertaByUsername(username string) (*model.Peserta, error)
	CreateUser(user *usermodel.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) GetUserByEmail(email string) (*usermodel.User, error) {
	var user usermodel.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByUsername(username string) (*usermodel.User, error) {
	var user usermodel.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetPesertaByUsername(username string) (*model.Peserta, error) {
	var peserta model.Peserta
	err := r.db.Where("username = ?", username).First(&peserta).Error
	if err != nil {
		return nil, err
	}
	return &peserta, nil
}

func (r *authRepository) CreateUser(user *usermodel.User) error {
	return r.db.Create(user).Error
}
