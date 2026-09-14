package configs

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func GetDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "fiber_backend"),
	}
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
	)
}

func (c *DatabaseConfig) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(c.GetDSN()), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
			},
		),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
