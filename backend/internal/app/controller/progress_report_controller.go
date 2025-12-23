package controller

import (
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"backend/config"
	"backend/internal/app/entity"
)

type ProgressReportController struct{}

func NewProgressReportController() *ProgressReportController {
	return &ProgressReportController{}
}

// 🟢 POST /progress_reports
func (ctrl *ProgressReportController) Create(c *gin.Context) {
	var report entity.ProgressReport

	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := govalidator.ValidateStruct(report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	// 🔍 แก้ไขจุดที่ 1: ใช้ AdvisorLogsID ให้ตรงกับ Entity
	var advisorLog entity.AdvisorLog
	if err := db.First(&advisorLog, report.AdvisorLogsID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advisor Log not found"})
		return
	}

	if advisorLog.Status != "PendingReport" { 
		c.JSON(http.StatusBadRequest, gin.H{"error": "This log is not waiting for a report"})
		return
	}

	report.SubmittedAt = time.Now()
	if err := db.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&advisorLog).Update("status", "ReportSubmitted").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update log status"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Report submitted successfully",
		"data":    report,
	})
}

// 🔵 GET /progress_reports/log/:log_id
func (ctrl *ProgressReportController) GetByLogID(c *gin.Context) {
	logID := c.Param("log_id")
	var reports []entity.ProgressReport

	// 🔍 แก้ไขจุดที่ 2: เปลี่ยนชื่อคอลัมน์เป็น advisor_logs_id
	if err := config.DB().Preload("Feedbacks").Where("advisor_logs_id = ?", logID).Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reports})
}