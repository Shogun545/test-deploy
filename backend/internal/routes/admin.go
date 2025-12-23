package routes

import (
	"backend/config"
	"backend/internal/app/controller"
	"backend/internal/app/repository"
	"backend/internal/middlewares"
	"backend/internal/service/academiccalendar" // หรือ backend/internal/app/service แล้วแต่โครงสร้างจริง
	"backend/internal/service/adminprofile"     // ใช้แพ็กเกจ service ของ admin
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(r *gin.Engine) {
	// 1. เชื่อมต่อ Database
	db := config.DB()

	// 2. Initialize Dependency Injection (เรียงลำดับการส่งต่อ)
	academicRepo := repository.NewAcademicCalendarRepository(db)
	academicService := service.NewAcademicCalendarService(academicRepo)
	academicCtrl := controller.NewAcademicCalendarController(academicService)

	adminRepo := repository.NewAdminProfileRepositoryImpl(db)
	adminSvc := adminprofile.NewAdminProfileService(adminRepo, db)
	adminCtrl := controller.NewAdminProfileController(adminSvc)

	// 3. สร้าง Group Route
	api := r.Group("/api")
	api.GET("/holidays", academicCtrl.GetHolidaysByYear)

	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/events", academicCtrl.CreateEvent)
		api.PUT("/events/:id", academicCtrl.UpdateEvent)
		api.DELETE("/events/:id", academicCtrl.DeleteEvent)

		api.GET("/admin/me/profile", adminCtrl.GetMyAdminProfile)
		api.PUT("/admin/me/profile", adminCtrl.UpdateMyAdminProfile)

		//NEW: Admin Management Endpoints
		api.GET("/admin/users", adminCtrl.GetManagedUsers)        // เพิ่มตรงนี้
		api.GET("/admin/users/:sut_id", adminCtrl.GetUserBySutID) // ดึงรายละเอียดผู้ใช้รายคน
		api.GET("/admin/majors", adminCtrl.GetMajors)             // ดึงรายการสาขาวิชา
		api.PUT("/admin/users/:sut_id/status", adminCtrl.UpdateUserStatus)
		api.GET("/admin/users/:sut_id/created-date", adminCtrl.GetUserCreatedDate)
		api.PUT("/admin/users/:sut_id", adminCtrl.UpdateUser) // อัปเดตข้อมูลผู้ใช้ที่ Admin ดูแล

	}
}
