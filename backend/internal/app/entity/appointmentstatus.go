package entity
import  
	("gorm.io/gorm"

	)
    const (
    StatusPendingID    uint = 1 // รอพิจารณา
    StatusApprovedID   uint = 2 // อนุมัติแล้ว
    StatusRescheduleID uint = 3 // เสนอเวลาใหม่
)

type AppointmentStatus struct {
	gorm.Model
    
    StatusCode   string    `gorm:"type:varchar(100)" json:"status_code"`
    StatusName   string    `gorm:"type:varchar(100)" json:"status_name"`
    IsTerminal   bool      `gorm:"type:boolean" json:"is_terminal"`
    DisplayOrder int       `gorm:"type:int" json:"display_order"`

}
