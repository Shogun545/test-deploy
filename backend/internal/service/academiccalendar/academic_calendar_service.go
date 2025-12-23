package service

import (
	"backend/internal/app/dto"
	"backend/internal/app/entity"
	"backend/internal/app/repository"
	
	"errors"
	"time"
)

type AcademicCalendarService struct {
	Repo repository.AcademicCalendarRepository
}

func NewAcademicCalendarService(repo repository.AcademicCalendarRepository) *AcademicCalendarService {
	return &AcademicCalendarService{Repo: repo}
}

// ---------------------------------------------------------
// Helper: จัดการ Timezone (Asia/Bangkok)
// ---------------------------------------------------------
func getLocation() *time.Location {
	// โหลด Timezone ไทย
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return time.Local // ถ้าหาไม่เจอให้ใช้ Local เครื่อง
	}
	return loc
}

// Helper: รวม "วันที่" (String) + "เวลา" (String) -> เป็น "time.Time" (Asia/Bangkok)
func combineDateAndTime(dateStr, timeStr string) (time.Time, error) {
	loc := getLocation()
	layout := "2006-01-02 15:04"

	if dateStr == "" {
		return time.Time{}, errors.New("date is required")
	}
	
	// ถ้าไม่ส่งเวลามา ให้ตั้งเป็น 00:00
	if timeStr == "" {
		timeStr = "00:00"
	}

	fullStr := dateStr + " " + timeStr

	// 🔥 Key Point: ใช้ ParseInLocation เพื่อระบุว่า string นี้คือเวลาไทย
	return time.ParseInLocation(layout, fullStr, loc)
}

// ---------------------------------------------------------
// Functions
// ---------------------------------------------------------

func (s *AcademicCalendarService) GetHolidaysByYear(yearStr string) ([]dto.HolidayResponse, error) {
	// Query ช่วงวันที่ 1 ม.ค. - 31 ธ.ค. ของปีนั้น
	startDate, err := time.Parse("2006-01-02", yearStr+"-01-01")
	if err != nil {
		return nil, errors.New("invalid year format")
	}
	endDate := startDate.AddDate(1, 0, 0)

	events, err := s.Repo.FindEventsByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	loc := getLocation() // เตรียมแปลงเวลาเป็นไทยตอนแสดงผล

	var response []dto.HolidayResponse
	for _, cal := range events {
		// แปลงเวลาจาก DB (ซึ่งอาจเป็น UTC) ให้เป็นเวลาไทยก่อนส่งกลับ
		startInThai := cal.StartDateTime.In(loc)
		endInThai := cal.EndDateTime.In(loc)

		response = append(response, dto.HolidayResponse{
			ID: cal.ID,
			
			// ส่งกลับแยกเป็น Date และ Time
			StartDate: startInThai.Format("2006-01-02"),
			EndDate:   endInThai.Format("2006-01-02"),
			StartTime: startInThai.Format("15:04"),
			EndTime:   endInThai.Format("15:04"),

			LocalName:   cal.EventName,
			Name:        cal.EventName,
			CountryCode: "TH",
			Types:       []string{cal.EventType},
		})
	}
	return response, nil
}

func (s *AcademicCalendarService) CreateEvent(req dto.EventRequest) (*entity.AcademicCalendar, error) {
	// 1. แปลง String เป็น Time Object (Zone ไทย)
	startDateTime, err := combineDateAndTime(req.StartDate, req.StartTime)
	if err != nil {
		return nil, errors.New("invalid start date/time format")
	}
	endDateTime, err := combineDateAndTime(req.EndDate, req.EndTime)
	if err != nil {
		return nil, errors.New("invalid end date/time format")
	}

	// 2. Validation: วันจบต้องไม่มาก่อนวันเริ่ม
	if endDateTime.Before(startDateTime) {
		return nil, errors.New("วันสิ้นสุดต้องมาหลังวันเริ่มต้น")
	}

	// 3. สร้าง Entity ลง DB
	event := entity.AcademicCalendar{
		EventName:     req.Title,
		EventType:     req.Type,
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
		AdminID:       req.AdminID,
		// IsHoliday: true/false (กำหนด Logic เพิ่มเติมได้ถ้าต้องการ)
	}

	if err := s.Repo.CreateEvent(&event); err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *AcademicCalendarService) UpdateEvent(id string, req dto.EventRequest) error {
	// 1. หาข้อมูลเดิม
	event, err := s.Repo.GetEventByID(id)
	if err != nil {
		return errors.New("event not found")
	}

	// 2. แปลงค่าใหม่เป็น Time Object (Zone ไทย)
	startDateTime, err := combineDateAndTime(req.StartDate, req.StartTime)
	if err != nil { return errors.New("invalid start date/time") }
	
	endDateTime, err := combineDateAndTime(req.EndDate, req.EndTime)
	if err != nil { return errors.New("invalid end date/time") }

	if endDateTime.Before(startDateTime) {
		return errors.New("วันสิ้นสุดต้องมาหลังวันเริ่มต้น")
	}

	// 3. Update fields
	event.EventName = req.Title
	event.EventType = req.Type
	event.StartDateTime = startDateTime
	event.EndDateTime = endDateTime

	return s.Repo.UpdateEvent(event)
}

func (s *AcademicCalendarService) DeleteEvent(id string) error {
	event, err := s.Repo.GetEventByID(id)
	if err != nil {
		return errors.New("event not found")
	}
	return s.Repo.DeleteEvent(event)
}