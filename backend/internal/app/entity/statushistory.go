package entity

import (
    "time"

    "gorm.io/gorm"
)

type StatusHistory struct {
    gorm.Model

    // นัดหมายที่ถูกเปลี่ยนสถานะ
    AppointmentID uint        `json:"appointment_id"`
    Appointment   Appointment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"appointment"`

    // คนที่เปลี่ยนสถานะ (อาจารย์หรือนักศึกษา)
    ChangedByUserID uint `json:"changed_by_user_id"`
    ChangedByUser   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"changed_by_user"`

    // จากสถานะเดิม → ไปสถานะใหม่
    FromStatusID uint              `json:"from_status_id"`
    FromStatus   AppointmentStatus `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"from_status"`

    ToStatusID uint              `json:"to_status_id"`
    ToStatus   AppointmentStatus `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"to_status"`

    // Action ที่ทำ เช่น APPROVE / REJECT / RESCHEDULE ฯลฯ
    ActionID uint          `json:"action_id"`
    Action   ApprovalAction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"action"`

    Reason    string    `gorm:"type:varchar(255)" json:"reason"`
    ChangedAt time.Time `json:"changed_at"`
}
