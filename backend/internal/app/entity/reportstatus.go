package entity

import (
	"gorm.io/gorm"
)

type ReportStatus struct {
	gorm.Model
	ReportStatusName string `json:"reportstatus_name"`
}


