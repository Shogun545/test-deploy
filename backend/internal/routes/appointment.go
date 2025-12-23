// internal/routes/appointment_routes.go
package routes

import (
	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/internal/app/controller"
	"backend/internal/app/repository"
	approvalService "backend/internal/service/approval"

	middleware "backend/internal/middlewares"
)

// route ของระบบนัดหมาย (ต้อง login แล้ว)
func SetupAppointmentRoutes(r *gin.Engine) {
	db := config.DB()

	apptRepo := repository.NewAppointmentRepository(db)
	apptService := approvalService.NewAppointmentService(apptRepo)
	apptController := controller.NewAppointmentController(apptService)

	appointments := r.Group("/api/appointments")
	appointments.Use(middleware.AuthMiddleware())

	// ✅ list: ดึงตามบัญชีที่ login (ไม่ต้องส่ง advisor_id)
	appointments.GET("", apptController.ListAll)       // GET /api/appointments
	appointments.GET("/pending", apptController.ListPending)
	appointments.GET("/done", apptController.ListDone)

	// ✅ detail/approve/reschedule
	appointments.GET("/:id", apptController.GetByID)
	appointments.PUT("/:id/approve", apptController.ApproveAppointment)
	appointments.PUT("/:id/reschedule", apptController.ProposeNewTime)
}
