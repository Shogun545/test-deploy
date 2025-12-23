package dto

import (
	"backend/internal/app/entity"
	"time"
)

// AdminManagedUserSummary เป็น DTO สรุปข้อมูลผู้ใช้ที่ Admin ดูแล
type AdminManagedUserSummary struct {
	SutID    string `json:"sutId"`
	Prefix   string `json:"prefix"` 
    FirstName string `json:"firstName"` 
    LastName string `json:"lastName"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
type MajorEntry struct {
    Major string `json:"major"`
}
func NewMajorEntry(m *entity.Major) MajorEntry {
    return MajorEntry{
        Major: m.Major,
    }
}

type UserCreatedDateResponse struct {
	SutID     string    `json:"sut_id"`
	CreatedAt time.Time `json:"created_at"`
}

type UserListFilter struct {
    Role       string `form:"role"`      
    Status     string `form:"status"`    // 'active' หรือ 'inactive'
	Major      string `form:"major"`
    Department string `form:"department"`
    Search     string `form:"search"`    // ค้นหา SutID, Name, Email
    // สามารถเพิ่ม Page, PageSize ได้ในอนาคต เช่น Page int `form:"page"`
}

// AdminProfileResponse (GET /api/admin/profile)
// AdminProfileResponse DTO สำหรับส่งข้อมูลโปรไฟล์ Admin กลับไป Frontend
type AdminProfileResponse struct {
	SutID          string                    `json:"sutId"`
	FirstName      string                    `json:"firstName"`
	LastName       string                    `json:"lastName"`
	Prefix         string                    `json:"prefix"` // จะมาจาก u.Prefix.Prefix
	Email          string                    `json:"email"`
	Phone          string                    `json:"phone"`
	Role           string                    `json:"role"`
	DepartmentName string                    `json:"departmentName"` // ชื่อหน่วยงานที่สังกัด
	Notes          string                    `json:"notes"`          // หน้าที่/หมายเหตุ
	Active         bool                      `json:"active"`         // สถานะบัญชี
	CreatedAt      time.Time                 `json:"createdAt"`
	LastLogin      *time.Time                `json:"lastLogin"`    // เวลาเข้าสู่ระบบล่าสุด
	ManagedUsers   []AdminManagedUserSummary `json:"managedUsers"` // เพิ่มตรงนี้
	TotalManaged   int                       `json:"totalManaged"`

}

// NewAdminProfileResponse สร้าง DTO จาก Entity User
func NewAdminProfileResponse(u *entity.User) AdminProfileResponse {
	deptName := ""
	if u.Department != nil {
		deptName = u.Department.Department
	}

	roleName := ""
	if u.Role != nil {
		roleName = u.Role.Role
	}

	prefix := ""
	if u.Prefix != nil {
		prefix = u.Prefix.Prefix
	}

	// สร้าง ManagedUsers slice
	managed := []AdminManagedUserSummary{}
	if u.ManagedUsers != nil {
		for _, mu := range u.ManagedUsers {
			fullName := mu.FirstName + " " + mu.LastName
			if mu.Prefix != nil {
				fullName = mu.Prefix.Prefix + " " + fullName
			}
			role := ""
			if mu.Role != nil {
				role = mu.Role.Role
			}
			managed = append(managed, AdminManagedUserSummary{
				SutID:    mu.SutId,
				FullName: fullName,
				Email:    mu.Email,
				Role:     role,
			})
		}
	}

	return AdminProfileResponse{
		SutID:          u.SutId,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Prefix:         prefix,
		Email:          u.Email,
		Phone:          u.Phone,
		Role:           roleName,
		DepartmentName: deptName,
		Notes:          u.Notes,
		Active:         u.Active,
		CreatedAt:      u.CreatedAt,
		LastLogin:      u.LastLogin,
		ManagedUsers:   managed,
		TotalManaged:   len(managed),
	}
}

// UpdateAdminProfileRequest (PUT/PATCH /api/admin/profile)
// UpdateAdminProfileRequest DTO สำหรับรับข้อมูลที่ Admin ส่งมาแก้ไข
type UpdateAdminProfileRequest struct {
	Prefix string `json:"prefix"`
	FirstName string `json:"firstName"` 
	LastName  string `json:"lastName"`
	Phone  string `json:"phone" `
	Notes  string `json:"notes" `
	Active *bool  `json:"active"`      
	Email  string `json:"email"` 
}

// AdminManagedUsersResponse (GET /api/admin/users)
// ManagedUserEntry Entry สำหรับแสดงในตารางรวม
type ManagedUserEntry struct {
	SutID          	string `json:"sutId"`
	Prefix   		string `json:"prefix"` 
    FirstName 		string `json:"firstName"` 
    LastName 		string `json:"lastName"`
	FullName       	string `json:"fullName"`
	Email          	string `json:"email"`
	Role           	string `json:"role"`
	MajorName      	string `json:"majorName"`      // สาขาวิชา (ถ้ามี)
	DepartmentName 	string `json:"departmentName"` // หน่วยงานที่สังกัด
	Active         bool   `json:"active"`
}

// AdminManagedUsersResponse DTO สำหรับรายการผู้ใช้ทั้งหมดในระบบ
type AdminManagedUsersResponse struct {
	TotalUsers int                `json:"totalUsers"`
	Users      []ManagedUserEntry `json:"users"`
}

// NewAdminManagedUsersResponse สร้าง DTO จาก []entity.User
func NewAdminManagedUsersResponse(users []*entity.User) AdminManagedUsersResponse {
	entries := make([]ManagedUserEntry, 0, len(users))
	for _, u := range users {
		majorName := ""
		if u.Major != nil {
			majorName = u.Major.Major
		}
		deptName := ""
		if u.Department != nil {
			deptName = u.Department.Department
		}
		roleName := ""
		if u.Role != nil {
			roleName = u.Role.Role
		}

		prefix := ""
		if u.Prefix != nil {
			prefix = u.Prefix.Prefix
		}
		fullName := u.FirstName + " " + u.LastName
		if prefix != "" {
			fullName = prefix + " " + fullName
		}

		entries = append(entries, ManagedUserEntry{
			SutID:          u.SutId,
			Prefix:       	prefix,
			FirstName:   	u.FirstName,
			LastName:    	u.LastName,
			FullName:       fullName,
			Email:          u.Email,
			Role:           roleName,
			MajorName:      majorName,
			DepartmentName: deptName,
			Active:         u.Active,
		})
	}

	return AdminManagedUsersResponse{
		TotalUsers: len(users),
		Users:      entries,
	}
}

// AdminUserDetailResponse (GET /api/admin/users/:sut_id)
// AdminUserDetailResponse DTO สำหรับแสดงรายละเอียดของผู้ใช้คนใดคนหนึ่ง
type AdminUserDetailResponse struct {
	SutID          string `json:"sutId"`
	Prefix         string `json:"prefix"` 
    FirstName      string `json:"firstName"` 
    LastName       string `json:"lastName"`
	FullName       string `json:"fullName"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Role           string `json:"role"`
	DepartmentName string `json:"departmentName"`
	MajorName      string `json:"majorName"`
	Active         bool   `json:"active"`

	// ข้อมูลเฉพาะทาง (ถ้ามี)
	OfficeRoom   string `json:"officeRoom,omitempty"`   // ถ้าเป็นอาจารย์
	Specialties  string `json:"specialties,omitempty"`  // ถ้าเป็นอาจารย์
	
	YearOfStudy int     `json:"yearOfStudy,omitempty"` // ถ้าเป็นนักศึกษา
	GPALatest   float64 `json:"gpaLatest,omitempty"`   // ถ้าเป็นนักศึกษา
}

