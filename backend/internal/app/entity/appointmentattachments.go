package entity

import (
	"gorm.io/gorm"
)

type AppointmentAttachment struct {
	gorm.Model

	FileAssetsID uint        `json:"file_assets_id"`
	FileAssets   *FileAssets `json:"file_assets" gorm:"foreignKey:FileAssetsID"`

	AppointmentID uint         `json:"appointment_id"`
	Appointment   *Appointment `json:"appointment" gorm:"foreignKey:AppointmentID"`
}
