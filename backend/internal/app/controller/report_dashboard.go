package controller

import (
	"net/http"

	"backend/config"
	"backend/internal/app/entity"
	"github.com/gin-gonic/gin"
)

type DashboardSummary struct {
	Total      int64 `json:"total"`
	Pending    int64 `json:"pending"`
	Inprogress int64 `json:"inprogress"`
	Resolved   int64 `json:"resolved"`
}

func GetDashboardSummary(c *gin.Context) {
	db := config.DB()

	var total, pending, inprogress, resolved int64

	db.Model(&entity.Report{}).Count(&total)
	db.Model(&entity.Report{}).Where("report_status_id = ?", 1).Count(&pending)
	db.Model(&entity.Report{}).Where("report_status_id = ?", 2).Count(&resolved)
	db.Model(&entity.Report{}).Where("report_status_id = ?", 3).Count(&inprogress)

	c.JSON(http.StatusOK, DashboardSummary{
		Total:      total,
		Pending:    pending,
		Inprogress: inprogress,
		Resolved:   resolved,
	})
}
