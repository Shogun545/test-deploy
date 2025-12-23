package entity

import "gorm.io/gorm"
import "time"

type AdvisorNonAvailabillity struct {
	gorm.Model
	// AdvisorNonAvailabillityID 	uint 			`json:"advisor_non_availabillity_id" gorm:"primaryKey;autoIncrement"`
	
	Description					string 			`json:"description"`
	TypeAvailabillity			string 			`json:"type_availabillity"`
	Day							time.Time 		`json:"day"`
	IsRecurring					bool			`json:"is_recurring"`		

	TimeID 						uint
	AdvisorID					uint 			`json:"advisor_id"`
	TimeNonAvailabillity 		TimeNonAvailabillity 	`json:"time_nonavailabillity" gorm:"foreignKey:TimeID"`

	User						User			`json:"user" gorm:"foreignKey:AdvisorID"`
}	