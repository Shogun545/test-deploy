package routes

import (
	"backend/config"
	"backend/internal/app/controller"
	"backend/internal/app/repository"
	"backend/internal/service/master"
	"github.com/gin-gonic/gin"
)

func SetupMasterRoutes(r *gin.Engine) {
	db := config.DB()

	prefixRepo := repository.NewPrefixRepository(db)
	prefixSvc := master.NewPrefixService(prefixRepo)
	masterCtrl := controller.NewMasterController(prefixSvc)

	api := r.Group("/api")
	{
		api.GET("/master/prefixes", masterCtrl.GetPrefixes)
	}
}
