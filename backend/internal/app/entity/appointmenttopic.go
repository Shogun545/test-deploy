package entity

import "gorm.io/gorm"

type AppointmentTopic struct {
	gorm.Model

	Topic       string `json:"topic" gorm:"type:varchar(100);not null"`
	Description string `json:"description" gorm:"type:text"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
}
