package handler

import (
	"invoice-ocr-backend/internal/service"
	"net/http"
	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	service *service.InvoiceService
}

func NewInvoiceHandler(service *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{service}
}

/**
 * Upload - HTTP Entry Point for Invoice Processing
 *
 * Flow:
 * - รับไฟล์จาก multipart request
 * - ส่งต่อไป service layer (Upload → OCR → DB)
 * - คืนผลลัพธ์ (invoice + OCR data) เป็น JSON
 *
 * Failure:
 * - ไม่มีไฟล์ → 400 Bad Request
 * - internal error → 500 Internal Server Error
 */
func (h *InvoiceHandler) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	// รับไฟล์จาก multipart form (key: "file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	invoice, ocrResult, err := h.service.UploadInvoice(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"invoice":    invoice,
		"ocr_result": ocrResult,
	})
}
