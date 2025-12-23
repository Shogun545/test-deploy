package seed

import (
	"gorm.io/gorm"
	"backend/internal/app/entity"
)

func SeedReportTopic(db *gorm.DB) {
	topics := []entity.ReportTopic{
		{
			ReportTopicName: "Login / Access",
			Description:   "ปัญหาเกี่ยวกับการเข้าสู่ระบบ",
		},
		{
			ReportTopicName: "Booking / Scheduling",
			Description:   "ปัญหาการจองหรือการนัดหมาย",
		},
		{
			ReportTopicName: "Profile / Advisor Data",
			Description:   "ปัญหาข้อมูลโปรไฟล์อาจารย์",
		},
		{
			ReportTopicName: "Notification / Communication",
			Description:   "ปัญหาระบบแจ้งเตือน",
		},
	}
		for _, t := range topics {
		db.FirstOrCreate(&t, entity.ReportTopic{
			ReportTopicName: t.ReportTopicName,
		})
	}
}
