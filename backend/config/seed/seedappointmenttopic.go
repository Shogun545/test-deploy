package seed

import (
	"log"

	"gorm.io/gorm"

	"backend/internal/app/entity"
)

func SeedAppointmentTopics(db *gorm.DB) {
	topics := []entity.AppointmentTopic{
		{
			Topic:       "ปรึกษาเรื่องการเรียน",
			Description: "ให้คำปรึกษาเกี่ยวกับแผนการเรียน การลงทะเบียน และผลการเรียน",
			IsActive:    true,
		},
		{
			Topic:       "ปรึกษาโปรเจกต์",
			Description: "ขอคำแนะนำเกี่ยวกับโปรเจกต์รายวิชาและโปรเจกต์จบ",
			IsActive:    true,
		},
		{
			Topic:       "ปรึกษาฝึกงาน / สหกิจ",
			Description: "ให้คำแนะนำเรื่องการเตรียมตัวฝึกงานและสหกิจศึกษา",
			IsActive:    true,
		},
		{
			Topic:       "ปัญหาส่วนตัว / การใช้ชีวิตในมหาวิทยาลัย",
			Description: "พูดคุยและขอคำแนะนำเรื่องการใช้ชีวิตและการปรับตัว",
			IsActive:    true,
		},
		{
			Topic:       "อื่น ๆ",
			Description: "หัวข้ออื่นที่ไม่อยู่ในรายการ",
			IsActive:    true,
		},
	}

	for _, t := range topics {
		var existing entity.AppointmentTopic
		err := db.
			Where("topic = ?", t.Topic).
			First(&existing).Error

		// ถ้ายังไม่มี → create
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&t).Error; err != nil {
				log.Printf("❌ Seed AppointmentTopic failed: %v\n", err)
			} else {
				log.Printf("✅ Seeded AppointmentTopic: %s\n", t.Topic)
			}
		}
	}
}
