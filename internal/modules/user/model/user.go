package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	Username  string    `gorm:"column:username;not null;uniqueIndex" json:"username"`
	Email     string    `gorm:"column:email;not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"column:password;not null" json:"-"`
	Role      string    `gorm:"column:role;not null;default:'user'" json:"role"`
	Status    string    `gorm:"column:status;not null;default:'active'" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

func (u *User) TableName() string {
	return "users"
}
