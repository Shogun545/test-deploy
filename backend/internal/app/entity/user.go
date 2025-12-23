package entity

import (
	"gorm.io/gorm"
	"time"
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
	"unicode"
)

type User struct {
	gorm.Model

	SutId        string     `json:"sut_id" gorm:"type:varchar(50);unique;not null"`
	FirstName    string     `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName     string     `json:"last_name" gorm:"type:varchar(100);not null"`
	Email        string     `json:"email" gorm:"type:varchar(150);unique;not null"`
	Phone        string     `json:"phone" gorm:"type:varchar(20)"`
	PasswordHash string     `json:"password_hash" gorm:"type:varchar(255);not null"`
	NationalId   string     `json:"national_id" gorm:"type:varchar(20)"`
	Birthday     string     `json:"birthday" gorm:"type:varchar(20)"`
	ProfileImage string     `json:"profile_image" gorm:"type:text"`
	Notes        string     `json:"notes" gorm:"type:text"`          //addเพิ่ม
	Active       bool       `json:"active" gorm:"default:true"`      //addเพิ่ม
	LastLogin    *time.Time `json:"lastLogin" gorm:"type:timestamp"` //addเพิ่ม

	RoleID       uint        `json:"role_id"`
	Role         *Role       `json:"role" gorm:"foreignKey:RoleID"`
	PrefixID     uint        `json:"prefix_id"`
	Prefix       *Prefix     `json:"prefix" gorm:"foreignKey:PrefixID"`
	MajorID      uint        `json:"major_id"`
	Major        *Major      `json:"major" gorm:"foreignKey:MajorID"`
	DepartmentID *uint       `json:"department_id"` // ทำให้เป็นค่า null ได้
	Department   *Department `gorm:"foreignKey:DepartmentID" json:"department"`

	ManagerID    *uint  `json:"manager_id"`           // admin ที่ดูแล user คนนี้
	ManagedUsers []User `gorm:"foreignKey:ManagerID"` // virtual field

	//อันนี้ไม่ได้ไปสร้าง column เพิ่มใน DB นะ แค่บอก GORM ว่า User 1 คน มี StudentProfile 1 อัน ผ่าน FK UserID ที่อยู่ฝั่ง StudentProfile
	StudentProfile *StudentProfile `json:"student_profile" gorm:"foreignKey:UserID"`
	AdvisorProfile *AdvisorProfile `json:"advisor_profile" gorm:"foreignKey:UserID"`
}
func (u *User) ValidateSutIdByRole() error {
	if u.SutId == "" {
		return errors.New("sut_id is required")
	}

	// ต้องยาว 8 ตัว
	if utf8.RuneCountInString(u.SutId) != 8 {
		return errors.New("sut_id must be 8 characters")
	}

	switch u.Role.Role {
	case "student":
		if u.SutId[0] != 'B' {
			return errors.New("student sut_id must start with B")
		}
	case "advisor":
		if u.SutId[0] != 'T' {
			return errors.New("advisor sut_id must start with T")
		}
	case "admin":
		if u.SutId[0] != 'A' {
			return errors.New("admin sut_id must start with A")
		}
	default:
		return errors.New("unknown role")
	}
	return nil
}

// ValidateStudentProfile ตรวจข้อมูลสำหรับนักศึกษา (student)
func (u *User) ValidateStudentProfile() error {
	// ตรวจ Birthday
	if u.Role == nil || u.Role.Role != "student" {
        return nil // role อื่นไม่ต้องเช็ก
    }

    if strings.TrimSpace(u.Birthday) == "" {
        return errors.New("birthday is required")
    }

	// ตรวจ Phone
	if strings.TrimSpace(u.Phone) == "" {
		return errors.New("phone is required")
	}

	// ตรวจว่าเบอร์โทร 10 หลัก
	match, _ := regexp.MatchString(`^\d{10}$`, u.Phone)
	if !match {
		return errors.New("phone must be 10 digits")
	}

	return nil
}

// ValidateAdminProfile ตรวจข้อมูลของ Admin
func (u *User) ValidateAdminProfile() error {
	// Email ต้องไม่ว่าง
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email is required")
	}

	// ตรวจ format ของ Email แบบง่ายๆ
	match, _ := regexp.MatchString(`^[\w._%+-]+@[\w.-]+\.[a-zA-Z]{2,}$`, u.Email)
	if !match {
		return errors.New("email is invalid")
	}

	// Phone ต้องไม่ว่าง
	if strings.TrimSpace(u.Phone) == "" {
		return errors.New("phone is required")
	}

	// Phone ต้องเป็นตัวเลข 10 หลัก
	match, _ = regexp.MatchString(`^\d{10}$`, u.Phone)
	if !match {
		return errors.New("phone must be 10 digits")
	}

	// Notes ไม่ต้อง validate เพราะ optional

	// Active เป็น boolean ไม่ต้อง validate
	return nil
}

// ValidateActiveStatus ตรวจค่า Active ของ User
func (u *User) ValidateActiveStatus() error {
	// Active เป็น boolean → ต้องมีค่า true หรือ false (แต่ Go boolean มีค่า default อยู่แล้ว)
	// ถ้าอยากเจาะจงว่า ต้องกรอกจริงๆ (ไม่เป็น null) ก็เช็คว่า u != nil
	if u == nil {
		return errors.New("User is nil")
	}
	return nil
}

func (u *User) ValidatePasswordStrength() error {
    pw := u.PasswordHash // สมมติ PasswordHash ใช้สำหรับ validate ก่อน hash

    if len(pw) < 8 {
        return errors.New("password must be at least 8 characters long")
    }

    var hasUpper, hasLower, hasNumber, hasSpecial bool
    for _, c := range pw {
        switch {
        case unicode.IsUpper(c):
            hasUpper = true
        case unicode.IsLower(c):
            hasLower = true
        case unicode.IsNumber(c):
            hasNumber = true
        case unicode.IsPunct(c) || unicode.IsSymbol(c):
            hasSpecial = true
        }
    }

    if !hasUpper {
        return errors.New("password must contain at least one uppercase letter")
    }
    if !hasLower {
        return errors.New("password must contain at least one lowercase letter")
    }
    if !hasNumber {
        return errors.New("password must contain at least one number")
    }
    if !hasSpecial {
        return errors.New("password must contain at least one special character")
    }

    return nil
}