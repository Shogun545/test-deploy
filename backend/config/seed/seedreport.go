package seed

import (
	"log"
	"math/rand"
	"time"

	"backend/internal/app/entity"
	"gorm.io/gorm"
)

func SeedReport(db *gorm.DB) {

	rand.Seed(time.Now().UnixNano())

	// ------------------------------
	// เช็กก่อน: ถ้ามี report แล้ว ไม่ต้อง seed ซ้ำ
	// ------------------------------
	var reportCount int64
	db.Model(&entity.Report{}).Count(&reportCount)
	if reportCount > 0 {
		log.Println("ℹ️ มี Report อยู่แล้ว ข้ามการ SeedReport")
		return
	}

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

	// map status name → status
	statusMap := make(map[string]entity.ReportStatus)
	for _, s := range statuses {
		statusMap[s.ReportStatusName] = s
	}

	requiredStatuses := []string{"Pending", "Inprogress", "Resolved"}
	for _, name := range requiredStatuses {
		if _, ok := statusMap[name]; !ok {
			log.Printf("❌ ไม่พบ ReportStatus: %s\n", name)
			return
		}
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
	// สร้าง Reports (ใช้ transaction)
	// ------------------------------
	err := db.Transaction(func(tx *gorm.DB) error {

		for i := 0; i < 10; i++ {

			reporter := reporters[rand.Intn(len(reporters))]
			topic := topics[rand.Intn(len(topics))]

			statusNames := []string{"Pending", "Inprogress", "Resolved"}
			status := statusMap[statusNames[rand.Intn(len(statusNames))]]

			report := entity.Report{
				Description:  descriptions[rand.Intn(len(descriptions))],
				ReportByID:     reporter.ID, // FK user
				ReportStatusID: status.ID,   // FK status
				ReportTopicID:  topic.ID,    // FK topic
			}

			if err := tx.Create(&report).Error; err != nil {
				return err // rollback
			}

			log.Printf("✅ Seed Report ID=%d by UserID=%d\n", report.ID, reporter.ID)
		}

		return nil
	})

	if err != nil {
		log.Println("❌ SeedReport ล้มเหลว (rollback):", err)
		return
	}

	log.Println("Seed Reports เสร็จสมบูรณ์ ✅")
}
