package dto

import "backend/internal/app/entity"

type StudentInfo struct {
	SutId       string  `json:"sut_id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	YearOfStudy int     `json:"year_of_study"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	MajorName   string  `json:"major_name"` // เพิ่ม Major
	GPALatest   float64 `json:"gpa_latest"` // เพิ่ม GPA ล่าสุด
	Birthday    string  `json:"birthday"`   // เพิ่ม วันเกิด
}

type AdvisorStudentsResponse struct {
	AdvisorSutId string        `json:"advisor_sut_id"`
	AdvisorName  string        `json:"advisor_name"`
	Students     []StudentInfo `json:"students"`
}

func NewAdvisorStudentsResponse(advisor *entity.User) AdvisorStudentsResponse {
	students := make([]StudentInfo, 0, len(advisor.AdvisorProfile.Students))
	for _, s := range advisor.AdvisorProfile.Students {
		// Major
		majorName := ""
		if s.User.Major != nil {
			majorName = s.User.Major.Major
		}

		// YearOfStudy
		year := 0
		if s.User.StudentProfile != nil {
			year = s.User.StudentProfile.YearOfStudy
		}

		// GPA ล่าสุด
		var gpa float64 = 0
		if s.User.StudentProfile != nil && len(s.User.StudentProfile.StudentAcademicRecords) > 0 {
			gpa = float64(s.User.StudentProfile.StudentAcademicRecords[0].CumulativeGPA)
		}

		students = append(students, StudentInfo{
			SutId:       s.User.SutId,
			FirstName:   s.User.FirstName,
			LastName:    s.User.LastName,
			YearOfStudy: year,
			Email:       s.User.Email,
			Phone:       s.User.Phone,
			MajorName:   majorName,
			GPALatest:   gpa,
			Birthday:    s.User.Birthday,
		})
	}

	return AdvisorStudentsResponse{
		AdvisorSutId: advisor.SutId,
		AdvisorName:  advisor.FirstName + " " + advisor.LastName,
		Students:     students,
	}
}

type AdvisorStudentDetailResponse struct {
	SutId       string  `json:"sut_id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	YearOfStudy int     `json:"year_of_study"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	MajorName   string  `json:"major_name"`
	GPALatest   float64 `json:"gpa_latest"`
	Birthday    string  `json:"birthday"`
}
