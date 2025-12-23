package entity

import (
	"gorm.io/gorm"
)

type Report struct {
	gorm.Model

	Description string `json:"description" valid:"required~description is required"`

	ReportByID uint  `json:"report_by_id" valid:"required~reportById is required"`
	User       *User `json:"user" gorm:"foreignKey:ReportByID"`

	ReportStatusID uint          `json:"report_status_id" valid:"required~reportStatusId is required"`
	Status         *ReportStatus `json:"status" gorm:"foreignKey:ReportStatusID"`

	ReportTopicID uint         `json:"report_topic_id" valid:"required~reportTopicId is required"` 
	Topic         *ReportTopic `json:"topic" gorm:"foreignKey:ReportTopicID"`
}
