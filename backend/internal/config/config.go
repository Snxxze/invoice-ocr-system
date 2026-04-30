package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	StorageEndpoint  string
	StorageAccessKey string
	StorageSecretKey string
	StorageBucket    string
	StorageUseSSL    bool
	OCRServiceURL    string
}

func LoadConfig() *Config {
	// โหลดไฟล์ .env
	_ = godotenv.Load()

	config := &Config{
		DBHost:     getEnv("DB_HOST", ""),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),
		DBPort:     getEnv("DB_PORT", "5432"),
		StorageEndpoint:  getEnv("STORAGE_ENDPOINT", ""),
    StorageAccessKey: getEnv("STORAGE_ACCESS_KEY", ""),
    StorageSecretKey: getEnv("STORAGE_SECRET_KEY", ""),
    StorageBucket:    getEnv("STORAGE_BUCKET", ""),
    StorageUseSSL:    getEnvAsBool("STORAGE_USE_SSL", false),
		OCRServiceURL: getEnv("OCR_SERVICE_URL", "http://localhost:8000/extract"),
	}

	if config.DBHost == "" || config.DBPassword == "" || config.StorageAccessKey == "" || config.StorageSecretKey == "" {
		log.Fatal("Environment variables are missing")
	}

	if config.OCRServiceURL == "" {
		log.Fatal("OCR_SERVICE_URL is required")
	}

	return config
}

// Helper ดึงค่า .env
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper เพิ่มเติมสำหรับอ่านค่าเป็น Boolean
func getEnvAsBool(key string, fallback bool) bool {
	val := getEnv(key, "")
	if val == "" { return fallback }
	return val == "true"
}