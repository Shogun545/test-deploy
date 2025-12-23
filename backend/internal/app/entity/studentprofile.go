package entity

import "gorm.io/gorm"

type StudentProfile struct {
    gorm.Model

    YearOfStudy int `json:"year_of_study"`

    // FK user ที่เป็นเจ้าของ profile
    UserID uint  `json:"user_id"`
    User   *User `json:"user" gorm:"foreignKey:UserID"`

	//StudentProfile มี advisor คนเดียว *AdvisorProfile
    AdvisorProfileID *uint           `json:"advisor_profile_id"`
    AdvisorProfile   *AdvisorProfile `json:"advisor_profile" gorm:"foreignKey:AdvisorProfileID"`

    StudentAcademicRecords []StudentAcademicRecord `gorm:"foreignKey:StudentProfileID"`

}
