package entity

import (
    "time"

    "gorm.io/gorm"
)

type NotificationsHistory struct {
    gorm.Model

    // FK ไป Appointment
    AppointmentID uint        `json:"appointment_id"`
    Appointment   Appointment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"appointment"`

    // FK ไป User (ผู้รับ)
    RecipientUserID uint `json:"recipient_user_id"`
    RecipientUser   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"recipient_user"`

    // ใช้ Go enum + varchar แทน enum ของ DB
    EventType NotificationEventType `gorm:"type:varchar(20)" json:"event_type"`

    TitleSnapshot   string `gorm:"type:varchar(255)" json:"title_snapshot"`
    MessageSnapshot string `gorm:"type:text" json:"message_snapshot"`

    // ถ้าอยากเก็บเวลาที่สร้าง snapshot จริงๆ
    CreatedAt time.Time `json:"created_at"` // หรือจะใช้ของ gorm.Model เฉยๆ ก็ได้ แล้วลบอันนี้ออก
}