// NewAdminUserDetailResponse สร้าง DTO จาก entity.User พร้อมความสัมพันธ์
func NewAdminUserDetailResponse(u *entity.User) AdminUserDetailResponse {
	prefix := ""
	if u.Prefix != nil {
		prefix = u.Prefix.Prefix
	}
	fullName := u.FirstName + " " + u.LastName
	if prefix != "" {
		fullName = prefix + " " + fullName
	}

	resp := AdminUserDetailResponse{
		SutID:    u.SutId,
		Prefix:   prefix,
		FirstName: u.FirstName,
		LastName: u.LastName,
		FullName: fullName,
		Email:    u.Email,
		Phone:    u.Phone,
		Active:   u.Active,
	}

	// Set Role and Department
	if u.Department != nil {
		resp.DepartmentName = u.Department.Department
	}
	if u.Role != nil {
		resp.Role = u.Role.Role
	}
	if u.Major != nil {
		resp.MajorName = u.Major.Major
	}

	// Set Advisor specific fields
	if u.AdvisorProfile != nil {
		resp.OfficeRoom = u.AdvisorProfile.OfficeRoom
		resp.Specialties = u.AdvisorProfile.Specialties
	}

	// Set Student specific fields
	if u.StudentProfile != nil {
		resp.YearOfStudy = u.StudentProfile.YearOfStudy

		// ดึง GPA ล่าสุด (สมมติว่าถูก Preload มาแล้ว)
		if len(u.StudentProfile.StudentAcademicRecords) > 0 {
			resp.GPALatest = float64(u.StudentProfile.StudentAcademicRecords[0].CumulativeGPA)
		}
	}

	return resp
}

// UpdateUserStatusRequest ใช้สำหรับรับสถานะใหม่เมื่อ Admin ต้องการเปลี่ยนสถานะผู้ใช้งาน
type UpdateUserStatusRequest struct {
    Status string `json:"status" binding:"required"` // เช่น "Active", "Inactive", "Suspended"
}

type UpdateManagedUserRequest struct {
	Phone  *string `json:"phone"`
}
