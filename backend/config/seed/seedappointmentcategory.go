package seed

import (
	"log"

	"gorm.io/gorm"

	"backend/internal/app/entity"
)

func SeedAppointmentCategories(db *gorm.DB) {
	categories := []entity.AppointmentCategory{
		{Category: "ขอคำปรึกษาทั่วไป"},
		{Category: "ปรึกษาเรื่องการเรียน"},
		{Category: "ปรึกษาเรื่องโปรเจกต์"},
		{Category: "ปรึกษาฝึกงาน / สหกิจ"},
		{Category: "อื่น ๆ"},
	}

	for _, c := range categories {
		var existing entity.AppointmentCategory
		err := db.
			Where("category = ?", c.Category).
			First(&existing).Error

		// ถ้ายังไม่มี → create
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&c).Error; err != nil {
				log.Printf("❌ Seed AppointmentCategory failed: %v\n", err)
			} else {
				log.Printf("✅ Seeded AppointmentCategory: %s\n", c.Category)
			}
		}
	}
}
