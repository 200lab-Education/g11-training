package image

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	TbName          = "images"
	ProviderAWSS3   = "aws_s3"
	StatusUploaded  = "uploaded"
	StatusActivated = "activated"
	StatusDeleted   = "deleted"
)

type Image struct {
	Id              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	FileName        string    `json:"file_name"`
	FileUrl         string    `json:"file_url" gorm:"-"`
	FileSize        int       `json:"file_size"`
	FileType        string    `json:"file_type"`
	StorageProvider string    `json:"storage_provider"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (img *Image) SetCDNDomain(domain string) {
	img.FileUrl = fmt.Sprintf("%s/%s", domain, img.FileName)
}
