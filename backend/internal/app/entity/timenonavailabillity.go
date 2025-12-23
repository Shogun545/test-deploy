package entity

import (
	"time"

	"gorm.io/gorm"
)

type TimeNonAvailabillity struct {
	gorm.Model
	TimeID 		uint		`json:"time_id" `
	Subjects	string		`json:"subjects"`
	SubjectsID	string		`json:"subjects_id"`
	StartTime	time.Time	`json:"start_time"`
	EndTime		time.Time	`json:"end_time"`

	// AdvisorNonAvailabillityID  		uint 	`json:"advisor_non_availabillity_id"`

	// AdvisorNonAvailabillity			AdvisorNonAvailabillity		`json:"advisor_non_availabillity" gorm:"foreignKey:AdvisorNonAvailabillityID"`
}