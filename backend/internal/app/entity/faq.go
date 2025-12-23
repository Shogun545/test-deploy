package entity

import (
	"gorm.io/gorm"
)

type FAQStatus string

const (
	StatusHide   FAQStatus = "Hide"
	StatusUnhide FAQStatus = "Unhide"
)

type FAQ struct {
	gorm.Model

	FAQQuestion string    `json:"faq_question" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:varchar(1024);not null"`

	FAQstatus   FAQStatus `json:"faq_status" gorm:"type:varchar(20);default:'Unhide'"`

	ViewCount   uint      `json:"view_count" gorm:"type:integer;default:0"`

	CraeteBy uint  `json:"create_by"`
	UserID   *User `json:"user_id" gorm:"foreignKey:CraeteBy"`

	FaqTopic           uint              `json:"faq_topic"`
	AppointmentTopicID *AppointmentTopic `json:"appointmenttopic_id" gorm:"foreignKey:FaqTopic"`
}
