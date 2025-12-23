package seed

import (
	"log"
	"math/rand"
	"time"

	"backend/internal/app/entity"
	"gorm.io/gorm"
)

func SeedReportttt(db *gorm.DB) {

	rand.Seed(time.Now().UnixNano())

	// ------------------------------
	// โหลด Users ที่สามารถ report ได้ (Student)
	// ------------------------------
	var reporters []entity.User
	if err := db.
		Joins("JOIN roles r ON r.id = users.role_id").
		Where("r.role = ?", "Student").
		Find(&reporters).Error; err != nil || len(reporters) == 0 {

		log.Println("❌ ไม่พบ Student สำหรับสร้าง Report")
		return
	}

	// ------------------------------
	// โหลด ReportStatus
	// ------------------------------
	var statuses []entity.ReportStatus
	if err := db.Find(&statuses).Error; err != nil || len(statuses) == 0 {
		log.Println("❌ ไม่พบ ReportStatus")
		return
	}

	// helper map
	statusMap := make(map[string]entity.ReportStatus)
	for _, s := range statuses {
		statusMap[s.ReportStatusName] = s
	}

	// ------------------------------
	// โหลด ReportTopic
	// ------------------------------
	var topics []entity.ReportTopic
	if err := db.Find(&topics).Error; err != nil || len(topics) == 0 {
		log.Println("❌ ไม่พบ ReportTopic")
		return
	}

	// ------------------------------
	// เตรียมข้อความตัวอย่าง
	// ------------------------------
	descriptions := []string{
		"ไม่สามารถเข้าสู่ระบบได้",
		"จองเวลานัดหมายแล้วไม่ขึ้นในระบบ",
		"ข้อมูลอาจารย์แสดงผลไม่ถูกต้อง",
		"ไม่ได้รับการแจ้งเตือนจากระบบ",
		"หน้าโปรไฟล์โหลดช้า",
	}

	// ------------------------------
	// สร้าง Reports
	// ------------------------------
	for i := 0; i < 10; i++ {

		reporter := reporters[rand.Intn(len(reporters))]
		topic := topics[rand.Intn(len(topics))]

		// สุ่ม status
		statusNames := []string{"Pending", "Inprogress", "Resolved"}
		status := statusMap[statusNames[rand.Intn(len(statusNames))]]

		report := entity.Report{
			Description:  descriptions[rand.Intn(len(descriptions))],
			ReportByID:     reporter.ID,   // ✅ user ที่มีจริง
			ReportStatusID: status.ID,     // ✅ FK ถูก
			ReportTopicID:  topic.ID,      // ✅ FK ถูก
		}

		if err := db.Create(&report).Error; err != nil {
			log.Println("❌ SeedReport error:", err)
		} else {
			log.Printf("✅ Seed Report ID=%d by UserID=%d\n", report.ID, reporter.ID)
		}
	}

	log.Println("Seed Reports เสร็จสมบูรณ์ ✅")
}
