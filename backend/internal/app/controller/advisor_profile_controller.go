package controller

import (
	"backend/internal/app/dto"
	"backend/internal/service/advisorprofile"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdvisorProfileController struct {
	Service *advisorprofile.AdvisorProfileService
}

func NewAdvisorProfileController(s *advisorprofile.AdvisorProfileService) *AdvisorProfileController {
	return &AdvisorProfileController{Service: s}
}

func (c *AdvisorProfileController) GetMyAdvisorProfile(ctx *gin.Context) {
	sutIdRaw, ok := ctx.Get("sut_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	sutID, _ := sutIdRaw.(string)

	resp, err := c.Service.GetMyProfile(sutID)
	if err != nil {
		if errors.Is(err, advisorprofile.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบผู้ใช้"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *AdvisorProfileController) UpdateMyAdvisorProfile(ctx *gin.Context) {
    sutIdRaw, ok := ctx.Get("sut_id")
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    sutID, _ := sutIdRaw.(string)

    var req dto.UpdateAdvisorProfileRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := c.Service.UpdateMyProfile(sutID, &req) // เรียก service เดิม
    if err != nil {
        if errors.Is(err, advisorprofile.ErrNotFound) {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบผู้ใช้"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "อัปเดตโปรไฟล์ไม่สำเร็จ"})
        return
    }

    ctx.JSON(http.StatusOK, resp)
}

func (c *AdvisorProfileController) GetMyStudents(ctx *gin.Context) {
    sutIdRaw, ok := ctx.Get("sut_id")
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    sutID, _ := sutIdRaw.(string)

    students, err := c.Service.GetMyStudents(sutID)
    if err != nil {
        if errors.Is(err, advisorprofile.ErrNotFound) {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบนักศึกษา"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return
    }

    ctx.JSON(http.StatusOK, students)
}
func (c *AdvisorProfileController) GetStudentBySutID(ctx *gin.Context) {
    advisorRaw, ok := ctx.Get("sut_id")
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    advisorSutID := advisorRaw.(string)
    studentSutID := ctx.Param("sut_id")

    resp, err := c.Service.GetStudentBySutID(advisorSutID, studentSutID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบนักศึกษา"})
        return
    }
    

    ctx.JSON(http.StatusOK, resp)
}
