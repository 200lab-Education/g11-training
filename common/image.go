package common

import (
	"github.com/google/uuid"
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
}
