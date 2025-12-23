package routes

import (
    "github.com/gin-gonic/gin"
    "backend/config"
    "backend/internal/app/controller" 
    "backend/internal/app/repository"
    "backend/internal/service/studentprofile"
    "backend/internal/middlewares"
)

func SetupStudentRoutes(r *gin.Engine) {
    db := config.DB()

    repo := repository.NewProfileRepository(db)
    svc := service.NewProfileService(repo)
    profileCtrl := controller.NewProfileController(svc)
    reportCtrl := controller.NewProgressReportController()

    api := r.Group("/api")
    api.Use(middleware.AuthMiddleware())

    api.GET("/student/me/profile", profileCtrl.GetMyStudentProfile)
    api.PUT("/student/me/profile", profileCtrl.UpdateMyStudentProfile)

    api.POST("/progress_reports", reportCtrl.Create)
}