// internal/routes/auth_routes.go
package routes

import (
	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/internal/app/controller"
	"backend/internal/app/repository"
	userservice "backend/internal/service/users"
)

// route ที่เกี่ยวกับ auth (ส่วนที่ไม่ต้อง login เช่น /login)
func SetupAuthRoutes(r *gin.Engine) {
	db := config.DB()

	userRepo := repository.NewUserRepository(db)
	authService := userservice.NewAuthService(userRepo)
	authController := controller.NewAuthController(authService)

	// login
	r.POST("/api/auth/login", authController.Login)
	// ถ้าจะมี register เพิ่มค่อยใส่ตรงนี้ก็ได้
}
