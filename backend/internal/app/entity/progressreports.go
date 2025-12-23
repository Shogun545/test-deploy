package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProgressReport struct {
	gorm.Model
	Body string `gorm:"type:text" valid:"required~Body is required"`
	Status string `gorm:"type:varchar(50)" valid:"required~Status is required"`
	AdvisorLogsID uint `gorm:"not null" valid:"required~AdvisorLogsID is required"`
	SubmittedAt time.Time
	FileName string `gorm:"type:varchar(255)"`
	FilePath string `gorm:"type:varchar(255)"`
	Feedbacks []ReportFeedback `gorm:"foreignKey:ProgressReportsID"`
}