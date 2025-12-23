package repository

import (
	"backend/internal/app/entity"
	"time"
	"gorm.io/gorm"
)

type AcademicCalendarRepository interface {
	// แก้ Signature ไม่ต้องรับ TimeCalendar แล้ว
	FindEventsByDateRange(startDate, endDate time.Time) ([]entity.AcademicCalendar, error)
	GetEventByID(id string) (*entity.AcademicCalendar, error)
	CreateEvent(event *entity.AcademicCalendar) error
	UpdateEvent(event *entity.AcademicCalendar) error
	DeleteEvent(event *entity.AcademicCalendar) error
}

type academicCalendarRepository struct {
	DB *gorm.DB
}

func NewAcademicCalendarRepository(db *gorm.DB) AcademicCalendarRepository {
	return &academicCalendarRepository{DB: db}
}

func (r *academicCalendarRepository) FindEventsByDateRange(queryStart, queryEnd time.Time) ([]entity.AcademicCalendar, error) {
	var calendars []entity.AcademicCalendar
	// ✅ Logic: หา Event ที่ช่วงเวลา "ทับซ้อน" (Overlap) กับช่วงที่ Query
	// (Start_Event <= Query_End) AND (End_Event >= Query_Start)
	err := r.DB.
		Where("start_date_time <= ? AND end_date_time >= ?", queryEnd, queryStart).
		Find(&calendars).Error
	return calendars, err
}

func (r *academicCalendarRepository) GetEventByID(id string) (*entity.AcademicCalendar, error) {
	var event entity.AcademicCalendar
	// ไม่ต้อง Preload TimeCalendar แล้ว
	err := r.DB.First(&event, "id = ?", id).Error
	return &event, err
}

func (r *academicCalendarRepository) CreateEvent(event *entity.AcademicCalendar) error {
	// ✅ ง่ายขึ้นเยอะ ไม่ต้องทำ Transaction 2 ตาราง
	return r.DB.Create(event).Error
}

func (r *academicCalendarRepository) UpdateEvent(event *entity.AcademicCalendar) error {
	// ใช้ Save เพื่อบันทึกทุก field รวมถึงกรณีเปลี่ยนค่าเป็น null/false
	return r.DB.Save(event).Error
}

func (r *academicCalendarRepository) DeleteEvent(event *entity.AcademicCalendar) error {
	return r.DB.Delete(event).Error
}