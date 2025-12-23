package controller

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"backend/config"
	"backend/internal/app/entity"
)

type ReportFeedbackController struct{}

func NewReportFeedbackController() *ReportFeedbackController {
	return &ReportFeedbackController{}
}

// 🟢 POST /report_feedbacks
// หน้าที่: อาจารย์บันทึกคอมเมนต์ (Feedback) ลงในรายงานของนักศึกษา
func (ctrl *ReportFeedbackController) Create(c *gin.Context) {
	var feedback entity.ReportFeedback

	// 1. รับค่าจาก JSON
	if err := c.ShouldBindJSON(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Validation
	if _, err := govalidator.ValidateStruct(feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	// 3. (Optional) เช็คว่า Report ที่จะคอมเมนต์มีอยู่จริงไหม
	var report entity.ProgressReport
	if err := db.First(&report, feedback.ProgressReportsID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Progress Report not found"})
		return
	}

	// 4. บันทึก Feedback ลงฐานข้อมูล
	if err := db.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Feedback submitted successfully",
		"data":    feedback,
	})
}