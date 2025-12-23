package controller

import (
	"backend/config"
	"backend/internal/app/entity"
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetAllFAQ(c *gin.Context) {

	db := config.DB()
	var FAQ []entity.FAQ
	results := db.Preload("AppointmentTopicID").Preload("UserID").Find(&FAQ)
	if results.Error != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": results.Error.Error()})

		return

	}
	c.JSON(http.StatusOK, &FAQ)
}

func GetFAQ(c *gin.Context) {

	ID := c.Param("id")
	var FAQ entity.FAQ
	db := config.DB()
	results := db.Preload("AppointmentTopicID").Preload("UserID").First(&FAQ, ID)
	if results.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": results.Error.Error()})
		return
	}
	if FAQ.ID == 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}
	c.JSON(http.StatusOK, FAQ)
}

func UpdateFAQ(c *gin.Context) {

	var FAQ entity.FAQ

	FAQID := c.Param("id")

	db := config.DB()

	result := db.First(&FAQ, FAQID)

	if result.Error != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "id not found"})

		return

	}

	if err := c.ShouldBindJSON(&FAQ); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, unable to map payload"})

		return

	}

	result = db.Save(&FAQ)

	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})

		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successful"})

}

func DeleteFAQ(c *gin.Context) {

	id := c.Param("id")

	db := config.DB()

	if tx := db.Exec("DELETE FROM FAQes WHERE id = ?", id); tx.RowsAffected == 0 {

		c.JSON(http.StatusBadRequest, gin.H{"error": "id not found"})

		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successful"})

}

func CreateFAQ(c *gin.Context) {
	var FAQ entity.FAQ

	if err := c.ShouldBindJSON(&FAQ); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, unable to map payload"})
		return
	}
	db := config.DB()
	result := db.Create(&FAQ)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "FAQ Created successful", "id": FAQ.ID})
}
