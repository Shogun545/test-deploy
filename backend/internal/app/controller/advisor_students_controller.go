package controller

import (
	"backend/config"
	"backend/internal/app/dto"
	"backend/internal/app/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
)

func GetMyStudents(c *gin.Context) {
	sutIdRaw, ok := c.Get("sut_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	sutId, _ := sutIdRaw.(string)

	db := config.DB()

	var advisor entity.User
	if err := db.
		Preload("AdvisorProfile.Students.User.Major").                          // preload Major ของแต่ละ student
		Preload("AdvisorProfile.Students.User.StudentProfile.StudentAcademicRecords", func(db *gorm.DB) *gorm.DB {
			return db.Order("academic_year desc, semester desc") // เรียง GPA ล่าสุดไว้บนสุด
		}).
		Where("sut_id = ?", sutId).
		First(&advisor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบ advisor หรือ ไม่มีนักศึกษา"})
		return
	}

	resp := dto.NewAdvisorStudentsResponse(&advisor)
	c.JSON(http.StatusOK, resp)
}