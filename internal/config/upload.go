package config

import (
	"os"
	"strings"
)

type UploadConfig struct {
	ImageUploadPath string
	MaxFileSize     int64
	AllowedFormats  []string
	ImageBaseURL    string
}

func GetUploadConfig() UploadConfig {
	minioCfg := GetMinioConfig()
	imageBaseURL := strings.TrimRight(minioCfg.PublicURL, "/") + "/" + minioCfg.Bucket

	return UploadConfig{
		ImageUploadPath: "./uploads/soal",
		MaxFileSize:     5 * 1024 * 1024,
		AllowedFormats:  []string{"jpg", "jpeg", "png", "gif", "webp"},
		ImageBaseURL:    imageBaseURL,
	}
}

func GetSoalImageSourceURL() string {
	url := os.Getenv("SOAL_IMAGE_SOURCE_URL")
	if url == "" {
		url = "https://apps.smkn2semarang.sch.id/ujian/soal/"
	}
	return strings.TrimRight(url, "/") + "/"
}
