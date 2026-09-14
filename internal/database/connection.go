package database

import (
	"fmt"

	"backend/configs"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	dbConfig := configs.GetDatabaseConfig()
	db, err := dbConfig.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
