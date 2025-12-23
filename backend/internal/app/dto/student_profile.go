package dto

import "backend/internal/app/entity"

type StudentProfileResponse struct {
	// ข้อมูลพื้นฐานจาก User
	SutId        string `json:"sut_id"`
	Prefix       string `json:"prefix"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	NationalId   string `json:"national_id"`
	Birthday     string `json:"birthday"`
	ProfileImage string `json:"profile_image"`

	// Major / Year
	MajorName   string `json:"major_name"`
	YearOfStudy int    `json:"year_of_study"`

	// Advisor
	AdvisorName string `json:"advisor_name"`

	// Academic (ล่าสุด)
	AcademicYear   string  `json:"academic_year"`
	Semester       int     `json:"semester"`
	TermGPA        float32 `json:"term_gpa"`
	CumulativeGPA  float32 `json:"cumulative_gpa"`
	AcademicStatus string  `json:"academic_status"`

	// สำหรับหน้า UI profile
	GPX          float32 `json:"gpx"`
	GPALatest    float32 `json:"gpa_latest"`
	GPATermLabel string  `json:"gpa_term_label"`
}

type UpdateStudentProfileRequest struct {
	// ฟิลด์ของ User
	Prefix       string `json:"prefix,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Phone        string `json:"phone"`
	Birthday     string `json:"birthday"`
	ProfileImage string `json:"profile_image"`
	Email        string `json:"email"`

	// ฟิลด์ของ StudentProfile (ถ้าอยากให้แก้)
	YearOfStudy *int `json:"year_of_study,omitempty"`

	// ปรับปรุงผลการเรียนล่าสุด (ถ้ามี record ล่าสุดอยู่แล้ว)
	AcademicYear  string   `json:"academic_year,omitempty"` 
    Semester      *int     `json:"semester,omitempty"`     
	TermGPA       *float32 `json:"term_gpa,omitempty"`
	CumulativeGPA *float32 `json:"cumulative_gpa,omitempty"`
}

// ฟังก์ชันช่วยสร้าง response จาก entity ต่าง ๆ
func NewStudentProfileResponse(
	user *entity.User,
	sp *entity.StudentProfile,
	record *entity.StudentAcademicRecord,
	advisorName string,
	termLabel string,
) StudentProfileResponse {

	resp := StudentProfileResponse{
		SutId:        user.SutId,
		Prefix:       user.Prefix.Prefix,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Phone:        user.Phone,
		NationalId:   user.NationalId,
		Birthday:     user.Birthday,
		ProfileImage: user.ProfileImage,

		MajorName: user.Major.Major,
		// กัน nil เผื่อไว้ เธอเช็คก่อนแล้วก็จริง แต่กันอีกชั้นเฉย ๆ
		YearOfStudy: 0,
		AdvisorName: advisorName,

		GPATermLabel: termLabel,
	}

	if sp != nil {
		resp.YearOfStudy = sp.YearOfStudy
	}

	if record != nil {
		resp.AcademicYear = record.AcademicYear
		resp.Semester = record.Semester
		resp.TermGPA = record.TermGPA
		resp.CumulativeGPA = record.CumulativeGPA
		resp.AcademicStatus = record.AcademicStatus

		resp.GPX = record.CumulativeGPA
		resp.GPALatest = record.TermGPA
	}

	return resp
}
