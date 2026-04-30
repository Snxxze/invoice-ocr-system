package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client     *minio.Client
	BucketName string
}

func InitMinio(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*MinioClient, error) {
	// สร้าง Client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// ตรวจสอบ/สร้าง Bucket
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		fmt.Printf("Successfully created bucket: %s\n", bucketName)
	}

	return &MinioClient{
		Client:     client,
		BucketName: bucketName,
	}, nil
}

var allowedTypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/png":       true,
}

func (s *MinioClient) PrepareFile(file *multipart.FileHeader) (io.Reader, minio.PutObjectOptions, func(), error) {
	f, err := file.Open()
	if err != nil {
		return nil, minio.PutObjectOptions{}, nil, err
	}

	buffer := make([]byte, 512)
	n, _ := io.ReadFull(f, buffer)

	contentType := http.DetectContentType(buffer[:n])

	if !allowedTypes[contentType] {
		f.Close()
		return nil, minio.PutObjectOptions{}, nil, fmt.Errorf("unsupported file type: %s", contentType)
	}

	reader := io.MultiReader(
		bytes.NewReader(buffer[:n]),
		f,
	)

	cleanup := func() {
		f.Close()
	}

	return reader, minio.PutObjectOptions{
		ContentType: contentType,
	}, cleanup, nil
}

// GenerateFileName สร้างชื่อไฟล์ใหม่เป็น UUID เพื่อไม่ให้ชื่อไฟล์ซ้ำกัน
func GenerateFileName(original string) string {
	ext := strings.ToLower(filepath.Ext(original))
	id := uuid.New().String()
	
	return id + ext
}