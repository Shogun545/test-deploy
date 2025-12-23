package controller

import (
	"net/http"
	"backend/internal/app/dto"
	"backend/config"
	"backend/internal/app/entity"
	"github.com/gin-gonic/gin"
)

/*
GET: /reports
*/
func GetAllReport(c *gin.Context) {
	db := config.DB()

	var reports []entity.Report
	db.Preload("User").
	   Preload("Status").
	   Preload("Topic").
	   Find(&reports)

	var response []dto.ReportDTO

	for _, r := range reports {
		item := dto.ReportDTO{
			ID:          r.ID,
			Description: r.Description,
		}

		// ---- User ----
		if r.User != nil {
			item.User = &dto.ReportUserDTO{
				ID:        r.User.ID,
				FirstName: r.User.FirstName,
				LastName:  r.User.LastName,
				Email:     r.User.Email,
			}
		}

		// ---- Status ----
		if r.Status != nil {
			item.Status = &dto.ReportStatusDTO{
				ID:   r.Status.ID,
				Name: r.Status.ReportStatusName,
			}
		}

		// ---- Topic ----
		if r.Topic != nil {
			item.Topic = &dto.ReportTopicDTO{
				ID:   r.Topic.ID,
				Name: r.Topic.ReportTopicName,
			}
		}

		response = append(response, item)
	}

	c.JSON(http.StatusOK, response)
}


/*
GET: /reports/:id
*/
func GetReport(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	var report entity.Report
	if err := db.
		Preload("User").
		Preload("Status").
		Preload("Topic").
		First(&report, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "report not found",
		})
		return
	}

	c.JSON(http.StatusOK, report)
}

/*
POST: /reports
*/
func CreateReport(c *gin.Context) {
	var report entity.Report

	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	db := config.DB()
	if err := db.Create(&report).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "report created successfully",
		"id":      report.ID,
	})
}

/*
PUT: /reports/:id
*/
func UpdateReport(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	var payload struct {
		ReportStatusID uint `json:"report_status_id"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid payload",
		})
		return
	}

	if payload.ReportStatusID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid status id",
		})
		return
	}

	// ✅ อัปเดตเฉพาะ column เดียว
	result := db.
		Model(&entity.Report{}).
		Where("id = ?", id).
		Update("report_status_id", payload.ReportStatusID)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "report not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "report status updated successfully",
	})
}

/*
DELETE: /reports/:id
*/
func DeleteReport(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	if err := db.Delete(&entity.Report{}, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "report deleted successfully",
	})
}
