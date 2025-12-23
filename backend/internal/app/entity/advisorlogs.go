package entity

import (
	"gorm.io/gorm"
)

type AdvisorLog struct {
	gorm.Model

	AppointmentID  uint   `gorm:"not null;uniqueIndex" valid:"required~appointmentId is required"`
	Title          string `gorm:"type:varchar(255)" valid:"required~title is required"`
	Body           string `gorm:"type:text" valid:"required~body is required"`
	Status         string `gorm:"type:varchar(50)" valid:"required~status is required"`
	RequiresReport bool   `gorm:"not null;default:false"`
	FileName string `gorm:"type:text"`
	FilePath string `gorm:"type:text"`

	Appointment *Appointment `gorm:"foreignKey:AppointmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}