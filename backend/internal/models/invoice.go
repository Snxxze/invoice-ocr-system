package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invoice struct {
	ID		string		`gorm:"primaryKey" json:"id"`
	FileUrl		string		`json:"file_url"`
	Status		string		`json:"status"`
	InvoiceNo		*string		`json:"invoice_no"`
	Vender			*string		`json:"vender"`
	Total				*string		`json:"total"`
	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
}

// Helper สร้าง UUID
func (i *Invoice) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = uuid.NewString()
	
	return
}