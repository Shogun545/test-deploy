package controller

import (
	"net/http"

	"backend/config"
	"backend/internal/app/entity"
	"github.com/gin-gonic/gin"
)

/*
GET: /report-topics
*/
func GetAllReportTopic(c *gin.Context) {
	db := config.DB()

	var topics []entity.ReportTopic
	if err := db.Find(&topics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, topics)
}

/*
GET: /report-topics/:id
*/
func GetReportTopic(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	var topic entity.ReportTopic
	if err := db.First(&topic, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "topic not found",
		})
		return
	}

	c.JSON(http.StatusOK, topic)
}

/*
POST: /report-topics
*/
func CreateReportTopic(c *gin.Context) {
	var topic entity.ReportTopic

	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	db := config.DB()
	if err := db.Create(&topic).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "report topic created successfully",
		"id":      topic.ID,
	})
}

/*
PUT: /report-topics/:id
*/
func UpdateReportTopic(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	var topic entity.ReportTopic
	if err := db.First(&topic, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "topic not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	if err := db.Save(&topic).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "report topic updated successfully",
	})
}

/*
DELETE: /report-topics/:id
*/
func DeleteReportTopic(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	if err := db.Delete(&entity.ReportTopic{}, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "report topic deleted successfully",
	})
}
