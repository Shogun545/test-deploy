package dto

import "backend/internal/app/entity"

type AdvisorProfileResponse struct {
	Prefix    string `json:"prefix"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`

	OfficeRoom string `json:"office_room"`

	MajorName    string `json:"major_name"`
	Department   string `json:"department_name"`
	AdviseeCount int    `json:"advisee_count"`

	Specialties string `json:"specialties"`
	IsActive    bool   `json:"is_active"`
	ProfileImage string `json:"profile_image"`

}

type AdviseeSummaryResponse struct {
	SutId          string  `json:"sut_id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	YearOfStudy    int     `json:"year_of_study"`
	MajorName      string  `json:"major_name"`
	GPX            float32 `json:"gpx"`
	AcademicStatus string  `json:"academic_status"`
}

func NewAdvisorProfileResponse(user *entity.User, advisor *entity.AdvisorProfile) AdvisorProfileResponse {
	prefix := ""
	if user.Prefix != nil {
		prefix = user.Prefix.Prefix
	}

	majorName := ""
	if user.Major != nil {
		majorName = user.Major.Major
	}

	adviseeCount := 0
	if advisor != nil && advisor.Students != nil {
		adviseeCount = len(advisor.Students)
	}

	officeRoom := ""
	specialties := ""
	isActive := false
	if advisor != nil {
		officeRoom = advisor.OfficeRoom
		specialties = advisor.Specialties
		isActive = advisor.IsActive
	}

	return AdvisorProfileResponse{
		Prefix:       prefix,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Phone:        user.Phone,
		OfficeRoom:   officeRoom,
		MajorName:    majorName,
		Department:   "", // ถ้ามี department ค่อยเพิ่ม
		AdviseeCount: adviseeCount,
		Specialties:  specialties,
		IsActive:     isActive,
		ProfileImage: user.ProfileImage,

	}
}

type UpdateAdvisorProfileRequest struct {
	Prefix       string `json:"prefix"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Phone        string `json:"phone"`
	ProfileImage string `json:"profile_image"`
	OfficeRoom   string `json:"office_room"`
	Specialties  string `json:"specialties"`
	IsActive     *bool  `json:"is_active"`
	Email        string `json:"email"`
	
}
