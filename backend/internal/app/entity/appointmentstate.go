package entity
import ( "gorm.io/gorm"
		"time"

)

type AppointmentState struct {
	gorm.Model
AppointmentID uint        `json:"appointment_id"`
    Appointment   Appointment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"appointment"`
CurrentStatusID uint              `json:"current_status_id"`
    CurrentStatus   AppointmentStatus `json:"current_status"`
    StatusChangeAt time.Time `gorm:"autoCreateTime" json:"status_change_at"`
}
