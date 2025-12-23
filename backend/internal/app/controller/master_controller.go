package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"backend/internal/service/master"
)

type MasterController struct {
	prefixService *master.PrefixService
}

func NewMasterController(prefixService *master.PrefixService) *MasterController {
	return &MasterController{
		prefixService: prefixService,
	}
}

func (c *MasterController) GetPrefixes(ctx *gin.Context) {
	prefixes, err := c.prefixService.GetAllPrefixes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot load prefixes",
		})
		return
	}

	ctx.JSON(http.StatusOK, prefixes)
}
