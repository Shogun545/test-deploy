package entity

import (
	"gorm.io/gorm"
)

type FileAssets struct {
	gorm.Model

	OriginalName   uint   `gorm:"primaryKey;autoIncrement;type:int" json:"original_name"`
	MimeType string `gorm:"type:varchar(100);not null" json:"mime_type"`
	FStoragePath		string `gorm:"type:text" json:"storage_path"`
}