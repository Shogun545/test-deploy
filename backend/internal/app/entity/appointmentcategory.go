package entity

import (
	"gorm.io/gorm"
)

type AppointmentCategory struct {
	gorm.Model
	Category string ` json:"category" gorm:"type:varchar(100);not null"`
}