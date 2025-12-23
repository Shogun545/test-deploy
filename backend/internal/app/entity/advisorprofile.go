package entity

import (
	"errors"
	"regexp"
	"strings"
	"gorm.io/gorm"
)

type AdvisorProfile struct {
    gorm.Model
    Specialties  string `json:"specialties"`    // ความเชี่ยวชาญ
    IsActive     bool   `json:"is_active"`      // เปิด/ปิดรับคำปรึกษา
    OfficeRoom   string `json:"office_room"`    // ห้องทำงาน

    UserID uint   `json:"user_id"`
    User   *User  `json:"user" gorm:"foreignKey:UserID"`

    Students []StudentProfile `json:"students" gorm:"foreignKey:AdvisorProfileID"`
}

// Validate ตรวจข้อมูลของ AdvisorProfile
func (a *AdvisorProfile) Validate() error {

	// ตรวจว่า User ต้องไม่เป็น nil เพราะเราต้องใช้ Phone
	if a.User == nil {
		return errors.New("User is required") // ถ้าไม่มี User จะ error
	}

	// ตรวจ Phone ของ Advisor
	if strings.TrimSpace(a.User.Phone) == "" {
		return errors.New("phone is required") // เบอร์โทรต้องกรอก
	}

	// ตรวจว่า Phone ต้องเป็นตัวเลข 10 หลัก
	match, _ := regexp.MatchString(`^\d{10}$`, a.User.Phone)
	if !match {
		return errors.New("phone must be 10 digits")
	}

	// ตรวจ OfficeRoom ต้องไม่ว่าง
	if strings.TrimSpace(a.OfficeRoom) == "" {
		return errors.New("officeRoom is required")
	}

	// ตรวจ Specialties ต้องไม่ว่าง
	if strings.TrimSpace(a.Specialties) == "" {
		return errors.New("specialties is required")
	}

	// IsActive เป็น boolean ไม่ต้องตรวจ เพราะค่า false valid ได้
	return nil
}