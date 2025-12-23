package entity

import (
	"gorm.io/gorm"
)

type ReportFeedback struct {
	gorm.Model
	ProgressReportsID uint `gorm:"not null" valid:"required~ProgressReportsID is required"`
	Body string `gorm:"type:text" valid:"required~Body is required"`
}