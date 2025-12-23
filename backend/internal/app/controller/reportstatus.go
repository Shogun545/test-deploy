package controller

import (
	"net/http"

	"backend/config"
	"backend/internal/app/entity"
	"github.com/gin-gonic/gin"
)

/*
GET: /report-status
*/
func GetAllReportStatus(c *gin.Context) {
	db := config.DB()

	var status []entity.ReportStatus
	if err := db.Find(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

/*
GET: /report-status/:id
*/
func GetReportStatus(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	var status entity.ReportStatus
	if err := db.First(&status, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "status not found",
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

/*
POST: /report-status
*/
func CreateReportStatus(c *gin.Context) {
	var status entity.ReportStatus

	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	db := config.DB()
	if err := db.Create(&status).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "report status created successfully",
		"id":      status.ID,
	})
}

/*
PUT: /report-status/:id
*/
func UpdateReportStatus(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	var status entity.ReportStatus
	if err := db.First(&status, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "status not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	if err := db.Save(&status).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "report status updated successfully",
	})
}

/*
DELETE: /report-status/:id
*/
func DeleteReportStatus(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	if err := db.Delete(&entity.ReportStatus{}, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "report status deleted successfully",
	})
}
