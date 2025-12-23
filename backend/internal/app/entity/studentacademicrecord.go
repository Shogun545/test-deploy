package entity

import "gorm.io/gorm"

type StudentAcademicRecord struct {
	gorm.Model
	AcademicYear 		string 	`json:"academic_year" binding:"required"`
	Semester	 		int    	`json:"semester" binding:"required"`
	TermGPA	 			float32 `json:"term_gpa" binding:"required"`
	CumulativeGPA 		float32 `json:"cumulative_gpa" binding:"required"`
	AcademicStatus 		string 	`json:"academic_status" binding:"required"`

	// FK StudentProfile ที่เป็นเจ้าของ record
	StudentProfileID 	uint           	`json:"student_profile_id"`
	StudentProfile   	*StudentProfile `json:"student_profile" gorm:"foreignKey:StudentProfileID"`
}