package service

import (
	"bytes"
	"context"
	"invoice-ocr-backend/internal/models"
	"invoice-ocr-backend/internal/repository"
	"invoice-ocr-backend/pkg/ocr"
	"invoice-ocr-backend/pkg/storage"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/minio/minio-go/v7"
)

type InvoiceService struct {
	repo      *repository.InvoiceRepository
	storage   *storage.MinioClient
	ocrClient *ocr.Client
}

func NewInvoiceService(repo *repository.InvoiceRepository, storage *storage.MinioClient, ocrClient *ocr.Client) *InvoiceService {
	return &InvoiceService{repo, storage, ocrClient}
}

/**
 * UploadInvoice - File Ingestion Pipeline
 *
 * Flow:
 * - Read file
 * - Buffer in memory
 * - Upload to MinIO
 * - Send to OCR service
 * - Save metadata to DB
 *
 * Design:
 * - ใช้ buffer กลางเพื่อ reuse data (ลด I/O และ stream bug)
 *   Trade-off: ใช้ RAM เพิ่ม (เหมาะกับไฟล์ขนาดเล็ก)
 *
 * - Detect Content-Type จาก content จริง (ไม่ trust client)
 * - ใช้ UUID เป็น filename (กัน collision / path attack)
 *
 * OCR:
 * - เรียกแบบ synchronous → latency ผูกกับ OCR service
 * - Future: เปลี่ยนเป็น async (queue + worker)
 *
 * Failure:
 * - ไม่มี transaction ครอบ MinIO + DB → อาจเกิด orphan file
 * - ปัจจุบัน fail-fast ถ้า OCR error
 *
 * Insight:
 * - read once, reuse ทั้ง pipeline
 * - แยก concern: storage / OCR / persistence ชัดเจน
 */
func (s *InvoiceService) UploadInvoice(file *multipart.FileHeader) (*models.Invoice, *models.OCRResponse, error) {
	ctx := context.Background()
	
	src, err := file.Open()
	if err != nil {
		return nil, nil, err
	}
	defer src.Close()

	buf := &bytes.Buffer{}
	if _,err := io.Copy(buf, src); err != nil {
		return nil, nil, err
	}

	data := buf.Bytes()
	reader := bytes.NewReader(data)

	contentType := http.DetectContentType(data)
	opts := minio.PutObjectOptions{ContentType: contentType}

	objectName := storage.GenerateFileName(file.Filename)

	_, err = s.storage.Client.PutObject(ctx, s.storage.BucketName, objectName, reader, int64(len(data)), opts)
	if err != nil {
		return nil, nil, err
	}

	ocrResult, err := s.ocrClient.Extract(data, file.Filename)
	if err != nil {
		return nil, nil, err
	}

	invoice := &models.Invoice{
		FileUrl: objectName,
		Status:  "processed",
	}
	
	if ocrResult != nil {
		total := ocrResult.Summary.Total
		invoice.Total = &total
	}

	if err := s.repo.Create(invoice); err != nil {
		return nil, nil, err
	}

	return invoice, ocrResult, nil
}