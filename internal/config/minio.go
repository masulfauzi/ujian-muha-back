package config

import (
	"os"
	"strings"
)

type MinioConfig struct {
	Endpoint  string // host:port tanpa scheme, untuk MinIO SDK
	PublicURL string // URL lengkap dengan scheme, untuk generate public link
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func GetMinioConfig() MinioConfig {
	endpoint := strings.TrimRight(os.Getenv("MINIO_ENDPOINT"), "/")
	if endpoint == "" {
		endpoint = "http://localhost:9000"
	}

	publicURL := endpoint

	// SDK membutuhkan host:port tanpa scheme
	sdkEndpoint := endpoint
	sdkEndpoint = strings.TrimPrefix(sdkEndpoint, "https://")
	sdkEndpoint = strings.TrimPrefix(sdkEndpoint, "http://")

	useSSL := strings.ToLower(os.Getenv("MINIO_USE_SSL")) == "true"

	bucket := os.Getenv("MINIO_BUCKET")
	if bucket == "" {
		bucket = "gambar-soal"
	}

	return MinioConfig{
		Endpoint:  sdkEndpoint,
		PublicURL: publicURL,
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		Bucket:    bucket,
		UseSSL:    useSSL,
	}
}
