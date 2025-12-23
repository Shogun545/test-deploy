package entity

import (
	"gorm.io/gorm"
)

type ReportTopic struct {
	gorm.Model
	ReportTopicName string `json:"reporttopic_name"`
	Description     string `json:"description"`
}
