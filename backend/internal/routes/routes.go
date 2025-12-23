package routes

import (
	"github.com/gin-gonic/gin"
	"backend/internal/app/controller"
)


// ฟังก์ชันรวม route ทั้งหมด
func SetupRoutes(r *gin.Engine) {
	// ลำดับไหนก่อนหลังก็ได้ แต่เพื่อความชัดก็ตามนี้
	SetupAuthRoutes(r)         // /api/auth/login (ไม่ต้อง login)
	SetupStudentRoutes(r)      // /api/student/me/profile (ต้อง login)
	SetupAdvisorRoutes(r)      // /api/advisor/me/profile (ต้อง login)
	SetupAppointmentRoutes(r)  // /api/appointments/... (ต้อง login)
	SetupAdminRoutes(r)
	SetupMasterRoutes(r)

	// ===== Report =====	
	r.GET("/reports", controller.GetAllReport)
	r.GET("/reports/:id", controller.GetReport)
	r.POST("/reports", controller.CreateReport)
	r.PUT("/reports/:id", controller.UpdateReport)
	r.DELETE("/reports/:id", controller.DeleteReport)

	// ===== Report Status =====
	r.GET("/report-status", controller.GetAllReportStatus)
	r.GET("/report-status/:id", controller.GetReportStatus)
	r.POST("/report-status", controller.CreateReportStatus)
	r.PUT("/report-status/:id", controller.UpdateReportStatus)
	r.DELETE("/report-status/:id", controller.DeleteReportStatus)
	
	// ===== Report Topic =====
	r.GET("/report-topics", controller.GetAllReportTopic)
	r.GET("/report-topics/:id", controller.GetReportTopic)
	r.POST("/report-topics", controller.CreateReportTopic)
	r.PUT("/report-topics/:id", controller.UpdateReportTopic)
	r.DELETE("/report-topics/:id", controller.DeleteReportTopic)

	r.GET("/report/summary", controller.GetDashboardSummary)
}
