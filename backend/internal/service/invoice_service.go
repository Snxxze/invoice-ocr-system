package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"invoice-ocr-backend/internal/config"
	"invoice-ocr-backend/internal/models"
	"invoice-ocr-backend/internal/repository"
	"invoice-ocr-backend/pkg/storage"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
)

type InvoiceService struct {
	repo    *repository.InvoiceRepository
	storage *storage.MinioClient
	cfg 		*config.Config
}

func NewInvoiceService(repo *repository.InvoiceRepository, storage *storage.MinioClient, cfg *config.Config) *InvoiceService {
	return &InvoiceService{repo, storage, cfg}
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
func (s *InvoiceService) UploadInvoice(file *multipart.FileHeader) (*models.Invoice, interface{}, error) {
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

	// เตรียมข้อมูลสำหรับ MinIO (ใช้ Buffer ที่มีอยู่)
	data := buf.Bytes()
	reader := bytes.NewReader(data)

	// Detect ContentType จาก bytes ตรงๆ
	contentType := http.DetectContentType(data)
	opts := minio.PutObjectOptions{ContentType: contentType}

	objectName := storage.GenerateFileName(file.Filename)

	// อัปโหลดขึ้น MinIO
	_, err = s.storage.Client.PutObject(ctx, s.storage.BucketName, objectName, reader, int64(len(data)), opts)
	if err != nil {
		return nil, nil, err
	}

	// เรียก OCR Service โดยใช้ Data ชุดเดียวกัน
	ocrResult, err := s.CallOCR(data, file.Filename)
	if err != nil {
		return nil, nil, fmt.Errorf("OCR failed: %v", err)
	}

	invoice := &models.Invoice{
		FileUrl: objectName,
		Status:	"processed",
	}

	if err := s.repo.Create(invoice); err != nil {
		return nil, nil, err
	}

	return invoice, ocrResult, nil
}

/**
 * CallOCR - External OCR Integration
 *
 * Flow:
 * - สร้าง multipart request (แนบไฟล์)
 * - ส่งไป OCR service
 * - รับและแปลง response เป็น JSON
 *
 * Design:
 * - ส่ง data เป็น byte (ไม่ reopen file / ลด I/O)
 * - ใช้ timeout กัน request ค้าง (network / OCR ช้า)
 *
 * Contract:
 * - OCR service ต้องรับ field "file" และตอบ JSON
 *
 * Failure:
 * - non-200 → treat เป็น error (fail-fast)
 * - ไม่มี retry → อาจเพิ่มในอนาคต (transient error)
 *
 * Insight:
 * - เป็น synchronous + network-bound call → latency ขึ้นกับ OCR service
 * - ควรแยกเป็น async job เมื่อ scale
 */
func (s *InvoiceService) CallOCR(data []byte, filename string) (interface{}, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err := part.Write(data); err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", s.cfg.OCRServiceURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OCR service error: %s", resp.Status)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}