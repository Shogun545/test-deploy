package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"backend/internal/app/entity"
	service "backend/internal/service/approval"

	"github.com/gin-gonic/gin"
)

type AppointmentController struct {
	service service.AppointmentService
}

// Constructor
func NewAppointmentController(service service.AppointmentService) *AppointmentController {
	return &AppointmentController{service: service}
}

// ---------------------------
// 1) อาจารย์อนุมัติคำขอนัดหมาย (ดึง actor/role จาก token)
// ---------------------------
func (ctr *AppointmentController) ApproveAppointment(c *gin.Context) {
	idStr := c.Param("id")
	AppointmentID64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || AppointmentID64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	actorID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}
	role, ok := getRoleFromContext(c)
	if !ok {
		return
	}
	if role != "ADVISOR" {
		c.JSON(http.StatusForbidden, gin.H{"error": "advisor only"})
		return
	}

	var request struct {
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appt, err := ctr.service.ApproveAppointment(
		uint(AppointmentID64),
		actorID,
		role,
		request.Description,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appt)
}

// ---------------------------
// 2) อาจารย์เสนอเวลาใหม่ (ดึง actor/role จาก token)
// ---------------------------
func (ctr *AppointmentController) ProposeNewTime(c *gin.Context) {
	idStr := c.Param("id")
	AppointmentID64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || AppointmentID64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	actorID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}
	role, ok := getRoleFromContext(c)
	if !ok {
		return
	}
	if role != "ADVISOR" {
		c.JSON(http.StatusForbidden, gin.H{"error": "advisor only"})
		return
	}

	var request struct {
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appt, err := ctr.service.ProposeNewTime(
		uint(AppointmentID64),
		actorID,
		role,
		request.Description,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appt)
}

func (ctr *AppointmentController) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	appt, err := ctr.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appt)
}

// ---------------------------
// 3) list: pending (ดึง advisorID จาก token)
// ---------------------------
func (ctr *AppointmentController) ListPending(c *gin.Context) {
	advisorID, ok := requireAdvisorFromContext(c)
	if !ok {
		return
	}

	appts, err := ctr.service.ListPendingByAdvisor(advisorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapToListRows(appts))
}

// ---------------------------
// 4) list: done (APPROVED + RESCHEDULE) (ดึง advisorID จาก token)
// ---------------------------
func (ctr *AppointmentController) ListDone(c *gin.Context) {
	advisorID, ok := requireAdvisorFromContext(c)
	if !ok {
		return
	}

	appts, err := ctr.service.ListDoneByAdvisor(advisorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapToListRows(appts))
}

// ---------------------------
// 5) list: all (ดึง advisorID จาก token)
// ---------------------------
func (ctr *AppointmentController) ListAll(c *gin.Context) {
	advisorID, ok := requireAdvisorFromContext(c)
	if !ok {
		return
	}

	appts, err := ctr.service.ListAllByAdvisor(advisorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapToListRows(appts))
}

// ---- helpers ----

// requireAdvisorFromContext: เช็ค role=ADVISOR และคืน user_id
func requireAdvisorFromContext(c *gin.Context) (uint, bool) {
	role, ok := getRoleFromContext(c)
	if !ok {
		return 0, false
	}
	if role != "ADVISOR" {
		c.JSON(http.StatusForbidden, gin.H{"error": "advisor only"})
		return 0, false
	}

	uid, ok := getUserIDFromContext(c)
	if !ok {
		return 0, false
	}
	return uid, true
}

// getUserIDFromContext: รองรับ user_id จาก jwt.MapClaims ได้หลายชนิด (float64/string/uint/int)
func getUserIDFromContext(c *gin.Context) (uint, bool) {
	v, ok := c.Get("user_id")
	if !ok || v == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user_id"})
		return 0, false
	}

	switch t := v.(type) {
	case float64:
		if t <= 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
			return 0, false
		}
		return uint(t), true
	case int:
		if t <= 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
			return 0, false
		}
		return uint(t), true
	case uint:
		if t == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
			return 0, false
		}
		return t, true
	case string:
		n, err := strconv.ParseUint(t, 10, 32)
		if err != nil || n == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
			return 0, false
		}
		return uint(n), true
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("invalid user_id type: %T", v)})
		return 0, false
	}
}

func getRoleFromContext(c *gin.Context) (string, bool) {
	v, ok := c.Get("role")
	if !ok || v == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing role"})
		return "", false
	}
	role := strings.ToUpper(fmt.Sprint(v))
	return role, true
}

func mapToListRows(appts []entity.Appointment) []gin.H {
	out := make([]gin.H, 0, len(appts))
	for _, a := range appts {
		out = append(out, gin.H{
			"id":          a.ID,
			"studentName": a.StudentUser.FirstName + " " + a.StudentUser.LastName,
			"sutId":       a.StudentUser.SutId,
			"topic":       a.Description,
			"submittedAt": a.CreatedAt.Format("02/01/2006"),
			"status":      a.AppointmentStatus.StatusCode, // PENDING/APPROVED/RESCHEDULE
		})
	}
	return out
}
