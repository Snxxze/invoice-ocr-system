package main

import (
	"invoice-ocr-backend/internal/config"
	"invoice-ocr-backend/internal/handler"
	"invoice-ocr-backend/internal/repository"
	"invoice-ocr-backend/internal/service"
	"invoice-ocr-backend/pkg/ocr"
	"invoice-ocr-backend/pkg/storage"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	cfg := config.LoadConfig()

	// Setup Database
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Setup MinIO
	minioClient, err := storage.InitMinio(
		cfg.StorageEndpoint,
		cfg.StorageAccessKey,
		cfg.StorageSecretKey,
		cfg.StorageBucket,
		cfg.StorageUseSSL,
	)
	if err != nil {
		panic("Failed to connect to MinIO: " + err.Error())
	}

	// Setup OCR Client
	ocrClient := ocr.NewClient(cfg.OCRServiceURL, 120*time.Second)

	// Repo -> Service -> Handler
	invoiceRepo := repository.NewInvoiceRepository(db)
	invoiceService := service.NewInvoiceService(invoiceRepo, minioClient, ocrClient)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong! Database is connected",
		})
	})

	// Route
	r.POST("/upload", invoiceHandler.Upload)

	r.Run(":8080")
}
