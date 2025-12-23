package entity

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type AcademicCalendar struct {
	gorm.Model

	// ✅ เพิ่ม Tags validation
	EventName string `json:"event_name" valid:"required~กรุณาระบุชื่อกิจกรรม,maxstringlength(100)~ชื่อกิจกรรมยาวเกินไป"`
	EventType string `json:"event_type" valid:"required~กรุณาระบุประเภทกิจกรรม,in(exam|activity|holiday)~ประเภทกิจกรรมไม่ถูกต้อง"`

	// ใช้ time.Time แล้ว ไม่ต้องเช็ค Format regex แต่เช็ค required ได้
	StartDateTime time.Time `json:"start_date_time" valid:"required~กรุณาระบุวันเวลาเริ่มต้น"`
	EndDateTime   time.Time `json:"end_date_time" valid:"required~กรุณาระบุวันเวลาสิ้นสุด"`

	IsHoliday bool   `json:"is_holiday"`
	AdminID   uint   `json:"admin_id"`
	User      User   `json:"user" gorm:"foreignKey:AdminID" valid:"-"`
}

// ✅ ฟังก์ชัน Validate ฝังอยู่ใน Entity
func (e *AcademicCalendar) Validate() error {
	// 1. ตรวจสอบ Basic Tags (ชื่อ, ประเภท, ค่าว่าง)
	if _, err := govalidator.ValidateStruct(e); err != nil {
		return err
	}

	// 2. ตรวจสอบ Logic เวลา (Time Logic)
	// เนื่องจากเป็น time.Time แล้ว สามารถเปรียบเทียบได้เลย ไม่ต้อง Parse String
	if e.EndDateTime.Before(e.StartDateTime) {
		return fmt.Errorf("วันเวลาสิ้นสุด ต้องอยู่หลังจากวันเวลาเริ่มต้น")
	}

	return nil
}