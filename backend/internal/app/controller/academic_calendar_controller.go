package controller

import (
	"backend/internal/app/dto"
	"backend/internal/app/entity" // ✅ Import Entity เพื่อเรียกใช้ Validate()
	"backend/internal/service/academiccalendar"
	"fmt"
	"net/http"
	"time" // ✅ Import time เพื่อแปลงวันที่

	"github.com/gin-gonic/gin"
)

type AcademicCalendarController struct {
	Service *service.AcademicCalendarService
}

func NewAcademicCalendarController(s *service.AcademicCalendarService) *AcademicCalendarController {
	return &AcademicCalendarController{Service: s}
}

func (ctrl *AcademicCalendarController) GetHolidaysByYear(c *gin.Context) {
	yearStr := c.Query("year")
	if yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Year is required"})
		return
	}

	response, err := ctrl.Service.GetHolidaysByYear(yearStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *AcademicCalendarController) CreateEvent(c *gin.Context) {
	var req dto.EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลที่ส่งมาไม่ถูกต้อง"})
		return
	}

	// -----------------------------------------------------
	// ✅ สร้าง Entity จำลองเพื่อเรียกใช้ฟังก์ชัน Validate()
	// -----------------------------------------------------
	startStr := fmt.Sprintf("%s %s", req.StartDate, req.StartTime)
	endStr := fmt.Sprintf("%s %s", req.EndDate, req.EndTime)
	
	startTime, errStart := time.Parse("2006-01-02 15:04", startStr)
	endTime, errEnd := time.Parse("2006-01-02 15:04", endStr)

	if errStart != nil || errEnd != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบวันเวลาไม่ถูกต้อง"})
		return
	}

	tempEvent := entity.AcademicCalendar{
		EventName:     req.Title,
		EventType:     req.Type,
		StartDateTime: startTime,
		EndDateTime:   endTime,
	}

	// เรียกฟังก์ชัน Validate ที่อยู่ใน Entity
	if err := tempEvent.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// -----------------------------------------------------

	userID, exists := c.Get("user_id") 
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in context"})
		return
	}

	switch v := userID.(type) {
	case uint:
		req.AdminID = v
	case float64:
		req.AdminID = uint(v)
	case int:
		req.AdminID = uint(v)
	}

	event, err := ctrl.Service.CreateEvent(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event created", "data": event})
}

func (ctrl *AcademicCalendarController) UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	var req dto.EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// -----------------------------------------------------
	// ✅ Validate Logic
	// -----------------------------------------------------
	startStr := fmt.Sprintf("%s %s", req.StartDate, req.StartTime)
	endStr := fmt.Sprintf("%s %s", req.EndDate, req.EndTime)
	
	startTime, errStart := time.Parse("2006-01-02 15:04", startStr)
	endTime, errEnd := time.Parse("2006-01-02 15:04", endStr)

	if errStart != nil || errEnd != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบวันเวลาไม่ถูกต้อง"})
		return
	}

	tempEvent := entity.AcademicCalendar{
		EventName:     req.Title,
		EventType:     req.Type,
		StartDateTime: startTime,
		EndDateTime:   endTime,
	}

	if err := tempEvent.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// -----------------------------------------------------

	if err := ctrl.Service.UpdateEvent(id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated"})
}

func (ctrl *AcademicCalendarController) DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.Service.DeleteEvent(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event and associated time deleted"})
}