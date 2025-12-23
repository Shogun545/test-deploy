package controller

import (
    "fmt"
    "net/http"
    "strconv"
    "strings"

    "backend/internal/app/dto"
    "backend/internal/service/advisorlog"
    "github.com/gin-gonic/gin"
)


type AdvisorLogController struct {
	svc advisorlog.Service
}

func NewAdvisorLogController(svc advisorlog.Service) *AdvisorLogController {
	return &AdvisorLogController{
		svc: svc,
	}
}
// ==========================================
// 🛡️ Helper: ดึง User Info (ฉบับแก้ปัญหา ID=0 โดยไม่ง้อ Middleware)
// ==========================================
func getUserFromContext(c *gin.Context) (uint, string) {
	var userID uint
	
	// 1. ดึงค่าดิบๆ ออกมาก่อน (ยังไม่รู้ว่าเป็น Type อะไร)
	v, exists := c.Get("user_id")
	
	if exists {
		// 2. เช็ค Type และแปลงให้ถูก
		switch val := v.(type) {
		case float64: 
			userID = uint(val) 
		case uint:
			userID = val
		case int:
			userID = uint(val)
		case int64:
			userID = uint(val)
		case string:
			// เผื่อเพื่อนส่งมาเป็น String
			if n, err := strconv.Atoi(val); err == nil {
				userID = uint(n)
			}
		default:
			fmt.Printf("[Check] user_id เป็น Type แปลกๆ: %T\n", val)
		}
	} 

	role := c.GetString("role")
	return userID, role
}

// ------------------------------
// CREATE LOG (multipart/form-data)
// ------------------------------
func (ctrl *AdvisorLogController) Create(c *gin.Context) {
	// 1. ดึง User ID และ Role จาก Context
	userID, role := getUserFromContext(c)

	const maxMemory = 32 << 20
	_ = c.Request.ParseMultipartForm(maxMemory)

	var req dto.AdvisorLogCreateReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if form, _ := c.MultipartForm(); form != nil {
		req.Files = form.File["files"]
	}

	// 2. ✅ ส่ง userID และ role ไปให้ Service
	out, err := ctrl.svc.Create(c.Request.Context(), req, userID, role)
	
	if err != nil {
		if err == advisorlog.ErrForbidden {
			 c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this appointment (neither student nor advisor)"})
			 return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": out})
}

// ------------------------------
// GET BY ID (🔒 แก้ IDOR + Draft Leakage)
// ------------------------------
func (ctrl *AdvisorLogController) GetByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    // 1. ดึง User Info
    userID, role := getUserFromContext(c)

    // 2. ✅ เรียก Service แบบใหม่ (ส่ง userID, role เข้าไป)
    // ให้ Service เป็นคนตัดสินใจเองว่า "เจอ" หรือ "ไม่เจอ" (Forbidden/NotFound)
    out, err := ctrl.svc.GetByID(c.Request.Context(), uint(id), userID, role)
    
    if err != nil {
        switch err {
        case advisorlog.ErrAdvisorLogNotFound:
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        case advisorlog.ErrForbidden:
            c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": out})
}

// ------------------------------
// LIST BY STUDENT (🔒 แก้ไขแล้ว: Clean & Secure)
// ------------------------------
func (ctrl *AdvisorLogController) ListByStudent(c *gin.Context) {
	// 1. ดึง User Info
	requesterID, role := getUserFromContext(c)
	
	var targetStudentID uint

	// 2. Logic เลือก Target ID (ยังคงไว้ เพื่อกัน IDOR)
	if role == "student" {
		targetStudentID = requesterID // นักศึกษาดูได้แค่ของตัวเอง
	} else {
		// อาจารย์ดูของใครก็ได้
		paramID, err := strconv.Atoi(c.Param("student_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student_id"})
			return
		}
		targetStudentID = uint(paramID)
	}

	// 3. ✅ เรียก Service: ส่ง role ไปด้วย!
	// Service จะเอา role ไปเช็ค: ถ้าเป็น student -> เติม WHERE status != 'Draft' ให้เอง
	out, err := ctrl.svc.ListByStudent(c.Request.Context(), targetStudentID, role)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. ส่งข้อมูลกลับได้เลย (ไม่ต้องมานั่งวนลูปกรองเองแล้ว)
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// ------------------------------
// LIST ALL (advisor)
// ------------------------------
func (ctrl *AdvisorLogController) ListAll(c *gin.Context) {
	out, err := ctrl.svc.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": out})
}

// ------------------------------
// UPDATE STATUS ONLY
// ------------------------------
func (ctrl *AdvisorLogController) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.AdvisorLogUpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.svc.UpdateStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		switch err {
		case advisorlog.ErrInvalidStatus:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case advisorlog.ErrAdvisorLogNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// ------------------------------
// UPDATE FULL LOG (Title, Body, RequiresReport, Files)
// ------------------------------
// ------------------------------
// UPDATE FULL LOG (Title, Body, RequiresReport, Files)
// ------------------------------
func (ctrl *AdvisorLogController) Update(c *gin.Context) {
    // 1. ดึง User Info
    userID, role := getUserFromContext(c)
    
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    const maxMemory = 32 << 20
    _ = c.Request.ParseMultipartForm(maxMemory)

    var req dto.AdvisorLogUpdateReq

    // Text fields (optional)
    if v := c.PostForm("title"); v != "" {
        req.Title = &v
    }
    if v := c.PostForm("body"); v != "" {
        req.Body = &v
    }
    if v := c.PostForm("requiresReport"); v != "" {
        b, err := strconv.ParseBool(v)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requiresReport"})
            return
        }
        req.RequiresReport = &b
    }

    // Files (optional)
    if form, _ := c.MultipartForm(); form != nil && len(form.File["files"]) > 0 {
        req.Files = form.File["files"]
    }

    // 2. ✅ ส่ง userID และ role ไปให้ Service เช็คสิทธิ์
    out, err := ctrl.svc.Update(c.Request.Context(), uint(id), req, userID, role)
    
    if err != nil {
        switch err {
        case advisorlog.ErrAdvisorLogNotFound:
            c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
        case advisorlog.ErrForbidden:
            // 3. ✅ เพิ่ม Case นี้: เมื่อพยายามแก้ของคนอื่น
            c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
        case advisorlog.ErrSaveFileFailed:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "save file failed"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "updated successfully",
        "data":    out,
    })
}
// ------------------------------
// DOWNLOAD FILE (🔒 แก้ไข Bug ประกาศตัวแปรซ้ำ)
// ------------------------------
func (ctrl *AdvisorLogController) DownloadFile(c *gin.Context) {
	// รับ id, index
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	idx, err := strconv.Atoi(c.Param("index"))
	if err != nil || idx < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid index"})
		return
	}

	sutAny, ok := c.Get("sut_id") 
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing sut_id in context"})
		return
	}
	sutID := strings.TrimSpace(fmt.Sprint(sutAny)) // ใช้ตัวแปรนี้ส่งไป

	fileName, absPath, err := ctrl.svc.GetFileForLog(c.Request.Context(), uint(id), idx, sutID)
	if err != nil {
		switch err {
		case advisorlog.ErrAdvisorLogNotFound, advisorlog.ErrFileNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case advisorlog.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.FileAttachment(absPath, fileName)
}
