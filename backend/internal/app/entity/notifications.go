package entity

import (
    "time"

    "gorm.io/gorm"
)

// ENUM TYPE
type NotificationEventType string
type NotificationStatusSnapshot string

// ENUM VALUE
const (
    EventApproved    NotificationEventType = "APPROVED"
    EventRescheduled NotificationEventType = "RESCHEDULED"
    EventFollowup    NotificationEventType = "FOLLOWUP"
    EventUniEvent    NotificationEventType = "UNIEVENT"
)

const (
    StatusSnapshotApproved    NotificationStatusSnapshot = "APPROVED"
    StatusSnapshotRescheduled NotificationStatusSnapshot = "RESCHEDULED"
)

type Notification struct {
    gorm.Model

    // FK ไป Appointment
    AppointmentID uint        `json:"appointment_id"`
    Appointment   Appointment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"appointment"`

    // FK ผู้รับ (Recipient)
    RecipientUserID uint `json:"recipient_user_id"`
    RecipientUser   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"recipient_user"`

    // FK ผู้ส่ง (Sender)
    SenderUserID uint `json:"sender_user_id"`
    SenderUser   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"sender_user"`

    // ใช้ Go enum + varchar แทน DB enum (ปลอดภัยกับ Postgres)
    EventType      NotificationEventType      `gorm:"type:varchar(20)" json:"event_type"`
    StatusSnapshot NotificationStatusSnapshot `gorm:"type:varchar(20)" json:"status_snapshot"`

    Topic   string `gorm:"type:varchar(255)" json:"topic"`
    Message string `gorm:"type:text" json:"message"`

    SentAt  *time.Time `json:"sent_at"`
    IsRead  bool       `gorm:"default:false" json:"is_read"`
    ReadAt  *time.Time `json:"read_at"`

    IsDelete   bool       `gorm:"default:false" json:"is_delete"`
    DeleteAt   *time.Time `json:"delete_at"`
    RestoredAt *time.Time `json:"restored_at"`
}
