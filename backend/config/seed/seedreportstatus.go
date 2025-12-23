package seed

import (
	"log"

	"gorm.io/gorm"
	"backend/internal/app/entity"
)

func SeedReportStatus(db *gorm.DB) {
	reportstatus := []entity.ReportStatus{
		{ReportStatusName: "Pending"},
		{ReportStatusName: "Resolved"},
		{ReportStatusName: "Inprogress"},
	}
	for _, reportstatus := range reportstatus {
		if err := db.FirstOrCreate(&reportstatus, &entity.ReportStatus{ReportStatusName: reportstatus.ReportStatusName}).Error; err != nil {
			log.Fatalf("failed to seed reportstatus %s: %v", reportstatus.ReportStatusName, err)
		}
	}
}
