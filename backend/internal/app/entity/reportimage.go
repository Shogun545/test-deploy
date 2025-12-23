package entity

import (
	"gorm.io/gorm"
)

type ReportImage struct {
	gorm.Model

	FileURL  string `json:"file_url" gorm:"type:varchar(255);not null"`
	FileType string `json:"file_type" gorm:"type:varchar(50);not null"`

	ReportID uint     `json:"report_id"`
	Report   *Report  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
