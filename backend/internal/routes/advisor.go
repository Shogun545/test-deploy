package routes

import (
	"backend/config"
	"backend/internal/app/controller"
	"backend/internal/app/repository"
	"backend/internal/middlewares"
	"backend/internal/service/advisorlog"
	"backend/internal/service/advisorprofile"
	
	"github.com/gin-gonic/gin"
)

func SetupAdvisorRoutes(r *gin.Engine) {
	db := config.DB()

	repo := repository.NewAdvisorProfileRepository(db)
	svc := advisorprofile.NewAdvisorProfileService(repo, db)
	profileCtrl := controller.NewAdvisorProfileController(svc)
	logSvc := advisorlog.New(db, "uploads")
	logCtrl := controller.NewAdvisorLogController(logSvc)
	reportCtrl := controller.NewProgressReportController()
	feedbackCtrl := controller.NewReportFeedbackController()

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	api.GET("/advisor/me/profile", profileCtrl.GetMyAdvisorProfile)
	api.PUT("/advisor/me/profile", profileCtrl.UpdateMyAdvisorProfile)
	api.GET("/advisor/me/students", profileCtrl.GetMyStudents)
	api.GET("/advisor/me/students/:sut_id", profileCtrl.GetStudentBySutID)

	// -------------------------
	// Advisor Logs (เรียงถูกต้อง)
	// -------------------------
	api.GET("/advisor_logs", logCtrl.ListAll)
	api.POST("/advisor_logs", logCtrl.Create)
	api.GET("/advisor_logs/student/:student_id", logCtrl.ListByStudent)
	api.PATCH("/advisor_logs/:id", logCtrl.UpdateStatus)
	api.GET("/advisor_logs/:id", logCtrl.GetByID)
	api.PATCH("/advisor_logs/:id/edit", logCtrl.Update)
	api.GET("/advisor_logs/:id/files/:index", logCtrl.DownloadFile)

	// -------------------------
	// Report & Feedback
	// -------------------------
	api.GET("/progress_reports/log/:log_id", reportCtrl.GetByLogID)
	api.POST("/report_feedbacks", feedbackCtrl.Create)

	
	



	

}
