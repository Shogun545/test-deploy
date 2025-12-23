package entity

import (
    "time"

    "gorm.io/gorm"
)

// enum type (ใน Go)
type ReminderStatus string

// enum values
const (
    ReminderPending   ReminderStatus = "pending"
    ReminderSent      ReminderStatus = "sent"
    ReminderCancelled ReminderStatus = "cancelled"
    ReminderFailed    ReminderStatus = "failed"
)

type AppointmentReminder struct {
    gorm.Model

    AppointmentID uint        `json:"appointment_id"`
    Appointment   Appointment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"appointment"`

    RecipientUserID uint `json:"recipient_user_id"`
    RecipientUser   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"recipient_user"`

    OffsetDays         int `json:"offset_days"`
    OffsetHours        int `json:"offset_hours"`
    OffsetMinutes      int `json:"offset_minutes"`
    TotalOffsetMinutes int `json:"total_offset_minutes"`

    ScheduledAt   time.Time      `json:"scheduled_at"`
    Status        ReminderStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
    AttemptCount  int            `json:"attempt_count"`
    LastAttemptAt *time.Time     `json:"last_attempt_at"`

    IsActive bool `gorm:"default:true" json:"is_active"`

}
