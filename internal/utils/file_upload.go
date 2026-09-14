package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"backend/internal/config"
	"backend/internal/storage"
)

func SaveImage(file *multipart.FileHeader, folder string) (string, error) {
	cfg := config.GetUploadConfig()

	if file.Size > cfg.MaxFileSize {
		return "", fmt.Errorf("file terlalu besar, maksimal %.2f MB", float64(cfg.MaxFileSize)/1024/1024)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	ext = strings.TrimPrefix(ext, ".")
	if !isAllowedFormat(ext, cfg.AllowedFormats) {
		return "", fmt.Errorf("format file tidak diizinkan: %s", ext)
	}

	return storage.UploadFile(context.Background(), folder, file)
}

func DeleteImage(filename string, folder string) error {
	if filename == "" {
		return nil
	}
	return storage.DeleteFile(context.Background(), folder, filename)
}

func isAllowedFormat(ext string, allowed []string) bool {
	for _, a := range allowed {
		if ext == a {
			return true
		}
	}
	return false
}
