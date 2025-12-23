package controller

import (
	"backend/internal/app/dto"
	"backend/internal/service/studentprofile"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProfileController struct {
	Service *service.ProfileService
}

func NewProfileController(s *service.ProfileService) *ProfileController {
	return &ProfileController{Service: s}
}

func (pc *ProfileController) GetMyStudentProfile(c *gin.Context) {
	sutIdRaw, ok := c.Get("sut_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	sutID, _ := sutIdRaw.(string)

	resp, err := pc.Service.GetMyProfile(sutID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบผู้ใช้"})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (pc *ProfileController) UpdateMyStudentProfile(c *gin.Context) {
	sutIdRaw, ok := c.Get("sut_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	sutID, _ := sutIdRaw.(string)

	var req dto.UpdateStudentProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := pc.Service.UpdateMyProfile(sutID, &req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบผู้ใช้"})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "อัปเดตโปรไฟล์ไม่สำเร็จ"})
		return
	}
	c.JSON(http.StatusOK, resp)
}