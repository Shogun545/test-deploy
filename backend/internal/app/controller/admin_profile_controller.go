package controller

import (
    "backend/internal/service/adminprofile" 
    "errors"
    "github.com/gin-gonic/gin"
    "net/http"
    "backend/internal/app/dto"
)

type AdminProfileController struct {
    Service *adminprofile.AdminProfileService
}

func NewAdminProfileController(s *adminprofile.AdminProfileService) *AdminProfileController {
    return &AdminProfileController{Service: s}
}

// GetMyAdminProfile ดึงข้อมูลโปรไฟล์ของ Admin ที่กำลังล็อกอิน
func (c *AdminProfileController) GetMyAdminProfile(ctx *gin.Context) {
    // ดึง sut_id (Admin ID) จาก Context (สมมติว่า Middleware จัดการการยืนยันตัวตนแล้ว)
    sutIdRaw, ok := ctx.Get("sut_id")
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    adminSutID, _ := sutIdRaw.(string)

    // เรียก Service เพื่อดึงโปรไฟล์
    resp, err := c.Service.GetMyProfile(adminSutID)
    if err != nil {
        if errors.Is(err, adminprofile.ErrNotFound) {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบผู้ดูแลระบบ"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: failed to fetch profile"})
        return
    }

    ctx.JSON(http.StatusOK, resp)
}

// GetManagedUsers ดึงรายการผู้ใช้งานทั้งหมดที่ Admin ดูแล (Students/Advisors)
func (c *AdminProfileController) GetManagedUsers(ctx *gin.Context) {
    filters := dto.UserListFilter{
        Role:      ctx.Query("role"),
        Status:    ctx.Query("status"),
        Department: ctx.Query("department"),
        Search:    ctx.Query("search"),
		Major:    ctx.Query("major"),
    }
    users, err := c.Service.GetAllUsers(filters) 
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: failed to fetch users"})
        return
    }

    ctx.JSON(http.StatusOK, users)
}

// GetUserBySutID: ดึงข้อมูลผู้ใช้งานคนใดคนหนึ่งในระบบ โดย Admin ระบุ SutID ของผู้ใช้งานนั้น
func (c *AdminProfileController) GetUserBySutID(ctx *gin.Context) {
    // Admin ID ที่ล็อกอิน (optional: ใช้เพื่อตรวจสอบสิทธิ์การเข้าถึง)
    targetSutID := ctx.Param("sut_id") // SutID ของ User ที่ต้องการดูข้อมูล
    resp, err := c.Service.GetUserBySutID(targetSutID) // เรียก Service เพื่อดึงข้อมูล User เป้าหมาย
    if err != nil {
        if errors.Is(err, adminprofile.ErrNotFound) {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบผู้ใช้งานที่ระบุ"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return
    }

    ctx.JSON(http.StatusOK, resp)
}
// UpdateMyAdminProfile เป็น Handler Function สำหรับการอัปเดต Profile ของ Admin ที่เข้าสู่ระบบ
// Endpoint PUT /api/v1/admin/profile
func (c *AdminProfileController) UpdateMyAdminProfile(ctx *gin.Context) {
    // ดึง sut_id (Admin ID) จาก Context
    sutIdRaw, ok := ctx.Get("sut_id")
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    adminSutID, ok := sutIdRaw.(string)
    if !ok || adminSutID == "" {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid sut_id in context"})
        return
    }

    // ใช้ dto type ตรง ๆ เพื่อ bind JSON
    var req dto.UpdateAdminProfileRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data", "details": err.Error()})
        return
    }

    // เรียก Service — ตรวจสอบให้แน่ใจว่า service มี signature แบบนี้:
    // func (s *AdminProfileService) UpdateMyProfile(sutID string, req *dto.UpdateAdminProfileRequest) (/*ผลลัพธ์*/, error)
    updatedProfile, err := c.Service.UpdateMyProfile(adminSutID, &req)
    if err != nil {
        if errors.Is(err, adminprofile.ErrNotFound) {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "admin profile not found"})
            return
        }
        if err.Error() == "email already exists" {
            ctx.JSON(http.StatusConflict, gin.H{"error": "email is already in use"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile", "details": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Admin profile updated successfully",
        "data":    updatedProfile,
    })
}

// func (c *AdminProfileController) UpdateUser(ctx *gin.Context) { ... }
func (c *AdminProfileController) GetMajors(ctx *gin.Context) {
    majors, err := c.Service.GetMasterMajors() 
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch majors"})
        return
    }

    ctx.JSON(http.StatusOK, majors)
}

// UpdateUserStatus อัปเดตสถานะของผู้ใช้งานคนใดคนหนึ่งในระบบ
func (c *AdminProfileController) UpdateUserStatus(ctx *gin.Context) {
    targetSutID := ctx.Param("sut_id") 
    var req dto.UpdateUserStatusRequest // สมมติว่ามี DTO นี้
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data", "details": err.Error()})
        return
    }
    err := c.Service.UpdateUserStatus(targetSutID, req.Status) 
    if err != nil {
        if errors.Is(err, adminprofile.ErrNotFound) {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบผู้ใช้งานที่ระบุ"})
            return
        }
        if errors.Is(err, adminprofile.ErrInvalidStatus) { 
            ctx.JSON(http.StatusBadRequest, gin.H{"error": "สถานะไม่ถูกต้อง"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user status", "details": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{
        "message": "User status updated successfully",
        "sut_id":  targetSutID,
        "new_status": req.Status,
    })
}

func (c *AdminProfileController) GetUserCreatedDate(ctx *gin.Context) {
	targetSutID := ctx.Param("sut_id")

	resp, err := c.Service.GetUserCreatedDate(targetSutID)
	if err != nil {
		if errors.Is(err, adminprofile.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบผู้ใช้งานที่ระบุ"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch created date"})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// PUT /api/admin/users/:sut_id
func (c *AdminProfileController) UpdateUser(ctx *gin.Context) {
	targetSutID := ctx.Param("sut_id")

	var req dto.UpdateManagedUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
			"details": err.Error(),
		})
		return
	}

	err := c.Service.UpdateUser(targetSutID, &req)
	if err != nil {
		if errors.Is(err, adminprofile.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบผู้ใช้งาน"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
	})
}
